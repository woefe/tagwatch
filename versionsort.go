/*
 * Copyright (c) 2024. Wolfgang Popp
 *
 * This file is part of tagwatch.
 *
 * tagwatch is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * tagwatch is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tagwatch.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"strconv"
	"unicode"
)

type ByVersion []string

func sliceGroup(str string, digit bool) (string, string) {
	for i, c := range str {
		if digit != unicode.IsDigit(c) {
			return str[:i], str[i:]
		}
	}
	return str, ""
}

// Less Debian-like version sorting
func (v ByVersion) Less(i, j int) bool {
	lhs, rhs := v[i], v[j]
	for lhs != "" && rhs != "" {
		lAlpha, rAlpha, lDigit, rDigit := "", "", "", ""

		// Compare non-numeric part
		lAlpha, lhs = sliceGroup(lhs, false)
		rAlpha, rhs = sliceGroup(rhs, false)
		if lAlpha != rAlpha {
			return lAlpha < rAlpha
		}

		// Compare numeric part
		lDigit, lhs = sliceGroup(lhs, true)
		rDigit, rhs = sliceGroup(rhs, true)
		lInt, _ := strconv.ParseUint(lDigit, 10, 0)
		rInt, _ := strconv.ParseUint(rDigit, 10, 0)
		if lInt != rInt {
			return lInt < rInt
		}
	}
	return len(lhs) < len(rhs)
}

func (v ByVersion) Len() int {
	return len(v)
}

func (v ByVersion) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
