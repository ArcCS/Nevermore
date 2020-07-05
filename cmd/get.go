package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(get{},
           "Usage:  get [container_name] itemName # \n \n Get the specified item.",
           permissions.Player,
           "GET")
}

type get cmd

func (get) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("You go to get.. uh??")
		return
	}

	argParse := 1
	targetStr := s.words[0]
	targetNum := 1
	whereStr := ""
	whereNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
		argParse += 1
	}

	if len(s.words) > 2 {
		whereStr = s.words[argParse]
		argParse += 1
	}

	if val, err := strconv.Atoi(s.words[argParse]); err == nil {
		whereNum = val
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target == nil {
		s.msg.Actor.SendInfo("What're you trying to put?")
		return
	}

	// Try to find the where if it's not the room
	if whereStr != "" {
		where := s.where.Items.Search(whereStr, whereNum)

		if where == nil {
			// Try to find it on us next.
			where = s.actor.Inventory.Search(whereStr, whereNum)
		}

		// Is where still nil?
		if where == nil {
			s.msg.Actor.SendInfo("Put it where?")
			return
		}

		// Do you specify itself?
		if target == where {
			s.msg.Actor.SendInfo("What're you even trying to do?")
			return
		}

		// is it a chest?
		if !where.Flags["chest"] {
			s.msg.Actor.SendInfo("That's not a chest")
			return
		}

		if (s.actor.Inventory.TotalWeight + target.GetWeight()) <= s.actor.MaxWeight() {
			where.Storage.Lock()
			s.actor.Inventory.Lock()
			where.Storage.Remove(target)
			s.actor.Inventory.Add(target)
			where.Storage.Unlock()
			s.actor.Inventory.Unlock()
		}else{
			s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
			return
		}

		s.msg.Actor.SendGood("You get ", target.Name, " from ", where.Name, ".")
		s.msg.Observers.SendInfo("You see ", s.actor.Name, " put ", target.Name, " into ", where.Name, ".")
	}else {
		where := s.where

		if (s.actor.Inventory.TotalWeight + target.GetWeight()) <= s.actor.MaxWeight() {
			where.Items.Lock()
			s.actor.Inventory.Lock()
			where.Items.Remove(target)
			s.actor.Inventory.Add(target)
			where.Items.Unlock()
			s.actor.Inventory.Unlock()
		}else{
			s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
			return
		}


		s.msg.Actor.SendGood("You take ", target.Name, ".")
		s.msg.Observers.SendInfo("You see ", s.actor.Name, " take ", target.Name, ".")
	}

	s.ok = true
}

