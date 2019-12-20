/*
 * modbuilder.go
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
	"path/filepath"

	"github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/parsexml"
	"github.com/isangeles/flame/core/data/res"
)

// NewModule creates new module directory
// in data/modules with all one chapter and
// one empty start area.
func NewModule(name string) error {
	path := filepath.FromSlash("data/modules/" + name)
	// Mod dir.
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("fail to create module dir: %v", err)
	}
	// Sub-dirs.
	err = os.MkdirAll(filepath.FromSlash(path+"/characters"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create characters dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/items"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create items dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/effects"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create effects dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/skills"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create skills dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/objects"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create objects dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/recipes"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create recipes dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/lang"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create lang dir: %v", err)
	}
	// Mod conf.
	confPath := filepath.FromSlash(path + "/mod.conf")
	confFile, err := os.Create(confPath)
	if err != nil {
		return fmt.Errorf("fail to create module conf file: %v", err)
	}
	defer confFile.Close()
	confFormat := "id:%s;\nstart-chapter:%s;\nchar-skills:%s;\nchar-items:%s;\n"
	conf := fmt.Sprintf(confFormat, name, "prologue", "", "")
	w := bufio.NewWriter(confFile)
	w.WriteString(conf)
	w.Flush()
	// Start chapter.
	err = createChapter(path+"/chapters/prologue", "prologue")
	if err != nil {
		return fmt.Errorf("fail to create chapter: %v", err)
	}
	return nil
}

// createChapter creates new chapter
// directory.
func createChapter(path, id string) error {
	// Dir.
	err := os.MkdirAll(filepath.FromSlash(path), 0755)
	if err != nil {
		return fmt.Errorf("fail to create dir: %v", err)
	}
	// Sub-dirs.
	err = os.MkdirAll(filepath.FromSlash(path+"/npc"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create npc dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/dialogs"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create dialogs dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/quests"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create quests dir: %v", err)
	}
	err = os.MkdirAll(filepath.FromSlash(path+"/lang"), 0755)
	if err != nil {
		return fmt.Errorf("fail to create lang dir: %v", err)
	}
	// Conf.
	confPath := filepath.FromSlash(path + "/chapter.conf")
	confFile, err := os.Create(confPath)
	if err != nil {
		return fmt.Errorf("fail to create conf file: %v", err)
	}
	defer confFile.Close()
	conf := fmt.Sprintf("start-area:%s;\n", "area1_main")
	w := bufio.NewWriter(confFile)
	w.WriteString(conf)
	w.Flush()
	// Start area.
	areaDirPath := filepath.FromSlash(path + "/areas/area1_main")
	err = os.MkdirAll(areaDirPath, 0755)
	if err != nil {
		return fmt.Errorf("fail to create areas dir: %v", err)
	}
	ad := res.AreaData{ID: "area1_main"}
	xmlArea, err := parsexml.MarshalArea(&ad)
	if err != nil {
		return fmt.Errorf("fail to marshal start area: %v", err)
	}
	
	areaPath := filepath.FromSlash(areaDirPath + "/main" + data.AreaFileExt)
	areaFile, err := os.Create(areaPath)
	if err != nil {
		return fmt.Errorf("fail to create start area file: %v", err)
	}
	defer areaFile.Close()
	w = bufio.NewWriter(areaFile)
	w.WriteString(xmlArea)
	w.Flush()
	return nil
}
