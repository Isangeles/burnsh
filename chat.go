/*
 * chat.go
 *
 * Copyright 2020-2023 Dariusz Sikora <ds@isangeles.dev>
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
	"sort"
	"time"

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/objects"
)

// Struct for log message.
type Message struct {
	author string
	time   time.Time
	text   string
}

// Struct for sorting messages by the messsage time.
type MessagesByTime []Message

func (mbt MessagesByTime) Len() int           { return len(mbt) }
func (mbt MessagesByTime) Swap(i, j int)      { mbt[i], mbt[j] = mbt[j], mbt[i] }
func (mbt MessagesByTime) Less(i, j int) bool { return mbt[i].time.UnixNano() < mbt[j].time.UnixNano() }

var chatOpen bool

// chatDialog starts chat CLI dialog.
func chatDialog() error {
	if activeGame == nil {
		msg := lang.Text("no_game_err")
		return fmt.Errorf(msg)
	}
	chatOpen = true
	go updateChat()
	scan := bufio.NewScanner(os.Stdin)
	for chatOpen {
		fmt.Printf("%s\n", lang.Text("chat_info"))
		scan.Scan()
		input := scan.Text()
		if input == "$close" {
			chatOpen = false
		}
	}
	return nil
}

// updateChat prints messages from players, nearby objects, and system
// log on the standard out.
func updateChat() {
	for chatOpen {
		// Add messages from players and nearby objects.
		messages := make([]Message, 0)
		for _, pc := range activeGame.Players() {
			// PC's private messages.
			for _, lm := range pc.Log().Messages() {
				m := Message{
					author: pc.ID(),
					time:   lm.Time,
					text:   fmt.Sprintf("%s\n", lm.String()),
				}
				if !lm.Translated {
					m.text = fmt.Sprintf("%s\n", lang.Text(lm.String()))
				}
				messages = append(messages, m)
			}
			// Near objects chat.
			area := activeGame.Chapter().ObjectArea(pc.Character)
			if area == nil {
				continue
			}
			charX, charY := pc.Position()
			for _, tar := range area.NearObjects(charX, charY, pc.SightRange()) {
				tar, ok := tar.(objects.Logger)
				if !ok {
					continue
				}
				for _, lm := range tar.ChatLog().Messages() {
					m := Message{
						author: tar.ID(),
						time:   lm.Time,
						text:   fmt.Sprintf("%s\n", lm.String()),
					}
					if !lm.Translated {
						m.text = fmt.Sprintf("%s\n", lang.Text(lm.String()))
					}
					messages = append(messages, m)
				}
			}
		}
		// Sort and print messages.
		sort.Sort(MessagesByTime(messages))
		for _, m := range messages {
			// Skip old messages.
			if m.time.UnixNano() < lastUpdate.UnixNano() {
				continue
			}
			fmt.Printf("%s: %s", lang.Text(m.author), m.text)
		}
	}
}
