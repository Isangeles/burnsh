/*
 * newgame.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module/character"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/burn"

	"github.com/isangeles/burnsh/game"
)

var (
	playableChars []flameres.CharacterData
)

// newGameDialog starts CLI dialog for new game.
func newGameDialog() error {
	if mod == nil {
		return fmt.Errorf("no module loaded")
	}
	if len(playableChars) < 1 {
		return fmt.Errorf(lang.Text("cli_newgame_no_chars_err"))
	}
	var playerData flameres.CharacterData
	scan := bufio.NewScanner(os.Stdin)
	for accept := false; !accept; {
		fmt.Printf("%s:\n", lang.Text("cli_newgame_chars"))
		for i, c := range playableChars {
			fmt.Printf("[%d]%v\n", i, charDataDisplayString(c))
		}
		fmt.Printf("%s:", lang.Text("cli_newgame_select_char"))
		for scan.Scan() {
			input := scan.Text()
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s:%s\n",
					lang.Text("cli_nan_error"), input)
			}
			if id >= 0 && id < len(playableChars) {
				playerData  = playableChars[id]
				break
			}
		}
		fmt.Printf("%s: %v\n", lang.Text("cli_newgame_summary"),
			charDataDisplayString(playerData))
		fmt.Printf("%s:", lang.Text("cli_accept_dialog"))
		scan.Scan()
		input := scan.Text()
		if input != "r" {
			accept = true
		}
	}
	activeGame = game.New(flame.NewGame(mod))
	burn.Game = activeGame.Game
	if server != nil {
		activeGame.SetServer(server)
		req := request.Request{NewChar: []flameres.CharacterData{playerData}}
		err := activeGame.Server().Send(req)
		if err != nil {
			return fmt.Errorf("Unable to send new character request: %v",
				err)
		}
		return nil
	}
	char := character.New(playerData)
	err := activeGame.AddPlayer(char)
	if err != nil {
		return fmt.Errorf("Unable to add player: %v", err)
	}
	return nil
}
