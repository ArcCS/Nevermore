// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package frontend

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
)

// login embeds a frontend instance adding fields and methods specific to
// account logins.
type login struct {
	*frontend
	inputName string
}

// NewLogin returns a login with the specified frontend embedded. The returned
// login can be used for processing the logging in of accounts.
func NewLogin(f *frontend) (l *login) {
	l = &login{frontend: f}
	l.accountDisplay()
	return
}

// accountDisplay asks for the player's account ID so that they can log into
// the system.
func (l *login) accountDisplay() {
	l.buf.Send("Enter your account ID or just press enter to create a new account, enter QUIT to leave the server:")
	l.nextFunc = l.accountProcess
}

// accountProcess processes the current input as an account ID.
func (l *login) accountProcess() {
	//log.Printf("String from: %s", string(l.input))
	switch {
	case len(l.input) == 0:
		NewAccount(l.frontend)
	case bytes.Equal(bytes.ToUpper(l.input), []byte("QUIT")):
		l.Close()
	default:
		l.inputName = string(l.input)
		l.passwordDisplay()
	}
}

// passwordDisplay asks for the player's password for their account ID.
func (l *login) passwordDisplay() {
	l.buf.Send("Enter the password for your account ID or just press enter to cancel:")
	l.nextFunc = l.passwordProcess
}

// passwordProcess takes the current input and treats is as the player's
// password for logging into the system. If no password is entered processing
// goes back to asking for the players account ID. If the account ID is valid
// and the password is correct we load the player data and move on to
// displaying the main menu. If either the account ID or password is invalid we
// go back to asking for an account ID.
func (l *login) passwordProcess() {

	// If no password given go back and ask for an account ID.
	if len(l.input) == 0 {
		l.buf.Send(text.Info, "Login cancelled.\n", text.Reset)
		NewLogin(l.frontend)
		return
	}

	acctData, err := data.LoadAcct(l.inputName)
	if err {
		l.buf.Send(text.Bad, "Account ID or password is incorrect. (Load acct failure) \n", text.Reset)
		NewLogin(l.frontend)
		return
	}

	password, _ := acctData["password"].(string)
	encPass := md5.Sum(l.input)
	// Check password is valid
	if hex.EncodeToString(encPass[:]) != password {
		l.buf.Send(text.Bad, "Account ID or password is incorrect.\n", text.Reset)
		NewLogin(l.frontend)
		return
	}

	if !acctData["active"].(bool) {
		l.buf.Send(text.Bad, "Account is not active.\n", text.Reset)
		NewLogin(l.frontend)
		return
	}

	accounts.Lock()
	l.frontend.account = acctData["name"].(string)
	l.frontend.permissions = permissions.Permissions(acctData["permissions"].(int64))
	accounts.Unlock()

	// Greet returning account
	l.buf.Send(text.Good, "Welcome back ", l.frontend.account, "!", text.Reset)

	NewStart(l.frontend)
}
