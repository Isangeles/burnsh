/*
 * creafting.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/isangeles/flame/module/character"
	"github.com/isangeles/flame/module/craft"
	"github.com/isangeles/flame/module/effect"
)

// craftingDialog starts CLI dialog for
// active PC crafting.
func craftingDialog() error {
	if activeGame == nil {
		msg := lang.Text("no_game_err")
		return fmt.Errorf(msg)
	}
	if activePC == nil {
		msg := lang.Text("no_pc_err")
		return fmt.Errorf(msg)
	}
	for {
		// Select recipe.
		recipe, err := recipeDialog(activePC)
		if err != nil {
			fmt.Printf("%v\n", err)
			break
		}
		// Print recipe details.
		fmt.Printf("%s:\t%s\n", lang.Text("crafting_recipe"),
			recipe.ID())
		fmt.Printf("%s:\t%s\n", lang.Text("crafting_category"),
			recipe.Category())
		fmt.Printf("%s:\t%s\n", lang.Text("crafting_reqs"),
			reqsInfo(recipe.UseAction().Requirements()...))
		fmt.Printf("%s:\n", lang.Text("crafting_result"))
		for _, m := range recipe.UseAction().UserMods() {
			m, ok := m.(*effect.AddItemMod)
			if ok {
				fmt.Printf("\t%s\tx%d\n", m.ItemID(), m.Amount())
			}
		}
		// Recipe options.
		ans := 0
		for ans == 0 {
			// Print options.
			fmt.Printf("[%d]%s\n", 1, lang.Text("crafting_make"))
			fmt.Printf("[%d]%s\n", 2, lang.Text("dialog_back"))
			fmt.Printf("[%d]%s\n", 3, lang.Text("dialog_cancel"))
			// Scan input.
			scan := bufio.NewScanner(os.Stdin)
			scan.Scan()
			input := scan.Text()
			n, err := strconv.Atoi(input)
			if err != nil {
				msg := lang.Text("nan_err")
				fmt.Printf("%s:%v\n", msg, err)
				continue
			}
			if n < 1 || n > 3 {
				msg := lang.Text("invalid_input_err")
				fmt.Printf("%s:%s\n", msg, input)
				continue
			}
			ans = n
		}
		if ans == 2 {
			continue
		}
		if ans == 3 {
			break
		}
		if ans == 1 {
			activePC.Use(recipe)
			break
		}
	}
	return nil
}

// recipeDialog starts recipe dialog for specified
// game character.
func recipeDialog(c *character.Character) (*craft.Recipe, error) {
	recipes := c.Crafting().Recipes()
	if len(recipes) < 1 {
		msg := lang.Text("crafting_no_recipes_err")
		return nil, fmt.Errorf(msg)
	}
	var recipe *craft.Recipe
	for recipe == nil {
		// List recipes.
		fmt.Printf("%s:\n", lang.Text("crafting_recipes"))
		for i, r := range recipes {
			fmt.Printf("[%d]%s\t%s\n", i, r.ID(), r.Category())
		}
		// Select ID.
		fmt.Printf("%s:", lang.Text("crafting_select_recipe"))
		// Scan input.
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		input := scan.Text()
		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s:%v\n", lang.Text("nan_err"), input)
			continue
		}
		if id < 0 || id > len(recipes)-1 {
			fmt.Printf("%s:%s\n", lang.Text("invalid_input_err"), input)
			continue
		}
		recipe = recipes[id]
	}
	return recipe, nil
}
