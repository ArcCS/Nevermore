// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package comms

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
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

// client contains state information about a client connection. The err field
// should not be manipulated directly. Instead call Error() and SetError().
//
// The current frontend in use is an anonymous interface as this lets us define
// what type frontend is - even though we don't have access to the unexported
// frontend struct from the frontend package.
type client struct {
	TCPConn    *net.TCPConn // The client's network connection
	remoteAddr string       // Client's remote address
	err        error        // Last error encountered

	frontend interface { // The current frontend in use
		Parse([]byte) error
		Close()
		GetCharacter() *objects.Character
		AccountCleanup()
	}
}

func (c *client) WriteError(err error) {
	c.err = err
}

// newClient returns an initialised client for the passed connection.
func newClient(conn *net.TCPConn) *client {

	// Setup connection parameters
	if cerr := conn.SetKeepAlive(true); cerr != nil {
		log.Printf("Error setting keep alive: %s", cerr)
	}
	if cerr := conn.SetKeepAlivePeriod(5 * time.Second); cerr != nil {
		log.Printf("Error setting keep alive period: %s", cerr)
	}
	if cerr := conn.SetLinger(10); cerr != nil {
		log.Printf("Error setting linger: %s", cerr)
	}
	if cerr := conn.SetNoDelay(false); cerr != nil {
		log.Printf("Error setting no delay: %s", cerr)
	}
	if cerr := conn.SetWriteBuffer((termColumns * termLines) * 5); cerr != nil {
		log.Printf("Error setting write buffer: %s", cerr)
	}
	if cerr := conn.SetReadBuffer(inputBuffer); cerr != nil {
		log.Printf("Error setting read buffer: %s", cerr)
	}

	c := &client{
		TCPConn:    conn,
		remoteAddr: conn.RemoteAddr().String(),
	}

	log.Printf("Acquired lease: %s", conn.RemoteAddr())
	c.leaseAcquire()

	// Setup frontend if no error acquiring a lease
	c.frontend = frontend.New(c, c.remoteAddr, c.WriteError)
	if err := c.frontend.Parse([]byte("")); err != nil {
		return nil
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
		log.Println("Post process, ending player loop")
		if c != nil {
			c.close()
		}
	}()

	// Main input processing loop, terminates on any error raised not just read
	// or Parse errors.
	{
		// Variables for use in the loop only hence the scoping outer braces
		var (
			s   = bufio.NewReaderSize(c.TCPConn, inputBuffer) // Sized network read buffer
			err error                                         // Local Error
			in  []byte                                        // Input string from buffer
		)

		log.Print("Starting game loop: ", c.TCPConn.RemoteAddr())
		for c.err == nil {
			if config.Server.Running == false {
				_ = c.TCPConn.Close()
			}
			// Time in seconds to wait for input
			pingTime := 45 * time.Second
			idleTime := config.Server.IdleTimeout
			if ok := c.frontend.GetCharacter(); ok != (*objects.Character)(nil) {
				if c.frontend.GetCharacter().Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					idleTime = 30
				} else if c.frontend.GetCharacter().Flags["AFK"] {
					idleTime = config.Server.AFKTimeout
				} else if c.frontend.GetCharacter().Flags["OOC"] {
					idleTime = config.Server.OOCTimeout
				}
			}
			c.err = c.TCPConn.SetReadDeadline(time.Now().Add(pingTime))

			if in, err = s.ReadSlice('\n'); err != nil {
				frontend.Zero(in)

				// Check if this is a timeout error
				if err != nil {
					var netErr net.Error
					ok := errors.As(err, &netErr)
					if ok && netErr.Timeout() {
						_, err = c.Write([]byte(">"))
						if err == nil {
							if c.frontend.GetCharacter() != (*objects.Character)(nil) {
								log.Println(time.Now().Sub(c.frontend.GetCharacter().LastAction).Seconds())
								if time.Now().Sub(c.frontend.GetCharacter().LastAction).Minutes() > idleTime {
									c.err = errors.New("idle Timeout")
									continue
								}
							}
							c.err = nil
							continue
						} else {
							log.Println("Failed to write to client, client actually DC'd?: ", err)
						}
					}
				}

				if !errors.Is(err, bufio.ErrBufferFull) {
					log.Println("Client Error " + err.Error())
					if c.frontend.GetCharacter() != (*objects.Character)(nil) {
						c.frontend.GetCharacter().SuppressWrites()
					}
					c.WriteError(err)
					continue
				}

				for errors.Is(err, bufio.ErrBufferFull) {
					in, err = s.ReadSlice('\n')
					frontend.Zero(in)
				}
				if _, werr := c.Write([]byte(text.Bad + "\nYou type too much.\n" + text.Prompt + ">")); werr != nil {
					log.Println("Error writing to player: ", werr)
				}
				continue
			}

			//log.Println(&in)
			fixDEL(&in)
			if err = c.frontend.Parse(in); err != nil {
				log.Println("Text parse error " + err.Error())
				c.WriteError(err)
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
	// Deallocate current frontend if we have one
	if c.frontend != nil {
		// Sometimes these disconnects are a little messy,  need to add some extra cleanup
		if c.frontend.GetCharacter() != (*objects.Character)(nil) {
			log.Println("Force Close from Client")
			c.frontend.GetCharacter().PrepareUnload()
			c.frontend.GetCharacter().Unload()
			c.frontend.AccountCleanup()
		}

		c.frontend.Close()
		c.frontend = nil
	}

	// Make sure connection closed down and deallocated
	if err := c.TCPConn.Close(); err != nil {
		log.Printf("Error closing connection: %s", err)
	} else {
		log.Printf("Connection closed: %s", c.remoteAddr)
	}
	c.TCPConn = nil
	c.leaseRelease()

}

// Write handles output for the network connection.
func (c *client) Write(d []byte) (n int, err error) {
	if c.TCPConn != nil {
		if n, err = c.TCPConn.Write(d); err != nil {
			log.Println("TCP Error" + err.Error())
		}
	}
	return
}
