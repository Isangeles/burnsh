/*
 * modbuilder.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"path/filepath"

	"github.com/isangeles/flame/data"
	"github.com/isangeles/flame/data/res"
)

// NewModule creates new module directory
// in data/modules with one chapter and one
// empty area.
func NewModule(name string) error {
	// Module data.
	modData := res.ModuleData{}
	modData.Config["id"] = []string{name}
	path := filepath.Join("data/modules", name)
	modData.Config["path"] = []string{path}
	modData.Config["chapter"] = []string{"ch1"}
	// Chapter data.
	modData.Chapter.Config["id"] = []string{"ch1"}
	chapterPath := filepath.Join(path, "chapters/ch1")
	modData.Chapter.Config["path"] = []string{chapterPath}
	// Start area data.
	startArea := res.AreaData{ID: "area1"}
	modData.Chapter.Areas = append(modData.Chapter.Areas, startArea)
	// Export module.
	err := data.ExportModule(path, modData)
	if err != nil {
		return fmt.Errorf("unable to export module: %v", err)
	}
	return nil
}
