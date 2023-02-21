// gorc project
// Copyright (C) 2023 IllusionMan1212
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see https://www.gnu.org/licenses.

package login

import (
	"errors"
	"strconv"
	"strings"
)

func NoSpacesValidation(input string) error {
	if strings.Contains(input, " ") {
		return errors.New("Field cannot contain spaces")
	}

	return nil
}

func ValidatePort(input string) error {
	pNum, err := strconv.ParseUint(input, 10, 16)
	if err != nil {
		return err
	}

	if pNum < 1 || pNum > 65535 {
		return errors.New("Invalid port number")
	}

	return nil
}
