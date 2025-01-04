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

