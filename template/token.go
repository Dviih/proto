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

package template

import (
	"fmt"
	"github.com/Dviih/Map"
	"reflect"
	"unicode"
)

type Token struct {
	data                  []byte
	i, j                  int
	found, noappend, sign bool

	v interface{}

	ret []byte
}

// Next loops and returns the next $.
func (token *Token) Next() string {
	for ; token.i < len(token.data); token.i++ {
		switch token.data[token.i] {
		case ' ', '\n', '<', '>':
			if !token.found {
				if !token.noappend {
					token.add(token.data[token.i])
				}
				continue
			}

			token.sign = false
			token.found = false
			return string(token.data[token.j:token.i])
		case '$':
			if token.sign {
				token.sign = false
				token.found = false
				if !token.noappend {
					token.add('$')
				}
				continue
			}

			token.sign = true
			token.found = true
			token.j = token.i
			continue
		default:
			if !token.found && !token.noappend {
				token.add(token.data[token.i])
			}
		}
	}

	return ""
}

// Info returns information associated with Next after its name.
func (token *Token) Info() string {
	k := token.i
	b := false

	for ; token.i < len(token.data); token.i++ {
		switch token.data[token.i] {
		case '\n', '<', '>':
			return string(token.data[k:token.i])
		default:
			if !b {
				k++
			}

			b = true
		}
	}

	return string(TrimSingle(token.data[k:]))
}

// End gets everything until finds `$end` tag.
func (token *Token) End() []byte {
	k := token.i
	token.noappend = true

	for {
		node := token.Next()

		switch node {
		case "":
			token.noappend = false
			return nil
		case "$end":
			token.noappend = false
			return TrimSingle(token.data[k:token.i])
		}
	}
}

func (token *Token) add(v interface{}) {
	switch v := v.(type) {
	case []byte:
		token.ret = append(token.ret, TrimSingle(v)...)
	case string:
		token.ret = append(token.ret, TrimSingle(v)...)
	case byte:
		token.ret = append(token.ret, v)
	case rune:
		token.ret = append(token.ret, string(v)...)
	case nil:
		return
	default:
		token.ret = append(token.ret, TrimSingle(fmt.Sprintf("%v", v))...)
	}
}

func (token *Token) get(name string) interface{} {
	switch data := token.v.(type) {
	case *Map.Map[string, interface{}]:
		v, err := data.Load(name)
		if err != nil {
			return nil
		}

		return v
	default:
		value := reflect.ValueOf(data)
		for value.Kind() == reflect.Pointer {
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.Struct:
			field := value.FieldByName(name)
			if field.Kind() == reflect.Invalid || field.IsZero() {
				return nil
			}

			return field.Interface()
		default:
			return nil
		}
	}
}

