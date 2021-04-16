/*
 * cli.go
 *
 * Copyright 2018-2021 Dariusz Sikora <dev@isangeles.pl>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either Version 2 of the License, or
 * (at your option) any later Version.
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

// Command line interface for flame engine.
// Uses Burn CI to handle user input and communicate with Flame Engine.
// All commands to be handled by CI must starts with dollar sign($),
// otherwise input is directly send to out(like 'echo').
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/isangeles/flame"
	flamedata "github.com/isangeles/flame/data"

	"github.com/isangeles/burn"
	"github.com/isangeles/burn/ash"
	"github.com/isangeles/burn/syntax"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/burnsh/config"
	"github.com/isangeles/burnsh/data"
	"github.com/isangeles/burnsh/game"
	"github.com/isangeles/burnsh/log"
)

const (
	Name           = "Burn Shell"
	Version        = "0.1.0-dev"
	CommandPrefix  = "$"
	ScriptPrefix   = "%"
	RunBGSuffix    = "&"
	CloseCmd       = "close"
	LoginCmd       = "login"
	NewCharCmd     = "newchar"
	NewGameCmd     = "newgame"
	NewModCmd      = "newmod"
	SaveGameCmd    = "savegame"
	LoadGameCmd    = "loadgame"
	ImportCharsCmd = "importchars"
	MoveCmd        = "move"
	LootTargetCmd  = "loot"
	TalkTargetCmd  = "talk"
	FindTargetCmd  = "target"
	TargetInfoCmd  = "tarinfo"
	QuestsCmd      = "quests"
	UseSkillCmd    = "useskill"
	CraftingCmd    = "crafting"
	TradeTargetCmd = "trade"
	TrainTargetCmd = "train"
	EquipCmd       = "equip"
	InventoryCmd   = "inventory"
	RepeatInputCmd = "!"
	InputIndicator = ">"
)

var (
	mod         *flame.Module
	server      *game.Server
	activeGame  *game.Game
	lastCommand string
	lastUpdate  time.Time
)

// Main function.
func main() {
	fmt.Printf("*%s(%s)@%s(%s)*\n", Name, Version,
		flame.Name, flame.Version)
	// Load CLI config.
	err := config.Load()
	if err != nil {
		log.Err.Printf("Unable to load config: %v", err)
	}
	log.PrintStdOut(config.Debug)
	// Load module.
	err = loadModule(config.ModulePath())
	if err != nil {
		log.Err.Printf("Unable to load module: %v", err)
	}
	// Load UI data.
	err = data.LoadUIData(filepath.Join(config.ModulePath(), data.UIDirPath))
	if err != nil {
		log.Err.Printf("Unable to load UI data: %v", err)
	}
	// Fire server.
	if config.Multiplayer() {
		serv, err := game.NewServer(config.ServerHost, config.ServerPort)
		if err != nil {
			panic(fmt.Errorf("Unable to create game server connection: %v",
				err))
		}
		server = serv
		server.SetOnResponseFunc(handleResponse)
		log.Inf.Printf("Connected to the game server at: %s", server.Address())
	}
	fmt.Print(InputIndicator)
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		input := scan.Text()
		if strings.HasPrefix(input, CommandPrefix) {
			cmd := strings.TrimPrefix(input, CommandPrefix)
			execute(cmd)
			lastCommand = cmd
		} else if strings.HasPrefix(input, ScriptPrefix) {
			input := strings.TrimPrefix(input, ScriptPrefix)
			scrArgs := strings.Split(input, " ")
			bgrun := false
			if strings.HasSuffix(scrArgs[0], RunBGSuffix) {
				bgrun = true
				scrArgs[0] = strings.TrimSuffix(scrArgs[0], RunBGSuffix)
			}
			executeFile(bgrun, scrArgs[0], scrArgs...)
		} else if activeGame != nil && activeGame.ActivePlayer() != nil {
			activeGame.ActivePlayer().AddChatMessage(input)
		} else {
			log.Inf.Println(input)
		}
		fmt.Print(InputIndicator)
	}
	if err := scan.Err(); err != nil {
		log.Err.Printf("unable to init input scanner : %v\n", err)
	}
}

// execute handles specified command or passes it to CI.
func execute(input string) {
	switch input {
	case CloseCmd:
		err := config.Save()
		if err != nil {
			log.Err.Printf("unable to save config: %v", err)
		}
		if server != nil {
			req := request.Request{Close: time.Now().UnixNano()}
			err := server.Send(req)
			if err != nil {
				log.Err.Printf("Unable to send close request: %v", err)
			}
		}
		os.Exit(0)
	case LoginCmd:
		err := loginDialog()
		if err != nil {
			log.Err.Printf("Login error: %v", err)
			break
		}
	case NewCharCmd:
		charData, err := newCharacterDialog(mod)
		if err != nil {
			log.Err.Printf("%s\n", err)
			break
		}
		playableChars = append(playableChars, charData)
	case NewGameCmd:
		err := newGameDialog()
		if err != nil {
			log.Err.Printf("%s: %v", NewGameCmd, err)
			break
		}
		go gameLoop(activeGame)
	case NewModCmd:
		err := newModDialog()
		if err != nil {
			log.Err.Printf("%s: %v", NewModCmd, err)
			break
		}
	case SaveGameCmd:
		err := saveGameDialog()
		if err != nil {
			log.Err.Printf("%s: %v", SaveGameCmd, err)
		}
	case LoadGameCmd:
		err := loadGameDialog()
		if err != nil {
			log.Err.Printf("%s: %v", LoadGameCmd, err)
			break
		}
		lastUpdate = time.Now()
	case ImportCharsCmd:
		if mod == nil {
			log.Err.Printf("%s: no module loaded", ImportCharsCmd)
		}
		log.Inf.Printf("Imported characters: %d\n", len(mod.Resources().Characters))
		for _, cd := range mod.Resources().Characters {
			playableChars = append(playableChars, cd)
		}
	case MoveCmd:
		err := moveDialog()
		if err != nil {
			log.Err.Printf("%s: %v", MoveCmd, err)
			break
		}
	case LootTargetCmd:
		err := lootDialog()
		if err != nil {
			log.Err.Printf("%s: %v", LootTargetCmd, err)
			break
		}
	case TalkTargetCmd:
		err := talkDialog()
		if err != nil {
			log.Err.Printf("%s: %v", TalkTargetCmd, err)
			break
		}
	case FindTargetCmd:
		err := targetDialog()
		if err != nil {
			log.Err.Printf("%s: %v", FindTargetCmd, err)
			break
		}
	case TargetInfoCmd:
		err := targetInfoDialog()
		if err != nil {
			log.Err.Printf("%s: %v", TargetInfoCmd, err)
			break
		}
	case QuestsCmd:
		err := questsDialog()
		if err != nil {
			log.Err.Printf("%s: %v", QuestsCmd, err)
		}
	case UseSkillCmd:
		err := useSkillDialog()
		if err != nil {
			log.Err.Printf("%s: %v", UseSkillCmd, err)
		}
	case CraftingCmd:
		err := craftingDialog()
		if err != nil {
			log.Err.Printf("%s: %v", CraftingCmd, err)
		}
	case TradeTargetCmd:
		err := tradeDialog()
		if err != nil {
			log.Err.Printf("%s: %v", TradeTargetCmd, err)
		}
	case TrainTargetCmd:
		err := trainDialog()
		if err != nil {
			log.Err.Printf("%s: %v", TrainTargetCmd, err)
		}
	case EquipCmd:
		err := equipDialog()
		if err != nil {
			log.Err.Printf("%s: %v", EquipCmd, err)
		}
	case InventoryCmd:
		err := inventoryDialog()
		if err != nil {
			log.Err.Printf("%s: %v", InventoryCmd, err)
		}
	case RepeatInputCmd:
		execute(lastCommand)
	default: // pass command to CI
		exp, err := syntax.NewSTDExpression(input)
		if err != nil {
			log.Err.Printf("command build error: %v", err)
			break
		}
		res, out := burn.HandleExpression(exp)
		log.Inf.Printf("burn[%d]: %s\n", res, out)
		if server != nil {
			req := request.Request{Command: []string{exp.String()}}
			server.Send(req)
		}
	}
}

// executeFile executes script from data/scripts dir.
func executeFile(bgrun bool, fileName string, args ...string) {
	path := fmt.Sprintf("%s/%s.ash", config.ScriptsPath(),
		fileName)
	file, err := os.Open(path)
	if err != nil {
		log.Err.Printf("unable to open file: %v", err)
		return
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		log.Err.Printf("unable to read file: %v", err)
		return
	}
	scriptName := filepath.Base(path)
	scr, err := ash.NewScript(scriptName, fmt.Sprintf("%s", text), args...)
	if err != nil {
		log.Err.Printf("unable to parse script: %v", err)
		return
	}
	if bgrun {
		go runScript(scr)
		return
	}
	runScript(scr)
}

// runScript runs sprecified Ash script.
func runScript(s *ash.Script) {
	err := ash.Run(s)
	if err != nil {
		log.Err.Printf("unable to run script: %v", err)
		return
	}
}

// gameLoop handles game updating.
func gameLoop(g *game.Game) {
	lastUpdate = time.Now()
	for {
		dtNano := time.Since(lastUpdate).Nanoseconds()
		delta := dtNano / int64(time.Millisecond) // delta to milliseconds
		g.Update(delta)
		updateChat()
		lastUpdate = time.Now()
		// Wait for 16 millis.
		time.Sleep(time.Duration(16) * time.Millisecond)
	}
}

// loadModule loads module with all module data
// from directory with specified path.
func loadModule(path string) error {
	modData, err := flamedata.ImportModule(config.ModulePath())
	if err != nil {
		return fmt.Errorf("unable to import module: %v", err)
	}
	mod = flame.NewModule()
	mod.Apply(modData)
	burn.Module = mod
	return nil
}
