/*
 * loadgame.go
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

	"github.com/isangeles/flame"
	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/core"
	"github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/res/lang"
)

// loadGameDialog starts CLI dialog for loading
// saved game.
func loadGameDialog() (*core.Game, error) {
	if flame.Mod() == nil {
		return nil, fmt.Errorf("no module loaded")
	}
	savePattern := fmt.Sprintf(".*%s", data.SavegameFileExt)
	saves, err := data.DirFilesNames(flameconf.ModuleSavegamesPath(), savePattern)
	if err != nil {
		return nil, fmt.Errorf("fail to retrieve save files: %v")
	}
	savename := ""
	scan := bufio.NewScanner(os.Stdin)
	for accept := false; !accept; {
		fmt.Printf("%s:\n", lang.Text("loadgame_saves"))
		for i, s := range saves {
			fmt.Printf("[%d]%v\n", i, s)
		}
		fmt.Printf("%s:", lang.Text("loadgame_select_save"))
		for scan.Scan() {
			input := scan.Text()
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s:%s\n", lang.Text("nan_err"), input)
			}
			if id >= 0 && id < len(saves) {
				savename = saves[id]
				break
			}
		}
		accept = true
	}
	game, err := data.ImportGame(flame.Mod(), flameconf.ModuleSavegamesPath(), savename)
	if err != nil {
		return nil, fmt.Errorf("fail to load saved game: %v", err)
	}
	flame.SetGame(game)
	return game, nil
}
