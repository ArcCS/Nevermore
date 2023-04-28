// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package comms

import (
	"bufio"
	"bytes"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
	"unicode"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/frontend"
	"github.com/ArcCS/Nevermore/text"
)

const (
	termColumns = 160
	termLines   = 24
	inputBuffer = 1024
)

// This interface lets us assert network or our own errors
type temporary interface {
	Temporary() bool
}

// client contains state information about a client connection. The err field
// should not be manipulated directly. Instead call Error() and SetError().
//
// The current frontend in use is an anonymous interface as this lets us define
// what type frontend is - even though we don't have access to the unexported
// frontend struct from the frontend package.
type client struct {
	*net.TCPConn            // The client's network connection
	remoteAddr   string     // Client's remote address
	err          chan error // Error channel to sync between input & output

	frontend interface { // The current frontend in use
		Parse([]byte) error
		Close()
		GetCharacter() *objects.Character
	}
}

// newClient returns an initialised client for the passed connection.
func newClient(conn *net.TCPConn) *client {

	// Setup connection parameters
	conn.SetKeepAlive(true)
	conn.SetLinger(10)
	conn.SetNoDelay(false)
	conn.SetWriteBuffer(termColumns * termLines)
	conn.SetReadBuffer(inputBuffer)

	c := &client{
		TCPConn:    conn,
		remoteAddr: conn.RemoteAddr().String(),
		err:        make(chan error, 1),
	}

	c.err <- nil

	log.Printf("Acquired lease: %s", conn.RemoteAddr())
	c.leaseAcquire()

	// Setup frontend if no error acquiring a lease

	if c.Error() == nil {

		c.frontend = frontend.New(c, c.remoteAddr)
		c.frontend.Parse([]byte(""))
	}

	return c
}

// process handles input from the network connection.
func (c *client) process() {

	// If a client goroutine panics try not to bring down the whole server down
	// unless the configuration option Debug.Panic is set to true.
	defer func() {
		if !config.Debug.Panic {
			if err := recover(); err != nil {
				log.Printf("CLIENT PANICKED: %s", c.remoteAddr)
				log.Printf("%s: %s", err, debug.Stack())
			}
		}
		log.Print("Ending game loop: ", c.RemoteAddr(), " for character: ", c.frontend.GetCharacter().Name)
		c.close()
	}()

	// Main input processing loop, terminates on any error raised not just read
	// or Parse errors.
	{
		// Variables for use in the loop only hence the scoping outer braces
		var (
			s   = bufio.NewReaderSize(c, inputBuffer) // Sized network read buffer
			err error                                 // function local errors
			in  []byte                                // Input string from buffer
		)

		log.Print("Starting game loop: ", c.RemoteAddr())
		for c.Error() == nil {
			if config.Server.Running == false {
				_ = c.Close()
			}
			useTime := config.Server.IdleTimeout
			if ok := c.frontend.GetCharacter(); ok != nil {
				if c.frontend.GetCharacter().Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					useTime = 30 * time.Minute
				} else if c.frontend.GetCharacter().Flags["AFK"] {
					useTime = config.Server.AFKTimeout
				} else if c.frontend.GetCharacter().Flags["OOC"] {
					useTime = config.Server.OOCTimeout
				}
			}
			err = c.SetReadDeadline(time.Now().Add(useTime))
			if in, err = s.ReadSlice('\n'); err != nil {
				frontend.Zero(in)

				if err != bufio.ErrBufferFull {
					c.SetError(err)
					continue
				}

				for err == bufio.ErrBufferFull {
					in, err = s.ReadSlice('\n')
					frontend.Zero(in)
				}
				_, _ = c.Write([]byte(text.Bad + "\nYou type too much.\n" + text.Prompt + ">"))
				continue
			}

			//log.Println(&in)
			fixDEL(&in)
			if err = c.frontend.Parse(in); err != nil {
				c.SetError(err)
			}
			frontend.Zero(in)
		}
	}
}

