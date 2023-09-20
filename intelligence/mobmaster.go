// Copyright 2023 Nevermore.

package intelligence

import (
	"github.com/ArcCS/Nevermore/objects"
	"log"
)

var ActiveMobs []*objects.Mob

func init() {
	objects.ActivateMob = ActivateMob
}

func ActivateMob(mob *objects.Mob) {
	log.Println("Adding mob to active mobs: ", mob.MobId)
	ActiveMobs = append(ActiveMobs, mob)
}

func DeactivateMob(mob *objects.Mob) {
	for c, p := range ActiveMobs {
		if p == mob {
			copy(ActiveMobs[c:], ActiveMobs[c+1:])
			ActiveMobs = ActiveMobs[:len(ActiveMobs)-1]
			break
		}
	}
}
