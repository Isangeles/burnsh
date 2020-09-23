/*
 * response.go
 *
 * Copyright 2020 Dariusz Sikora <dev@isangeles.pl>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
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
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/module/character"
	"github.com/isangeles/flame/module/serial"

	"github.com/isangeles/fire/response"

	"github.com/isangeles/burnsh/log"
)

// handleResponse handles specified response from Fire server.
func (g *Game) handleResponse(resp response.Response) {
	if !resp.Logon && g.onLoginFunc != nil {
		g.onLoginFunc(g)
	}
	g.handleUpdateResponse(resp.Update)
	g.handleNewCharResponse(resp.NewChar)
	for _, r := range resp.Error {
		log.Err.Printf("Game server error: %s", r)
	}
}

// handleNewCharResponse handles new characters from server response.
func (g *Game) handleNewCharResponse(resp []flameres.CharacterData) {
	for _, cd := range resp {
		serial.Reset()
		char := character.New(cd)
		player := Player{char, g}
		g.players = append(g.players, &player)
		g.SetActivePlayer(&player)
	}
}

// handleUpdateRespone handles update response.
func (g *Game) handleUpdateResponse(resp response.Update) {
	serial.Reset()
	g.Module().Apply(resp.Module)
}
