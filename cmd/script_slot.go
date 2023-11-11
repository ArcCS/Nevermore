package cmd

import (
	//	"log"
	//"strconv"
	"strings"

	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	//	"github.com/jinzhu/copier"
)

func init() {
	addHandler(slot{},
		"",
		permissions.Player,
		"$SLOT")
}

type slot cmd

var symbols = map[int]string{
	1:  "The head of a unicorn against a flame and crossed by a ribbon of cerulean blue.",
	2:  "The silhouette of a griffon flying toward a black and crimson moon against a violet twilight sky.",
	3:  "Two rings, one gold, one platinum, overlapping on the edge to form the symbol of infinity.",
	4:  "Three Pointed Star.",
	5:  "A bunch of grapes surrounded by a wreath of wildflowers.",
	6:  "A cracked crystal goblet.",
	7:  "A Crescent Moon wrapped around an 8-pointed star, with the star's four primary points extending past the crescent.",
	8:  "A slender black dagger against a grey background.",
	9:  "Doves, White Roses, Lilies.",
	10: "A Longbow diagonally over a vertical Claymore",
	11: "An eye, set against an open tomb.",
	12: "A white four-pointed star obscured by cloud.",
	13: "An angel in silver robes placing a golden crown on his head, with the hand holding the crown scaled and Daemonic.",
	14: "A set of unbalanced scales.",
	15: "Crossed Swords, one of Steel (Valor), one of Silver (Honor).",
	// Some Gods don't have symbols
	//	16:
	//	17:
	//	18:

}

func (slot) process(s *state) {
	/*if len(s.words) < 1 {
		s.msg.Actor.SendBad("Buy what?")
		return
	}

	targetStr := ""
	targetNum := 1

	for _, word := range s.words {
		if val, err := strconv.Atoi(word); err == nil {
			targetNum = val
		} else {
			targetStr += " " + word
		}
	}

	targetStr = strings.Trim(targetStr, " ")

	if len(s.where.StoreInventory.Contents) > 0 {
		purchaseItem := s.where.StoreInventory.Search(targetStr, targetNum)
		if purchaseItem != nil {
			if s.actor.Gold.Value >= purchaseItem.StorePrice {
				if (s.actor.GetCurrentWeight() + purchaseItem.GetWeight()) <= s.actor.MaxWeight() {
					s.actor.RunHook("act")
					s.actor.Gold.Subtract(purchaseItem.StorePrice)
					if purchaseItem.Flags["infinite"] {
						newItem := objects.Item{}
						if err := copier.CopyWithOption(&newItem, objects.Items[purchaseItem.ItemId], copier.Option{DeepCopy: true}); err != nil {
							log.Println("Error copying item: ", err)
						}
						s.actor.Inventory.Add(&newItem)
					} else {
						if err := s.where.StoreInventory.Remove(purchaseItem); err != nil {
							log.Println("Error removing item from store: ", err)
						}
						s.actor.Inventory.Add(purchaseItem)
					}
					s.msg.Actor.SendGood("You purchase ", purchaseItem.Name, ".")
					s.msg.Observers.SendInfo(s.actor.Name, " purchases ", purchaseItem.Name, ".")
					return
				} else {
					s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
					return
				}
			} else {
				s.msg.Actor.SendInfo("You don't have enough gold to purchase that item.")
				return
			}
		} else {
			s.msg.Actor.SendInfo("You don't see that item to purchase.")
			return
		}
	} else {
		s.msg.Actor.SendBad("There's nothing to purchase at this store.")
		return
	}
	*/
	if s.actor.Gold.Value >= 100 {
		s.actor.Gold.Subtract(100)
		var results [4]int
		for i := 0; i < 4; i++ {
			results[i] = utils.Roll(15, 1, 0)
			//s.msg.Actor.SendInfo(strconv.Itoa(results[i]))
		}
		s.msg.Observers.SendInfo(utils.Title(strings.ToLower(strings.Join(s.words[1:], " "))))
		s.msg.Actor.SendInfo("You hear some whirling as the wheels start to spin. As they come to a stop symbols begin to glow.")
		s.msg.Actor.SendGood("The first wheel shows: ", text.Cyan+symbols[results[0]])
		s.msg.Actor.SendGood("The second wheel shows: ", text.Cyan+symbols[results[1]])
		s.msg.Actor.SendGood("The third wheel shows: ", text.Cyan+symbols[results[2]])
		s.msg.Actor.SendGood("The fourth wheel shows: ", text.Cyan+symbols[results[3]])

		jackpot := s.where.Items.Search(s.words[0], 1)
		if jackpot == nil {
			s.msg.Actor.SendInfo("Machine is under maintenance")
			return
		}

		uniqueValues := make(map[int]bool)
		for _, v := range results {
			uniqueValues[v] = true
		}
		switch len(uniqueValues) {
		case 1:
			{
				s.msg.Actor.SendGood("Jackpot!!!!")
				s.actor.Gold.Add(jackpot.Armor)
				objects.ActiveCharacters.MessageAll("###: You hear the distinctive Cha-Ching of the slot machine and know somebody has won the jackpot!")
				jackpot.Armor = 10000
			}
		case 2:
			{
				s.msg.Actor.SendGood("Multiple matches!! Payout:$1000")
				s.actor.Gold.Add(1000)
				jackpot.Armor -= 1000
				if jackpot.Armor < 10000 {
					jackpot.Armor = 10000
				}
			}
		case 3:
			{
				s.msg.Actor.SendGood("A match! Payout:$100")
				s.actor.Gold.Add(100)
				if jackpot.Armor < 10000 {
					jackpot.Armor = 10000
				}

			}
		case 4:
			{
				s.msg.Actor.SendGood("No matches, better luck next time")
				jackpot.Armor += 75
			}
		}
	} else {
		s.msg.Actor.SendInfo("You don't have enough gold to play")
	}

}
