/*
 * talk.go
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

	"github.com/isangeles/flame/core/data/res/lang"
	"github.com/isangeles/flame/core/module/character"
	"github.com/isangeles/flame/core/module/dialog"
	"github.com/isangeles/flame/core/module/effect"
)

// talkDialog starts CLI dialog for dialog with
// current target of active PC.
func talkDialog() error {
	if game == nil {
		msg := lang.Text("no_game_err")
		return fmt.Errorf(msg)
	}
	if activePC == nil {
		msg := lang.Text("no_pc_err")
		return fmt.Errorf(msg)
	}
	tar := activePC.Targets()[0]
	if tar == nil {
		msg := lang.Text("no_tar_err")
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
	d.Restart()
	// Dialog.
	for {
		fmt.Printf("%s:\n", lang.Text("talk_dialog"))
		// Dialog phase.
		phase := dialogPhase(d.Phases(), activePC)
		if phase == nil {
			return fmt.Errorf(lang.Text("talk_no_phase_err"))
		}
		// Dialog phase text.
		dlgText := lang.Text(phase.ID())
		fmt.Printf("[%s]:%s\n", d.Owner().Name(), dlgText)
		// Phase modifiers.
		if owner, ok := d.Owner().(effect.Target); ok {
			owner.TakeModifiers(activePC, phase.OwnerModifiers()...)
		}
		if owner, ok := d.Owner().(effect.Target); ok {
			activePC.TakeModifiers(owner, phase.TalkerModifiers()...)
		}
		// Answer.
		var (
			ans     *dialog.Answer
			ansText string
		)
		for ans == nil {
			// Select answers.
			answers := make([]*dialog.Answer, 0)
			for _, a := range phase.Answers() {
				if !activePC.MeetReqs(a.Requirements()...) {
					continue
				}
				answers = append(answers, a)
			}
			// Print answers.
			fmt.Printf("%s:\n", lang.Text("talk_answers"))
			for i, a := range answers {
				ansText = lang.Text(a.ID())
				fmt.Printf("[%d]%s\n", i, ansText)
			}
			// Select answer.
			fmt.Printf("%s:", lang.Text("talk_answers_select"))
			scan.Scan()
			input := scan.Text()
			id, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s:%s\n", lang.Text("nan_err"), input)
				continue
			}
			if id < 0 || id > len(phase.Answers())-1 {
				fmt.Printf("%s\n", lang.Text("talk_no_answer_id_err"))
				continue
			}
			ans = answers[id]
			ansText = lang.Text(ans.ID())
			// Answer modifiers.
			if owner, ok := d.Owner().(effect.Target); ok {
				owner.TakeModifiers(activePC, ans.OwnerModifiers()...)
			}
			if owner, ok := d.Owner().(effect.Target); ok {
				activePC.TakeModifiers(owner, ans.TalkerModifiers()...)
			}
		}
		fmt.Printf("[%s]:%s\n", activePC.Name(), ansText)
		// Dialog progress.
		d.Next(ans)
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

// dialogPhase selects dialog phase with requirements
// met by specified character.
func dialogPhase(phases []*dialog.Phase, char *character.Character) *dialog.Phase {
	for _, p := range phases {
		if char.MeetReqs(p.Requirements()...) {
			return p
		}
	}
	return nil
}
