package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

func init() {
	addHandler(give{},
           "Usage:  give [person] itemName # \n \n Give the specific person an item.",
           permissions.Player,
           "GIVE")
}

type give cmd

func (give) process(s *state) {

	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Give who what???")
		return
	}

	whoStr := s.words[0]
	targetStr := s.words[1]
	targetNum := 1

	var who *objects.Character
	if s.actor.Permission.HasFlag(permissions.Dungeonmaster) || s.actor.Permission.HasFlag(permissions.Gamemaster) {
		who = s.where.Chars.Search(whoStr, true)
	}else{
		who = s.where.Chars.Search(whoStr, false)
	}
	if who == nil {
		s.msg.Actor.SendInfo("Give who what???")
		return
	}

	// We're going to process a money transaction.
	if strings.HasPrefix("$", targetStr) {
		if amount, err := strconv.ParseInt(strings.Trim(targetStr, "$"), 10, 64); err==nil{
			if s.actor.Gold.CanSubtract(amount){
				s.actor.Gold.SubIfCan(amount)
				who.Gold.Add(amount)
				s.msg.Actor.SendGood("You give ", targetStr ,  " to ", who.Name, ".")
				s.msg.Observer.SendInfo("You see ", s.actor.Name, " give ", who.Name, " some gold.")
			}else{
				s.msg.Actor.SendInfo("You don't have that much gold.")
				return
			}
		}
	}


	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
	}


	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target == nil {
		s.msg.Actor.SendInfo("What're you trying to give away?")
		return
	}

	if (who.Inventory.TotalWeight + target.GetWeight()) <= who.MaxWeight() {
		s.actor.Inventory.Lock()
		who.Inventory.Lock()
		s.actor.Inventory.Remove(target)
		who.Inventory.Add(target)
		s.actor.Inventory.Unlock()
		who.Inventory.Unlock()
	}else{
		s.msg.Actor.SendInfo("They can't carry anymore.")
		return
	}

	s.msg.Actor.SendGood("You give ", target.Name, " to ", who.Name, ".")
	s.msg.Observer.SendInfo("You see ", s.actor.Name, " give ", target.Name, " to ", who.Name, ".")

	s.ok = true
}
