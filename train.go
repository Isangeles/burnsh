/*
 * train.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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

	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/core/data/text/lang"
	char "github.com/isangeles/flame/core/module/object/character"
	"github.com/isangeles/flame/core/module/train"
)

// tradeDialog starts CLI dialog for train
// with current PC target.
func trainDialog() error {
	langPath := flameconf.LangPath()
	if game == nil {
		msg := lang.TextDir(langPath, "no_game_err")
		return fmt.Errorf(msg)
	}
	if activePC == nil {
		msg := lang.TextDir(langPath, "no_pc_err")
		return fmt.Errorf(msg)
	}
	tar := activePC.Targets()[0]
	if tar == nil {
		msg := lang.TextDir(langPath, "no_tar_err")
		return fmt.Errorf(msg)
	}
	tarChar, ok := tar.(*char.Character)
	if !ok {
		msg := lang.TextDir(langPath, "tar_invalid")
		return fmt.Errorf(msg)
	}
	t := selectTraining(tarChar.Trainings())
	if t == nil {
		msg := lang.TextDir(langPath, "train_no_train_sel")
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
	langPath := flameconf.LangPath()
	if len(trainings) < 1 {
		msg := lang.TextDir(langPath, "train_no_trainings")
		fmt.Printf("%s\n", msg)
		return nil
	}
	var training train.Training
	for training == nil {
		// List available trainings.
		fmt.Printf("%s:\n", lang.TextDir(langPath, "train_trainings"))
		for i, t := range trainings {
			fmt.Printf("\t[%d]%v\n", i, t)
		}
		fmt.Printf("%s:", lang.TextDir(langPath, "train_select_training"))
		// Scan input.
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		input := scan.Text()
		if input == "" {
			break
		}
		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s:%v\n", lang.TextDir(langPath, "nan_err"), err)
			continue
		}
		if id < 0 || id > len(trainings)-1 {
			fmt.Printf("%s:%s\n", lang.TextDir(langPath, "invalid_input_err"),
				input)
			continue
		}
		training = trainings[id]
	}
	return training
}
