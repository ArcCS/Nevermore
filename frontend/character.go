package frontend

import (
	"fmt"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/text"
)

// account embeds a frontend instance adding fields and methods specific to
// account and player creation.
type newCharacter struct {
	*frontend
	name   string
	gender string
	class  int
	race   int
	str    int
	con    int
	dex    int
	pie    int
	intel  int
}

// verifyName is used to test that a players name only uses the letters A-Z,a-z.
var verifyName = regexp.MustCompile(`^[a-zA-Z]+$`)

func NewCharacter(f *frontend) (a *newCharacter) {
	a = &newCharacter{frontend: f}
	a.explainCharName()
	return
}

// Character Name
func (a *newCharacter) explainCharName() {
	a.buf.Send("Welcome to Aalynor's Nexus character creation.  The first step is to choose your character name. It needs to use characters only. No numbers, no spaces. ")
	a.newCharacterDisplay()
}

// newAccountDisplay asks the player for a new account ID
func (a *newCharacter) newCharacterDisplay() {
	a.buf.Send("Enter your characters name or just press enter to cancel:")
	a.nextFunc = a.charNameProcess
}

// Process a character name
func (a *newCharacter) charNameProcess() {
	switch l := len(a.input); {
	case l == 0:
		a.buf.Send(text.Info, "Character creation cancelled.\n", text.Reset)
		NewStart(a.frontend)
	case l < config.Login.AccountLength:
		l := strconv.Itoa(config.Login.AccountLength)
		a.buf.Send(text.Bad, "Character name is too short. Needs to be ", l, " characters or longer.\n", text.Reset)
		a.newCharacterDisplay()
	case verifyName.Find(a.input) == nil:
		a.buf.Send(text.Bad, "A character's name must only contain the upper or lower cased letters 'a' through 'z'. \n", text.Reset)
		a.newCharacterDisplay()
	case utils.StringInLike(strings.ToLower(string(a.input)), config.BlockedNames):
		a.buf.Send(text.Bad, "The requested name is unavailable. \n", text.Reset)
		a.newCharacterDisplay()
	case data.CharacterExists(string(a.input)):
		a.buf.Send(text.Bad, "That character already exists in this world. \n", text.Reset)
		a.newCharacterDisplay()
	default:
		a.name = string(a.input)
		a.creationSpeedDisplay()
	}
}

func (a *newCharacter) creationSpeedDisplay() {
	a.buf.Send(`Welcome to Aalynor's Nexus character builder!

Please choose one of the following options:
===========================================
Quick (q), a couple of questions to get you started.
Normal (n), story guidance generating your character.
Restart (r), return to name selection.
Cancel (c), return to the main menu.`)

	a.nextFunc = a.creationSpeedProcess
}

func (a *newCharacter) creationSpeedProcess() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.creationSpeedProcess
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case inputVal == "q":
		a.buf.Send(text.Info, "Quick menu requested. \n", text.Reset)
		a.fastStep1Display()
	case inputVal == "n":
		a.buf.Send(text.Info, "Story menu requested. \n", text.Reset)
		a.selectGenderDisplay()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.creationSpeedProcess
	}
}

// ******   Story Processing Options
func (a *newCharacter) selectGenderDisplay() {
	a.buf.SendInfo(`You find yourself on a pier.  The air is dense and foggy with mists swirling about you.
An old woman wanders up to you seemingly formed of mist herself.
The Old Woman says:  'Your spirit has come alive and set you on a new journey... and here you are.. hero...'
The Old Woman pauses and looks you over and asks: "My vision isn't what it used to be..  are you a man or a woman?'

(Your gender will have some effects on how certain creatures will react to you in the world, and additionally your ability to handle some magic weapons)
(Your gender doesn't affect any combat traits)"

Choose one of the following options:
===========================================
Male (m, male)
Female (f, female)
Back (b), return to speed selection.
Restart (r), return to name selection.
Cancel (c), return to the main menu.`)
	a.nextFunc = a.selectGenderProcess
}

