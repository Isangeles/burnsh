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
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flameconf "github.com/isangeles/flame/config"
	flamedata "github.com/isangeles/flame/data"
	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/burn"

	"github.com/isangeles/burnsh/log"
)

// loadGameDialog starts CLI dialog for loading
// saved game.
func loadGameDialog() error {
	if mod == nil {
		return fmt.Errorf("no module loaded")
	}
	savePattern := fmt.Sprintf(".*%s", flamedata.SavegameFileExt)
	saves, err := flamedata.DirFilesNames(flameconf.ModuleSavegamesPath(), savePattern)
	if err != nil {
		return fmt.Errorf("unable to retrieve save files: %v")
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
	// Game.
	savepath := filepath.Join(flameconf.ModuleSavegamesPath(), savename)
	g, err := flamedata.ImportGame(mod, savepath)
	if err != nil {
		return fmt.Errorf("unable to load saved game: %v", err)
	}
	game = g
	burn.Game = g
	// CLI.
	savename = strings.TrimSuffix(savename, flamedata.SavegameFileExt)
	cliSavePath := fmt.Sprintf("%s/%s%s", flameconf.ModuleSavegamesPath(), savename, CLISaveExt)
	cliSave, err := loadCLI(cliSavePath)
	if err != nil {
		return fmt.Errorf("unable to load CLI state: %v", err)
	}
	for _, pcSave := range cliSave.Players {
		c := game.Module().Chapter().Character(pcSave.ID, pcSave.Serial)
		if c == nil {
			log.Err.Printf("load game: unable to find pc: %s%s", pcSave.ID, pcSave.Serial)
		}
		players = append(players, c)
	}
	activePC = players[0]
	return nil
}

// loadCLI loads CLI save file from specified path.
func loadCLI(path string) (*CLISave, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open save file: %v", err)
	}
	data, _ := ioutil.ReadAll(file)
	cliSave := new(CLISave)
	err = xml.Unmarshal(data, cliSave)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal xml data: %v", err)
	}
	return cliSave, nil
}
