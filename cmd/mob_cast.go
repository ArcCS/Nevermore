package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(mobCast{},
		"",
		permissions.None,
		"$MCAST")
}

type mobCast cmd

func (mobCast) process(s *state) {
	/*
			if m.CurrentTarget != "" && m.ChanceCast > 0 {
				// Try to cast a spell first
				log.Println("High chance to cast, trying to cast a spell")
				target := Rooms[m.ParentId].Chars.MobSearch(m.CurrentTarget, m)
				spellSelected := false
				selectSpell := ""
				if utils.Roll(100, 1, 0) <= m.ChanceCast {
					log.Println("Successful Roll, trying to cast a spell")
					for range m.Spells {
						rand.Seed(time.Now().Unix())
						selectSpell = m.Spells[rand.Intn(len(m.Spells))]
						if selectSpell != "" {
							if utils.StringIn(selectSpell, OffensiveSpells) {
								if m.Mana.Current > Spells[selectSpell].Cost {
									spellSelected = true
								}
							}
						}
					}

					if spellSelected {
						spellInstance, ok := Spells[selectSpell]
						if !ok {
							spellSelected = false
						}
						Rooms[m.ParentId].MessageAll(m.Name + " casts a " + spellInstance.Name + " spell on " + target.Name + "\n")
						target.RunHook("attacked")
						m.Mana.Subtract(spellInstance.Cost)
						result := Cast(m, target, spellInstance.Effect, spellInstance.Magnitude)
						if strings.Contains(result, "$SCRIPT") {
							m.MobScript(result)
						}
						target.DeathCheck("was slain by a " + m.Name + ".")
						return
					}
				}
			}
		}
	*/
	return
}
