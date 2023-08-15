/*
 * talk.go
 *
 * Copyright 2019-2023 Dariusz Sikora <ds@isangeles.dev>
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
	"math"
	"os"
	"strconv"

	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/dialog"
)

var TalkRange float64 = 50

// talkDialog starts CLI dialog for dialog with
// current target of active PC.
func talkDialog() error {
	if activeGame == nil {
		msg := lang.Text("no_game_err")
		return fmt.Errorf(msg)
	}
	if activeGame.ActivePlayer() == nil {
		msg := lang.Text("no_pc_err")
		return fmt.Errorf(msg)
	}
	if len(activeGame.ActivePlayer().Targets()) < 1 {
		msg := lang.Text("no_tar_err")
		return fmt.Errorf(msg)
	}
	// Range check.
	tar := activeGame.ActivePlayer().Targets()[0]
	tarX, tarY := tar.Position()
	pcX, pcY := activeGame.ActivePlayer().Position()
	if math.Hypot(tarX-pcX, tarY-pcY) > TalkRange {
		msg := lang.Text("out_of_range_err")
		return fmt.Errorf(msg)
	}
	tarChar, ok := tar.(*character.Character)
	if !ok {
		return fmt.Errorf("invalid_target")
	}
	if len(tarChar.Dialogs()) < 1 {
		return fmt.Errorf("no_target_dialogs")
	}
	d := tarChar.Dialogs()[0]
	scan := bufio.NewScanner(os.Stdin)
	activeGame.StartDialog(d, activeGame.ActivePlayer())
	// Dialog.
	for {
		fmt.Printf("%s:\n", lang.Text("talk_dialog"))
		// Dialog stage.
		if d.Stage() == nil {
			return fmt.Errorf(lang.Text("talk_no_stage_err"))
		}
		// Dialog stage text.
		fmt.Printf("[%s]: %s\n", lang.Text(d.Owner().ID()),
			dialogText(d, d.Stage().ID()))
		// Answer.
		var answer *dialog.Answer
		for answer == nil {
			// Select answers.
			answers := make([]*dialog.Answer, 0)
			for _, a := range d.Stage().Answers() {
				if !activeGame.ActivePlayer().MeetReqs(a.Requirements()...) {
					continue
				}
				answers = append(answers, a)
			}
			// Print answers.
			fmt.Printf("%s:\n", lang.Text("talk_answers"))
			for i, a := range answers {
				fmt.Printf("[%d]%s\n", i, dialogText(d, a.ID()))
			}
			// Select answer.
			fmt.Printf("%s:", lang.Text("talk_answers_select"))
			scan.Scan()
			input := scan.Text()
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s: %s\n", lang.Text("nan_err"), input)
				continue
			}
			if id < 0 || id > len(d.Stage().Answers())-1 {
				fmt.Printf("%s\n", lang.Text("talk_no_answer_id_err"))
				continue
			}
			answer = answers[id]
		}
		fmt.Printf("[%s]: %s\n", lang.Text(activeGame.ActivePlayer().ID()),
			dialogText(d, answer.ID()))
		// Dialog progress.
		activeGame.AnswerDialog(d, answer)
		if d.Trading() {
			err := tradeDialog()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
		if d.Training() {
			err := trainDialog()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
		if d.Finished() {
			break
		}
	}
	return nil
}

// dialogText returns translated text for dialog with specified ID.
func dialogText(d *dialog.Dialog, textID string) string {
	return d.DialogText(lang.Text(textID))
}
