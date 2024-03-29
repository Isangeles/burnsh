/*
 * newcharacter.go
 *
 * Copyright 2018-2023 Dariusz Sikora <ds@isangeles.dev>
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
	"github.com/isangeles/flame/character"
)

const playerIDPrefix = "player_"

// newCharacterDialog starts CLI dialog to create new playable
// game character.
func newCharacterDialog(mod *flame.Module) (flameres.CharacterData, error) {
	var data flameres.CharacterData
	if mod == nil {
		return data, fmt.Errorf("no module loaded")
	}
	// Character creation dialog
	name := ""
	scan := bufio.NewScanner(os.Stdin)
	for mainAccept := false; !mainAccept; {
		// Name
		fmt.Printf("%s:", lang.Text("cli_newchar_name"))
		for scan.Scan() {
			name = scan.Text()
			if !charNameValid(name) {
				fmt.Printf("%s\n", lang.Text("cli_newchar_invalid_name_err"))
				fmt.Printf("%s:", lang.Text("cli_newchar_name"))
				continue
			}
			break
		}
		// Race.
		race := raceDialog()
		// Gender.
		sex := genderDialog()
		// Attributes.
		attrs := character.Attributes{}
		attrsPts := mod.Chapter().Conf().StartAttrs
		for accept := false; !accept; {
			attrs = newAttributesDialog(attrsPts)
			fmt.Printf("%s: %v\n", lang.Text("cli_newchar_attrs_summary"), attrs)
			fmt.Printf("%s:", lang.Text("cli_accept_dialog"))
			scan.Scan()
			input := scan.Text()
			if input != "r" {
				accept = true
			}
		}
		// Summary.
		charID := fmt.Sprintf("%s%s", playerIDPrefix, name)
		charData := flameres.CharacterData{
			ID:        charID,
			Level:     1,
			Sex:       string(sex),
			Race:      race,
			Attitude:  string(character.Friendly),
			Alignment: string(character.TrueNeutral),
		}
		charData.Attributes = flameres.AttributesData{
			Str: attrs.Str,
			Con: attrs.Con,
			Dex: attrs.Dex,
			Int: attrs.Int,
			Wis: attrs.Wis,
		}
		fmt.Printf("%s: %s\n", lang.Text("cli_newchar_summary"),
			charDataDisplayString(charData))
		fmt.Printf("%s:", lang.Text("cli_accept_dialog"))
		scan.Scan()
		input := scan.Text()
		if input != "r" {
			data = charData
			mainAccept = true
		}
	}
	// Add translation for new character name.
	nameTrans := flameres.TranslationData{data.ID, []string{name}}
	lang.AddTranslation(nameTrans)
	// Add player skills & items from interface config.
	for _, sid := range mod.Chapter().Conf().StartSkills {
		skill := flameres.ObjectSkillData{ID: sid}
		data.Skills = append(data.Skills, skill)
	}
	for _, iid := range mod.Chapter().Conf().StartItems {
		item := flameres.InventoryItemData{ID: iid}
		data.Inventory.Items = append(data.Inventory.Items, item)
	}
	return data, nil
}

// raceDialog starts CLI dialog for game character race.
// Returns character race.
func raceDialog() string {
	scan := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s:", lang.Text("cli_newchar_race"))
	races := make([]flameres.RaceData, 0)
	for _, r := range flameres.Races {
		if !r.Playable {
			continue
		}
		races = append(races, r)
	}
	race := ""
	for len(race) < 1 {
		fmt.Printf("[")
		for i, r := range races {
			fmt.Printf("%d - %s ", i, lang.Text(r.ID))
		}
		fmt.Printf("]:")
		scan.Scan()
		input := scan.Text()
		i, err := strconv.Atoi(input)
		if err != nil || i < 0 || i > len(races)-1 {
			fmt.Printf("%s: %s\n", lang.Text("cli_newchar_invalid_value_err"),
				input)
			continue
		}
		race = races[i].ID
	}
	return race
}

// genderDialog starts CLI dialog for game character gender.
// Returns character gender.
func genderDialog() character.Gender {
	scan := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s:", lang.Text("cli_newchar_gender"))
	genderNames := make([]string, 2)
	genderNames[0] = lang.Text(string(character.Male))
	genderNames[1] = lang.Text(string(character.Female))
	s := make([]interface{}, 0)
	for _, v := range genderNames {
		s = append(s, v)
	}
	for true {
		fmt.Printf("[1 - %s, 2 - %s]:", s...)
		scan.Scan()
		input := scan.Text()
		switch input {
		case "1":
			return character.Male
		case "2":
			return character.Female
		default:
			fmt.Printf("%s: %s\n", lang.Text("cli_newchar_invalid_value_err"),
				input)
		}
	}
	return character.Male
}

// newAttributesDialog Starts CLI dialog for game character attributes.
// Returns character attributes.
func newAttributesDialog(attrsPoints int) (attrs character.Attributes) {
	scan := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s:\n", lang.Text("cli_newchar_attrs"))
	for attrsPoints > 0 {
		// Strenght.
		for true {
			fmt.Printf("%s[%s = %d, %s = %d]+", lang.Text("attr_str"),
				lang.Text("cli_newchar_value"), attrs.Str,
				lang.Text("cli_newchar_points"), attrsPoints)
			scan.Scan()
			input := scan.Text()
			attr, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s:%s\n",
					lang.Text("cli_newchar_nan_error"), input)
			} else {
				if attrsPoints-attr >= 0 {
					attrs.Str += attr
					attrsPoints -= attr
					break
				} else {
					fmt.Printf("%s\n", lang.Text("cli_newchar_no_pts_error"))
				}
			}
		}
		// Constitution.
		for true {
			fmt.Printf("%s[%s = %d, %s = %d]+", lang.Text("attr_con"),
				lang.Text("cli_newchar_value"), attrs.Con,
				lang.Text("cli_newchar_points"), attrsPoints)
			scan.Scan()
			input := scan.Text()
			attr, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s:%s\n", lang.Text("cli_newchar_nan_error"), input)
			} else {
				if attrsPoints-attr >= 0 {
					attrs.Con += attr
					attrsPoints -= attr
					break
				} else {
					fmt.Printf("%s\n", lang.Text("cli_newchar_no_pts_error"))
				}
			}

		}
		// Dexterity.
		for true {
			fmt.Printf("%s[%s = %d, %s = %d]+", lang.Text("attr_dex"),
				lang.Text("cli_newchar_value"), attrs.Dex,
				lang.Text("cli_newchar_points"), attrsPoints)
			scan.Scan()
			input := scan.Text()
			attr, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s:%s\n", lang.Text("cli_newchar_nan_error"), input)
			} else {
				if attrsPoints-attr >= 0 {
					attrs.Dex += attr
					attrsPoints -= attr
					break
				} else {
					fmt.Printf("%s\n", lang.Text("cli_newchar_no_pts_error"))
				}
			}
		}
		// Wisdom.
		for true {
			fmt.Printf("%s[%s = %d, %s = %d]+", lang.Text("attr_wis"),
				lang.Text("cli_newchar_value"), attrs.Wis,
				lang.Text("cli_newchar_points"), attrsPoints)
			scan.Scan()
			input := scan.Text()
			attr, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s:%s\n", lang.Text("cli_newchar_nan_error"), input)
			} else {
				if attrsPoints-attr >= 0 {
					attrs.Wis += attr
					attrsPoints -= attr
					break
				} else {
					fmt.Printf("%s\n", lang.Text("cli_newchar_no_pts_error"))
				}
			}
		}
		// Inteligence.
		for true {
			fmt.Printf("%s[%s = %d, %s = %d]+", lang.Text("attr_int"),
				lang.Text("cli_newchar_value"), attrs.Int,
				lang.Text("cli_newchar_points"), attrsPoints)
			scan.Scan()
			input := scan.Text()
			attr, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("%s:%s\n", lang.Text("cli_newchar_nan_error"),
					input)
			} else {
				if attrsPoints-attr >= 0 {
					attrs.Int += attr
					attrsPoints -= attr
					break
				} else {
					fmt.Printf("%s\n",
						lang.Text("cli_newchar_no_pts_error"))
				}
			}
		}

	}
	return
}

// charNameVaild Checks if specified name
// is valid character name.
func charNameValid(name string) bool {
	return len(name) > 0
}
