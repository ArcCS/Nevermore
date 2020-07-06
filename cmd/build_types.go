package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(types{},
           "Usage:  types  \n Print all of the integer to type relationships. ",
           permissions.Builder,
           "types")
}

type types cmd

func (types) process(s *state) {

	for key, value := range config.ItemTypes {
		s.msg.Actor.SendInfo(strconv.Itoa(key) + ": " + value + "\n")
	}

	s.ok = true
	return
}