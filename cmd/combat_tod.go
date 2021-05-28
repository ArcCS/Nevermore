package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"math"
	"strconv"
)

func init() {
	addHandler(tod{},
		"Usage:  touch target # \n\n Attempt the secret art of a touch of death on your living target",
		permissions.Monk,
		"touch", "tod", "touch-of-death")
}

type tod cmd

func (tod) process(s *state) {
	//TODO Finish TOD command
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Attack what exactly?")
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob != nil {
		// This is an override for a GM to delete a mob
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			s.msg.Actor.SendInfo("You smashed ", whatMob.Name, " out of existence.")
			objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			whatMob = nil
			return
		}

		s.actor.RunHook("combat")

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == (*objects.Item)(nil) && s.actor.Class != 8 {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		if _, err := whatMob.ThreatTable[s.actor.Name]; !err {
			s.msg.Actor.Send(text.White + "You engaged " + whatMob.Name + " #" + strconv.Itoa(s.where.Mobs.GetNumber(whatMob)) + " in combat.")
			whatMob.AddThreatDamage(0, s.actor.Name)
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Equipment.Main.ItemType != 4 && (s.actor.Placement != whatMob.Placement) {
			s.msg.Actor.SendBad("You are too far away to attack.")
			return
		}

		// Lets use a list of attacks,  so we can expand this later if other classes get multi style attacks
		attacks := []float64{
			1.0,
		}

		skillLevel := config.WeaponLevel(s.actor.Skills[s.actor.Equipment.Main.ItemType].Value)

		// Kill is really the fighters realm for specialty..
		if s.actor.Permission.HasAnyFlags(permissions.Fighter) {
			// Did this mofo lethal?
			if config.RollLethal(skillLevel) {
				// Sure did.  Kill this fool and bail.
				s.msg.Actor.SendInfo("You landed a lethal blow on the " + whatMob.Name)
				s.msg.Observers.SendInfo(s.actor.Name + " landed a lethal blow on " + whatMob.Name)
				// Mob died
				//TODO Calculate experience
				stringExp := strconv.Itoa(whatMob.Experience)
				for k := range whatMob.ThreatTable {
					s.where.Chars.Search(k, s.actor).Write([]byte(text.Cyan + "You earn " + stringExp + " exp for the defeat of the " + whatMob.Name + "\n" + text.Reset))
					s.where.Chars.Search(k, s.actor).Experience.Add(whatMob.Experience)
				}
				s.msg.Observers.SendInfo(whatMob.Name + " dies.")
				s.msg.Actor.SendInfo(whatMob.DropInventory())
				objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
				whatMob = nil
				s.actor.SetTimer("combat", 8)
				return
			}

			if skillLevel >= 4 {
				attacks = append(attacks, .15)
				if skillLevel >= 5 {
					attacks[1] = .3
					if skillLevel >= 6 {
						attacks = append(attacks, .15)
						if skillLevel >= 7 {
							attacks[2] = .30
							if skillLevel >= 8 {
								attacks = append(attacks, .15)
								if skillLevel >= 9 {
									attacks[3] = .3
									if skillLevel >= 10 {
										attacks = append(attacks, .30)
									}
								}
							}
						}
					}
				}
			}
		}

		// Lets start executing the attacks
		for _, mult := range attacks {
			// Lets try to crit:
			//TODO: Parry/Miss?
			if config.RollCritical(skillLevel) {
				mult *= float64(config.CombatModifiers["critical"])
				s.msg.Actor.SendGood("Critical Strike!")
				// TODO: Something something shattered weapons something or other
			} else if config.RollDouble(skillLevel) {
				mult *= float64(config.CombatModifiers["double"])
				s.msg.Actor.SendGood("Double Damage!")
			}
			actualDamage, _ := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * mult)))
			whatMob.AddThreatDamage(actualDamage, s.actor.Name)
			s.msg.Actor.SendInfo("You hit the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
			if whatMob.Stam.Current <= 0 {
				s.msg.Actor.SendInfo("You killed " + whatMob.Name + text.Reset)
				s.msg.Observers.SendInfo(s.actor.Name + " killed " + whatMob.Name + text.Reset)
				//TODO Calculate experience
				stringExp := strconv.Itoa(whatMob.Experience)
				for k := range whatMob.ThreatTable {
					s.where.Chars.Search(k, s.actor).Write([]byte(text.Cyan + "You earn " + stringExp + " exp for the defeat of the " + whatMob.Name + "\n" + text.Reset))
					s.where.Chars.Search(k, s.actor).Experience.Add(whatMob.Experience)
				}
				s.msg.Observers.SendInfo(whatMob.Name + " dies.")
				s.msg.Actor.SendInfo(whatMob.DropInventory())
				objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
				whatMob = nil
			}
		}
		s.actor.SetTimer("combat", config.CombatCooldown)
		return

	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}
