package cmd

import (
	"strings"
)

// Syntax: SAY <message> | " <message>
func init() {
	addHandler(say{}, "SAY", "\"")
	addHelp("Usage:  say [character] \n \n Say something out loud!", 0, "say")
}

type say cmd

func (say) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	/*// TODO: Finish Get all location inventories within 1 move of current location
	for _, loc := range s.where.Exits {
		room := rooms.Rooms[loc.ToId]
		s.AddCharLock(int(room.RoomId))
	}
	*/
	who := s.actor.Name

	if s.actor.Flags["invisible"] {
		who = "Someone"
	}

	msg := strings.Join(s.input, " ")

	if msg[len(msg)-1:] == "?" {
		s.msg.Actor.SendGood("You ask: ", msg)
		s.msg.Observer.SendInfo(who, " asks: ", msg)
	}else if msg[len(msg)-1:] == "!" {
		s.msg.Actor.SendGood("You exclaim: ", msg)
		s.msg.Observer.SendInfo(who, " exclaims: ", msg)
	}else {
		s.msg.Actor.SendGood("You say: ", msg)
		s.msg.Observer.SendInfo(who, " says: ", msg)
	}

	// We need to calculate nearby locations in order to do this.
	// Notify observers in near by locations
	//s.msg.Observers.SendInfo("You hear talking nearby.")

	s.ok = true
	return
}