func (a *newCharacter) selectGenderProcess() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.selectGenderProcess
	case inputVal == "b":
		a.buf.Send(text.Info, "Returning to previous step. \n", text.Reset)
		a.creationSpeedDisplay()
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case strings.Contains("male", strings.ToLower(inputVal)):
		a.buf.Send(text.Info, "Your character is male. \n", text.Reset)
		a.gender = "m"
		a.selectRaceDisplay()
	case strings.Contains("female", strings.ToLower(inputVal)):
		a.buf.Send(text.Info, "Your character is female. \n", text.Reset)
		a.gender = "f"
		a.selectRaceDisplay()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.selectGenderProcess
	}
}

// ******   Story Processing Options
func (a *newCharacter) selectRaceDisplay() {
	a.buf.SendInfo(`The Old Woman tips her head to the side:  "I suppose I see it now""
The Old Woman reaches up and puts her hands on your face; they are cool, aged, and dry as she feels over your features.
THe Old Woman asks: "And from what race of people do you descend?"
(Type 'help' to get a list of available races,  type "help <race>" to learn more about a specific race)

===========================================
Enter your chosen race name.
	or:
Back (b), return to gender selection.
Restart (r), return to name selection.
Cancel (c), return to the main menu.`)
	a.nextFunc = a.selectRaceProcess
}

func (a *newCharacter) selectRaceProcess() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.selectRaceProcess
	case strings.Contains(inputVal, "help"):
		if strings.Contains(inputVal, " ") {
			topic := strings.Split(inputVal, " ")[1]
			a.helpDisplay(topic)
		} else {
			a.helpDisplay("races")
		}
		a.nextFunc = a.selectRaceProcess
	case inputVal == "b":
		a.buf.Send(text.Info, "Returning to previous step. \n", text.Reset)
		a.selectGenderDisplay()
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case utils.StringIn(inputVal, config.AvailableRaces):
		a.buf.Send(text.Info, "Your character is a: ", inputVal, ".", text.Reset)
		a.race = utils.IndexOf(inputVal, config.AvailableRaces)
		a.selectClassDisplay()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.selectRaceProcess
	}
}

func (a *newCharacter) selectClassDisplay() {
	a.buf.SendInfo(`The Old Woman says, "All heroes specialize in certain combat abilities,
This combat specialization will define how you fight, and how you work with others.
Tell me, how will you contribute to this world?"
(This will be your characters class, this cannot be changed later and determines the abilities available to your character.
Type 'help' to get a list of classes, or 'help <class name>' to find out more about a class)

===========================================
Enter your chosen class name.
	or:
Back (b), return to race selection.
Restart (r), return to name selection.
Cancel (c), return to the main menu.`)
	a.nextFunc = a.selectClassProcess
}

func (a *newCharacter) selectClassProcess() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.selectClassProcess
	case strings.Contains(inputVal, "help"):
		if strings.Contains(inputVal, " ") {
			topic := strings.Split(inputVal, " ")[1]
			a.helpDisplay(topic)
		} else {
			a.helpDisplay("classes")
		}
		a.nextFunc = a.selectClassProcess
	case inputVal == "b":
		a.buf.Send(text.Info, "Returning to previous step. \n", text.Reset)
		a.selectRaceDisplay()
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case utils.StringIn(inputVal, config.AvailableClasses):
		a.buf.Send(text.Info, "Your character class is ", inputVal, ".", text.Reset)
		a.class = utils.IndexOf(inputVal, config.AvailableClasses)
		a.selectStatsDisplay()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.selectClassProcess
	}
}

func (a *newCharacter) selectStatsDisplay() {
	a.buf.SendInfo(fmt.Sprintf(`The Old Woman hears your answer and appears to be visually sizing up your suitability for the desired profession.

Please select your stats.  You have 50 points to spend, and based on your
race selection here are your respective minimums and maximums:

%[1]s
Strength: %[2]s/%[3]s
Dexterity: %[4]s/%[5]s
Constitution: %[6]s/%[7]s
Intelligence: %[8]s/%[9]s
Piety: %[10]s/%[11]s"

===========================================
(Enter all 5 numbers with a space between each number, in order of -> STR DEX CON INT PIE)

	or:
Back (b), return to race selection.
Restart (r), return to name selection.
Cancel (c), return to the main menu.`,
		config.AvailableRaces[a.race],
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].StrMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].StrMax),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].DexMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].DexMax),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].ConMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].ConMax),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].IntMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].IntMax),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].PieMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].PieMax),
	))
	a.nextFunc = a.selectStatsProcess
}

