/*
 * useskill.go
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
	"github.com/isangeles/flame/module/skill"
)

// useSkillDialog starts CLI dialog for
// using skills.
func useSkillDialog() error {
	if activePC == nil {
		msg := lang.Text("no_pc_err")
		return fmt.Errorf(msg)
	}
	// List skills.
	fmt.Printf("%s:\n", lang.Text("useskill_skills"))
	skills := activePC.Skills()
	for i, s := range skills {
		fmt.Printf("[%d]%s\n", i, s.Name())
	}
	// Select skill.
	scan := bufio.NewScanner(os.Stdin)
	var skill *skill.Skill
	for skill == nil {
		fmt.Printf("%s:", lang.Text("useskill_select"))
		scan.Scan()
		input := scan.Text()
		id, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("%s:%s\n", lang.Text("nan_err"), input)
			continue
		}
		if id < 0 || id > len(skills)-1 {
			fmt.Printf("%s:%s\n", lang.Text("invalid_input_err"), input)
			continue
		}
		skill = skills[id]
	}
	activePC.UseSkill(skill)
	return nil
}
