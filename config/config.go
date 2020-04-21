/*
 * config.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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

package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/isangeles/flame/data/text"

	"github.com/isangeles/burnsh/log"
)

const (
	ConfigFileName = ".burnsh"
)

var (
	NewCharAttrs  = 10
	NewCharSkills []string
	NewCharItems  []string
)

// LoadConfig Loads CLI config file.
func LoadConfig() error {
	file, err := os.Open(ConfigFileName)
	if err != nil {
		return fmt.Errorf("unable to open config file: %v", err)
	}
	defer file.Close()
	conf, err := text.UnmarshalConfig(file)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config file: %v", err)
	}
	if len(conf["new-char-attrs"]) > 0 {
		NewCharAttrs, _ = strconv.Atoi(conf["new-char-attrs"][0])
	}
	NewCharSkills = conf["new-char-skills"]
	NewCharItems = conf["new-char-items"]
	log.Dbg.Println("Config file loaded")
	return nil
}

// SaveConfig Saves current config values in config file.
func SaveConfig() error {
	// Create file.
	file, err := os.Create(ConfigFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// Marshal config.
	conf := make(map[string][]string)
	conf["new-char-attrs"] = []string{fmt.Sprintf("%d", NewCharAttrs)}
	conf["new-char-skills"] = NewCharSkills
	conf["new-char-items"] = NewCharItems
	confText := text.MarshalConfig(conf)
	// Write values.
	w := bufio.NewWriter(file)
	w.WriteString(confText)
	// Save.
	w.Flush()
	log.Dbg.Println("Config file saved")
	return nil
}

// ScriptsPath returns path to
// scripts directory.
func ScriptsPath() string {
	return "data/scripts"
}