func (a *newCharacter) selectStatsProcess() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.selectStatsProcess
	case inputVal == "b":
		a.buf.Send(text.Info, "Returning to previous step. \n", text.Reset)
		a.selectClassDisplay()
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case a.validateStats(inputVal):
		stats := strings.Split(inputVal, " ")
		a.buf.SendInfo(fmt.Sprintf("Your stats will be Str: %[1]s Dex: %[2]s Con: %[3]s Int: %[4]s Pie: %[5]s", stats[0], stats[1], stats[2], stats[3], stats[4]))
		a.str, _ = strconv.Atoi(stats[0])
		a.dex, _ = strconv.Atoi(stats[1])
		a.con, _ = strconv.Atoi(stats[2])
		a.intel, _ = strconv.Atoi(stats[3])
		a.pie, _ = strconv.Atoi(stats[4])
		a.confirmSelections()
	default:
		a.buf.Send(text.Info, "Invalid stats were entered, please review.  \n", text.Reset)
		a.nextFunc = a.selectStatsProcess
	}
}

func (a *newCharacter) confirmSelections() {
	a.buf.SendInfo(fmt.Sprintf(`Here is what you selected:
		Gender:  %[1]s
		Race:  %[2]s
		Class:  %[3]s
		Str:  %[4]d
		Dex:  %[5]d
		Con:  %[6]d
		Int:  %[7]d
		Pie: %[8]d
	
Are you satisfied with these options?
===========================================
Yes (y) Finish the Character Builder
No (n) Go back to the last step
Restart (r) Restart the character builder
Cancel (c) Leave the character builder
`, a.gender,
		config.AvailableRaces[a.race],
		config.AvailableClasses[a.class],
		a.str,
		a.dex,
		a.con,
		a.intel,
		a.pie,
	))
	a.nextFunc = a.confirmProcess
}

func (a *newCharacter) confirmProcess() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.confirmProcess
	case inputVal == "n":
		a.buf.Send(text.Info, "Returning to previous step. \n", text.Reset)
		a.selectStatsDisplay()
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case inputVal == "y" || inputVal == "yes":
		a.storyFinish()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.selectRaceProcess
	}
}

// Finalize, save, and place hte character into the world.
func (a *newCharacter) storyFinish() {
	a.buf.SendInfo(`The Old Woman says, "You'll fit into this world just fine.
You should know that Altin saw nearly a thousand years of something that people might call peace.
However, during that time evil has been allowed to grow in the dark corners of the world.
Altin has need of you again. I wish you luck in your endeavours.. 
Welcome to the fortified training city: Rymek."
The Old Woman turns and walks away.  Her form distorting into the mists from which she came.

The mists begin to clear...

(Your character has been created.  Send any input to begin journey.`)
	a.nextFunc = a.completeBuilder
}

// ******   Fast Processing Options
func (a *newCharacter) fastStep1Display() {
	a.buf.SendInfo(`Fast Step 1, Gender, Race, Class

	Enter your gender (m|f), race, and class separated by spaces"
(type help <class>/<race> for information.)

	or:
Restart (r), return to name selection.
Cancel (c), return to the main menu.`)
	a.nextFunc = a.fastStep1Process
}

func (a *newCharacter) fastStep1Process() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.fastStep1Process
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case strings.Contains(inputVal, "help"):
		if strings.Contains(inputVal, " ") {
			topic := strings.Split(inputVal, " ")[1]
			a.helpDisplay(topic)
		} else {
			a.helpDisplay("classes")
		}
		a.nextFunc = a.fastStep1Process
	case validateFastStep(inputVal):
		items := strings.Split(inputVal, " ")
		a.buf.SendInfo(fmt.Sprintf("Gender: %[1]s, Race: %[2]s, Class: %[3]s", items[0], items[1], items[2]))
		a.gender = items[0]
		a.race = utils.IndexOf(items[1], config.AvailableRaces)
		a.class = utils.IndexOf(items[2], config.AvailableClasses)
		a.fastStep2Display()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.fastStep1Process
	}
}

