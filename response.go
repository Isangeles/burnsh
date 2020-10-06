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

package main

import (
	"github.com/isangeles/flame/module"

	"github.com/isangeles/fire/response"

	"github.com/isangeles/burnsh/log"
)

// handleResponse handles response from the Fire server.
func handleResponse(resp response.Response) {
	if !resp.Logon {
		log.Inf.Printf("Logged at: %s", server.Address())
		handleUpdateResponse(resp.Update)
	}
	for _, r := range resp.Error {
		log.Err.Printf("Server error response: %s", r)
	}
}

// handleUpdateResponse handles update response from the server.
func handleUpdateResponse(resp response.Update) {
	if mod == nil {
		mod = module.New()
		mod.Apply(resp.Module)
		return
	}
	mod.Apply(resp.Module)
}
