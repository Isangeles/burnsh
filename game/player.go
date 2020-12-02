/*
 * player.go
 *
 * Copyright 2020 Dariusz Sikora <dev@isangeles.pl>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either Version 2 of the License, or
 * (at your option) any later Version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
 * MA 02110-1301, USA.
 *
 *
 */

package game

import (
	"github.com/isangeles/flame/module/character"
	"github.com/isangeles/flame/module/serial"
	"github.com/isangeles/flame/module/useaction"

	"github.com/isangeles/burnsh/log"
)

// Wrapper struct for player character.
type Player struct {
	*character.Character
	game *Game
}

// SetDestPoint sets a specified XY position as current
// as a character destination point.
func (p *Player) SetDestPoint(x, y float64) {
	p.Character.SetDestPoint(x, y)
	if p.game.Server() == nil {
		return
	}
	err := p.game.Server().Move(p.ID(), p.Serial(), x, y)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send move request: %v",
			p.ID(), p.Serial(), err)
	}
}

// AddChatMessage adds new message to player chat log.
func (p *Player) AddChatMessage(message string) {
	p.ChatLog().Add(message)
	if p.game.Server() == nil {
		return
	}
	err := p.game.Server().Chat(p.ID(), p.Serial(), message)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send chat request: %v",
			p.ID(), p.Serial(), err)
	}
}

// Use uses specified usable object.
func (p *Player) Use(ob useaction.Usable) {
	p.Character.Use(ob)
	if p.game.Server() == nil {
		return
	}
	obSerial := ""
	if ob, ok := ob.(serial.Serialer); ok {
		obSerial = ob.Serial()
	}
	err := p.game.Server().Use(p.ID(), p.Serial(), ob.ID(), obSerial)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send use request: %v",
			p.ID(), p.Serial(), err)
	}
}