func (a *newCharacter) fastStep2Display() {
	a.buf.SendInfo(fmt.Sprintf(`Please select your stats.  You have 50 points to spend, and based on your
race selection here are your respective minimums and maximums:

%[1]s
Strength: %[2]s/%[3]s
Dexterity: %[4]s/%[5]s
Constitution: %[6]s/%[7]s
Intelligence: %[8]s/%[9]s
Piety: %[10]s/%[11]s"

(Enter all 5 numbers with a space between each number, in order of -> STR DEX CON INT PIE)

	or:
Back (b), return to gender, race and class select. 
Restart (r), return to name selection.
Cancel (c), return to the main menu.`,
		config.AvailableRaces[a.race],
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].StrMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].StrMax),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].DexMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].DexMax),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].ConMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].ConMax),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].IntMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].IntMax),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].PieMin),
		strconv.Itoa(config.RaceDefs[config.AvailableRaces[a.race]].PieMax),
	))
	a.nextFunc = a.fastStep2Process
}

func (a *newCharacter) fastStep2Process() {
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.fastStep2Process
	case inputVal == "b":
		a.buf.Send(text.Info, "Returning to previous step. \n", text.Reset)
		a.fastStep1Display()
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case a.validateStats(inputVal):
		stats := strings.Split(inputVal, " ")
		a.buf.SendInfo(fmt.Sprintf("Your stats will be Str: %[1]s Dex: %[2]s Con: %[3]s Int: %[4]s Pie: %[5]s", stats[0], stats[1], stats[2], stats[3], stats[4]))
		a.str, _ = strconv.Atoi(stats[0])
		a.dex, _ = strconv.Atoi(stats[1])
		a.con, _ = strconv.Atoi(stats[2])
		a.intel, _ = strconv.Atoi(stats[3])
		a.pie, _ = strconv.Atoi(stats[4])
		a.confirmFastSelections()
	default:
		a.buf.Send(text.Info, "Invalid stats were entered, please review. \n", text.Reset)
		a.nextFunc = a.fastStep2Process
	}
}

func (a *newCharacter) confirmFastSelections() {
	a.buf.SendInfo(fmt.Sprintf(`Here is what you selected:
		Gender:  %[1]s
		Race:  %[2]s
		Class:  %[3]s
		Str:  %[4]d
		Dex:  %[5]d
		Con:  %[6]d
		Int:  %[7]d
		Pie: %[8]d
	
Are you satisfied with these options?
===========================================
Yes (y) Finish the Character Builder
No (n) Go back to the last step
Restart (r) Restart the character builder
Cancel (c) Leave the character builder
`, a.gender,
		config.AvailableRaces[a.race],
		config.AvailableClasses[a.class],
		a.str,
		a.dex,
		a.con,
		a.intel,
		a.pie,
	))
	a.nextFunc = a.confirmFastProcess
}

func (a *newCharacter) confirmFastProcess() {
	log.Println("entering confirmFastProcess with input: ", string(a.input))
	inputVal := strings.ToLower(string(a.input))
	switch l := len(inputVal); {
	case l == 0:
		a.buf.Send(text.Info, "No input given. Please try again. \n", text.Reset)
		a.nextFunc = a.confirmFastProcess
	case inputVal == "n" || inputVal == "no":
		a.buf.Send(text.Info, "Returning to previous step. \n", text.Reset)
		a.fastStep1Display()
	case inputVal == "c":
		a.buf.Send(text.Info, "Character creation cancelled. \n", text.Reset)
		NewStart(a.frontend)
	case inputVal == "r":
		a.buf.Send(text.Info, "Restart requested. \n", text.Reset)
		a.newCharacterDisplay()
	case inputVal == "y" || inputVal == "yes":
		a.completeBuilder()
	default:
		a.buf.Send(text.Info, "Unrecognized input, please try again. \n", text.Reset)
		a.nextFunc = a.confirmFastProcess
	}
}

func (a *newCharacter) completeBuilder() {
	charData := make(map[string]interface{})
	charData["account"] = a.account
	charData["gender"] = a.gender
	charData["name"] = a.name
	charData["class"] = a.class
	charData["race"] = a.race
	charData["str"] = a.str
	charData["dex"] = a.dex
	charData["con"] = a.con
	charData["intel"] = a.intel
	charData["pie"] = a.pie
	charData["darkvision"] = config.RaceDefs[config.AvailableRaces[a.race]].Darkvision
	charData["birthday"] = objects.CurrentDay
	charData["birthdate"] = objects.DayOfMonth
	charData["birthmonth"] = objects.CurrentMonth
	if data.CreateChar(charData) {
		a.buf.Send(text.Info, "# New character created,  entering Altin. \n", text.Reset)
		FirstTimeStartGame(a.frontend, a.name)
	} else {
		a.buf.SendBad(text.Info, "Error, try again later. \n", text.Reset)
		NewStart(a.frontend)
	}
}

