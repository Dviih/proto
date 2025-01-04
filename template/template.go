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
	"errors"
	"github.com/Dviih/Map"
	"io"
	"io/fs"
	"reflect"
)

type Template struct {
	templates *Map.Map[string, []byte]
	data      *Map.Map[string, interface{}]
}

func (template *Template) Add(name string, data []byte) {
	template.templates.Store(name, data)
}

func (template *Template) Templates() []string {
	var templates []string

	template.templates.Range(func(template string, _ []byte) bool {
		templates = append(templates, template)
		return true
	})

	return templates
}

func (template *Template) Set(name string, v interface{}) {
	template.data.Store(name, v)
}

func (template *Template) Join(m map[string]interface{}) {
	for k, v := range m {
		template.data.Store(k, v)
	}
}

func (template *Template) Get(name string) interface{} {
	v, err := template.data.Load(name)
	if err != nil {
		return nil
	}

	return v
}

func (template *Template) Execute(name string) ([]byte, error) {
	data, err := template.templates.Load(name)
	if err != nil {
		return nil, err
	}

	return template.execute(data, template.data)
}

func (template *Template) execute(data []byte, v interface{}) ([]byte, error) {
	token := NewToken(data, v)

	for {
		node := token.Next()

		switch node {
		case "":
			return TrimSpaceRight(token.ret), nil
		case "$template":
			info := token.Info()

			t2, err := template.templates.Load(info)
			if err != nil {
				return nil, err
			}

			data, err := template.execute(t2, template.data)
			if err != nil {
				return nil, err
			}

			token.add(data)
		case "$range":
			info := token.Info()

			k := 0
			for ; k < len(info); k++ {
				if info[k] == ':' {
					info = info[:k] + info[k+1:]
					break
				}
			}

			if k == len(info) {
				k = 0
			}

			t := token.End()

			m := template.Get(info[k:])
			if m == nil {
				return nil, errors.New("nil")
			}

			value := reflect.ValueOf(m)

			switch value.Kind() {
			case reflect.Array, reflect.Slice:
				for i := 0; i < value.Len(); i++ {
					var v interface{}

					if k != 0 {
						v = Map.New[string, interface{}]()
						v.(*Map.Map[string, interface{}]).Store(info[:k], value.Index(i).Interface())
					} else {
						v = value.Index(i).Interface()
					}

					data, err := template.execute(t, v)
					if err != nil {
						return nil, err
					}

					token.add(data)
				}
			default:
				return nil, errors.New("invalid range")
			}
		default:
			s, i := Trim(node[1:])
			token.i -= i

			if s == "" {
				token.add(v)
				continue
			}

			token.add(token.get(s))
		}
	}
}