// fixDEL is used to delete characters when the input contains literal DEL
// characters (ASCII 0x7f or "\b"). This is the case when using a client that
// does not support line editing, for example a plain Windows TELNET client.
//
// For example if you type "ABD" then delete the "D" and enter "C" the data
// sent to the server would be "ABD\bC" if there is no line editing support.
// With line editing "ABC" would be sent to the server.
//
// Calling fixDEL on the data will interpret the DEL characters so that, for
// example, "ABD\bC" becomes "ABC".
//
// fixDEL can work on ASCII or UTF-8 and handles Unicode diacritics in addition
// to precomposed characters. For example 'à' or 'a\u0300'.
//
// It should be noted that this function modifies the slice passed to it.
func fixDEL(in *[]byte) {

	i := 0
	for j, v := range *in {
		(*in)[j] = '\x00'
		if v != '\b' {
			(*in)[i] = v
			i++
			continue
		}

		// Remove previous rune which may be Unicode, maybe combining diacritic
		for l, combi := 0, true; combi == true; {
			switch {
			case i > 0 && (*in)[i-1]&128 == 0:
				l, combi = 1, false
			case i > 1 && (*in)[i-2]&192 == 192:
				l = 2
			case i > 2 && (*in)[i-3]&192 == 192:
				l = 3
			case i > 3 && (*in)[i-4]&192 == 192:
				l = 4
			default:
				l, combi = 0, false
			}
			if l == 1 {
				(*in)[i-1] = '\x00'
			}
			if l > 1 {
				combi = unicode.In(bytes.Runes((*in)[i-l : i])[0], unicode.Mn, unicode.Me)
				copy((*in)[i-l:i], "\x00\x00\x00\x00")
			}
			i = i - l
		}
	}

	*in = (*in)[:i]

	return
}

// close shuts down a client cleanly, closes network connections and
// deallocates resources.
func (c *client) close() {

	idle, busy := false, false

	// Idle timeout?
	if oe, ok := c.Error().(*net.OpError); ok && oe.Timeout() {
		idle = true
	}

	// Server busy?
	if _, ok := c.Error().(noLeaseError); ok {
		busy = true
	}

	// Deallocate current frontend if we have one
	if c.frontend != nil {
		// Sometimes these disconnects are a little messy,  need to add some extra cleanup
		if c.frontend.GetCharacter() != nil {
			c.frontend.GetCharacter().Save()
			c.frontend.GetCharacter().Unfollow()
			c.frontend.GetCharacter().LoseParty()
			c.frontend.GetCharacter().PurgeEffects()
			objects.ActiveCharacters.Remove(c.frontend.GetCharacter())
			objects.Rooms[c.frontend.GetCharacter().ParentId].Chars.Remove(c.frontend.GetCharacter())
			c.frontend.GetCharacter().Unload()
		}

		if idle {
			_, _ = c.Write([]byte("\n")) // Move off prompt line
		}
		c.frontend.Close()
		c.frontend = nil
	}

	// If connection timed out notify the client
	if idle {
		_, _ = c.Write([]byte(text.Bad + "\nIdle connection terminated by server.\n"))
	}

	// Notify if server too busy to accept more players
	if busy {
		_, _ = c.Write([]byte(text.Bad + "\nServer too busy. Please come back in a short while.\n"))
	}

	// Say goodbye to client
	_, _ = c.Write([]byte(text.Info + "\nBye bye...\n\n"))

	// Revert to default colors
	_, _ = c.Write([]byte(text.Reset))

	// io.EOF does not give address info so handle specially, otherwise just
	// report the error
	if c.Error() == io.EOF {
		log.Printf("Connection dropped by remote client: %s", c.remoteAddr)
	} else {
		log.Printf("Connection error: %s, %s", c.Error(), c.remoteAddr)
	}

	// Make sure connection closed down and deallocated
	if err := c.Close(); err != nil {
		log.Printf("Error closing connection: %s", err)
	} else {
		log.Printf("Connection closed: %s", c.remoteAddr)
	}
	c.TCPConn = nil
	c.leaseRelease()

	// Close and drain error channel
	close(c.err)
	<-c.err
}

// Write handles output for the network connection.
func (c *client) Write(d []byte) (n int, err error) {

	// If we already have a non-temporary error do nothing
	if e := c.Error(); e != nil {
		if e, ok := e.(temporary); !ok || !e.Temporary() {
			return
		}
	}

	/*var t []byte

	if len(d) != 0 {
		t = text.Fold(d, termColumns)
	}*/

	if n, err = c.TCPConn.Write(d); err != nil {
		c.SetError(err)
	}
	return
}

// Error returns the first error raised or nil if there is no error. An error
// can be set by calling SetError.
func (c *client) Error() (err error) {
	err = <-c.err
	c.err <- err
	return err
}

// SetError is used to record the first error condition that occurs. Subsequent
// calls will not over write the initial error raised. The current error can be
// checked by calling Error.
func (c *client) SetError(err error) {
	e := <-c.err
	if e == nil {
		e = err
	}
	c.err <- e
}
