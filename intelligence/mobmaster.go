// Copyright 2023 Nevermore.

package intelligence

import (
	"github.com/ArcCS/Nevermore/objects"
)

var ActiveMobs []*objects.Mob

func init() {
	//objects.ActivateMob = ActivateMob
}

func ActivateMob(mob *objects.Mob) {
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
