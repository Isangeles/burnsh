/*
 * response.go
 *
 * Copyright 2020-2021 Dariusz Sikora <dev@isangeles.pl>
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

	"github.com/isangeles/fire/response"

	"github.com/isangeles/burnsh/data/res"
	"github.com/isangeles/burnsh/log"
)

// handleResponse handles specified response from Fire server.
func (g *Game) handleResponse(resp response.Response) {
	if !resp.Logon && g.onLoginFunc != nil {
		g.onLoginFunc(g)
	}
	g.handleUpdateResponse(resp.Update)
	for _, r := range resp.Character {
		g.handleCharacterResponse(r)
	}
	for _, r := range resp.Error {
		log.Err.Printf("Game server error: %s", r)
	}
}

// handleCharacterResponse handles new characters from server response.
func (g *Game) handleCharacterResponse(resp response.Character) {
	for _, p := range g.Players() {
		if p.ID() == resp.ID && p.Serial() == resp.Serial {
			return
		}
	}
	char := g.Module().Chapter().Character(resp.ID, resp.Serial)
	if char == nil {
		log.Err.Printf("Game server: handle new-char response: unable to find character in module: %s %s",
			resp.ID, resp.Serial)
		return
	}
	player := NewPlayer(char, g)
	g.AddPlayer(player)
}

// handleUpdateRespone handles update response.
func (g *Game) handleUpdateResponse(resp response.Update) {
	flameres.Clear()
	flameres.Add(flameres.ResourcesData{TranslationBases: res.TranslationBases})
	g.Module().Apply(resp.Module)
}
