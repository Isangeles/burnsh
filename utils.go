/*
 * utils.go
 *
 * Copyright 2018-2022 Dariusz Sikora <ds@isangeles.dev>
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

	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/req"
)

// charDataDisplayString returns string with character
// stats and info.
func charDataDisplayString(data flameres.CharacterData) string {
	return fmt.Sprintf("%s:%s,%s,%s:%d,%d,%d,%d,%d",
		lang.Text(data.ID), data.Race, data.Sex,
		"Stats", data.Attributes.Str, data.Attributes.Con,
		data.Attributes.Dex, data.Attributes.Wis,
		data.Attributes.Int)
}

// reqsInfo returns text with info to display
// about specified requirements.
func reqsInfo(reqs ...req.Requirement) string {
	out := ""
	for _, r := range reqs {
		switch r := r.(type) {
		case *req.Level:
			out = fmt.Sprintf("%s\n%s\t%d", out, lang.Text("req_level"),
				r.MinLevel())
		case *req.Item:
			out = fmt.Sprintf("%s\n%s:%s\tx%d", out, lang.Text("req_item"),
				r.ItemID(), r.ItemAmount())
		default:
			out = fmt.Sprintf("%s\n%s", out, lang.Text("req_unknown"))
		}
	}
	return out
}
