/*
 *     Proto is a minimal tool for real time HTML rendering.
 *     Copyright (C) 2024  Dviih
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package proto

import (
	"syscall/js"
)

type Value interface {
	Name() string
	Value() js.Value
}

type namedValue struct {
	name  string
	value js.Value
}

func (value *namedValue) Name() string {
	return value.name
}

func (value *namedValue) Value() js.Value {
	return value.value
}

func NewNamedValue(name string, value js.Value) Value {
	return &namedValue{
		name:  name,
		value: value,
	}
}

// NewEmptyValue is used for iterations.
func NewEmptyValue(value js.Value) Value {
	return &namedValue{
		name:  "",
		value: value,
	}
}
