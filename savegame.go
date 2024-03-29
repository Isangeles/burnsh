/*
 * savegame.go
 *
 * Copyright 2020-2023 Dariusz Sikora <ds@isangeles.dev>
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
	"os"
	"path/filepath"

	flamedata "github.com/isangeles/flame/data"
	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/burnsh/config"
	"github.com/isangeles/burnsh/log"
)

var (
	SaveExt         = ".savecli"
	ModuleSavesPath = "burnsh/saves"
)

// Struct for CLI save node.
type CLISave struct {
	XMLName xml.Name     `xml:"save"`
	Name    string       `xml:"name,attr"`
	Players []PlayerSave `xml:"players>player"`
}

// Struct for CLI player node.
type PlayerSave struct {
	XMLName xml.Name `xml:"player"`
	ID      string   `xml:"id,attr"`
	Serial  string   `xml:"serial,attr"`
}

// saveGameDialog starts CLI dialog for saving
// current game state.
func saveGameDialog() error {
	if activeGame == nil {
		return fmt.Errorf("no game started")
	}
	// CLI.
	save := new(CLISave)
	scan := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s:", lang.Text("savegame_save_name"))
	for scan.Scan() {
		save.Name = scan.Text()
		if len(save.Name) > 0 {
			break
		}
	}
	for _, pc := range activeGame.Players() {
		pcSave := PlayerSave{
			ID:     pc.ID(),
			Serial: pc.Serial(),
		}
		save.Players = append(save.Players, pcSave)
	}
	cliSavepath := filepath.Join(mod.Conf().Path, ModuleSavesPath)
	err := saveCLI(save, cliSavepath)
	if err != nil {
		return fmt.Errorf("unable to save cli: %v", err)
	}
	// Game.
	if activeGame.Server() != nil {
		req := request.Request{Save: []string{save.Name}}
		err := activeGame.Server().Send(req)
		if err != nil {
			return fmt.Errorf("unable to send save request: %v",
				err)
		}
		return nil
	}
	savepath := filepath.Join(config.ModulesPath, save.Name)
	err = flamedata.ExportModule(savepath, activeGame.Data())
	if err != nil {
		return fmt.Errorf("unable to export module: %v", err)
	}
	return nil
}

// saveCLI saves CLI state in file under specified path.
func saveCLI(save *CLISave, path string) error {
	out, err := xml.Marshal(save)
	if err != nil {
		return fmt.Errorf("unable to marshal save: %v", err)
	}
	xml := string(out[:])
	savePath := fmt.Sprintf("%s/%s%s", path, save.Name, SaveExt)
	savePath = filepath.FromSlash(savePath)
	err = os.MkdirAll(filepath.Dir(savePath), 0755)
	if err != nil {
		return fmt.Errorf("unable to create save directory: %v", err)
	}
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("unable to create save file: %v", err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.WriteString(xml)
	w.Flush()
	log.Dbg.Printf("cli state saved in: %s", savePath)
	return nil
}
