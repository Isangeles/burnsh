/*
 * loot.go
 *
 * Copyright 2019-2025 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/item"
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
	if len(activeGame.ActivePlayer().Targets()) < 1 {
		return fmt.Errorf("no target")
	}
	tar := activeGame.ActivePlayer().Targets()[0]
	ob, ok := tar.(*character.Character)
	if ok && ob.Live() && !ob.OpenLoot() {
		return fmt.Errorf("target is not lootable")
	}
	err := activeGame.TransferItems(activeGame.ActivePlayer(), ob, lootItems(ob.Inventory().Items())...)
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
