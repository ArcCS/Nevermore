package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"math/rand"
	"strings"
)

func init() {
	addHandler(split{},
		"",
		permissions.Player,
		"$SPLIT")
}

type split cmd

func (split) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Where do I output the split list?")
		return
	}

	var charList []string
	for _, char := range s.where.Chars.ListAllNoGM() {
		charList = append(charList, char.Name)
	}

	//Randomize the characterList
	rand.Shuffle(len(charList), func(i, j int) {
		charList[i], charList[j] = charList[j], charList[i]
	})

	targetString := strings.Join(s.words[:len(s.words)-2], " ")
	targetItem := s.where.Items.Search(targetString, 1)

	if targetItem != nil {
		s.msg.Actor.SendGood("The item shifts and shows you a new loot order.")
		targetItem.Description = objects.Items[targetItem.ItemId].Description + "\n\nSplit Order:\n====\n" + strings.Join(charList, "\n")
	} else {
		s.msg.Actor.SendBad("Couldn't find the split output.")
		return
	}
}
