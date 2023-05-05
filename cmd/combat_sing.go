package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"math"
	"strconv"
	"strings"
)

func init() {
	addHandler(sing{},
		"Usage:  sing song_name # \n\n Use your instrument to sing a song!  \n\n Use 'sing stop' to conclude your performance.",
		permissions.Bard,
		"sing")
}

type sing cmd

func (sing) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("What do you want to sing?")
		return
	}

	// Stop the song
	if s.words[0] == "STOP" {
		if s.actor.CheckFlag("singing") {
			s.actor.SongTickerUnload <- true
			return
		}
		s.msg.Actor.SendBad("You aren't singing!")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
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

	song := strings.ToLower(s.words[0])

	// Check if the song exists
	songInstance, ok := objects.Songs[song]
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

	// Calculate duration of the song, as well as the tick rate
	duration := 45 + (s.actor.GetStat("con") * config.DurationPerCon)

	tickRate := 8 - int(math.Floor(float64(s.actor.Tier/5)))

	s.actor.ApplyEffect("sing", strconv.Itoa(duration), 0, 0,
		func(triggers int) {
			s.actor.SingSong(song, tickRate)
		},
		func() {
			s.actor.SongTickerUnload <- true
		})

	s.msg.Actor.SendGood("You begin singing " + song + "!")
	s.msg.Observers.SendInfo(s.actor.Name + " begins a performance, singing: " + songInstance["verse"] + "!")
	s.actor.SetTimer("combat_sing", 60*2)
	s.ok = true
	return
}
