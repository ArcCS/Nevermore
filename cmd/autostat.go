package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(reroll{}, "reroll ## ## ## ## ## \n reroll str, con, dex, int, pie use this command to adjust your stats, must enter all of them and all must be values 50+2per level or less", permissions.Player, "reroll")
}

type reroll cmd

func (reroll) process(s *state) {
	if s.actor.Rerolls == 0 {
		s.msg.Actor.SendBad("You have no rerolls left.")
		return
	}
	if len(s.words) < 5 {
		s.msg.Actor.SendBad("You must enter all of the stats to adjust")
		return
	}
	var str, con, dex, intel, pie int
	for i := 0; i < 5; i++ {
		if val, err := strconv.Atoi(s.words[i]); err == nil {
			if i == 0 {
				str = val
			} else if i == 1 {
				con = val
			} else if i == 2 {
				dex = val
			} else if i == 3 {
				intel = val
			} else if i == 4 {
				pie = val
			}
		} else {
			s.msg.Actor.SendBad("Value must be an integer")
			return
		}
	}

	validateStats(s, str, con, dex, intel, pie)
	s.actor.Rerolls--
	s.actor.Str.Current = str
	s.msg.Actor.SendGood("Strength changed to " + strconv.Itoa(str))
	s.actor.Con.Current = con
	s.msg.Actor.SendGood("Constitution changed to " + strconv.Itoa(con))
	s.actor.Dex.Current = dex
	s.msg.Actor.SendGood("Dexterity changed to " + strconv.Itoa(dex))
	s.actor.Int.Current = intel
	s.msg.Actor.SendGood("Intelligence changed to " + strconv.Itoa(intel))
	s.actor.Pie.Current = pie
	s.msg.Actor.SendGood("Piety changed to " + strconv.Itoa(pie))
	s.actor.Stam.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
	s.actor.Stam.Current = s.actor.Stam.Max
	s.actor.Vit.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
	s.actor.Vit.Current = s.actor.Vit.Max
	s.actor.Mana.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
	s.actor.Mana.Current = s.actor.Mana.Max
	return
}

func validateStats(s *state, str int, con int, dex int, intel int, pie int) bool {

	if str+con+dex+intel+pie != 50+((s.actor.Tier*2)-2) {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[s.actor.Race]].StrMin > str || str > config.RaceDefs[config.AvailableRaces[s.actor.Race]].StrMax {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[s.actor.Race]].DexMin > dex || dex > config.RaceDefs[config.AvailableRaces[s.actor.Race]].DexMax {
		return false

	}
	if config.RaceDefs[config.AvailableRaces[s.actor.Race]].ConMin > con || con > config.RaceDefs[config.AvailableRaces[s.actor.Race]].ConMax {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[s.actor.Race]].IntMin > intel || intel > config.RaceDefs[config.AvailableRaces[s.actor.Race]].IntMax {
		return false
	}
	if config.RaceDefs[config.AvailableRaces[s.actor.Race]].PieMin > pie || pie > config.RaceDefs[config.AvailableRaces[s.actor.Race]].PieMax {
		return false
	}
	return true

}
