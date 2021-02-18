/*
 * equip.go
 *
 * Copyright 2021 Dariusz Sikora <dev@isangeles.pl>
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
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module/item"
)

// equipDialog starts CLI dialog for equip action.
func equipDialog() error {
	if activeGame == nil {
		msg := lang.Text("no_game_err")
		return fmt.Errorf(msg)
	}
	if activeGame.ActivePlayer() == nil {
		msg := lang.Text("no_pc_err")
		return fmt.Errorf(msg)
	}
	// List items.
	fmt.Printf("%s:\n", lang.Text("equip_items"))
	items := make([]item.Equiper, 0)
	for _, it := range activeGame.ActivePlayer().Inventory().Items() {
		if eit, ok := it.(item.Equiper); ok {
			items = append(items, eit)
		}
	}
	for i, it := range items {
		if activeGame.ActivePlayer().Equipment().Equiped(it) {
			fmt.Printf("[%d]%s[e]\n", i, lang.Text(it.ID()))
		} else {
			fmt.Printf("[%d]%s\n", i, lang.Text(it.ID()))
		}
	}
	// Select skill.
	scan := bufio.NewScanner(os.Stdin)
	var item item.Equiper
	for item == nil {
		fmt.Printf("%s:", lang.Text("equip_select"))
		scan.Scan()
		input := scan.Text()
		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s:%s\n", lang.Text("nan_err"), input)
			continue
		}
		if id < 0 || id > len(items)-1 {
			fmt.Printf("%s:%s\n", lang.Text("invalid_input_err"), input)
			continue
		}
		item = items[id]
	}
	// Equip/unequip item.
	if activeGame.ActivePlayer().Equipment().Equiped(item) {
		activeGame.ActivePlayer().Equipment().Unequip(item)
	} else {
		err := equip(item)
		if err != nil {
			msg := lang.Text("equip_error")
			return fmt.Errorf("%s: %s", msg, err)
		}
	}
	return nil
}

// equip inserts specified equipable item to all
// compatible slots in active PC equipment.
func equip(it item.Equiper) error {
	if !activeGame.ActivePlayer().MeetReqs(it.EquipReqs()...) {
		return fmt.Errorf(lang.Text("reqs_not_meet"))
	}
	for _, itSlot := range it.Slots() {
		equiped := false
		for _, eqSlot := range activeGame.ActivePlayer().Equipment().Slots() {
			if eqSlot.Item() != nil {
				continue
			}
			if eqSlot.Type() == itSlot {
				eqSlot.SetItem(it)
				equiped = true
				break
			}
		}
		if !equiped {
			activeGame.ActivePlayer().Equipment().Unequip(it)
			return fmt.Errorf(lang.Text("equip_no_free_slot_error"))
		}
	}
	if !activeGame.ActivePlayer().Equipment().Equiped(it) {
		return fmt.Errorf(lang.Text("equip_no_valid_slot_error"))
	}
	return nil
}