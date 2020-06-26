// Copyright 2015 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import (
	"github.com/ArcCS/Nevermore/message"
	"github.com/ArcCS/Nevermore/objects"
	"io"
	"strings"
)

func init() {
	//event.Script = Script
}

// state contains the current parsing state for commands. The state fields may
// be modified directly except for locks. The AddLocks method should be used to
// add locks, CanLock can be called to see if a lock has already been added.
//
// NOTE: the where field is only set when the state is creafdispatchted. If the actor
// moves to another location the where field should be updated as well. See the
// move command for such an example.
//
type state struct {
	actor       *objects.Character // The Thing executing the command
	where       *objects.Room 		// Where the character currently is
	participant *objects.Character  // The other Character participating in the command
	input       []string            // The original input of the actor minus cmd
	cmd         string              // The current command being processed
	words       []string            // Input as uppercased words, less stopwords
	ok          bool                // Flag to indicate if command was successful
	scripting   bool                // Is state in scripting mode?

	// DO NOT MANIPULATE LOCKS DIRECTLY - use AddLock and see it's comments
	cLocks []int
	mLocks []int
	iLocks []int

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

// Script processes the input string the same as Parse. However Script runs
// with the state.scripting flag set to true, permitting scripting specific
// commands to be executed.
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

	s.tokenizeInput(input)
	//log.Println("Received command ", input)
	s.where = objects.Rooms[o.ParentId]
	s.AddAllLocks(s.where.RoomId)

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
//
func (s *state) tokenizeInput(input string) {
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

	s.msg.Allocate(s.where.RoomId, s.cLocks)
	l := s.TotalLocks()

	dispatchHandler(s)

	if l-s.TotalLocks() == 0 {
		inSync = true
		s.messenger()
	}
	return
}

func (s *state) script(actor, participant, observers bool, inputs ...string) {
	//log.Println("Scripted input..." + strings.Join(inputs, " ") )
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
		players := []io.Writer{}
		for _, c := range objects.Rooms[where].Chars.Contents {
			if c != s.actor && c != s.participant {
				players = append(players, c)
			}
		}
		buffer.Deliver(players...)
	}

	s.msg.Deallocate()
}

func (s *state) AddAllLocks(r int){
	s.AddMobLock(r)
	s.AddItemLock(r)
	s.AddCharLock(r)
}

func (s *state) TotalLocks() int {
	return len(s.cLocks) + len(s.iLocks) + len(s.mLocks)
}

func (s *state) LockAll(){
	for _, l := range s.cLocks {
		objects.Rooms[l].Chars.Lock()
	}
	for _, l := range s.mLocks {
		objects.Rooms[l].Mobs.Lock()
	}
	for _, l := range s.iLocks {
		objects.Rooms[l].Items.Lock()
	}
}

func (s *state) UnlockAll(){
	for _, l := range s.cLocks {
		objects.Rooms[l].Chars.Unlock()
	}
	for _, l := range s.mLocks {
		objects.Rooms[l].Mobs.Unlock()
	}
	for _, l := range s.iLocks {
		objects.Rooms[l].Items.Unlock()
	}
}

func (s *state) AddItemLock(i int) {

	if i == 0{
		return
	}

	s.iLocks = append(s.iLocks, i)
	l := len(s.iLocks)

	if l == 1 {
		return
	}

	for x := 0; x < l; x++ {
		copy(s.iLocks[x+1:l], s.iLocks[x:l-1])
		s.iLocks[x] = i
		break
	}
	// After adding the lock to the context, lock the item
	objects.Rooms[i].Items.Lock()
}

func (s *state) AddMobLock(i int) {

	if i == 0{
		return
	}

	s.mLocks = append(s.mLocks, i)
	l := len(s.mLocks)

	if l == 1 {
		return
	}

	for x := 0; x < l; x++ {
		copy(s.mLocks[x+1:l], s.mLocks[x:l-1])
		s.mLocks[x] = i
		break
	}
	// After adding the lock to the context, lock the item
	objects.Rooms[i].Mobs.Lock()
}

func (s *state) AddCharLock(i int) {

	if i == 0{
		return
	}

	s.cLocks = append(s.cLocks, i)
	l := len(s.cLocks)

	if l == 1 {
		return
	}

	// After adding the lock to the context, lock the item
	objects.Rooms[i].Chars.Lock()

	for x := 0; x < l; x++ {
		copy(s.cLocks[x+1:l], s.cLocks[x:l-1])
		s.cLocks[x] = i
		break
	}


}