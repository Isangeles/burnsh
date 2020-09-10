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

// Package with CLI configuration values.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/isangeles/flame/data/text"

	"github.com/isangeles/burnsh/log"
)

const (
	ConfigFileName = ".burnsh"
)

var (
	Module      = ""
	Lang        = "english"
	Fire        = false
	ServerHost  = ""
	ServerPort  = "8000"
	ServerLogin = ""
	ServerPass  = ""
	Debug       = false
)

// Load loads the CLI config file.
func Load() error {
	file, err := os.Open(ConfigFileName)
	if err != nil {
		return fmt.Errorf("unable to open config file: %v", err)
	}
	defer file.Close()
	conf, err := text.UnmarshalConfig(file)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config file: %v", err)
	}
	if len(conf["module"]) > 0 {
		Module = conf["module"][0]
	}
	if len(conf["lang"]) > 0 {
		Lang = conf["lang"][0]
	}
	if len(conf["fire"]) > 0 {
		Fire = conf["fire"][0] == "true"
	}
	if len(conf["server"]) > 1 {
		ServerHost = conf["server"][0]
		ServerPort = conf["server"][1]
	}
	if len(conf["server-user"]) > 1 {
		ServerLogin = conf["server-user"][0]
		ServerPass = conf["server-user"][1]
	}
	if len(conf["debug"]) > 0 {
		Debug = conf["debug"][0] == "true"
	}
	log.Dbg.Println("Config file loaded")
	return nil
}

// Save saves current config values in the config file.
func Save() error {
	// Create file.
	file, err := os.Create(ConfigFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// Marshal config.
	conf := make(map[string][]string)
	conf["module"] = []string{Module}
	conf["lang"] = []string{Lang}
	conf["fire"] = []string{fmt.Sprintf("%v", Fire)}
	conf["server"] = []string{ServerHost, ServerPort}
	conf["server-user"] = []string{ServerLogin, ServerPass}
	conf["debug"] = []string{fmt.Sprintf("%v", Debug)}
	confText := text.MarshalConfig(conf)
	// Write to file.
	w := bufio.NewWriter(file)
	w.WriteString(confText)
	// Save.
	w.Flush()
	log.Dbg.Println("Config file saved")
	return nil
}

// ModulePath returns path to the directory of current module.
func ModulePath() string {
	return filepath.Join("data/modules", Module)
}

// LangPath returns path to the CLI lang directory.
func LangPath() string {
	return filepath.Join("data/lang", Lang)
}

// ScriptsPath returns path to the scripts directory.
func ScriptsPath() string {
	return filepath.FromSlash("data/scripts")
}
