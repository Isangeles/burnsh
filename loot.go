/*
 * loot.go
 *
 * Copyright 2019-2023 Dariusz Sikora <ds@isangeles.dev>
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
	"fmt"

	"github.com/isangeles/flame/item"
	"github.com/isangeles/flame/objects"
)

// lootDialog start CLI dialog current
// PC target loot.
func lootDialog() error {
	if activeGame == nil {
		return fmt.Errorf("no game started")
	}
	if activeGame.ActivePlayer() == nil {
		return fmt.Errorf("no active player")
	}
	tar := activeGame.ActivePlayer().Targets()[0]
	if tar == nil {
		return fmt.Errorf("no target")
	}
	if tar, ok := tar.(objects.Killable); ok && tar.Live() {
		return fmt.Errorf("tar not lootable")
	}
	con, ok := tar.(item.Container)
	if !ok {
		return fmt.Errorf("target have no inventory")
	}
	err := activeGame.TransferItems(activeGame.ActivePlayer(), con, lootItems(con.Inventory().Items())...)
	if err != nil {
		return fmt.Errorf("unable to transfer items: %v", err)
	}
	return nil
}

// lootItems returns all lootable items from specified
// inventory items list.
func lootItems(invItems []*item.InventoryItem) (items []item.Item) {
	for _, it := range invItems {
		if it.Loot {
			items = append(items, it.Item)
		}
	}
	return
}
