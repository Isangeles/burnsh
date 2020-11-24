/*
 * data.go
 *
 * Copyright 2020 Dariusz Sikora <dev@isangeles.pl>
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

package data

import (
	"fmt"
	"path/filepath"

	flamedata "github.com/isangeles/flame/data"
	flameres "github.com/isangeles/flame/data/res"

	"github.com/isangeles/burnsh/data/res"
)

const (
	UIDirPath = "burnsh"
)

// LoadUIData loads UI data directory with specified path.
func LoadUIData(path string) error {
	langPath := filepath.Join(path, "lang")
	lang, err := flamedata.ImportLangDirs(langPath)
	if err != nil {
		return fmt.Errorf("Unable to load translations: %v", err)
	}
	res.TranslationBases = lang
	flameres.Add(flameres.ResourcesData{TranslationBases: lang})
	return nil
}
