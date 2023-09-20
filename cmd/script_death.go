package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
	"log"
	"strconv"
	"strings"
	"time"
)

func init() {
	addHandler(scriptDeath{},
		"",
		permissions.Anyone,
		"$DEATH")
}

type scriptDeath cmd

func (scriptDeath) process(s *state) {

	healingHand := objects.Rooms[config.HealingHand]
	if !utils.IntIn(healingHand.RoomId, s.rLocks) {
		s.AddLocks(healingHand.RoomId)
		s.ok = false
		return
	}

	if time.Now().Sub(s.actor.LastAction).Seconds() < 60 {
		deathString := "### " + s.actor.Name + " has died."
		if len(s.words[0]) > 0 {
			deathString = "### " + s.actor.Name + " " + strings.Join(s.input[0:], " ")
		}

		objects.ActiveCharacters.MessageAll("### An otherworldly bell sounds once, the note echoing in your soul")
		objects.ActiveCharacters.MessageAll(deathString)

		if s.actor.Tier > config.FreeDeathTier {
			equipment := s.actor.Equipment.UnequipAll()

			var tempStore []*objects.Item
			for _, item := range s.actor.Inventory.Contents {
				tempStore = append(tempStore, item)
			}

			newItem := objects.Item{}
			copier.CopyWithOption(&newItem, objects.Items[1], copier.Option{DeepCopy: true})
			newItem.Name = "corpse of " + s.actor.Name
			newItem.Description = "It's the corpse of " + s.actor.Name + "."
			newItem.Placement = s.actor.Placement
			if len(tempStore) != 0 {
				for _, item := range tempStore {
					if !item.Flags["permanent"] {
						s.actor.Inventory.Remove(item)
						newItem.Storage.Add(item)
					}
				}
			}
			if len(equipment) != 0 {
				for _, item := range equipment {
					if !item.Flags["permanent"] {
						newItem.Storage.Add(item)
					}
				}
			}
			if s.actor.Gold.Value > 0 {
				newGold := objects.Item{}
				copier.CopyWithOption(&newGold, objects.Items[3456], copier.Option{DeepCopy: true})
				newGold.Name = strconv.Itoa(s.actor.Gold.Value) + " gold marks"
				newGold.Value = s.actor.Gold.Value
				newItem.Storage.Add(&newGold)
				s.actor.Gold.Value = 0
			}
			s.where.MessageAll("The lifeless body of " + s.actor.Name + " falls to the ground.\n\n" + text.Reset)
			s.where.Items.Add(&newItem)
		} else {
			s.msg.Actor.Send(text.Green + "An apprentice aura protects you from the worst of this death and ferries you and your gear safely to the healing hand...")
		}

		s.where.Chars.Remove(s.actor)
		healingHand.Chars.Add(s.actor)
		s.actor.Placement = 3
		s.actor.ParentId = healingHand.RoomId

		s.actor.Write([]byte(text.Cyan + "In what seems like a dream, an imposing black gate shrouded in fog speeds into view.. There is nothing else here to greet you, except a sorrowful sense of loneliness and longing... A chilling thought claws at the inside of your skull, behind your eyes, that this scene isn't right.. and just as swiftly as you arrived, the gate races past... and you awaken in another place..\n\n\n " + text.Reset))
		s.actor.RemoveEffect("blind")
		s.actor.RemoveEffect("poison")
		s.actor.RemoveEffect("disease")
		s.actor.Stam.Current = s.actor.Stam.Max
		s.actor.Vit.Current = s.actor.Vit.Max
		s.actor.Mana.Current = s.actor.Mana.Max

		totalExpNeeded := config.MaxLoss(s.actor.Tier)
		finalMin := config.TierExpLevels[s.actor.Tier+1] - int(float64(totalExpNeeded)*1.2)
		// Determine the death penalty
		if s.actor.Tier > config.FreeDeathTier {
			deathRoll := utils.Roll(100, 1, 0)
			switch {
			case deathRoll <= 20: // Light Passage
				s.actor.Write([]byte(text.Green + "You've pass through this death with minimal effects. (10% xp loss) \n\n" + text.Reset))
				s.actor.Experience.SubMax(int(float64(totalExpNeeded)*.10), finalMin)
				break
			case deathRoll <= 100: // Medium Passage
				s.actor.Write([]byte(text.Green + "The death did not come easy. (30% xp loss)\n\n" + text.Reset))
				s.actor.Experience.SubMax(int(float64(totalExpNeeded)*.30), finalMin)
				break
			}
		}

		s.actor.DeathInProgress = false
		s.scriptActor("LOOK")

	} else {
		deathString := "### " + s.actor.Name + " died a lag death."

		objects.ActiveCharacters.MessageAll("### An otherworldly bell attempts to ring but is abruptly muffled.")
		objects.ActiveCharacters.MessageAll(deathString)

		go func() {
			log.Println("Lag Death: Clean Room")
			s.where.Chars.Remove(s.actor)
			healingHand.Chars.Add(s.actor)
			s.actor.RemoveEffect("blind")
			s.actor.RemoveEffect("poison")
			s.actor.RemoveEffect("disease")
			s.actor.Stam.Current = s.actor.Stam.Max
			s.actor.Vit.Current = s.actor.Vit.Max
			s.actor.Mana.Current = s.actor.Mana.Max
			s.actor.Placement = 3
			s.actor.ParentId = healingHand.RoomId
		}()
	}

}
