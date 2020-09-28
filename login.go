/*
 * login.go
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
	"bufio"
	"fmt"
	"os"

	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/burnsh/config"
)

var logged bool

// login start CLI dialog for game server login.
func loginDialog() error {
	if server == nil {
		return fmt.Errorf("No server connection")
	}
	id, pass := config.ServerLogin, config.ServerPass
	if len(id) < 1 || len(pass) < 1 {
		scan := bufio.NewScanner(os.Stdin)
		fmt.Printf("%s:", lang.Text("cli_login_id"))
		for scan.Scan() {
			id = scan.Text()
			if len(id) > 0 {
				break
			}
		}
		fmt.Printf("%s:", lang.Text("cli_login_pass"))
		for scan.Scan() {
			pass = scan.Text()
			if len(pass) > 0 {
				break
			}
		}
	}
	err := server.Login(id, pass)
	if err != nil {
		return fmt.Errorf("Unable to send login request: %v",
			err)
	}
	return nil
}
