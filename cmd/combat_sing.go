package cmd

import (
	"math"
	"strconv"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(sing{},
		"Usage:  sing song_name # \n\n Use your instrument to sing a song!  \n\n"+
			" Use 'sing mix song_name' to weave your performance into another song without stopping.  "+
			"The song keeps the time it has left, and you can only mix once every 5 minutes.\n\n"+
			" Use 'sing stop' to conclude your performance.  Stopping early gives back part of the wait before your next song.",
		permissions.Bard,
		"sing")
}

type sing cmd

func (sing) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to sing?")
		return
	}

	switch s.words[0] {
	case "STOP":
		singStop(s)
		return
	case "MIX":
		singMix(s)
		return
	}

	if !s.requireStamina() {
		return
	}

	if s.actor.CheckFlag("singing") {
		s.msg.Actor.SendBad("You are already singing!")
		return
	}

	singReady, msg := s.actor.TimerReady("combat_sing")
	if !singReady {
		s.msg.Actor.SendBad(msg)
		return
	}

	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	// Check if the song exists, allowing a unique prefix like commands do
	song, songInstance, ok := objects.ResolveSong(s.words[0])
	if !ok {
		s.msg.Actor.SendBad("That song doesn't exist!")
		return
	}

	if !utils.StringIn(song, s.actor.Spells) {
		s.msg.Actor.SendBad("You don't know that song!")
		return
	}

	if s.actor.Equipment.Off != (*objects.Item)(nil) {
		if s.actor.Equipment.Off.ItemType != 16 {
			s.msg.Actor.SendBad("You need to be holding an instrument to sing!")
			return
		}
	} else {
		s.msg.Actor.SendBad("You need to be holding an instrument to sing!")
		return
	}

	duration := config.SongDuration(s.actor.Tier)
	tickRate := config.SongTickSeconds(s.actor.Tier)
	s.actor.SongLength = duration

	s.actor.ApplyEffect("sing", strconv.Itoa(duration), 0, 0,
		func(triggers int) {
			s.actor.SingSong(song, tickRate)
		},
		func() {
			// This runs with the room locked (death, sing stop). The song goroutine
			// may be waiting on that same lock for its tick, so a blocking send here
			// deadlocks both; leave the signal in the buffer for it to pick up.
			select {
			case s.actor.SongTickerUnload <- true:
			default:
			}
			// A tick that slips in before the goroutine reads the signal finds no
			// song and does nothing.
			s.actor.CurrentSong = ""
			// Damage Instrument
			weapMsg := s.actor.Equipment.DamageWeapon("off", 1)
			if weapMsg != "" {
				s.msg.Actor.SendInfo(weapMsg)
			}
		})

	s.msg.Actor.SendGood("You begin singing " + song + "!")
	s.msg.Observers.SendInfo(s.actor.Name + " begins a performance, singing: " + songInstance["verse"] + "!")
	s.actor.SetTimer("combat_sing", config.SongCooldown)
	s.ok = true
	return
}

// singStop ends the song early and hands back part of the song cooldown for
// the share of the song left unsung.
func singStop(s *state) {
	if !s.actor.CheckFlag("singing") {
		s.msg.Actor.SendBad("You aren't singing!")
		return
	}
	if song, ok := s.actor.Effects["sing"]; ok {
		s.actor.RefundTimer("combat_sing", config.SongStopRefundSeconds(song.TimeRemaining(), float64(s.actor.SongLength)))
	}
	s.actor.RemoveEffect("sing")
	if wait := int(math.Ceil(s.actor.TimerRemaining("combat_sing"))); wait > 0 {
		s.msg.Actor.SendInfo("You can sing again in " + strconv.Itoa(wait) + " seconds.")
	}
	s.ok = true
}

// singMix switches a running song to another one without stopping it: the
// song keeps the time it has left and its first pulse lands at once.
func singMix(s *state) {
	if !s.actor.CheckFlag("singing") || s.actor.CurrentSong == "" {
		s.msg.Actor.SendBad("You aren't singing!")
		return
	}
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("What do you want to mix into?")
		return
	}

	if !s.requireStamina() {
		return
	}

	mixReady, msg := s.actor.TimerReady("combat_sing_mix")
	if !mixReady {
		s.msg.Actor.SendBad(msg)
		return
	}

	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	song, songInstance, ok := objects.ResolveSong(s.words[1])
	if !ok {
		s.msg.Actor.SendBad("That song doesn't exist!")
		return
	}

	if !utils.StringIn(song, s.actor.Spells) {
		s.msg.Actor.SendBad("You don't know that song!")
		return
	}

	if song == s.actor.CurrentSong {
		s.msg.Actor.SendBad("You are already singing that!")
		return
	}

	s.actor.MixSong(song)
	s.msg.Actor.SendGood("You weave your performance into " + song + "!")
	s.msg.Observers.SendInfo(s.actor.Name + "'s song shifts into " + s.actor.SongDescriptor() + ": " + songInstance["verse"] + "!")
	s.actor.SetTimer("combat_sing_mix", config.SongMixCooldown)
	s.ok = true
}
