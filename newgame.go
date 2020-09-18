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
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module/character"

	"github.com/isangeles/burn"

	"github.com/isangeles/burnsh/config"
	"github.com/isangeles/burnsh/game"
	"github.com/isangeles/burnsh/log"
)

var (
	playableChars []*character.Character
	newGamePlayer *character.Character
)

// newGameDialog starts CLI dialog for new game.
func newGameDialog() error {
	if mod == nil {
		return fmt.Errorf("no module loaded")
	}
	if len(playableChars) < 1 {
		return fmt.Errorf(lang.Text("cli_newgame_no_chars_err"))
	}
	scan := bufio.NewScanner(os.Stdin)
	for accept := false; !accept; {
		fmt.Printf("%s:\n", lang.Text("cli_newgame_chars"))
		for i, c := range playableChars {
			fmt.Printf("[%d]%v\n", i, charDisplayString(c))
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
				newGamePlayer = playableChars[id]
				break
			}
		}

		fmt.Printf("%s: %v\n", lang.Text("cli_newgame_summary"),
			charDisplayString(newGamePlayer))
		fmt.Printf("%s:", lang.Text("cli_accept_dialog"))
		scan.Scan()
		input := scan.Text()
		if input != "r" {
			accept = true
		}
	}
	activeGame = game.New(flame.NewGame(mod), server)
	if activeGame.Server() != nil {
		activeGame.SetOnLoginFunc(onServerLogin)
		err := activeGame.Server().Login(config.ServerLogin,
			config.ServerPass)
		if err != nil {
			return fmt.Errorf("Unable to send login request: %v",
				err)
		}
		return nil
	}
	activeGame.AddPlayer(newGamePlayer)
	activeGame.SetActivePlayer(activeGame.Players()[0])
	burn.Game = activeGame.Game
	return nil
}

// onServerLogin callback function called after successful login.
func onServerLogin(game *game.Game) {
	game.SetOnLoginFunc(nil)
	err := game.AddPlayer(newGamePlayer)
	if err != nil {
		log.Err.Printf("New game: Unable to add player: %v",
			err)
		return
	}
	burn.Game = game.Game
	fmt.Printf("Game started\n")
}
