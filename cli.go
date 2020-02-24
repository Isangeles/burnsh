/*
 * cli.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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
	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/core"
	"github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/module/character"

	"github.com/isangeles/burn"
	"github.com/isangeles/burn/ash"
	"github.com/isangeles/burn/syntax"

	"github.com/isangeles/burnsh/config"
	"github.com/isangeles/burnsh/log"
)

const (
	Name           = "Burn Shell"
	Version        = "0.0.0"
	CommandPrefix  = "$"
	ScriptPrefix   = "%"
	RunBGSuffix    = "&"
	CloseCmd       = "close"
	NewCharCmd     = "newchar"
	NewGameCmd     = "newgame"
	NewModCmd      = "newmod"
	SaveGameCmd    = "savegame"
	LoadGameCmd    = "loadgame"
	ImportCharsCmd = "importchars"
	LootTargetCmd  = "loot"
	TalkTargetCmd  = "talk"
	FindTargetCmd  = "target"
	TargetInfoCmd  = "tarinfo"
	QuestsCmd      = "quests"
	UseSkillCmd    = "useskill"
	CraftingCmd    = "crafting"
	TradeTargetCmd = "trade"
	TrainTargetCmd = "train"
	RepeatInputCmd = "!"
	InputIndicator = ">"
)

var (
	game        *core.Game
	players     []*character.Character
	activePC    *character.Character
	lastCommand string
	lastUpdate  time.Time
)

// On init.
func init() {
	// Load flame config.
	err := flameconf.LoadConfig()
	if err != nil {
		log.Err.Printf("fail to load flame config: %v", err)
	}
	// Load UI translation.
	err = data.LoadTranslationData(flameconf.LangPath())
	// Load module.
	err = loadModule(flameconf.ModulePath(), flameconf.LangID())
	if err != nil {
		log.Err.Printf("fail to load ui translation data: %v", err)
	}
	if err != nil {
		log.Err.Printf("fail to load module: %v", err)
	}
	// Load CLI config.
	err = config.LoadConfig()
	if err != nil {
		log.Err.Printf("fail to load config: %v", err)
	}
}

func main() {
	fmt.Printf("*%s(%s)@%s(%s)*\n", Name, Version,
		flame.Name, flame.Version)
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
		} else if activePC != nil {
			activePC.SendChat(input)
		} else {
			log.Inf.Println(input)
		}
		fmt.Print(InputIndicator)
		// Game update on input.
		if game != nil {
			go gameLoop(game)
		}
	}
	if err := scan.Err(); err != nil {
		log.Err.Printf("input scanner init fail : %v\n", err)
	}
}

// execute passes specified command to CI.
func execute(input string) {
	switch input {
	case CloseCmd:
		err := flameconf.SaveConfig()
		if err != nil {
			log.Err.Printf("engine config save fail: %v",
				err)
		}
		err = config.SaveConfig()
		if err != nil {
			log.Err.Printf("config save fail: %v", err)
		}
		os.Exit(0)
	case NewCharCmd:
		createdChar, err := newCharacterDialog(flame.Mod())
		if err != nil {
			log.Err.Printf("%s\n", err)
			break
		}
		playableChars = append(playableChars, createdChar)
	case NewGameCmd:
		err := newGameDialog()
		if err != nil {
			log.Err.Printf("%s:%v", NewGameCmd, err)
			break
		}
		lastUpdate = time.Now()
	case NewModCmd:
		err := newModDialog()
		if err != nil {
			log.Err.Printf("%s:%v", NewModCmd, err)
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
			log.Err.Printf("%s:%v", LoadGameCmd, err)
			break
		}
		lastUpdate = time.Now()
	case ImportCharsCmd:
		chars, err := data.ImportCharactersDir(flame.Mod().Conf().CharactersPath())
		if err != nil {
			log.Err.Printf("%s:%v", ImportCharsCmd, err)
			break
		}
		log.Inf.Printf("imported chars: %d\n", len(chars))
		for _, c := range chars {
			playableChars = append(playableChars, c)
		}
	case LootTargetCmd:
		err := lootDialog()
		if err != nil {
			log.Err.Printf("%s:%v", LootTargetCmd, err)
			break
		}
	case TalkTargetCmd:
		err := talkDialog()
		if err != nil {
			log.Err.Printf("%s:%v", TalkTargetCmd, err)
			break
		}
	case FindTargetCmd:
		err := targetDialog()
		if err != nil {
			log.Err.Printf("%s:%v", FindTargetCmd, err)
			break
		}
	case TargetInfoCmd:
		err := targetInfoDialog()
		if err != nil {
			log.Err.Printf("%s:%v", TargetInfoCmd, err)
			break
		}
	case QuestsCmd:
		err := questsDialog()
		if err != nil {
			log.Err.Printf("%s:%v", QuestsCmd, err)
		}
	case UseSkillCmd:
		err := useSkillDialog()
		if err != nil {
			log.Err.Printf("%s:%v", UseSkillCmd, err)
		}
	case CraftingCmd:
		err := craftingDialog()
		if err != nil {
			log.Err.Printf("%s:%v", CraftingCmd, err)
		}
	case TradeTargetCmd:
		err := tradeDialog()
		if err != nil {
			log.Err.Printf("%s:%v", TradeTargetCmd, err)
		}
	case TrainTargetCmd:
		err := trainDialog()
		if err != nil {
			log.Err.Printf("%s:%v", TrainTargetCmd, err)
		}
	case RepeatInputCmd:
		execute(lastCommand)
		return
	default: // pass command to CI
		exp, err := syntax.NewSTDExpression(input)
		if err != nil {
			log.Err.Printf("command build error: %v", err)
		}
		res, out := burn.HandleExpression(exp)
		log.Inf.Printf("burn[%d]: %s\n", res, out)
	}
}

// executeFile executes script from data/scripts dir.
func executeFile(bgrun bool, fileName string, args ...string) {
	path := fmt.Sprintf("%s/%s.ash", config.ScriptsPath(),
		fileName)
	file, err := os.Open(path)
	if err != nil {
		log.Err.Printf("fail to open file: %v", err)
		return
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		log.Err.Printf("fail to read file: %v", err)
		return
	}
	scriptName := filepath.Base(path)
	scr, err := ash.NewScript(scriptName, fmt.Sprintf("%s", text), args...)
	if err != nil {
		log.Err.Printf("fail to parse script: %v", err)
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
		log.Err.Printf("script run fail: %v", err)
		return
	}
}

// gameLoop handles game updating.
func gameLoop(g *core.Game) {
	// Delta.
	dtNano := time.Since(lastUpdate).Nanoseconds()
	delta := dtNano / int64(time.Millisecond) // delta to milliseconds
	// Game.
	g.Update(delta)
	// Update time.
	lastUpdate = time.Now()
}

// loadModule loads module with all module data
// from directory with specified path.
func loadModule(path, langID string) error {
	m, err := data.ImportModule(flameconf.ModulePath(), flameconf.LangID())
	if err != nil {
		return fmt.Errorf("fail to dir: %v", err)
	}
	// Load module data.
	err = data.LoadModuleData(m)
	if err != nil {
		return fmt.Errorf("fail to load data: %v", err)
	}
	flame.SetModule(m)
	burn.Module = m
	return nil
}
