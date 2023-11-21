// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

// Package frontend provides interactive processing of front end pages and
// access to the backend game. Note that a 'page' can be anything from a menu
// of options to choose from to a simple request for a password.
//
// The frontend is responsible for coordinating the display of pages to a user
// account creation, player creation and other non in-game activities. When the
// player is in-game the frontend will simply pass any input through to the
// game backend for processing.
//
// Pages typically have a pair of methods - a display part and a processing
// part. For example accountDisplay and accountProcess. Sometimes there is only
// a display part, for example greetingDisplay.
//
// The current state is held in an instance of frontend. With frontend.nextFunc
// being the next method to call when input is received - usually an xxxProcess
// method.
//
// Each time input is received Parse will be called. The method in nextFunc
// will be called to handle the input. nextFunc should then call the next
// xxxDisplay method to send a response to the input processing and setup
// nextFunc with the method that will process the next input received. Any
// buffered response will then be sent back before Parse exits. Parse will then
// be called again when more input is received.
package frontend

import (
	"bytes"
	"github.com/ArcCS/Nevermore/cmd"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"io"
	"log"
	"sync"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/message"
	"github.com/ArcCS/Nevermore/text"
)

// accounts is used to track which (valid) accounts are logged in and in use.
// It's main purpose is to track logged in account IDs to prevent duplicate
// logins.
var accounts struct {
	sync.Mutex
	inuse map[string]struct{}
}

// init is used to initialise the map used in account ID tracking.
func init() {
	accounts.inuse = make(map[string]struct{})

}

// closedError represents the fact that Close has been called on a frontend
type closedError struct{}

// Error implements the error interface for errors and returns descriptive text
// for the closedError error.
func (closedError) Error() string {
	return "Game interface closed."
}

// Temporary always returns true for a frontend.Error. A frontend.Error is
// considered temporary as recovery is easy - create a new frontend instance.
func (closedError) Temporary() bool {
	return true
}

// frontend represents the current frontend state for a given io.Writer - this
// is typically from a player's network connection.
type frontend struct {
	output      io.Writer          // Writer to send output text to
	buf         *message.Buffer    // Buffered messages written with next prompt
	input       []byte             // The input text we are currently processing
	nextFunc    func()             // The next frontend function called by Parse
	remoteAddr  string             // IP Address
	character   *objects.Character // The current player instance (ingame or not)
	account     string
	permissions permissions.Permissions
	err         error       // First error to occur else nil
	writeError  func(error) // Network error to return to player
	ClientClose func()      // Disconnect up the stack
}

func (f *frontend) GetCharacter() *objects.Character {
	if f.character != (*objects.Character)(nil) {
		return f.character
	} else {
		return (*objects.Character)(nil)
	}
}

// New returns an instance of frontend initialised with the given io.Writer.
// The io.Writer is used to send responses back from calling Parse. The new
// frontend is initialised with a message buffer and nextFunc setup to call
// greetingDisplay.
func New(output io.Writer, address string, errorWriter func(error), clientclose func()) *frontend {
	f := &frontend{
		buf:        message.AcquireBuffer(),
		output:     output,
		remoteAddr: address,
		writeError: errorWriter,
	}
	f.buf.OmitLF(true)
	f.nextFunc = f.greetingDisplay
	f.ClientClose = clientclose
	return f
}

func (f *frontend) Disconnect() {
	f.character = (*objects.Character)(nil)
	f.ClientClose()
}

func (f *frontend) AccountCleanup() {
	delete(accounts.inuse, f.account)
	log.Println(accounts.inuse)
}

// Close makes sure the player is no longer 'in game' and frees up resources
func (f *frontend) Close() {

	// Just return if we already have an error
	if f.err != nil {
		return
	}
	f.err = closedError{}

	// If player is still in the game force them to quit
	if f.character != (*objects.Character)(nil) {
		if objects.ActiveCharacters.Find(f.character.Name) != nil {
			cmd.Parse(f.character, "QUIT")
		}
	}

	// Free up resources
	message.ReleaseBuffer(f.buf)
	f.buf = (*message.Buffer)(nil)

	f.output = nil
	f.nextFunc = nil

	f.character = (*objects.Character)(nil)

	f = (*frontend)(nil)
}

// Parse is the main input/output processing method for frontend. The input is
// stripped of leading and trailing whitespace before being stored in the
// frontend state. Any response from processing the input is written to the
// io.Writer passed to the initial New function that created the frontend. If
// the frontend is closed during processing a frontend.Error will be returned
// else nil.
func (f *frontend) Parse(input []byte) error {

	// If we already have an error just return it
	if f.err != nil {
		return f.err
	}

	// Trim whitespace from input and process it
	f.input = bytes.TrimSpace(input)
	f.nextFunc()

	// If we have a message buffer write out its content and a new prompt
	if f.buf != nil {
		f.buf.Deliver(f)
	}
	return f.err
}

// greetingDisplay displays the welcome message to players when they first
// connect to the server.
func (f *frontend) greetingDisplay() {
	f.buf.Send(config.DragonAscii)
	NewLogin(f)
}

// Write writes the specified byte slice to the associated client.
func (f *frontend) Write(b []byte) (n int, err error) {
	b = append(b, text.Prompt...)
	b = append(b, '>')
	n, err = f.output.Write(b)
	if err != nil {
		// This might be a broken pipe. if so, we need to close the connection
		f.writeError(err)
	}
	return
}

// Zero writes zero bytes into the passed slice
func Zero(data []byte) {
	if len(data) > 0 {
		data[0] = 0
		for i := 1; i < len(data); i *= 2 {
			copy(data[i:], data[:i])
		}
	}
}
