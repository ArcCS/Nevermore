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
	addHandler(kill{},
           "Usage:  kill target # \n\n Try to attack something! Can also use attack, or shorthand k",
           permissions.Player,
           "kill")
}

type kill cmd

func (kill) process(s *state) {
	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum,true)
	if whatMob != nil {
		// This is an override for a GM to delete a mob
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			s.msg.Actor.SendInfo("You smashed ", whatMob.Name, " out of existence.")
			objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			whatMob = nil
			return
		}

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil && s.actor.Class != 7 {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		if _, err := whatMob.ThreatTable[s.actor.Name]; err {
			s.msg.Actor.Send(text.White + "You engaged " + whatMob.Name + " #" + strconv.Itoa(s.where.Mobs.GetNumber(whatMob)) + " in combat.")
			whatMob.AddThreatDamage(0, s.actor.Name)
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Equipment.Main.Type != 4 && (s.actor.Placement != whatMob.Placement) {
			s.msg.Actor.SendBad("You are too far away to attack.")
			return
		}

		// Lets use a list of attacks,  so we can expand this later if other classes get multi style attacks
		attacks := []float64{
			1.0,
		}

		skillLevel := config.WeaponLevel(s.actor.Equipment.Main.Type)

		// Kill is really the fighters realm for specialty..
		if s.actor.Permission.HasAnyFlags(permissions.Fighter) {
			// Did this mofo lethal?
			if config.RollLethal(skillLevel) {
				// Sure did.  Kill this fool and bail.
				s.msg.Actor.SendInfo("You landed a lethal blow on the " + whatMob.Name)
				s.msg.Observers.SendInfo(s.actor.Name + " landed a lethal blow on " + whatMob.Name)
				whatMob.Died()
				objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
				whatMob = nil
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
				s.msg.Actor.SendGood("Crital Strike!")
				// TODO: Something something shattered weapons something or other
			}else if config.RollDouble(skillLevel) {
				mult *= float64(config.CombatModifiers["double"])
				s.msg.Actor.SendGood("Double Damage!")
			}
			actualDamage := whatMob.ReceiveDamage(int64(math.Ceil(float64(s.actor.InflictDamage()) * mult)))
			s.msg.Actor.SendInfo("You hit the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!")
			if whatMob.Stam.Current <= 0 {
				s.msg.Actor.SendInfo("You killed  " + whatMob.Name)
				s.msg.Observers.SendInfo(s.actor.Name + " killed " + whatMob.Name)
				whatMob.Died()
				objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			}
		}

	}

	s.msg.Actor.SendInfo("Attack what?")
	s.ok = true
}
