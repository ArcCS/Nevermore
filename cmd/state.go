// Copyright 2015 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import (
	"github.com/ArcCS/Nevermore/message"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/utils"
	"io"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func init() {
	objects.Script = Script
}

// state contains the current parsing state for commands. The state fields may
// be modified directly except for locks. The AddLocks method should be used to
// add locks, CanLock can be called to see if a lock has already been added.
//
// NOTE: the where field is only set when the state is creafdispatchted. If the actor
// moves to another location the where field should be updated as well. See the
// move command for such an example.
type state struct {
	actor       *objects.Character // The Thing executing the command
	where       *objects.Room      // Where the character currently is
	participant *objects.Character // The other Character participating in the command
	original    string             // The original input string with command removed
	input       []string           // The original input of the actor minus cmd, parsed and cleaned
	cmd         string             // The current command being processed
	words       []string           // Input as uppercased words, less stopwords
	ok          bool               // Flag to indicate if command was successful
	scripting   bool               // Is state in scripting mode?

	// DO NOT MANIPULATE LOCKS DIRECTLY - use AddLock and see it's comments
	rLocks []int

	// msg contains the message buffers for sending data to different recipients
	msg message.Msg
}

// Parse initiates processing of the input string for the specified Thing. The
// input string is expected to be input from a player. The actual command
// processed will be returned. For example GET or DROP.
//
// Parse runs with state.scripting set to false, disallowing scripting specific
// commands from being executed by players directly.
//
// When sync handles a command the command may determine it needs to hold
// additional locks. In this case sync will return false and should be called
// again. This repeats until the list of locks is complete, the command
// processed and sync returns true.
func Parse(o *objects.Character, input string) string {
	s := newState(o, input)
	for !s.sync() {
	}
	return s.cmd
}

// Script processes the input string the same as Parse.
func Script(o *objects.Character, input string) string {
	s := newState(o, input)
	s.scripting = true
	for !s.sync() {
	}
	return s.cmd
}

// newState returns a *state initialised with the passed Thing and input. If
// the passed Thing is locatable the containing Inventory is added to the lock
// list, but the lock is not taken at this point.
func newState(o *objects.Character, input string) *state {

	s := &state{actor: o}

	s.original = input
	s.tokenizeInput(input)
	s.where = objects.Rooms[o.ParentId]
	s.AddLocks(s.where.RoomId)

	return s
}

// tokenizeInput takes the given string and breaks it into uppercased words
// which are stored in the current state. After processing s.cmd will contain
// the leading command, uppercased. s.input will contain the original input
// minus the leading s.cmd. s.words will contain the input, uppercased with
// stopwords and the leading s.cmd removed. For example:
//
//	input = "Say I'm in need of help!"
//	s.cmd = "SAY"
//	s.input = []string{"I'm", "in", "need", "of", "help!"}
//	s.words = []string{"I'M", "NEED", "HELP!"}
func (s *state) tokenizeInput(input string) {
	if len(strings.Fields(s.original)) > 1 {
		s.original = strings.Join(strings.Fields(s.original)[1:], " ")
	}
	quoteReg := regexp.MustCompile("`([^`]*)`")
	for _, match := range quoteReg.FindStringSubmatch(input) {
		input = strings.ReplaceAll(input, match, strings.ReplaceAll(strings.ReplaceAll(match, " ", "%_R%"), "`", ""))
	}
	s.input = strings.Fields(input)
	s.words = make([]string, 0)
	if len(s.input) > 0 {
		if len(s.input) == 0 {
			s.cmd = "Eh?"
			return
		}

		for _, o := range s.input {
			s.words = append(s.words, strings.ToUpper(o))
		}

		s.cmd, s.words = s.words[0], s.words[1:]
		s.input = s.input[1:]
	}
	// Clean up words
	for i := range s.words {
		s.words[i] = strings.ReplaceAll(s.words[i], "%_R%", " ")
	}
	for i := range s.input {
		s.input[i] = strings.ReplaceAll(s.input[i], "%_R%", " ")
	}
}

// sync is called to do the actual locking/unlocking for commands. Having this
// separate from takes advantage of unwinding the locks using defer. This makes
// sync very simple. If the list of locks before and after handling a command
// are the same we are 'in sync' and had all the locks we needed to process the
// command. In this case we return true. If more locks need to be acquired we
// return false and should be called again.
//
// NOTE: There is usually at least one lock, added by newState, which is the
// containing Inventory of the current actor - if it is locatable.
//
// NOTE: At the moment locks are only added - using AddLock. A change in the
// lock list can therefore be detected by simply checking the length of the
// list. If at a later time we need to be able to remove locks as well this
// simple length check will not be sufficient.
func (s *state) sync() (inSync bool) {
	s.LockAll()
	defer s.UnlockAll()

	s.msg.Allocate(s.rLocks)
	l := s.TotalLocks()

	dispatchHandler(s)

	if l-s.TotalLocks() == 0 {
		inSync = true
		s.messenger()
	}
	return
}

