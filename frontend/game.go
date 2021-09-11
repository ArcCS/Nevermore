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
	"strings"
)

// game embeds a frontend instance adding fields and methods specific to
// communicating with the game.
type game struct {
	*frontend
}

// NewGame returns a game with the specified frontend embedded. The returned
// game can be used for processing communication to the actual game.
func StartGame(f *frontend, charName string) (g *game) {
	g = &game{frontend: f}
	g.character, _ = objects.LoadCharacter(charName, f.output)
	g.gameInit()
	return
}

// NewGame returns a game with the specified frontend embedded. The returned
// game can be used for processing communication to the actual game.
func FirstTimeStartGame(f *frontend, charName string) (g *game) {
	g = &game{frontend: f}
	g.character, _ = objects.LoadCharacter(charName, f.output)
	for _, item_id := range config.StartingGear[g.character.Class] {
		newItem := objects.Item{}
		copier.Copy(&newItem, objects.Items[item_id])
		g.character.Inventory.Add(&newItem)
	}
	g.gameInit()
	return
}

// gameInit is used to place the player into the game world. As the game
// backend has it's own output handling we remove the frontend.buf buffer to
// prevent duplicate output. The buffer is restored by gameProcess when the
// player quits the game world.
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
	objects.Rooms[g.character.ParentId].Chars.Lock()
	objects.Rooms[g.character.ParentId].Chars.Add(g.character)
	objects.ActiveCharacters.Add(g.character, g.remoteAddr)
	objects.Rooms[g.character.ParentId].Chars.Unlock()

	cmd.Script(g.character, "$POOF")
	// Initialize this characters ticker
	g.nextFunc = g.gameProcess
}

// gameProcess hands input to the game backend for processing while the player
// is in the game. When the player quits the game the frontend.buf buffer is
// restored - see gameInit.
func (g *game) gameProcess() {
	c := cmd.Parse(g.character, string(g.input))
	if c == "QUIT" || c == strings.ToUpper("DELETE"+g.GetCharacter().Name){
		g.character.Unload()
		g.character = nil
		g.buf = message.AcquireBuffer()
		g.buf.OmitLF(true)
		NewStart(g.frontend)
	}
}
