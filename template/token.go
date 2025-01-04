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
		fields := token.getFields(name)

		if len(fields) < 2 {
			v, err := data.Load(name)
			if err != nil {
				return nil
			}

			return v
		}

		v, err := data.Load(fields[0])
		if err != nil {
			return nil
		}

		return token.getStruct(v, fields)
	default:
		value := reflect.ValueOf(data)
		for value.Kind() == reflect.Pointer {
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.Struct:
			return token.getStruct(data, token.getFields(name))
		default:
		}

		return data
	}
}

func (token *Token) getFields(name string) []string {
	var fields []string

	j := 0

	for i := 0; i < len(name); i++ {
		if name[i] == '.' {
			fields = append(fields, name[j:i])
			j = i + 1
		}
	}

	return append(fields, name[j:])
}

func (token *Token) getStruct(v interface{}, fields []string) interface{} {
	value := reflect.ValueOf(v)

	for _, fname := range fields[1:] {
		for value.Kind() == reflect.Pointer {
			value = value.Elem()
		}

		value = value.FieldByName(fname)
		if value.Kind() != reflect.Struct {
			break
		}
	}

	if value.Kind() == reflect.Invalid || value.IsZero() {
		return nil
	}

	return value.Interface()
}

func NewToken(data []byte, v interface{}) *Token {
	return &Token{
		data: data,
		v:    v,
	}
}

func TrimSpaceLeft[T string | []byte](t T) T {
	data := string(t)

	for i, b := range data {
		if !unicode.IsSpace(b) {
			return T(data[i:])
		}
	}

	return t
}

func TrimSpaceRight[T string | []byte](t T) T {
	data := string(t)

	for i := len(data) - 1; i > 0; i-- {
		if !unicode.IsSpace(rune(data[i])) {
			return T(data[:i+1])
		}
	}

	return t
}

func TrimSingle[T string | []byte](t T) T {
	var ls, rs, ln, rn bool

	switch t[0] {
	case ' ':
		ls = true
	case '\n':
		ln = true
	}

	switch t[len(t)-1] {
	case ' ':
		rs = true
	case '\n':
		rn = true
	}

	t = TrimSpaceRight(TrimSpaceLeft(t))

	if ls {
		t = T(" " + string(t))
	}

	if rs {
		t = T(string(t) + " ")
	}

	if ln {
		t = T("\n" + string(t))
	}

	if rn {
		t = T(string(t) + "\n")
	}

	return t
}

func Trim(s string) (string, int) {
	i := len(s) - 1
	for ; i >= 0; i-- {
		switch s[i] {
		case '"', '\'', ' ':
			continue
		default:
			return s[:i+1], len(s) - i - 1
		}
	}

	return s[:i+1], len(s) - i - 1
}
