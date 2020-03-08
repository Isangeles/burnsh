/*
 * train.go
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

	"github.com/isangeles/flame/data/res/lang"
	char "github.com/isangeles/flame/module/character"
	"github.com/isangeles/flame/module/req"
	"github.com/isangeles/flame/module/train"
)

// tradeDialog starts CLI dialog for train
// with current PC target.
func trainDialog() error {
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
	tarChar, ok := tar.(*char.Character)
	if !ok {
		msg := lang.Text("tar_invalid")
		return fmt.Errorf(msg)
	}
	t := selectTraining(tarChar.Trainings())
	if t == nil {
		msg := lang.Text("train_no_train_sel")
		fmt.Printf("%s\n", msg)
		return nil
	}
	err := activePC.Train(t)
	if err != nil {
		return err
	}
	return nil
}

// selectTrainings starts dialog for selecting training from
// specified trainings.
func selectTraining(trainings []train.Training) train.Training {
	if len(trainings) < 1 {
		msg := lang.Text("train_no_trainings")
		fmt.Printf("%s\n", msg)
		return nil
	}
	var training train.Training
	for training == nil {
		// List available trainings.
		fmt.Printf("%s:\n", lang.Text("train_trainings"))
		for i, t := range trainings {
			fmt.Printf("\t[%d]%s\n", i, trainingInfo(t))
			// List training reqs.
			fmt.Printf("\t%s:\n", lang.Text("train_reqs"))
			for _, r := range t.Reqs() {
				fmt.Printf("\t%s\n", reqInfo(r))
			}
		}
		fmt.Printf("%s:", lang.Text("train_select_training"))
		// Scan input.
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		input := scan.Text()
		if input == "" {
			break
		}
		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s:%v\n", lang.Text("nan_err"), err)
			continue
		}
		if id < 0 || id > len(trainings)-1 {
			fmt.Printf("%s:%s\n", lang.Text("invalid_input_err"), input)
			continue
		}
		training = trainings[id]
	}
	return training
}

// trainingInfo returns information about
// training do display.
func trainingInfo(t train.Training) string {
	// Train info.
	info := ""
	switch t := t.(type) {
	case *train.AttrsTraining:
		info = fmt.Sprintf("%s", lang.Text("train_attrs_training"))
		if t.Strenght() > 0 {
			strLabel := lang.Text("attr_str")
			info = fmt.Sprintf("%s: %s(%d)", info, strLabel, t.Strenght())
		}
		if t.Constitution() > 0 {
			conLabel := lang.Text("attr_con")
			info = fmt.Sprintf("%s: %s(%d)", info, conLabel, t.Constitution())
		}
		if t.Dexterity() > 0 {
			dexLabel := lang.Text("attr_dex")
			info = fmt.Sprintf("%s: %s(%d)", info, dexLabel, t.Dexterity())
		}
		if t.Wisdom() > 0 {
			wisLabel := lang.Text("attr_wis")
			info = fmt.Sprintf("%s: %s(%d)", info, wisLabel, t.Wisdom())
		}
		if t.Intelligence() > 0 {
			intLabel := lang.Text("attr_int")
			info = fmt.Sprintf("%s: %s(%d)", info, intLabel, t.Intelligence())
		}
	default:
		info = fmt.Sprintf("%s%s", info, "unknown")
	}
	return info
}

// reqInfo returns information about specified
// requirement.
func reqInfo(r req.Requirement) string {
	info := ""
	switch r := r.(type) {
	case *req.ItemReq:
		reqLabel := lang.Text("req_item")
		info = fmt.Sprintf("%s: %s x%d", reqLabel, r.ItemID(),
			r.ItemAmount())
	case *req.CurrencyReq:
		reqLabel := lang.Text("req_currency")
		info = fmt.Sprintf("%s: %d", reqLabel, r.Amount())
	default:
		return "unknown"
	}
	return info
}
