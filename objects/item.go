package objects

import (
	"strconv"
	"strings"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
)

type Item struct {
	Object
	ParentItemId int
	ItemId       int
	ItemType     int
	Flags        map[string]bool
	Creator      string
	NumDice      int
	PlusDice     int
	SidesDice    int
	Armor        int
	ArmorClass   int
	MaxUses      int
	Value        int
	Spell        string
	StorePrice   int
	Adjustment   int
	Repairs      int    // repairs since the last overhaul, see config.WeaponRepairsBeforeOverhaul
	BoundTo      string // the character a worn ring bound itself to, see BindOnWear

	Storage *ItemInventory
	Weight  int
}

func LoadItem(itemData map[string]interface{}) (*Item, bool) {
	description := ""
	var ok bool
	if description, ok = itemData["description"].(string); !ok {
		description = "This one was missing a description.. maybe fix that?"
	}
	newItem := &Item{
		Object{
			Name:        itemData["name"].(string),
			Description: description,
			Placement:   3,
			Commands:    DeserializeCommands(itemData["commands"].(string)),
		},
		0,
		int(itemData["item_id"].(int64)),
		int(itemData["type"].(int64)),
		make(map[string]bool),
		itemData["creator"].(string),
		int(itemData["ndice"].(int64)),
		int(itemData["pdice"].(int64)),
		int(itemData["sdice"].(int64)),
		int(itemData["armor"].(int64)),
		int(itemData["armor_class"].(int64)),
		int(itemData["max_uses"].(int64)),
		int(itemData["value"].(int64)),
		itemData["spell"].(string),
		0,
		int(itemData["adjustment"].(int64)),
		0,
		"",
		&ItemInventory{},
		int(itemData["weight"].(int64)),
	}
	for k, v := range itemData["flags"].(map[string]interface{}) {
		if v == nil {
			newItem.Flags[k] = false
		} else {
			newItem.Flags[k] = v.(int64) != 0
		}
	}
	return newItem, true
}

func (i *Item) GetWeight() int {
	if i.ItemType == 9 && !i.Flags["weightless_chest"] {
		return i.Weight + i.Storage.GetTotalWeight()
	} else {
		return i.Weight
	}
}

// Look describes the item to viewer, whose class decides the tier a weapon reads
// as usable from.  A nil viewer reads the tier without any class adjustment.
func (i *Item) Look(viewer *Character) string {
	resString := i.appraisal(viewer) + strings.ReplaceAll(i.Description, "\\n", "\n") + "\n\n"
	if utils.IntIn(i.ItemType, config.WeaponTypes) {
		resString = "It is a " + config.ItemTypes[i.ItemType] + " weapon, and it" + i.ReturnState() + "\n" + resString
	}
	if utils.IntIn(i.ItemType, config.ArmorTypes) {
		resString = "It is a " + config.ArmorDescription(i.ItemType, i.ArmorClass) + ", and it" + i.ReturnState() + "\n" + resString
	}
	if i.ItemType == 9 {
		items := i.Storage.ReducedList()
		if len(items) > 0 {
			resString += "The " + i.Name + " contains " + strconv.Itoa(len(i.Storage.Contents)) + " items: \n" + items
		}
	}
	return resString
}

// appraisal is the qualitative lines a look gives on a weapon, armor, device or
// wand.  Ratings read the template, so enchanting or binding a copy doesn't
// change them.
func (i *Item) appraisal(viewer *Character) string {
	base := i
	if template, ok := Items[i.ItemId]; ok && template != nil {
		base = template
	}
	switch {
	case utils.IntIn(i.ItemType, []int{0, 1, 2, 3, 4}):
		if base.NumDice <= 0 {
			return ""
		}
		roll := config.WeaponRoll{Dice: base.NumDice, Sides: base.SidesDice, Plus: base.PlusDice}
		twoHanded := base.Flags["two_handed"]
		special := base.Flags["no_shatter"] || base.Flags["quest_loot"] || base.Flags["permanent"] ||
			base.Flags["magic"] || base.Flags["light"]
		class := -1
		if viewer != nil {
			class = viewer.Class
		}
		return qualityLine(config.WeaponQuality(roll, twoHanded, base.Flags["always_crit"], special)) +
			tierLine(config.WeaponTierFor(roll.Max(), twoHanded, class))
	case utils.IntIn(i.ItemType, config.ArmorTypes):
		armor := i.BaseArmor()
		return qualityLine(config.ArmorQuality(i.ItemType, armor)) + tierLine(config.ArmorTier(i.ItemType, armor))
	case i.ItemType == 6 || i.ItemType == 8:
		resString := "This is a " + config.ItemTypes[i.ItemType]
		if kind := SpellEffectKind(i.Spell); kind != "" {
			resString += " with " + kind + " effects"
		}
		resString += ".\n"
		if charges := config.ChargeDescription(i.MaxUses, base.MaxUses); charges != "" && !i.Flags["infinite"] {
			resString += charges + "\n"
		}
		return resString
	}
	return ""
}

