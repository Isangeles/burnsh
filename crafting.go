/*
 * creafting.go
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

	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/craft"
	"github.com/isangeles/flame/core/module/object/character"
	flameconf "github.com/isangeles/flame/config"
)

// craftingDialog starts CLI dialog for
// active PC crafting.
func craftingDialog() error {
  langPath := flameconf.LangPath()
  if game == nil {
		msg := lang.TextDir(langPath, "no_game_err")
		return fmt.Errorf(msg)
  }
  if activePC == nil {
    msg := lang.TextDir(langPath, "no_pc_err")
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
    fmt.Printf("%s:\t%s\n", lang.TextDir(langPath, "crafting_recipe"),
      recipe.ID())
    fmt.Printf("%s:\t%s\n", lang.TextDir(langPath, "crafting_category"),
      recipe.CategoryID())
    // TODO: print result & amount.
    // Recipe options.
    ans := 0
    for ans == 0 {
      // Print options.
      fmt.Printf("[%d]%s\n", 1, lang.TextDir(langPath, "crafting_make"))
      fmt.Printf("[%d]%s\n", 2, lang.TextDir(langPath, "dialog_back"))
      fmt.Printf("[%d]%s\n", 3, lang.TextDir(langPath, "dialog_cancel"))
      // Scan input.
      scan := bufio.NewScanner(os.Stdin)
      scan.Scan()
      input := scan.Text()
      n, err := strconv.Atoi(input)
      if err != nil {
        msg := lang.TextDir(langPath, "nan_err")
        fmt.Printf("%s:%v\n", msg, err)
        continue
      }
      if n < 1 || n > 3 {
        msg := lang.TextDir(langPath, "invalid_input_err")
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
      err := pcMakeRecipe(recipe)
      if err != nil {
        fmt.Printf("%v\n", err)
        continue
      }
      break
    }
  }
  return nil
}

// recipeDialog starts recipe dialog for specified
// game character.
func recipeDialog(c *character.Character) (*craft.Recipe, error) {
  langPath := flameconf.LangPath()
  recipes := c.Recipes()
  if len(recipes) < 1 {
    msg := lang.TextDir(langPath, "crafting_no_recipes_err")
    return nil, fmt.Errorf(msg)
  }
  var recipe *craft.Recipe
  for recipe == nil {
    // List recipes.
    fmt.Printf("%s:\n", lang.TextDir(langPath, "crafting_recipes"))
    for i, r := range recipes {
      fmt.Printf("[%d]%s\t%s\n", i, r.ID(), r.CategoryID())
    }
    // Select ID.
    fmt.Printf("%s:", lang.TextDir(langPath, "crafting_select_recipe"))
    // Scan input.
    scan := bufio.NewScanner(os.Stdin)
    scan.Scan()
    input := scan.Text()
    id, err := strconv.Atoi(input)
    if err != nil {
      fmt.Printf("%s:%v\n", lang.TextDir(langPath, "nan_err"), input)
      continue
    }
  	if id < 0 || id > len(recipes)-1 {
  		fmt.Printf("%s:%s\n", lang.TextDir(langPath, "invalid_input_err"), input)
  		continue
  	}
    recipe = recipes[id]
  }
  return recipe, nil
}

// pcMakeRecipe checks requirements and makes
// items from specified recipe for active PC.
func pcMakeRecipe(r *craft.Recipe) error {
  langPath := flameconf.LangPath()
  if activePC == nil {
    msg := lang.TextDir(langPath, "no_pc_err")
    return fmt.Errorf(msg)
  }
  if !activePC.MeetReqs(r.Reqs()...) {
    msg := lang.TextDir(langPath, "reqs_not_meet")
    return fmt.Errorf(msg)
  }
  for _, i := range r.Make() {
    activePC.Inventory().AddItem(i)
  }
  return nil
}
