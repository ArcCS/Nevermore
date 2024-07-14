// Copyright 2016 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package frontend

import (
	"github.com/ArcCS/Nevermore/cmd"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/message"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/jinzhu/copier"
	"log"
	"time"
)

// Game embeds a Frontend instance adding fields and methods specific to
// communicating with the Game.
type Game struct {
	*Frontend
}

func StartGame(f *Frontend, charName string) (g *Game) {
	accounts.inuse[f.account] = struct{}{}
	g = &Game{Frontend: f}
	g.character, _ = objects.LoadCharacter(charName, f.output, g.Disconnect)
	g.gameInit()
	return
}

func FirstTimeStartGame(f *Frontend, charName string) (g *Game) {
	g = &Game{Frontend: f}
	g.character, _ = objects.LoadCharacter(charName, f.output, f.Disconnect)
	for _, itemId := range config.StartingGear[g.character.Class] {
		newItem := objects.Item{}
		if err := copier.CopyWithOption(&newItem, objects.Items[itemId], copier.Option{DeepCopy: true}); err != nil {
			log.Println("Error copying character: ", err)
		}
		g.character.Inventory.Add(&newItem)
	}
	g.gameInit()
	return
}

// gameInit is used to place the player into the Game world.
func (g *Game) gameInit() {

	message.ReleaseBuffer(g.buf)
	g.buf = nil

	if _, ok := objects.Rooms[g.character.ParentId]; !ok {
		g.character.ParentId = config.StartingRoom
	}
	if g.character.Class == 100 || g.character.Class == 99 {
		g.character.Permission.ToggleFlag(g.permissions)
	} else {
		g.character.Permission.ToggleFlag(permissions.Anyone)
		g.character.Permission.ToggleFlag(permissions.Player)
		g.character.Permission.ToggleFlag(config.ClassPerms[g.character.Class])
	}
	g.character.Unloader = g.CharUnloader
	g.character.Disconnect = g.Disconnect
	objects.Rooms[g.character.ParentId].LockRoom("GameInit", false)
	objects.Rooms[g.character.ParentId].Chars.Add(g.character)
	objects.ActiveCharacters.Add(g.character, g.remoteAddr)
	objects.Rooms[g.character.ParentId].UnlockRoom("GameInit", false)

	cmd.Script(g.character, "$POOF")
	objects.LastActivity[g.character.Name] = time.Now()
	// Initialize this characters ticker
	g.nextFunc = g.gameProcess
}

func ResumeGame(f *Frontend, charRef *objects.Character) (g *Game) {
	g = &Game{Frontend: f}
	g.character = charRef
	g.gameResumeInit()
	return
}

// gameInit is used to place the player into the Game world.
func (g *Game) gameResumeInit() {
	message.ReleaseBuffer(g.buf)
	g.buf = nil
	g.character.Writer = g.output
	g.character.Unloader = g.CharUnloader
	g.character.Disconnect = g.Disconnect
	cmd.Script(g.character, "LOOK")
	// Initialize this characters ticker
	g.nextFunc = g.gameProcess
}

// gameProcess hands input to the Game backend for processing while the player
// is in the Game. When the player quits the Game the Frontend.buf buffer is
// restored - see gameInit.
func (g *Game) gameProcess() {
	c := cmd.Parse(g.character, string(g.input))
	if c == "QUIT" {
		g.CharUnloader()
	}
}

func (g *Game) CharUnloader() {
	g.AccountCleanup()
	g.character.Unload()
	g.character = nil
	g.buf = message.AcquireBuffer()
	g.buf.OmitLF(true)
	NewStart(g.Frontend)
}
