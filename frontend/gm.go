package frontend

import (
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
	"strings"
)

// account embeds a frontend instance adding fields and methods specific to
// account and player creation.
type newPCharacter struct {
	*frontend
	name   string
	gender string
	race   int
}

func NewPCharacter(f *frontend) (a *newPCharacter) {
	a = &newPCharacter{frontend: f}
	a.explainPChar()
	return
}

// Character Name
func (a *newPCharacter) explainPChar() {
	a.buf.Send("Welcome to Aalynor's Nexus GM creator.  Just a couple of steps to get you running.")
	a.newPCharacterDisplay()
}

// newAccountDisplay asks the player for a new account ID
func (a *newPCharacter) newPCharacterDisplay() {
	a.buf.Send("Enter your GM name or just press enter to cancel:")
	a.nextFunc = a.charPNameProcess
}

// Process a character name
func (a *newPCharacter) charPNameProcess() {
	switch l := len(a.input); {
	case l == 0:
		a.buf.Send(text.Info, "GM creation cancelled.\n", text.Reset)
		NewStart(a.frontend)
	case l < config.Login.AccountLength:
		l := strconv.Itoa(config.Login.AccountLength)
		a.buf.Send(text.Bad, "GM name is too short. Needs to be ", l, " characters or longer.\n", text.Reset)
		a.newPCharacterDisplay()
	case verifyName.Find(a.input) == nil:
		a.buf.Send(text.Bad, "A character's name must only contain the upper or lower cased letters 'a' through 'z'. \n", text.Reset)
		a.newPCharacterDisplay()
	default:
		a.name = string(a.input)
		a.fastPCharDisplay()
	}
}

// ******   Fast Processing Options
func (a *newPCharacter) fastPCharDisplay() {
	a.buf.SendInfo(`Fast Step 1, Gender, Race

	Enter your gender (m|f) and race separated by spaces.  Don't worry, if you have a GM account interacting with players
you'll be able to change this later.  Otherwise it doesn't matter.'
(type help races for a list of races.)

	or:
Restart (r), return to name selection.
Cancel (c), return to the main menu.`)
	a.nextFunc = a.fastPCharProcess
}

func (a *newPCharacter) fastPCharProcess() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.fastPCharProcess
	case strings.HasPrefix("help", inputVal):
		if strings.Contains(" ", inputVal) {
			topic := strings.Split(inputVal, " ")[1]
			a.helpDisplay(topic)
		} else {
			a.helpDisplay("races")
		}
		a.nextFunc = a.fastPCharProcess
	case inputVal == "c":
		a.buf.Send(text.Info, "GM creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newPCharacterDisplay()
	case validateFastPStep(inputVal):
		items := strings.Split(inputVal, " ")
		a.buf.SendInfo(fmt.Sprintf("Gender: %[1]s, Race: %[2]s", items[0], items[1]))
		a.gender = items[0]
		a.race = utils.IndexOf(items[1], config.AvailableRaces)
		a.confirmFastSelections()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.fastPCharDisplay
	}
}

func (a *newPCharacter) confirmFastSelections() {
	a.buf.SendInfo(fmt.Sprintf(`Here is what you selected:
		Name: %[1]s
		Gender:  %[2]s
		Race:  %[3]s
	
Are you satisfied with these options?
===========================================
Yes (y) Finish your GM creation. 
No (n) Go back to selections
Restart (r) Restart the GM builder
Cancel (c) Leave the GM builder
`, a.name,
		a.gender,
		config.AvailableRaces[a.race],
	))
	a.nextFunc = a.confirmFastProcess
}

func (a *newPCharacter) confirmFastProcess() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.confirmFastProcess
	case inputVal == "n":
		a.buf.Send(text.Info, "Returning to previous step. \n", text.Reset)
		a.fastPCharDisplay()
	case inputVal == "c":
		a.buf.Send(text.Info, "GM creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newPCharacterDisplay()
	case inputVal == "y":
		a.completeBuilder()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.confirmFastProcess
	}
}

func (a *newPCharacter) completeBuilder() {
	charData := make(map[string]interface{})
	charData["name"] = a.name
	charData["class"] = 100
	charData["race"] = a.race
	charData["account"] = a.account
	charData["gender"] = a.gender
	charData["str"] = 30
	charData["dex"] = 30
	charData["con"] = 30
	charData["intel"] = 30
	charData["pie"] = 30
	charData["birthday"] = objects.CurrentDay
	charData["birthdate"] = objects.DayOfMonth
	charData["birthmonth"] = objects.CurrentMonth
	charData["darkvision"] = config.RaceDefs[config.AvailableRaces[a.race]].Darkvision
	if data.CreateChar(charData) {
		a.buf.Send(text.Info, "New GM created,  entering Altin. \n", text.Reset)
		StartGame(a.frontend, a.name)
	} else {
		a.buf.SendBad(text.Info, "Error, try again later. \n", text.Reset)
		NewStart(a.frontend)
	}
}

func (a *newPCharacter) helpDisplay(subject string) {
	// Print Race
	subject = strings.ToLower(subject)
	if subject == "races" {
		a.buf.Send("Available races: \n")
		a.buf.Send(strings.Join(config.AvailableRaces, ", "))
	} else {
		a.buf.Send("No help on that topic found")
	}
}

func validateFastPStep(choiceInput string) bool {
	inputs := strings.Split(choiceInput, " ")
	if len(inputs) != 2 {
		return false
	}
	if inputs[0] != "f" {
		if inputs[0] != "m" {
			return false
		}
	}
	if !utils.StringIn(inputs[1], config.AvailableRaces) {
		return false
	}
	return true
}
