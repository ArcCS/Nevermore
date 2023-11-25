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
)

// game embeds a frontend instance adding fields and methods specific to
// communicating with the game.
type game struct {
	*frontend
}

func StartGame(f *frontend, charName string) (g *game) {
	accounts.inuse[f.account] = struct{}{}
	g = &game{frontend: f}
	g.character, _ = objects.LoadCharacter(charName, f.output, g.Disconnect)
	g.gameInit()
	return
}

func FirstTimeStartGame(f *frontend, charName string) (g *game) {
	g = &game{frontend: f}
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

// gameInit is used to place the player into the game world.
func (g *game) gameInit() {

	message.ReleaseBuffer(g.buf)
	g.buf = nil

	if _, ok := objects.Rooms[g.character.ParentId]; !ok {
		g.character.ParentId = config.StartingRoom
	}
	if g.character.Class == 100 {
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
	// Initialize this characters ticker
	g.nextFunc = g.gameProcess
}

func ResumeGame(f *frontend, charRef *objects.Character) (g *game) {
	g = &game{frontend: f}
	g.character = charRef
	g.gameResumeInit()
	return
}

// gameInit is used to place the player into the game world.
func (g *game) gameResumeInit() {
	message.ReleaseBuffer(g.buf)
	g.buf = nil
	g.character.Writer = g.output
	g.character.Unloader = g.CharUnloader
	g.character.Disconnect = g.Disconnect
	cmd.Script(g.character, "LOOK")
	// Initialize this characters ticker
	g.nextFunc = g.gameProcess
}

// gameProcess hands input to the game backend for processing while the player
// is in the game. When the player quits the game the frontend.buf buffer is
// restored - see gameInit.
func (g *game) gameProcess() {
	c := cmd.Parse(g.character, string(g.input))
	if c == "QUIT" {
		g.CharUnloader()
	}
}

func (g *game) CharUnloader() {
	g.AccountCleanup()
	g.character.Unload()
	g.character = nil
	g.buf = message.AcquireBuffer()
	g.buf.OmitLF(true)
	NewStart(g.frontend)
}