// qualityLine colors only the rating, then picks back up the look's info color.
func qualityLine(quality int) string {
	return "This item appears to be " + config.QualityColors[quality] + config.QualityNames[quality] + text.Info + ".\n"
}

// tierLine is empty for a tier no character can reach.
func tierLine(tier int) string {
	if tier < 1 || tier > config.LevelCap {
		return ""
	}
	return "It can be used from tier " + strconv.Itoa(tier) + ".\n"
}

// DisplayName Return a display name with numerics
func (i *Item) DisplayName() string {
	typeReturn := 0
	preName := ""
	if i.Flags["magic"] {
		preName += "magic "
	}
	// Mapping value definitions
	switch i.ItemType {
	case 0:
		typeReturn = 1
	case 1:
		typeReturn = 1
	case 2:
		typeReturn = 1
	case 3:
		typeReturn = 1
	case 4:
		typeReturn = 1
	case 5:
		typeReturn = 2
	case 15:
		typeReturn = 1
	case 16:
		typeReturn = 1
	case 19:
		typeReturn = 2
	case 20:
		typeReturn = 2
	case 21:
		typeReturn = 2
	case 22:
		typeReturn = 2
	case 23:
		typeReturn = 2
	case 24:
		typeReturn = 2
	case 25:
		typeReturn = 2
	case 26:
		typeReturn = 2
	}
	if typeReturn == 1 {
		return preName + i.Name + " (" + strconv.Itoa(i.Adjustment) + ")"
	} else if typeReturn == 2 {
		if i.Armor > 0 {
			return preName + i.Name + " (" + strconv.Itoa(i.Armor) + ")"
		} else {
			return preName + i.Name
		}
	}
	return preName + i.Name
}

// BindOnWear binds worn gear the first time it takes damage: quest loot is
// soulbound (permanent), and a ring is bound to its wearer, who alone can pick
// it up again.  It returns the message for the wearer, or "" when nothing
// changed.
func (i *Item) BindOnWear(owner string) string {
	if i.Flags["quest_loot"] {
		if i.Flags["permanent"] {
			return ""
		}
		i.Flags["permanent"] = true
		return "Your " + i.DisplayName() + " binds itself to your soul."
	}
	if i.ItemType == 24 && i.BoundTo == "" && owner != "" {
		i.BoundTo = owner
		return "Your " + i.DisplayName() + " binds itself to you."
	}
	return ""
}

// Bound reports whether the item is bound to its owner, either soulbound or a
// ring bound by wear, and so can't be dropped, given, put away or sold.
func (i *Item) Bound() bool {
	return i.Flags["permanent"] || i.BoundTo != ""
}

func (i *Item) ToggleFlag(flagName string) bool {
	if val, exists := i.Flags[flagName]; exists {
		i.Flags[flagName] = !val
		return true
	} else {
		return false
	}
}

func (i *Item) Save() {
	itemData := make(map[string]interface{})
	itemData["item_id"] = i.ItemId
	itemData["creator"] = i.Creator
	itemData["ndice"] = i.NumDice
	itemData["weight"] = i.Weight
	itemData["description"] = i.Description
	itemData["type"] = i.ItemType
	itemData["pdice"] = i.PlusDice
	itemData["armor"] = i.Armor
	itemData["armor_class"] = i.ArmorClass
	itemData["max_uses"] = i.MaxUses
	itemData["name"] = i.Name
	itemData["sdice"] = i.SidesDice
	itemData["value"] = i.Value
	itemData["spell"] = i.Spell
	itemData["always_crit"] = utils.Btoi(i.Flags["always_crit"])
	itemData["quest_loot"] = utils.Btoi(i.Flags["quest_loot"])
	itemData["no_shatter"] = utils.Btoi(i.Flags["no_shatter"])
	itemData["permanent"] = utils.Btoi(i.Flags["permanent"])
	itemData["magic"] = utils.Btoi(i.Flags["magic"])
	itemData["light"] = utils.Btoi(i.Flags["light"])
	itemData["no_take"] = utils.Btoi(i.Flags["no_take"])
	itemData["weightless_chest"] = utils.Btoi(i.Flags["weightless_chest"])
	itemData["two_handed"] = utils.Btoi(i.Flags["two_handed"])
	itemData["adjustment"] = i.Adjustment
	itemData["commands"] = i.SerializeCommands()
	data.UpdateItem(itemData)
	return
}

// BaseArmor is the armor a piece had before any enchantment: its template's
// armor. Only a cleric's enchant raises a copy above its template, and that
// bonus is allowed past the wearer's armor cap, so caps are checked here.
func (i *Item) BaseArmor() int {
	if template, ok := Items[i.ItemId]; ok && template != nil {
		return template.Armor
	}
	return i.Armor
}

