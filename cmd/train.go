package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
	"strings"
)

func init() {
	addHandler(train{}, "train stat stat \n Use this command at the trainer to train to the next level. \n Stat options of str dex con pie int to advance into the next level.", permissions.Player, "TRAIN")
}

type train cmd

func (train) process(s *state) {
	if !s.where.Flags["train"] {
		s.msg.Actor.SendBad("You must find a training location in order to advance your tier.")
		return
	}
	if !(s.actor.Experience.Value >= config.TierExpLevels[s.actor.Tier+1]) {
		s.msg.Actor.SendBad("You don't have enough experience earned to train to the next tier.")
		return
	}
	if !(s.actor.Gold.Value >= config.GoldPerLevel[s.actor.Tier+1]) {
		s.msg.Actor.SendBad("You don't have enough gold to train to the next tier. (" + strconv.Itoa(config.GoldPerLevel[s.actor.Tier+1]) + ")")
		return
	}
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("You must enter both of the stat points you wish to advance into your coming tier.")
		return
	} else if len(s.words) > 2 {
		s.msg.Actor.SendBad("You can only train two stats at a time.")
		return
	}
	message := ""

	if !validateStats(s, s.actor.Str.Current, s.actor.Con.Current, s.actor.Dex.Current, s.actor.Int.Current, s.actor.Pie.Current) {
		s.msg.Actor.SendBad("Stats are not valid, you cannot train this character")
		return
	}

	if !utils.StringIn(strings.ToLower(s.words[0]), []string{"str", "dex", "con", "int", "pie"}) || !utils.StringIn(strings.ToLower(s.words[1]), []string{"str", "dex", "con", "int", "pie"}) {
		s.msg.Actor.SendBad("You must enter a valid stat to train. (pie, int, con, dex, str)")
		return
	}
	for _, val := range s.input {
		proc := strings.ToLower(val)
		// Verify Moves
		if proc == "str" {
			if s.actor.Str.Current == s.actor.Str.Max {
				s.msg.Actor.SendBad("You've already maxed out that stat.")
				return
			}
		} else if proc == "dex" {
			if s.actor.Dex.Current == s.actor.Dex.Max {
				s.msg.Actor.SendBad("You've already maxed out that stat.")
				return
			}
		} else if proc == "con" {
			if s.actor.Con.Current == s.actor.Con.Max {
				s.msg.Actor.SendBad("You've already maxed out that stat.")
				return
			}
		} else if proc == "int" {
			if s.actor.Int.Current == s.actor.Int.Max {
				s.msg.Actor.SendBad("You've already maxed out that stat.")
				return
			}
		} else if proc == "pie" {
			if s.actor.Pie.Current == s.actor.Pie.Max {
				s.msg.Actor.SendBad("You've already maxed out that stat.")
				return
			}
		}
	}

	for count, val := range s.input {
		log.Println(count)
		if message != "" {
			message += " and "
		}
		proc := strings.ToLower(val)
		// Process Moves
		if proc == "str" {
			s.actor.Str.Current += 1
			message += "your strength"
		} else if proc == "dex" {
			s.actor.Dex.Current += 1
			message += "your dexterity"
		} else if proc == "con" {
			s.actor.Con.Current += 1
			message += "your constitution"
		} else if proc == "int" {
			s.actor.Int.Current += 1
			message += "your intelligence"
		} else if proc == "pie" {
			s.actor.Pie.Current += 1
			message += "your piety"
		}
	}

	s.actor.Tier += 1
	s.actor.Gold.Subtract(config.GoldPerLevel[s.actor.Tier])
	s.actor.Stam.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
	s.actor.Stam.Current = s.actor.Stam.Max
	s.actor.Vit.Max = config.CalcHealth(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
	s.actor.Vit.Current = s.actor.Vit.Max
	s.actor.Mana.Max = config.CalcMana(s.actor.Tier, s.actor.Int.Current, s.actor.Class)
	s.actor.Mana.Current = s.actor.Mana.Max
	s.actor.ClassTitle = config.ClassTitle(s.actor.Class, s.actor.Gender, s.actor.Tier)
	s.msg.Actor.SendGood(utils.Title(message + " were increased by 1 and tier increased to " + strconv.Itoa(s.actor.Tier)))
	return
}
