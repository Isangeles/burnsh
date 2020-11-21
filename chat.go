/*
 * chat.go
 *
 * Copyright 2020 Dariusz Sikora <dev@isangeles.pl>
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
	"sort"
	"time"

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module/objects"
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

// updateChat prints messages from players, nearby objects, and system
// log on the standard out.
func updateChat() {
	// Add messages from players and nearby objects.
	messages := make([]Message, 0)
	for _, pc := range activeGame.Players() {
		// PC's private messages.
		for _, m := range pc.PrivateLog().Messages() {
			m := Message{
				author: pc.ID(),
				time:   m.Time(),
				text:   fmt.Sprintf("%s\n", m.String()),
			}
			messages = append(messages, m)
		}
		// Near objects chat & combat.
		area := activeGame.Module().Chapter().CharacterArea(pc.Character)
		if area == nil {
			continue
		}
		for _, tar := range area.NearTargets(pc.Character, pc.SightRange()) {
			tar, ok := tar.(objects.Logger)
			if !ok {
				continue
			}
			for _, m := range tar.ChatLog().Messages() {
				m := Message{
					author: tar.ID(),
					time:   m.Time(),
					text:   fmt.Sprintf("%s\n", m.String()),
				}
				messages = append(messages, m)
			}
			for _, m := range tar.CombatLog().Messages() {
				m := Message{
					author: tar.ID(),
					time:   m.Time(),
					text:   fmt.Sprintf("%s\n", m.String()),
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
