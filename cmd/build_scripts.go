package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(scripts{},
		"Usage:  types  \n Print all of the command scripts that can be attached to things. ",
		permissions.Builder,
		"scripts")
}

type scripts cmd

var ScriptList = map[string]string{
	"$TELEPORTTO": "Usage: $TELEPORTO room_id, message; send player to a different room",
	"$TELEPORT":   "Usage: $TELEPORT message; invokes a randomized teleport with a message",
	"$REPAIR":     "Usage: $REPAIR, allows a player to repair weapons and armor",
	"$POOF":       "Usage: $POOF, announces an arrival of a player",
	"$BALANCE":    "Usage: $BALANCE, provides the player their bank information",
	"$DEPOSIT":    "Usage: $DEPOSIT, Deposits a value from primary gold to bank gold",
	"$WITHDRAW":   "Usage: $WITHDRAW, Withdraws gold from the bank and places it in the users gold pouch",
	"$BUY":        "Usage: $BUY item_name, will exchange a store list item for users gold",
	"$LIST":       "Usage: $LIST, lists all of the items for sale in the store",
	"$SELL":       "Usage: $SELL,  ALlows the user to pawn items.",
	"$ECHO":       "Usage: $ECHO, sends a message to the actor using the command",
	"$ECHOALL":    "Usage: $ECHOALL, sends a messae to the actor and everyone in the room",
	"$TEACH":      "Usage: $TEACH, teaches a spell to the actor",
	"$MELD":       "Usage: $MELD, melds 2 like items together",
	"$SOULBIND":   "Usage: $SOULBIND, binds an item to a player",
}

func (scripts) process(s *state) {
	for key, value := range ScriptList {
		s.msg.Actor.SendInfo(key + "| " + value + "\n")
	}

	s.ok = true
	return
}
