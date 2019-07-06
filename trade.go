/*
 * trade.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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
	"strings"
	
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/character"
	"github.com/isangeles/flame/core/module/object/item"
	flameconf "github.com/isangeles/flame/config"
)

// tradeDialog starts CLI dialog for trade with
// current PC target.
func tradeDialog() error {
	langPath := flameconf.LangPath()
	if game == nil {
		msg := lang.TextDir(langPath, "no_game_err")
		return fmt.Errorf(msg)
	}
	if activePC == nil {
		msg := lang.TextDir(langPath, "no_pc_err")
		return fmt.Errorf(msg)
	}
	tar := activePC.Targets()[0]
	if tar == nil {
		msg := lang.TextDir(langPath, "no_tar_err")
		return fmt.Errorf(msg)
	}
	tarChar, ok := tar.(*character.Character)
	if !ok {
		return fmt.Errorf("invalid_target")
	}
	fmt.Printf("%s:\n", lang.TextDir(langPath, "trade_buy_items"))
	buyItems := selectItems(tarChar.Inventory())
	fmt.Printf("%s:\n", lang.TextDir(langPath, "trade_sell_items"))
	sellItems := selectItems(activePC.Inventory())
	// Check trade value.
	buyValue := 0
	for _, it := range buyItems {
		buyValue += it.Value()
	}
	sellValue := 0
	for _, it := range sellItems {
		sellValue += it.Value()
	}
	valueLabel := lang.TextDir(langPath, "trade_item_value")
	fmt.Printf("%s[%s:%d]:\n", lang.TextDir(langPath, "trade_buy_items"), valueLabel, buyValue)
	for _, it := range buyItems {
		fmt.Printf("\t%s\n", it.ID())
	}
	fmt.Printf("%s[%s:%d]:\n", lang.TextDir(langPath, "trade_sell_items"), valueLabel, sellValue)
	for _, it := range sellItems {
		fmt.Printf("\t%s\n", it.ID())
	}
	fmt.Printf("%s[y/N]:", lang.TextDir(langPath, "trade_accept"))
	// Scan input.
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	input := scan.Text()
	if strings.ToLower(input) != "y" {
		return nil
	}
	if sellValue < buyValue {
		fmt.Printf("%s\n", lang.TextDir(langPath, "trade_sell_value_small"))
		return nil
	}
	// Trade items.
	for _, it := range buyItems {
		activePC.Inventory().AddItem(it)
		tarChar.Inventory().RemoveItem(it)
	}
	for _, it := range sellItems {
		tarChar.Inventory().AddItem(it)
		activePC.Inventory().RemoveItem(it)
	}	
	return nil
}

// selectItems starts dialog for selecting
// items from specified inventory.
func selectItems(inv *item.Inventory) (items []item.Item) {
	langPath := flameconf.LangPath()
	selectItems := make(map[string]item.Item)
	if len(inv.Items()) < 1 {
		fmt.Printf("%s\n", lang.TextDir(langPath, "trade_no_items"))
		return
	}
	for {
		invItems := make([]item.Item, 0)
		for _, it := range inv.Items() {
			if selectItems[it.ID()+it.Serial()] != nil {
				continue
			}
			invItems = append(invItems, it)
		}
		// List items to select.
		fmt.Printf("%s:\n", lang.TextDir(langPath, "trade_select_items"))
		valueLabel := lang.TextDir(langPath, "trade_item_value")
		for i, it := range invItems {
			fmt.Printf("\t[%d]%s\t%s:%d\n", i, it.ID(), valueLabel, it.Value()) 
		}
		// Scan input.
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		input := scan.Text()
		if input == "" {
			break
		}
		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s:%v\n", lang.TextDir(langPath, "nan_err"), input)
			continue
		}
		if id < 0 || id > len(invItems)-1 {
			fmt.Printf("%s:%s\n", lang.TextDir(langPath, "invalid_input_err"), input)
			continue
		}
		it := invItems[id]
		selectItems[it.ID()+it.Serial()] = it
	}
	for _, it := range selectItems {
		items = append(items, it)
	}
	return
}
