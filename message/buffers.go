// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package message

import (
	"github.com/ArcCS/Nevermore/text"
)

// buffers are a collection of Buffer indexed by location.
type buffers map[int]*Buffer

// Send calls Buffer.Send for each Buffer in the receiver buffers.
//
// See also Buffer.Send for more details.
func (b buffers) Send(s ...string) {
	for _, b := range b {
		b.Send(s...)
	}
}

// SendGood is convenient for sending a message to all Buffer in buffers using
// text.Good for the color.
func (b buffers) SendGood(s ...string) {
	for _, b := range b {
		b.sendColor(text.Good, s...)
	}
}

// SendBad is convenient for sending a message to all Buffer in buffers using
// text.Bad for the color.
func (b buffers) SendBad(s ...string) {
	for _, b := range b {
		b.sendColor(text.Bad, s...)
	}
}

// SendInfo is convenient for sending a message to all Buffer in buffers using
// text.Info for the color.
func (b buffers) SendInfo(s ...string) {
	for _, b := range b {
		b.sendColor(text.Info, s...)
	}
}

// Append calls Buffer.Append for each Buffer in the receiver buffers.
//
// See also Buffer.Append for more details.
func (b buffers) Append(s ...string) {
	for _, b := range b {
		b.Append(s...)
	}
}

// Silent calls Buffer.Silent with the passed new flag for each Buffer in the
// receiver buffers. Silent returns two sets of buffers, one for all buffers
// that were true and one for all buffers that were false. The previous silent
// state of buffers can be restored by calling Silent with true or false on the
// returned buffers. For example:
//
//	t,f := s.msg.Observers.Silent(true)
//	:
//	: // do something
//	:
//	t.Silent(true)
//	f.silent(false)
//
// See also Buffer.Silent for more details.

func (b buffers) Silent(new bool) (t buffers, f buffers) {
	t = make(map[int]*Buffer)
	f = make(map[int]*Buffer)
	for where, b := range b {
		if old := b.Silent(new); old {
			t[where] = b
		} else {
			f[where] = b
		}
	}
	return
}

// Len returns the number of messages for each Buffer in buffers as a
// [has.Inventory]int map.
func (b buffers) Len() (l map[int]int) {
	l = make(map[int]int)
	for where, b := range b {
		l[where] = b.count
	}
	return
}

// Filter takes a list of Room Id's and filter them out
func (b buffers) Filter(limit ...int) (filtered buffers) {
	filtered = make(map[int]*Buffer)
	for _, l := range limit {
		if _, ok := b[l]; ok {
			filtered[l] = b[l]
		}
	}
	return
}