func (s *state) script(actor, participant, observers bool, inputs ...string) {
	input := strings.Join(inputs, " ")

	i, w, c, sc := s.input, s.words, s.cmd, s.scripting // Save state

	// Set silent mode on buffers storing old modes
	a := s.msg.Actor.Silent(!actor)
	p := s.msg.Participant.Silent(!participant)
	ot, of := s.msg.Observers.Silent(!observers)

	s.tokenizeInput(strings.TrimSpace(input))
	s.ok = false
	s.scripting = true
	dispatchHandler(s)

	// Restore old silent modes
	s.msg.Actor.Silent(a)
	s.msg.Participant.Silent(p)
	ot.Silent(true)
	of.Silent(false)

	s.input, s.words, s.cmd, s.scripting = i, w, c, sc // Restore state
}

// scriptAll is a helper method that is equivalent to calling script with no
// messages suppressed for the actor, participant or observers.
func (s *state) scriptAll(input ...string) {
	s.script(true, true, true, input...)
}

// scriptAll is a helper method that is equivalent to calling script with all
// messages suppressed.
func (s *state) scriptNone(input ...string) {
	s.script(false, false, false, input...)
}

// scriptAll is a helper method that is equivalent to calling script with
// messages suppressed for any participant or observers. Only the actor will
// receive any messages.
func (s *state) scriptActor(input ...string) {
	s.script(true, false, false, input...)
}

// messenger is used to send buffered messages to the actor, participant and
// observers. The participant may be in another location to the actor - such as
// when throwing something at someone or shooting someone.
//
// For the actor we don't check the buffer length to see if there is anything
// in it to send. We always send to the actor so that we can redisplay the
// prompt even if they just hit enter.
func (s *state) messenger() {

	if s.actor != nil {
		s.msg.Actor.Deliver(s.actor)
	}

	if s.participant != nil && s.msg.Participant.Len() > 0 {
		s.msg.Participant.Deliver(s.participant)
	}

	for where, buffer := range s.msg.Observers {
		if buffer.Len() == 0 {
			continue
		}
		var players []io.Writer
		for _, c := range objects.Rooms[where].Chars.Contents {
			if c != s.actor && c != s.participant {
				players = append(players, c)
			}
		}
		buffer.Deliver(players...)
	}

	s.msg.Deallocate()
}

func (s *state) AddLocks(r int) {
	if !utils.IntIn(r, s.rLocks) {
		if r == 0 {
			return
		}

		s.rLocks = append(s.rLocks, r)
		l := len(s.rLocks)

		if l == 1 {
			return
		}

		for x := 0; x < l; x++ {
			_ = copy(s.rLocks[x+1:l], s.rLocks[x:l-1])
			s.rLocks[x] = r
			break
		}
		objects.Rooms[r].LockRoom("StateHandler("+s.actor.Name+":"+s.cmd+")", false)
	}
}

func (s *state) TotalLocks() int {
	return len(s.rLocks)
}

func (s *state) LockAll() {
	s.AcquireLockPriority()
	for _, l := range s.rLocks {
		objects.Rooms[l].LockRoom("StateHandler("+s.actor.Name+":"+s.cmd+")", false)
	}
}

func (s *state) UnlockAll() {
	s.RemoveLockPriority()
	for _, l := range s.rLocks {
		objects.Rooms[l].UnlockRoom("StateHandler("+s.actor.Name+":"+s.cmd+")", false)
	}
}

func (s *state) AcquireLockPriority() {
	ready := false
	for !ready {
		ready = true
		for _, l := range s.rLocks {
			if objects.Rooms[l].LockPriority == "" {
				objects.Rooms[l].LockPriority = s.actor.Name
			} else if objects.Rooms[l].LockPriority != s.actor.Name {
				ready = false
			}
		}
		if !ready {
			for _, l := range s.rLocks {
				if objects.Rooms[l].LockPriority == s.actor.Name {
					objects.Rooms[l].LockPriority = ""
				}
			}
			r := rand.Int()
			t, _ := time.ParseDuration(string(rune(r)) + "ms")
			time.Sleep(t)
		}
	}
}

func (s *state) RemoveLockPriority() {
	for _, l := range s.rLocks {
		objects.Rooms[l].LockPriority = ""
	}
}
