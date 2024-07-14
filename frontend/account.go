package frontend

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"log"
	"strconv"
)

// Account embeds a frontend instance adding fields and methods specific to
type Account struct {
	*Frontend
	account    string
	password   [16]byte
	permission permissions.Permissions
}

// NewAccount returns an account with the specified frontend embedded. The
// returned account can be used for processing the creation of new accounts and
// players.
func NewAccount(f *Frontend) (a *Account) {
	a = &Account{Frontend: f, permission: 1}
	a.explainAccountDisplay()
	return
}

// explainAccountDisplay displays the requirements for new account IDs. It is
// separated from newAccountDisplay so that if there is a problem we can ask
// for the new account ID again without having to have the explanation as well.
func (a *Account) explainAccountDisplay() {
	l := strconv.Itoa(config.Login.AccountLength)
	a.buf.Send("Your account ID can be anything you can remember: an email address, a book title, a film title, a quote. You can use upper and lower case characters, numbers and symbols. The only restriction is it has to be at least ", l, " characters long.\n\nThis is NOT your character's name it is for your account ID for logging in only.\n")
	a.newAccountDisplay()
}

// newAccountDisplay asks the player for a new account ID
func (a *Account) newAccountDisplay() {
	a.buf.Send("Enter text to use for your new account ID or just press enter to cancel:")
	a.nextFunc = a.newAccountProcess
}

// newAccountProcess takes the current input and stores it as an account ID
// hash. We don't know if it's already taken yet, we are just storing it.
func (a *Account) newAccountProcess() {
	switch l := len(a.input); {
	case l == 0:
		a.buf.Send(text.Info, "Account creation cancelled.\n", text.Reset)
		NewLogin(a.Frontend)
	case l < config.Login.AccountLength:
		l := strconv.Itoa(config.Login.AccountLength)
		a.buf.Send(text.Bad, "Account ID is too short. Needs to be ", l, " characters or longer.\n", text.Reset)
		a.newAccountDisplay()
	default:
		a.account = string(a.input)
		a.newPasswordDisplay()
	}
}

// newPasswordDisplay asks for a password to associate with the account ID.
func (a *Account) newPasswordDisplay() {
	a.buf.Send("Enter a password to use for your account ID or just press enter to cancel:")
	a.nextFunc = a.newPasswordProcess
}

// newPasswordProcess takes the current input and stores it in the current
// state as a hash. The hash is calculated with a random salt that is also
// stored in the current state.
func (a *Account) newPasswordProcess() {
	switch l := len(a.input); {
	case l == 0:
		a.buf.Send(text.Info, "Account creation cancelled.\n", text.Reset)
		NewLogin(a.Frontend)
	case l < config.Login.PasswordLength:
		l := strconv.Itoa(config.Login.PasswordLength)
		a.buf.Send(text.Bad, "Password is too short. Needs to be ", l, " characters or longer.\n", text.Reset)
		a.newPasswordDisplay()
	default:
		a.password = md5.Sum(a.input)
		a.confirmPasswordDisplay()
	}
}

// confirmPasswordDisplay asks for the password to be typed again for confirmation.
func (a *Account) confirmPasswordDisplay() {
	a.buf.Send("Enter your password again to confirm or just press enter to cancel:")
	a.nextFunc = a.confirmPasswordProcess
}

// confirmPasswordProcess verifies that the confirmation password matches the
// new password already stored in the current state.
func (a *Account) confirmPasswordProcess() {
	switch l := len(a.input); {
	case l == 0:
		a.buf.Send(text.Info, "Account creation cancelled.\n", text.Reset)
		NewLogin(a.Frontend)
	default:
		if md5.Sum(a.input) != a.password {
			a.buf.Send(text.Bad, "Passwords do not match, please try again.\n", text.Reset)
			a.newPasswordDisplay()
			return
		}
		a.write()
	}
}

func (a *Account) write() {

	// Check if account ID is already registered
	if data.AccountExists(a.account) {
		a.buf.Send(text.Bad, "The account ID you used is not available.\n", text.Reset)
		NewLogin(a.Frontend)
		return
	}

	newAcct := make(map[string]interface{})
	newAcct["name"] = a.account
	newAcct["password"] = hex.EncodeToString(a.password[:])

	if data.NewAcct(newAcct) {
		log.Printf("New account failed to create: %s", a.account)
		NewLogin(a.Frontend)
		return
	}

	log.Printf("New account created: %s", a.account)
	a.Frontend.account = a.account
	accounts.inuse[a.Frontend.account] = struct{}{}

	// Greet new player
	a.buf.Send(text.Good, "Welcome ", a.account, "!", text.Reset)

	NewStart(a.Frontend)
}
