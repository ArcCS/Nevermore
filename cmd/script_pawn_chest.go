package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
)

func init() {
	addHandler(sellchest{},
		"",
		permissions.Player,
		"$SELLCHEST")
	addHandler(sellchest_confirm{},
		"",
		permissions.Player,
		"$SELLCHEST_CONFIRM")
}

type sellchest cmd

func (sellchest) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Sell what????")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if len(s.words) == 2 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		}
	}

	// Search the room for the chest
	target := s.where.Items.Search(targetStr, targetNum)

	// Search the player for the chest
	if target == nil {
		target = s.actor.Inventory.Search(targetStr, targetNum)
	}

	if target != nil {
		if target.ItemType != 9 {
			s.msg.Actor.SendInfo("This command is intended to sell the entire contents of a chest.")
			return
		}

		if len(target.Storage.Contents) >= 1 {
			if s.actor.PartyFollow != "" || len(s.actor.PartyFollowers) > 0 {
				s.msg.Actor.SendInfo("The pawn broker nods to the ", target.Name, " and says, 'I'll give you a fair price for the contents, but only for the stuff that hasn't been used.. and you can't back out, if you agree I'll go through and everything and then I'll split the gold to everyone in your party, are you still interested?'")
			} else {
				s.msg.Actor.SendInfo("The pawn broker nods to the ", target.Name, " and says, 'I'll give you a fair price for the contents, but only for the stuff that hasn't been used.. and you can't back out, if you agree I'll go through and everything and toss you the gold, are you still interested?'")
			}
			s.actor.AddCommands("yes", "$SELLCHEST_CONFIRM "+targetStr+" "+strconv.Itoa(targetNum))
			s.actor.AddCommands("y", "$SELLCHEST_CONFIRM "+targetStr+" "+strconv.Itoa(targetNum))
		} else {
			s.msg.Actor.SendBad("There's nothing to sell in there.")
		}
	} else {
		s.msg.Actor.SendInfo("What're you trying to sell from??")
		return
	}

}

type sellchest_confirm cmd

func (sellchest_confirm) process(s *state) {
	log.Println("length of words: ", len(s.words))
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Not enough arguements on sellchest command.")
		return
	}

	targetStr := s.words[0]
	targetNum := 1

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
	}

	// Search the room for the chest
	target := s.where.Items.Search(targetStr, targetNum)

	// Search the player for the chest
	if target == nil {
		target = s.actor.Inventory.Search(targetStr, targetNum)
	}

	if target != nil {
		if target.ItemType != 9 {
			s.msg.Actor.SendInfo("This command is intended to sell the entire contents of a chest.")
			return
		}

		finalValue := 0
		itemValue := 0

		for _, item := range target.Storage.ListItems() {
			log.Println("Item: ", item.Name, " Value: ", item.Value, " MaxUses: ", objects.Items[item.ItemId].MaxUses)
			if target.MaxUses != objects.Items[target.ItemId].MaxUses {
				s.msg.Actor.SendInfo("The pawn broker places your ", item.Name, " back in the ", target.Name, " and says, 'I don't buy used items.'")
				continue
			}

			if s.actor.GetStat("int") < 5 {
				itemValue = int(.10 * float64(item.Value))
			} else if s.actor.GetStat("int") < 10 {
				itemValue = int((.25 * float64(item.Value)) +
					(float64(utils.Roll(10, 1, 0))/float64(10))*(.25*float64(item.Value)))
			} else if s.actor.GetStat("int") >= 10 {
				itemValue = int((.5 * float64(item.Value)) +
					((float64(s.actor.Int.Current) / 45) * (.25 * float64(item.Value))) +
					(float64(utils.Roll(10, 1, 0))/float64(10))*(.25*float64(item.Value)))
			}

			if ok := target.Storage.Remove(item); ok != nil {
				s.msg.Actor.SendBad("Issue completing sell. Issuing accrued gold from sales, and stopping routine")
				log.Println("Error removing item from chest: ", ok)
				s.actor.Gold.Add(finalValue)
				s.msg.Actor.SendGood("The pawn broker gives you ", strconv.Itoa(finalValue), " for the items he was able to take from ", target.Name, ".")
				return
			} else {
				finalValue += itemValue
			}

			data.StoreItemSale(target.ItemId, s.actor.CharId, s.actor.Tier, itemValue)
			data.StoreItemTotals(target.ItemId, 1, itemValue)
		}

		if s.actor.PartyFollow == "" && len(s.actor.PartyFollowers) == 0 {
			s.actor.Gold.Add(finalValue)
			s.msg.Actor.SendGood("The pawn broker gives you ", strconv.Itoa(finalValue), " for the contents of ", target.Name, ".")
		}
		if s.actor.PartyFollow != "" {
			leadChar := objects.ActiveCharacters.Find(s.actor.PartyFollow)
			if leadChar != nil {
				s.participant = leadChar
				goldSplit := finalValue / (len(leadChar.PartyFollowers) + 1)
				msg := "The pawn broker gives you " + strconv.Itoa(goldSplit) + " for the contents of " + target.Name + "."
				s.participant.Gold.Add(goldSplit)
				s.msg.Participant.Send(msg)
				leadChar.MessageParty(msg, s.actor)
				// Give the gold...
				for _, follower := range leadChar.PartyFollowers {
					followerChar := objects.ActiveCharacters.Find(follower)
					if followerChar != nil {
						followerChar.Gold.Add(goldSplit)
					}
				}
				return
			}
		}
		if len(s.actor.PartyFollowers) > 0 {
			goldSplit := finalValue / (len(s.actor.PartyFollowers) + 1)
			msg := "The pawn broker gives you " + strconv.Itoa(goldSplit) + " for the contents of " + target.Name + "."
			s.actor.Gold.Add(goldSplit)
			s.msg.Actor.Send(msg)
			s.actor.MessageParty(msg, s.actor)
			// Give the gold...
			for _, follower := range s.actor.PartyFollowers {
				followerChar := objects.ActiveCharacters.Find(follower)
				if followerChar != nil {
					followerChar.Gold.Add(goldSplit)
				}
			}
			return
		}
		s.actor.Gold.Add(finalValue)
		s.msg.Actor.SendGood("The pawn broker gives you ", strconv.Itoa(finalValue), " for the contents of ", target.Name, ".")

	} else {
		s.msg.Actor.SendInfo("What're you trying to sell??")
		return
	}

}
