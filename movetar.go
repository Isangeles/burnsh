/*
 * movetar.go
 *
 * Copyright 2023 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/isangeles/flame/data/res/lang"
)

// moveTarDialog starts CLI dialog for moving active player to
// the positionf of the current target.
func moveTarDialog() error {
	if activeGame == nil {
		return fmt.Errorf("%s\n", lang.Text("no_game_err"))
	}
	if activeGame.ActivePlayer() == nil {
		return fmt.Errorf("%s\n", lang.Text("no_pc_err"))
	}
	tar := activeGame.ActivePlayer().Targets()[0]
	if tar == nil {
		return fmt.Errorf("%s\n", lang.Text("no_tar_err"))
	}
	tarX, tarY := tar.Position()
	activeGame.ActivePlayer().SetDestPoint(tarX, tarY)
	info := fmt.Sprintf("%s: %fx%f", lang.Text("movetar_info"), tarX, tarY)
	fmt.Printf("%s\n", info)
	return nil
}
