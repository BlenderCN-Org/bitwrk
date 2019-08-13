//  BitWrk - A Bitcoin-friendly, anonymous marketplace for computing power
//  Copyright (C) 2013-2019  Jonas Eschenburg <jonas@bitwrk.net>
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrNotADirectory        = errors.New("Not a directory")
	ErrNotFoundInSearchPath = errors.New("Resource directory not found")
)

type infoStruct struct {
	Info    string `json:"info"`
	Version string `json:"version"`
}

func TestResourceDir(dir, name, version string) error {
	infoPath := filepath.Join(dir, "info.json")
	var infoFile *os.File
	if fi, err := os.Stat(dir); err != nil {
		return err
	} else if !fi.IsDir() {
		return ErrNotADirectory
	} else if f, err := os.Open(infoPath); err != nil {
		return err
	} else {
		infoFile = f
	}

	decoder := json.NewDecoder(infoFile)
	info := infoStruct{}
	if err := decoder.Decode(&info); err != nil {
		return err
	}

	if info.Version != version {
		return errors.New("Wrong version: " + info.Version)
	}

	if info.Info != name+" resource files" {
		return errors.New("Not a resource directory")
	}

	return nil
}

// Function AutoFindResourceDir checks a number of candidate directories for whether an info.json file
// exists in them and whether that file can be parsed (using TestResourceDir) as a resource directory
// info file. In case of success, it returns the first successful candidate directory. Otherwise, it
// returns an error.
func AutoFindResourceDir(name, version string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	cmdDir := filepath.Dir(execPath)
	candidates := []string{
		filepath.Join(cmdDir, "share/", name),
		filepath.Join(cmdDir, "../share/", name),
		filepath.Join(cmdDir, "rsc/"),
		filepath.Join(cmdDir, "resources/"),
		cmdDir,
	}
	errs := make([]error, len(candidates))

	for i, dir := range candidates {
		if err := TestResourceDir(dir, name, version); err != nil {
			errs[i] = err
		} else {
			return dir, nil
		}
	}

	// Log a custom reason that explains why each candidate failed the check
	for i, dir := range candidates {
		log.Printf("No resource directory at location [%v]: %v", dir, errs[i])
	}
	return "", ErrNotFoundInSearchPath
}
