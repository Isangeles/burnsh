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
