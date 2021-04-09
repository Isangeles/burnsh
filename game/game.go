/*
 * game.go
 *
 * Copyright 2020-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/isangeles/flame/dialog"
	"github.com/isangeles/flame/flag"
	"github.com/isangeles/flame/item"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/ignite/ai"

	"github.com/isangeles/burnsh/log"
)

const (
	aiCharFlag = flag.Flag("igniteNpc")
)

// Struct for game wrapper.
type Game struct {
	*flame.Module
	server       *Server
	players      []*Player
	activePlayer *Player
	localAI      *ai.AI
	onLoginFunc  func(g *Game)
}

// New creates new game wrapper for specified module.
func New(module *flame.Module) *Game {
	g := Game{Module: module}
	g.localAI = ai.New(ai.NewGame(module))
	return &g
}

// Update updates game.
func (g *Game) Update(delta int64) {
	g.Module.Update(delta)
	if g.Server() != nil {
		return
	}
	g.updateAIChars()
	g.localAI.Update(delta)
}

// Players returns player characters.
func (g *Game) Players() []*Player {
	return g.players
}

// AddPlayer adds new player character.
func (g *Game) AddPlayer(player *Player) {
	g.players = append(g.players, player)
	g.SetActivePlayer(player)
}

// ActivePlayer returns active player.
func (g *Game) ActivePlayer() *Player {
	return g.activePlayer
}

// SetActivePlayer sets specified player as active player.
func (g *Game) SetActivePlayer(player *Player) {
	g.activePlayer = player
}

// SetServer sets remote game server.
func (g *Game) SetServer(server *Server) {
	g.server = server
	g.Server().SetOnResponseFunc(g.handleResponse)
	err := g.Server().Update()
	if err != nil {
		log.Err.Printf("Game: unable to send update request to the server: %v",
			err)
	}
}

// Server retruns game server.
func (g *Game) Server() *Server {
	return g.server
}

// SetOnLoginFunc sets function triggered on login.
func (g *Game) SetOnLoginFunc(f func(g *Game)) {
	g.onLoginFunc = f
}

// SpawnPlayer places specified player in the area and on the position specified in
// game module configuration.
func (g *Game) SpawnPlayer(player *Player) error {
	// Set start position.
	player.SetPosition(g.Chapter().Conf().StartPosX, g.Chapter().Conf().StartPosY)
	// Set start area.
	startArea := g.Chapter().Area(g.Chapter().Conf().StartArea)
	if startArea == nil {
		return fmt.Errorf("game: start area not found: %s",
			g.Chapter().Conf().StartArea)
	}
	startArea.AddCharacter(player.Character)
	return nil
}

// TransferItems transfer items between specified objects.
// Items are in the form of a map with IDs as keys and serial values as values.
func (g *Game) TransferItems(from, to item.Container, items ...item.Item) error {
	for _, i := range items {
		if from.Inventory().Item(i.ID(), i.Serial()) == nil {
			return fmt.Errorf("Item not found: %s %s",
				i.ID(), i.Serial())
		}
		from.Inventory().RemoveItem(i)
		err := to.Inventory().AddItem(i)
		if err != nil {
			return fmt.Errorf("Unable to add item inventory: %v",
				err)
		}
	}
	if g.Server() == nil {
		return nil
	}
	transferReq := request.TransferItems{
		ObjectFromID:     from.ID(),
		ObjectFromSerial: from.Serial(),
		ObjectToID:       to.ID(),
		ObjectToSerial:   to.Serial(),
		Items:            make(map[string][]string),
	}
	for _, i := range items {
		transferReq.Items[i.ID()] = append(transferReq.Items[i.ID()], i.Serial())
	}
	req := request.Request{TransferItems: []request.TransferItems{transferReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: transfer items: unable to send transfer items request: %v",
			err)
	}
	return nil
}

// Trade exchanges items between specified containers.
func (g *Game) Trade(seller, buyer item.Container, sellItems, buyItems []item.Item) {
	for _, it := range sellItems {
		buyer.Inventory().RemoveItem(it)
		err := seller.Inventory().AddItem(it)
		if err != nil {
			log.Err.Printf("Game: trade items: unable to add sell item: %s %s: %v",
				it.ID(), it.Serial(), err)
		}
	}
	for _, it := range buyItems {
		seller.Inventory().RemoveItem(it)
		err := buyer.Inventory().AddItem(it)
		if err != nil {
			log.Err.Printf("Game: trade items: unable to add buy item: %s %s: %v",
				it.ID(), it.Serial(), err)
		}
	}
	if g.Server() == nil {
		return
	}
	transferReqSell := request.TransferItems{
		ObjectFromID:     buyer.ID(),
		ObjectFromSerial: buyer.Serial(),
		ObjectToID:       seller.ID(),
		ObjectToSerial:   seller.Serial(),
		Items:            make(map[string][]string),
	}
	for _, i := range sellItems {
		transferReqSell.Items[i.ID()] = append(transferReqSell.Items[i.ID()], i.Serial())
	}
	transferReqBuy := request.TransferItems{
		ObjectFromID:     seller.ID(),
		ObjectFromSerial: seller.Serial(),
		ObjectToID:       buyer.ID(),
		ObjectToSerial:   buyer.Serial(),
		Items:            make(map[string][]string),
	}
	for _, i := range buyItems {
		transferReqBuy.Items[i.ID()] = append(transferReqBuy.Items[i.ID()], i.Serial())
	}
	tradeReq := request.Trade{Sell: transferReqSell, Buy: transferReqBuy}
	req := request.Request{Trade: []request.Trade{tradeReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: trade items: unable to send trade request: %v",
			err)
	}
}

// StartDialog starts dialog with specified object as dialog target.
func (g *Game) StartDialog(dialog *dialog.Dialog, target dialog.Talker) {
	dialog.Restart()
	dialog.SetTarget(target)
	if g.Server() == nil || dialog.Owner() == nil {
		return
	}
	dialogReq := request.Dialog{
		TargetID:     target.ID(),
		TargetSerial: target.Serial(),
		OwnerID:      dialog.Owner().ID(),
		OwnerSerial:  dialog.Owner().Serial(),
		DialogID:     dialog.ID(),
	}
	req := request.Request{Dialog: []request.Dialog{dialogReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: start dialog: unable to send dialog request: %v",
			err)
	}
}

// AnswerDialog answers dialog with specified answer.
func (g *Game) AnswerDialog(dialog *dialog.Dialog, answer *dialog.Answer) {
	dialog.Next(answer)
	if g.Server() == nil || dialog.Owner() == nil || dialog.Target() == nil {
		return
	}
	dialogReq := request.Dialog{
		TargetID:     dialog.Target().ID(),
		TargetSerial: dialog.Target().Serial(),
		OwnerID:      dialog.Owner().ID(),
		OwnerSerial:  dialog.Owner().Serial(),
		DialogID:     dialog.ID(),
	}
	dialogAnswerReq := request.DialogAnswer{
		Dialog:   dialogReq,
		AnswerID: answer.ID(),
	}
	req := request.Request{DialogAnswer: []request.DialogAnswer{dialogAnswerReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: answer dialog: unable to send dialog answer: %v",
			err)
	}
}

// updateAIChars updates list of characters controlled by the AI.
func (g *Game) updateAIChars() {
outer:
	for _, c := range g.Chapter().Characters() {
		for _, aic := range g.localAI.Game().Characters() {
			if aic.ID() == c.ID() && aic.Serial() == c.Serial() {
				continue outer
			}
		}
		if !c.HasFlag(aiCharFlag) {
			continue
		}
		aiChar := ai.NewCharacter(c, g.localAI.Game())
		g.localAI.Game().AddCharacter(aiChar)
	}
}
