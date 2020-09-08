/*
 * target.go
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
	"github.com/isangeles/flame/module/effect"
)

// targetDialog starts target CLI dialog for
// active player.
func targetDialog() error {
	if activeGame == nil {
		return fmt.Errorf("%s\n", lang.Text("no_game_err"))
	}
	mod := activeGame.Module()
	area := mod.Chapter().CharacterArea(activePC)
	scan := bufio.NewScanner(os.Stdin)
	var tar effect.Target
	for tar == nil {
		fmt.Printf("%s:\n", lang.Text("target_near_targets"))
		targets := area.NearTargets(activePC, activePC.SightRange())
		if len(targets) < 1 {
			return nil
		}
		for i, t := range targets {
			fmt.Printf("[%d]%s\n", i, lang.Text(t.ID()))
		}
		fmt.Printf("%s:", lang.Text("target_select_target"))
		scan.Scan()
		input := scan.Text()
		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s:%s\n", lang.Text("cli_nan_error"),
				input)
			continue
		}
		if id < 0 || id > len(targets)-1 {
			fmt.Printf("%s\n", lang.Text("invalid_input_err"))
			continue
		}
		tar = targets[id]
	}
	activePC.SetTarget(tar)
	return nil
}
