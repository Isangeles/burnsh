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

	"github.com/isangeles/flame/core"
	"github.com/isangeles/flame/core/data/res/lang"
	"github.com/isangeles/flame/core/module/character"

	"github.com/isangeles/burn"
)

var (
	playableChars []*character.Character
)

// newGameDialog starts CLI dialog for new game.
func newGameDialog() error {
	if mod == nil {
		return fmt.Errorf("no module loaded")
	}
	if len(playableChars) < 1 {
		return fmt.Errorf(lang.Text("cli_newgame_no_chars_err"))
	}
	var pc *character.Character
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
				pc = playableChars[id]
				break
			}
		}

		fmt.Printf("%s: %v\n", lang.Text("cli_newgame_summary"),
			charDisplayString(pc))
		fmt.Printf("%s:", lang.Text("cli_accept_dialog"))
		scan.Scan()
		input := scan.Text()
		if input != "r" {
			accept = true
		}
	}
	players = append(players, pc)
	game = core.NewGame(mod)
	// All players to start area.
	chapter := mod.Chapter()
	startArea := chapter.Area(chapter.Conf().StartArea)
	if startArea == nil {
		return fmt.Errorf("start area not found: %s",
			chapter.Conf().StartArea)
	}
	for _, pc := range players {
		startArea.AddCharacter(pc)
	}
	// Set start positions for players.
	for _, pc := range players {
		pc.SetPosition(chapter.Conf().StartPosX, chapter.Conf().StartPosY)
	}
	burn.Game = game
	activePC = players[0]
	return nil
}
