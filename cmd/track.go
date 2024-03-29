package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(track{},
		"Usage:  Try to determine if there are tracks of creatures the frequent this area",
		permissions.Ranger,
		"track")
}

type track cmd

func (track) process(s *state) {
	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}
	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")

		return
	}
	if s.actor.Tier < config.MinorAbilityTier {
		s.msg.Actor.SendBad("You must be at least tier " + strconv.Itoa(config.MinorAbilityTier) + " to use this skill.")
		return
	}

	ready, msg := s.actor.TimerReady("track")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	s.msg.Actor.SendInfo("You begin searching for traces of recent activity")
	if s.actor.Flags["hidden"] != true {
		s.msg.Observers.SendInfo(s.actor.Name, " begins searching for tracks")
	}

	s.actor.SetTimer("track", config.TrackCooldown)
	knownTracks := 0
	unknownTracks := 0
	if len(s.where.EncounterTable) == 0 {
		s.msg.Actor.SendGood("You are unable to find any tracks")
	} else {
		for k := range s.where.EncounterTable {
			whatMob := objects.Mobs[k]
			curChance := config.TrackChance + (s.actor.Int.Current * config.TrackChancePerPoint) + (config.TrackChancePerLevel * (s.actor.Tier - whatMob.Level))
			if _, ok := objects.Mobs[k]; ok {
				if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {
					knownTracks += 1
					s.msg.Actor.SendGood("You find the tracks of a ", objects.Mobs[k].Name)
				} else {
					unknownTracks += 1
				}

			}
		}
		if knownTracks == 0 && unknownTracks > 0 {
			s.msg.Actor.SendGood("You find some tracks but they are unknown to you")
		} else if unknownTracks > 0 {
			s.msg.Actor.SendGood("You also find some tracks that are unknown to you")
		}
	}
}
