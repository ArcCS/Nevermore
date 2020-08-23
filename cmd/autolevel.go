package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(autolevel{}, "autolevel ## \n PLAYTEST ONLY,  use this command to auto adjust your level to the chosen level, all attributes are reset to default and attribute points are assigned for you to apply with autostat", permissions.Player,  "autolevel", "autotier")
}

type autolevel cmd

func (autolevel) process(s *state) {
		if len(s.words) < 1 {
			s.msg.Actor.SendBad("You must enter the level that you want")
			return
		}

		if val, err := strconv.Atoi(s.words[0]); err == nil {
			if val <= 25 {
				s.actor.Tier = val
			}else{
				s.msg.Actor.SendBad("Value must be 25 or below.")
				return
			}
		}else{
			s.msg.Actor.SendBad("Value must be an integer")
			return
		}
		s.actor.Stam.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
		s.actor.Stam.Current = s.actor.Stam.Max
		s.actor.Vit.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
		s.actor.Vit.Current = s.actor.Vit.Max
		s.actor.Mana.Max = config.CalcStamina(s.actor.Tier, s.actor.Con.Current, s.actor.Class)
		s.actor.Mana.Current = s.actor.Mana.Max
		s.msg.Actor.SendGood("Tier changed to " + strconv.Itoa(s.actor.Tier))
		return
}
