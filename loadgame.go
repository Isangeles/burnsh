/*
 * loadgame.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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

	"github.com/isangeles/flame"
	flamedata "github.com/isangeles/flame/data"
	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/burnsh/game"
	"github.com/isangeles/burnsh/log"
)

// loadGameDialog starts CLI dialog for loading
// saved game.
func loadGameDialog() error {
	if mod == nil {
		return fmt.Errorf("no module loaded")
	}
	savePattern := fmt.Sprintf(".*%s", SaveExt)
	path := filepath.Join(mod.Conf().Path, ModuleSavesPath)
	saves, err := flamedata.DirFilesNames(path, savePattern)
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
	// Handle game server.
	if server != nil {
		savename = strings.ReplaceAll(savename, SaveExt, "")
		req := request.Request{Load: savename}
		err := server.Send(req)
		if err != nil {
			return fmt.Errorf("unable to send load request: %v", err)
		}
		return nil
	}
	// CLI.
	savename = strings.TrimSuffix(savename, flamedata.ModuleFileExt)
	cliSavePath := filepath.Join(mod.Conf().Path, ModuleSavesPath, savename+SaveExt)
	cliSave, err := loadCLI(cliSavePath)
	if err != nil {
		return fmt.Errorf("unable to load CLI state: %v", err)
	}
	for _, pcSave := range cliSave.Players {
		c := activeGame.Chapter().Character(pcSave.ID, pcSave.Serial)
		if c == nil {
			log.Err.Printf("load game: unable to find pc: %s%s", pcSave.ID, pcSave.Serial)
		}
		activeGame.AddPlayer(game.NewPlayer(c, activeGame))
	}
	if len(activeGame.Players()) > 0 {
		activeGame.SetActivePlayer(activeGame.Players()[0])
	}
	// Game.
	savepath := filepath.Join(mod.Conf().SavesPath(), savename)
	modData, err := flamedata.ImportModule(savepath)
	if err != nil {
		return fmt.Errorf("unable to import module file: %v", err)
	}
	m := flame.NewModule(modData)
	activeGame = game.New(m)
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
