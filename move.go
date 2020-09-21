/*
 * move.go
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
)

// moveDialog starts dialog for setting a destination point
// for the active player.
func moveDialog() error {
	if activeGame == nil {
		return fmt.Errorf("%s\n", lang.Text("no_game_err"))
	}
	if activeGame.ActivePlayer() == nil {
		return fmt.Errorf("%s\n", lang.Text("no_pc_err"))
	}
	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s:", lang.Text("move_enter_x_position"))
		scan.Scan()
		input := scan.Text()
		x, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Printf("%s: %s\n", lang.Text("cli_nan_error"),
				input)
			continue
		}
		fmt.Printf("%s:", lang.Text("move_enter_y_position"))
		scan.Scan()
		input = scan.Text()
		y, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Printf("%s: %s\n", lang.Text("cli_nan_error"),
				input)
			continue
		}
		activeGame.ActivePlayer().SetDestPoint(x, y)
		break
	}
	return nil
}
