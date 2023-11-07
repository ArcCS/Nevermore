// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package message

// Msg is a collection of buffers for gathering messages to send back as a
// result of processing a command.
type Msg struct {
	Actor              *Buffer
	GM                 *Buffer
	ActorVerbose       *Buffer
	Participant        *Buffer
	ParticipantVerbose *Buffer
	Observer           *Buffer
	ObserverVerbose    *Buffer
	Observers          buffers
	ObserversVerbose   buffers
}

// Allocate sets up the message buffers for the actor, participant and
// observers.
func (m *Msg) Allocate(obsRooms []int) {
	if m.Actor == nil {
		m.Actor = AcquireBuffer()
		m.Actor.omitLF = true
		m.Participant = AcquireBuffer()
		m.Observers = make(map[int]*Buffer)
	}

	for _, l := range obsRooms {
		if _, ok := m.Observers[l]; !ok {
			m.Observers[l] = AcquireBuffer()
		}
	}

}

// Deallocate releases the references to message buffers for the actor,
// participant and observers. Specific deallocation can help with garbage
// collection.
func (m *Msg) Deallocate() {
	ReleaseBuffer(m.Actor)
	m.Actor = nil
	ReleaseBuffer(m.Participant)
	m.Participant = nil
	m.Observer = nil
	for where := range m.Observers {
		ReleaseBuffer(m.Observers[where])
		m.Observers[where] = nil
		delete(m.Observers, where)
	}
}
