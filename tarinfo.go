/*
 * tarinfo.go
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
	"fmt"

	"github.com/isangeles/flame/data/res/lang"
)

// Interface for objects with info for
// targetinfo command.
type InfoTarget interface {
	ID() string
	Health() int
	Mana() int
}

// targetInfoDialog starts CLI dialog that prints informations
// about current PC target.
func targetInfoDialog() error {
	if activeGame == nil {
		return fmt.Errorf("%s\n", lang.Text("no_game_err"))
	}
	if activeGame.ActivePlayer() == nil {
		return fmt.Errorf("%s\n", lang.Text("no_pc_err"))
	}
	pcTar := activeGame.ActivePlayer().Targets()[0]
	if pcTar == nil {
		return fmt.Errorf("%s\n", lang.Text("no_tar_err"))
	}
	tar, ok := pcTar.(InfoTarget)
	if !ok {
		return fmt.Errorf("%s\n", lang.Text("invalid_tar"))
	}
	// Name.
	info := fmt.Sprintf("%s:%s", lang.Text("ob_name"),
		lang.Text(tar.ID()))
	// Health.
	info += fmt.Sprintf("\n%s:%d", lang.Text("ob_health"),
		tar.Health())
	// Mana.
	info += fmt.Sprintf("\n%s:%d", lang.Text("ob_mana"),
		tar.Mana())
	// Print.
	fmt.Printf("%s\n", info)
	return nil
}
