package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(autostat{}, "autostat ## ## ## ## ## \n PLAYTEST ONLY, str, con, dex, int, pie use this command to adjust your stats, must enter all of them and all must be values 50 or less", permissions.Player, "autostat")
}

type autostat cmd

//TODO: DISABLE THIS,  do not go into production

func (autostat) process(s *state) {
	if len(s.words) < 5 {
		s.msg.Actor.SendBad("You must enter all of the stats to adjust")
		return
	}
	for i := 0; i < 5; i++ {
		if val, err := strconv.Atoi(s.words[i]); err == nil {
			if val <= 50 {
				if i == 0 {
					s.actor.Str.Current = val
					s.msg.Actor.SendGood("Strength changed to " + strconv.Itoa(val))
				} else if i == 1 {
					s.actor.Con.Current = val
					s.msg.Actor.SendGood("Constitution changed to " + strconv.Itoa(val))
				} else if i == 2 {
					s.actor.Dex.Current = val
					s.msg.Actor.SendGood("Dexterity changed to " + strconv.Itoa(val))
				} else if i == 3 {
					s.actor.Int.Current = val
					s.msg.Actor.SendGood("Intelligence changed to " + strconv.Itoa(val))
				} else if i == 4 {
					s.actor.Pie.Current = val
					s.msg.Actor.SendGood("Piety changed to " + strconv.Itoa(val))
				}
			} else {
				s.msg.Actor.SendBad("Value must be 25 or below.")
				return
			}
		} else {
			s.msg.Actor.SendBad("Value must be an integer")
			return
		}
	}
	s.actor.Stam.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
	s.actor.Stam.Current = s.actor.Stam.Max
	s.actor.Vit.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
	s.actor.Vit.Current = s.actor.Vit.Max
	s.actor.Mana.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
	s.actor.Mana.Current = s.actor.Mana.Max
	return
}
