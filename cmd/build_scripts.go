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
	"$TELEPORT": "Usage: $teleport room_id, send player to a different room",
	"$POOF": "Usage: $POOF, announces an arrival of a player",
	"$BALANCE": "Usage: $BALANCE, provides the player their bank information",
	"$DEPOSIT": "Usage: $DEPOSIT, Deposits a value from primary gold to bank gold",
	"$WITHDRAW": "Usage: $WITHDRAW, Withdraws gold from the bank and places it in the users gold pouch",
	"$BUY": "Usage: $BUY item_name, will exchange a store list item for users gold",
	"$LIST": "Usage: $LIST, lists all of the items for sale in the store",
	"$SELL": "Usage: $SELL,  ALlows the user to pawn items.",
	"$ECHO": "Usage: $ECHO, sends a message to the actor using the command",
	"$ECHOALL": "Usage: $ECHOALL, sends a messae to the actor and everyone in the room",
	"$TEACH": "Usage: $TEACH, teaches a spell to the actor",
}


func (scripts) process(s *state) {
	for key, value := range ScriptList {
		s.msg.Actor.SendInfo(key + "| " + value + "\n")
	}

	s.ok = true
	return
}