func (i *Item) ReturnState() string {
	stamStatus := "slightly used"
	if i.MaxUses == Items[i.ItemId].MaxUses {
		stamStatus = "pristine"
	} else if i.MaxUses < (Items[i.ItemId].MaxUses - int(.90*float32(Items[i.ItemId].MaxUses))) {
		stamStatus = "about to break"
	} else if i.MaxUses < (Items[i.ItemId].MaxUses - int(.75*float32(Items[i.ItemId].MaxUses))) {
		stamStatus = "badly damaged"
	} else if i.MaxUses < (Items[i.ItemId].MaxUses - int(.5*float32(Items[i.ItemId].MaxUses))) {
		stamStatus = "well used"
	} else if i.MaxUses < (Items[i.ItemId].MaxUses - int(.25*float32(Items[i.ItemId].MaxUses))) {
		stamStatus = "used"
	}
	return " looks " + stamStatus
}

// ReturnItemInstanceProps Function to return only the modifiable properties
func ReturnItemInstanceProps(item *Item) map[string]interface{} {
	serialList := map[string]interface{}{
		"itemId":     item.ItemId,
		"name":       item.Name,
		"uses":       item.MaxUses,
		"adjustment": item.Adjustment, //  Adjustable by Mages
		"magic":      utils.Btoi(item.Flags["magic"]),
		"spell":      item.Spell,
		"light":      utils.Btoi(item.Flags["light"]),
		"armor":      item.Armor, // Adjustable by Paladins
	}
	if _, ok := item.Flags["infinite"]; ok {
		serialList["infinite"] = utils.Btoi(item.Flags["infinite"])
	}
	if item.StorePrice != 0 {
		serialList["store_price"] = item.StorePrice
	}
	if item.Repairs != 0 {
		serialList["repairs"] = item.Repairs
	}
	if item.Flags["permanent"] {
		serialList["permanent"] = 1
	}
	if item.BoundTo != "" {
		serialList["bound_to"] = item.BoundTo
	}
	if item.ItemType == 9 {
		serialList["contents"] = item.Storage.Jsonify()
	}
	return serialList
}

func (i *Item) Eval() string {

	stringOut := "You study the " + i.Name + " in your minds eye.... \n\n"

	if utils.IntIn(i.ItemType, []int{0, 1, 2, 3, 4}) { // Weapons
		stringOut += "It is a " + config.ItemTypes[i.ItemId] + " weapon. \n" +
			"It deals between " + strconv.Itoa(utils.RollMin(i.NumDice, i.PlusDice)+i.Adjustment) + " and " + strconv.Itoa(utils.RollMax(i.SidesDice, i.NumDice, i.PlusDice)+i.Adjustment) + " damage. \n" +
			"It has " + strconv.Itoa(i.MaxUses) + " uses before it breaks \n."
		if i.Flags["two_handed"] {
			stringOut += "It takes both hands to wield. \n"
		}
	} else if utils.IntIn(i.ItemType, []int{5, 26, 25, 24, 23, 22, 21, 20, 19}) { // Armor
		stringOut += "It is a " + config.ArmorDescription(i.ItemType, i.ArmorClass) + ". \n" +
			"It has " + strconv.Itoa(i.MaxUses) + " uses before it breaks. \n"
	} else if i.ItemType == 17 { // Beverage
		stringOut += "It is a beverage. \n" +
			"It has " + strconv.Itoa(i.MaxUses) + " sips remaining. \n"
	} else if i.ItemType == 18 { // Music
		stringOut += "It is sheet music. \n" +
			"It contains " + i.Spell + ".\n"
	} else if i.ItemType == 16 { // Instrument
		stringOut += "It is an instrument. \n"
	} else if i.ItemType == 15 { // Ammo  //TODO Implement if we do ammo.
		stringOut += "It is ammunition, and not implemented. \n"
	} else if i.ItemType == 6 || i.ItemType == 8 { //device/wand
		stringOut += "It is a " + config.ItemTypes[i.ItemType] + ". \n" +
			"It is charged with " + i.Spell + ".\n" +
			"It has " + strconv.Itoa(i.MaxUses) + " uses remaining. \n"
	} else if i.ItemType == 7 {
		stringOut += "It is a scroll. \n" +
			"It contains " + i.Spell + ".\n"
	} else if i.ItemType == 9 { //chest
		stringOut += "It is a container. \n" +
			"It can hold " + strconv.Itoa(i.MaxUses) + " items. \n"
		if i.Flags["weightless_chest"] {
			stringOut += "It holds it's contents weightlessly. \n"
		}
	} else if i.ItemType == 10 { //gold
		stringOut += "It is gold. \n"
	} else if i.ItemType == 11 { //key
		stringOut += "It is a key. \n"
	} else if i.ItemType == 12 { //light
		stringOut += "It is a light source. \n"
	} else if i.ItemType == 13 { //object
		stringOut += "It's just an object. \n"
	}

	if i.Flags["permanent"] {
		stringOut += "It is permanent. \n"
	}
	if i.Flags["no_take"] {
		stringOut += "It cannot be picked up \n"
	}

	stringOut += "\n You determine its weight to be " + strconv.Itoa(i.Weight) + "lbs. \nYou judge its value to be " + strconv.Itoa(i.Value) + " gold marks."
	return stringOut
}
