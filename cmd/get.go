package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(get{},
		"Usage:  get itemName [container_name] # \n \n Get the specified item.",
		permissions.Player,
		"GET", "TAKE", "G")
}

type get cmd

func (get) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("You go to get.. uh??")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	checkWeight := true
	targetStr := s.words[0]
	targetNum := 1
	whereStr := ""
	whereNum := 1

	if len(s.words) > 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		} else {
			whereStr = s.words[1]
		}
	}
	if len(s.words) > 2 {
		if whereStr != "" {
			if val2, err2 := strconv.Atoi(s.words[2]); err2 == nil {
				whereNum = val2
			} else {
				whereStr = s.words[1]
			}
		}
	}

	if len(s.words) > 3 {
		if val3, err3 := strconv.Atoi(s.words[3]); err3 == nil {
			whereNum = val3
		}
	}

	if whereStr == "" {
		roomInventory := s.where.Items.Search(targetStr, targetNum)
		if roomInventory != nil {
			if roomInventory.Placement != s.actor.Placement {
				s.msg.Actor.SendBad("You must be next to the item to get it.")
				return
			}
			if roomInventory.Flags["no_take"] && !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				s.msg.Actor.SendBad("You cannot take that!")
				return
			}
			// Check the mobs in the room to see if any of them block
			for _, mob := range s.where.Mobs.Contents {
				if mob.Flags["guard_treasure"] && mob.Placement == s.actor.Placement {
					s.msg.Actor.SendBad(mob.Name + " prevents you from picking up " + roomInventory.DisplayName())
					return
				}
			}
			if roomInventory.ItemType == 10 {
				s.actor.RunHook("act")
				s.where.Items.Remove(roomInventory)
				s.actor.Gold.Add(roomInventory.Value)
				s.msg.Actor.SendGood("You picked up ", strconv.Itoa(roomInventory.Value), " gold pieces.")
				if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					s.msg.Observers.SendInfo("You see ", s.actor.Name, " get ", roomInventory.Name, ".")
				}
				return
			} else if (s.actor.GetCurrentWeight()+roomInventory.GetWeight()) <= s.actor.MaxWeight() || s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
				s.actor.RunHook("act")
				s.where.Items.Remove(roomInventory)
				s.actor.Inventory.Add(roomInventory)
				s.msg.Actor.SendGood("You get ", roomInventory.Name, ".")
				if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					s.msg.Observers.SendInfo(s.actor.Name, " takes ", roomInventory.Name, ".")
				}
				if roomInventory.Flags["permanent"] {
					s.where.Save()
				}
				return
			} else {
				s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
				return
			}
		}
	} else {
		where := s.where.Items.Search(whereStr, whereNum)

		if where != nil && where.ItemType == 9 {
			if where.Placement != s.actor.Placement {
				s.msg.Actor.SendBad("You must be next to the chest to get items from it.")
				return
			}
		}

		// If we didn't find it in the room, look on the person.
		if where == nil {
			where = s.actor.Inventory.Search(whereStr, whereNum)
			if where != nil {
				if !where.Flags["weightless_chest"] {
					checkWeight = false
				}
			}
		}

		if where != nil {
			whereInventory := where.Storage.Search(targetStr, targetNum)
			if whereInventory != nil {
				if whereInventory.ItemType == 10 {
					s.actor.RunHook("act")
					where.Storage.Remove(whereInventory)
					s.actor.Gold.Add(whereInventory.Value)
					s.msg.Actor.SendGood("You take ", whereInventory.Name, " from ", where.Name, " and put it in your gold pouch.")
					if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
						s.msg.Observers.SendInfo("You see ", s.actor.Name, " take ", whereInventory.Name, " from ", where.Name, ".")
					}
					return
				} else if (!checkWeight) || (checkWeight && (s.actor.GetCurrentWeight()+whereInventory.GetWeight() <= s.actor.MaxWeight())) || s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					s.actor.RunHook("act")
					where.Storage.Remove(whereInventory)
					s.actor.Inventory.Add(whereInventory)
					s.msg.Actor.SendGood("You take ", whereInventory.Name, " from ", where.Name, ".")
					if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
						s.msg.Observers.SendInfo("You see ", s.actor.Name, " take ", whereInventory.Name, " from ", where.Name, ".")
					}
					return
				} else {
					s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
					return
				}
			}
		}
	}

	s.msg.Actor.SendBad("Get what?")
	s.ok = true
}