func (a *newCharacter) helpDisplay(subject string) {
	// Print Race
	subject = strings.ToLower(subject)
	if subject == "races" {
		a.buf.Send("Available races: \n")
		a.buf.Send(strings.Join(config.AvailableRaces, ", "))
	} else if subject == "classes" {
		a.buf.Send("Available classes: \n")
		a.buf.Send(strings.Join(config.AvailableClasses, ", "))
	} else if utils.StringIn(subject, config.AvailableRaces) {
		outLine := fmt.Sprintf("Race: %[1]s \n"+
			"Desc: %[2]s \n"+
			"Attribute: Min/Max \n"+
			"Strength: %[3]s/%[4]s, \n"+
			"Dexterity: %[5]s/%[6]s, \n"+
			"Constitution: %[7]s/%[8]s, \n"+
			"Intelligence: %[9]s/%[10]s, \n"+
			"Piety: %[11]s/%[12]s \n",
			subject,
			config.RaceDefs[subject].Desc,
			strconv.Itoa(config.RaceDefs[subject].StrMin),
			strconv.Itoa(config.RaceDefs[subject].StrMax),
			strconv.Itoa(config.RaceDefs[subject].DexMin),
			strconv.Itoa(config.RaceDefs[subject].DexMax),
			strconv.Itoa(config.RaceDefs[subject].ConMin),
			strconv.Itoa(config.RaceDefs[subject].ConMax),
			strconv.Itoa(config.RaceDefs[subject].IntMin),
			strconv.Itoa(config.RaceDefs[subject].IntMax),
			strconv.Itoa(config.RaceDefs[subject].PieMin),
			strconv.Itoa(config.RaceDefs[subject].PieMax),
		)
		a.buf.Send(outLine)
	} else if utils.StringIn(subject, config.AvailableClasses) {
		outLine := fmt.Sprintf(
			"Class: %[1]s \n"+
				"Desc: %[2]s \n"+
				"Abilities: %[3]s \n"+
				"Recommended Races: %[4]s \n"+
				"Recommended Stat Focuses: %[5]s \n",
			subject,
			config.Classes[subject].Desc,
			config.Classes[subject].Skills,
			config.Classes[subject].Races,
			config.Classes[subject].Stats,
		)
		a.buf.Send(outLine)
	} else {
		a.buf.Send("No help on that topic found")
	}
}

func parseStats(statInput string) []int {
	stats := strings.Split(statInput, " ")
	statOut := make([]int, 5)
	for i, v := range stats {
		statOut[i], _ = strconv.Atoi(v)
	}
	return statOut
}

func validateFastStep(choiceInput string) bool {
	inputs := strings.Split(choiceInput, " ")
	if len(inputs) < 3 {
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
	if !utils.StringIn(inputs[2], config.AvailableClasses) {
		return false
	}
	return true
}

func (a *newCharacter) validateStats(statInput string) bool {

	stats := parseStats(statInput)
	if stats[0]+stats[1]+stats[2]+stats[3]+stats[4] != 50 {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[a.race]].StrMin > stats[0] || stats[0] > config.RaceDefs[config.AvailableRaces[a.race]].StrMax || stats[0] > 20 {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[a.race]].DexMin > stats[1] || stats[1] > config.RaceDefs[config.AvailableRaces[a.race]].DexMax || stats[1] > 20 {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[a.race]].ConMin > stats[2] || stats[2] > config.RaceDefs[config.AvailableRaces[a.race]].ConMax || stats[2] > 20 {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[a.race]].IntMin > stats[3] || stats[3] > config.RaceDefs[config.AvailableRaces[a.race]].IntMax || stats[3] > 20 {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[a.race]].PieMin > stats[4] || stats[4] > config.RaceDefs[config.AvailableRaces[a.race]].PieMax || stats[4] > 20 {
		return false
	}
	return true

}
