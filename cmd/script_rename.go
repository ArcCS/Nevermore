package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

func init() {
	addHandler(scriptRename{},
		"",
		permissions.Player,
		"$RENAME")
	addHandler(confirmName{},
		"",
		permissions.Player,
		"$CONFIRMNAME")
	addHandler(confirmDesc{},
		"",
		permissions.Player,
		"$CONFIRMDESC")
	addHandler(confirmRename{},
		"",
		permissions.Player,
		"$CONFIRMRENAME")
}

type scriptRename cmd

func (scriptRename) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Rename what????")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if len(s.words) == 2 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target != nil {
		s.msg.Actor.SendGood("It costs ", strconv.Itoa(config.RenameCost), " gold marks to rename "+target.Name+".")
		s.msg.Actor.SendGood("Accept offer? (y, yes to confirm)")
		s.actor.AddCommands("yes", "$CONFIRMNAME "+targetStr+" "+strconv.Itoa(targetNum))
		s.actor.AddCommands("y", "$CONFIRMNAME "+targetStr+" "+strconv.Itoa(targetNum))
	} else {
		s.msg.Actor.SendInfo("What're you trying to rename??")
		return
	}
}

type confirmName cmd

func (confirmName) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Issue completing rename")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target != nil {
		if s.actor.Gold.Value < config.RenameCost {
			s.msg.Actor.SendBad("You do not have enough gold to rename this item.")
			return
		}
		s.msg.Actor.SendGood("Please enter the description of your new item first with the `description` directive,  \nExample: description my lovely customized item. \n Or simply stop to cancel. \n !! You must stop NOW in order to not be charged. !!")
		s.actor.AddCommands("description", "$CONFIRMDESC "+targetStr+" "+strconv.Itoa(targetNum))
	} else {
		s.msg.Actor.SendInfo("What're you trying to rename??")
		return
	}
}

type confirmDesc cmd

func (confirmDesc) process(s *state) {
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Issue completing rename")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
	}

	newDescription := strings.Join(s.input[2:], " ")

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target != nil {
		if s.actor.Gold.Value < config.RenameCost {
			s.msg.Actor.SendBad("You do not have enough gold to rename this item.")
			return
		}
		s.msg.Actor.SendGood("The description of your item is now " + newDescription + ".")
		target.Description = newDescription
		s.msg.Actor.SendGood("Please enter the name of your item first with the `name` directive \nExample: name my shiny item. \n Or simply stop to cancel.")
		s.actor.AddCommands("name", "$CONFIRMRENAME "+targetStr+" "+strconv.Itoa(targetNum))
	} else {
		s.msg.Actor.SendInfo("What're you trying to rename??")
		return
	}
}

type confirmRename cmd

func (confirmRename) process(s *state) {
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Issue completing rename")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
	}

	newName := strings.Join(s.input[2:], " ")

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target != nil {
		if s.actor.Gold.Value < config.RenameCost {
			s.msg.Actor.SendBad("You do not have enough gold to rename this item.")
			return
		}
		s.actor.Gold.Subtract(config.RenameCost)
		target.Name = newName
		s.msg.Actor.SendGood("The name of your item is " + newName + ".")
	} else {
		s.msg.Actor.SendInfo("What're you trying to rename??")
		return
	}
}
