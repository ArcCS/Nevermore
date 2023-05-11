package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(put{},
		"Usage:  put itemName # [chest] \n \n Put the specified item name in a chest.",
		permissions.Player,
		"PUT")
}

type put cmd

func (put) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("You go to put something into something else...")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	if len(s.words) == 1 {
		s.msg.Actor.SendInfo("Put it where?")
		return
	}

	// We have at least 2 items here so lets move forward with that
	argParse := 1
	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
		argParse = 2
	}

	if argParse == 2 && len(s.words) <= 2 {
		s.msg.Actor.SendInfo("Put it where?")
		return
	}

	whereStr := s.words[argParse]
	whereNum := 1

	if len(s.words) >= argParse+2 {
		if val, err := strconv.Atoi(s.words[argParse+1]); err == nil {
			whereNum = val
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target == nil {
		s.msg.Actor.SendInfo("What're you trying to put?")
		return
	}

	if target.Flags["permament"] {
		s.msg.Actor.SendBad("You cannot get rid of this item.. it is bound to you.")
		return
	}

	where := s.where.Items.Search(whereStr, whereNum)

	if where == nil {
		// Try to find it on us next.
		where = s.actor.Inventory.Search(whereStr, whereNum)
		// Is where still nil?
		if where == nil {
			s.msg.Actor.SendInfo("Put it where?")
			return
		}
	}

	// Do you specify itself?
	if target == where {
		s.msg.Actor.SendInfo("You can't put something inside of itself...")
		return
	}

	// is it a chest?
	if where.ItemType != 9 {
		s.msg.Actor.SendInfo("You can't put anything in that.")
		return
	}

	if target.ItemType == 9 {
		s.msg.Actor.SendInfo("You cannot put a container into another container.")
		return
	}

	if len(where.Storage.Contents) >= where.MaxUses {
		s.msg.Actor.SendBad("That container is full.")
		return
	}

	s.actor.RunHook("act")

	s.actor.Inventory.Remove(target)
	where.Storage.Add(target)

	s.msg.Actor.SendGood("You put ", target.Name, " into ", where.Name, ".")
	s.msg.Observers.SendInfo("You see ", s.actor.Name, " put ", target.Name, " into ", where.Name, ".")

	s.ok = true
}
