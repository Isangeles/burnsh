/*
 * quests.go
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

// questsDialog starts quests journal CLI dialog.
func questsDialog() error {
	if activeGame == nil {
		return fmt.Errorf("No active game")
	}
	if activeGame.ActivePlayer() == nil {
		return fmt.Errorf("No active PC")
	}
	fmt.Printf("%s:\n", lang.Text("quests_list"))
	for i, q := range activeGame.ActivePlayer().Journal().Quests() {
		questInfo := lang.Texts(q.ID())
		fmt.Printf("[%d]%s\n", i, questInfo[0])
		if len(questInfo) > 1 {
			fmt.Printf("\t%s\n", questInfo[1])
		}
		if q.Completed() {
			completeInfo := lang.Text("quests_q_completed")
			fmt.Printf("\t%s\n", completeInfo)
			continue
		}
		if q.ActiveStage() == nil {
			continue
		}
		stageInfo := lang.Texts(q.ActiveStage().ID())
		fmt.Printf("\t%s\n", stageInfo[0])
	}
	return nil
}
