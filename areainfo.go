/*
 * areainfo.go
 *
 * Copyright 2021-2023 Dariusz Sikora <ds@isangeles.dev>
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
	"time"

	"github.com/isangeles/flame/data/res/lang"
)

// areaInfoDialog starts CLI dialog that prints information
// about current area.
func areaInfoDialog() error {
	if activeGame == nil {
		return fmt.Errorf("%s\n", lang.Text("no_game_err"))
	}
	if activeGame.ActivePlayer() == nil {
		return fmt.Errorf("%s\n", lang.Text("no_pc_err"))
	}
	area := activeGame.Chapter().ObjectArea(activeGame.ActivePlayer().Character)
	if area == nil {
		return fmt.Errorf("%s\n", lang.Text("no_pc_area_err"))
	}
	// Name.
	info := fmt.Sprintf("%s: %s", lang.Text("area_name"),
		lang.Text(area.ID()))
	// Time.
	info += fmt.Sprintf("\n%s: %s", lang.Text("area_time"),
		area.Time.Format(time.Kitchen))
	// Weather.
	info += fmt.Sprintf("\n%s: %s", lang.Text("area_weather"),
		lang.Text(string(area.Weather().Conditions)))
	// Print.
	fmt.Printf("%s\n", info)
	return nil
}
