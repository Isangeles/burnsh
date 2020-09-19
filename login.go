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

	"github.com/isangeles/fire/response"

	"github.com/isangeles/burnsh/log"
)

var logged bool

// login start CLI dialog for game server login.
func loginDialog() error {
	if server == nil {
		return fmt.Errorf("No server connection")
	}
	scan := bufio.NewScanner(os.Stdin)
	for reqSend := false; !reqSend; {
		id := ""
		fmt.Printf("%s:", lang.Text("cli_login_id"))
		for scan.Scan() {
			id = scan.Text()
			if len(id) > 0 {
				break
			}
		}
		pass := ""
		fmt.Printf("%s:", lang.Text("cli_login_pass"))
		for scan.Scan() {
			pass = scan.Text()
			if len(pass) > 0 {
				break
			}
		}
		server.SetOnResponseFunc(handleResponse)
		err := server.Login(id, pass)
		if err != nil {
			return fmt.Errorf("Unable to send login request: %v",
				err)
		} else {
			break
		}
	}
	return nil
}

// handleResponse handles response from the Fire server.
func handleResponse(resp response.Response) {
	if !resp.Logon {
		log.Inf.Printf("Logged at: %s", server.Address())
		server.SetOnResponseFunc(nil)
	}
	for _, r := range resp.Error {
		log.Err.Printf("Server error response: %s", r)
	}
}
