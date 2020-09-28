/*
 * game.go
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

// Package with interface wrapper for game.
package game

import (
	"fmt"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/module/character"
)

// Struct for game wrapper.
type Game struct {
	*flame.Game
	server       *Server
	players      []*Player
	activePlayer *Player
	onLoginFunc  func(g *Game)
}

// New creates new game wrapper for specified game.
func New(game *flame.Game, server *Server) *Game {
	g := Game{Game: game}
	g.server = server
	if g.server != nil {
		g.server.SetOnResponseFunc(g.handleResponse)
	}
	return &g
}

// Players returns player characters.
func (g *Game) Players() []*Player {
	return g.players
}

// AddPlayer adds new player character.
func (g *Game) AddPlayer(char *character.Character) error {
	if g.server != nil {
		err := g.server.NewCharacter(char.Data())
		if err != nil {
			return fmt.Errorf("Unable to send new character request: %v",
				err)
		}
		return nil
	}
	player, err := g.newPlayer(char)
	if err != nil {
		return fmt.Errorf("Unable to create player: %v", err)
	}
	g.players = append(g.players, player)
	g.SetActivePlayer(player)
	return nil
}

// ActivePlayer returns active player.
func (g *Game) ActivePlayer() *Player {
	return g.activePlayer
}

// SetActivePlayer sets specified player as active player.
func (g *Game) SetActivePlayer(player *Player) {
	g.activePlayer = player
}

// Server retruns game server.
func (g *Game) Server() *Server {
	return g.server
}

// SetOnLoginFunc sets function triggered on login.
func (g *Game) SetOnLoginFunc(f func(g *Game)) {
	g.onLoginFunc = f
}

// newPlayer creates new character avatar for the player from specified data and
// places this character in the start area of game module.
func (g *Game) newPlayer(char *character.Character) (*Player, error) {
	// Set start position.
	char.SetPosition(g.Module().Chapter().Conf().StartPosX,
		g.Module().Chapter().Conf().StartPosY)
	// Set start area.
	startArea := g.Module().Chapter().Area(g.Module().Chapter().Conf().StartArea)
	if startArea == nil {
		return nil, fmt.Errorf("game: start area not found: %s",
			g.Module().Chapter().Conf().StartArea)
	}
	startArea.AddCharacter(char)
	player := Player{char, g}
	return &player, nil
}