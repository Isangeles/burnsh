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

package main

import (
	"github.com/isangeles/flame"
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/serial"

	"github.com/isangeles/burn"

	"github.com/isangeles/fire/response"

	"github.com/isangeles/burnsh/data/res"
	"github.com/isangeles/burnsh/game"
	"github.com/isangeles/burnsh/log"
)

// handleResponse handles response from the Fire server.
func handleResponse(resp response.Response) {
	if !resp.Logon {
		if len(resp.Load.Save) > 0 {
			handleLoadResponse(resp.Load)
		}
		handleUpdateResponse(resp.Update)
	}
	for _, r := range resp.Error {
		log.Err.Printf("Server error response: %s", r)
	}
}

// handleUpdateResponse handles update response from the server.
func handleUpdateResponse(resp response.Update) {
	flameres.Clear()
	flameres.Add(flameres.ResourcesData{TranslationBases: res.TranslationBases})
	if mod == nil {
		serial.Reset()
		mod = flame.NewModule(resp.Module)
		return
	}
	mod.Apply(resp.Module)
}

// handleLoadResponse handles load response.
func handleLoadResponse(resp response.Load) {
	serial.Reset()
	flameres.Clear()
	mod = flame.NewModule(resp.Module)
	burn.Module = mod
	activeGame = game.New(mod)
	activeGame.SetServer(server)
}
