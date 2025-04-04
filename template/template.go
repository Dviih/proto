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
	"github.com/Dviih/proto"
	"github.com/Dviih/proto/state"
	"html/template"
	"io/fs"
	"reflect"
	"sync/atomic"
)

type Template struct {
	*template.Template
	states proto.Store
}

var funcMap = map[string]interface{}{
	"state": func(s string) string {
		return "<state state=\"" + s + "\"></state>"
	},
}

func (template *Template) StoreState(id string, v interface{}) {
	template.states.Set(id, v)
}

func (template *Template) NewState(_ reflect.Type, id string) state.Virtual {
	return &State{
		id:       id,
		template: template,
		current:  atomic.Pointer[string]{},
		store:    &proto.MapStore{},
		c:        make(chan struct{}),
	}
}

func (template *Template) LoadState(id string) interface{} {
	return template.states.Get(id)
}

func ParseFS(fs fs.FS, patterns ...string) (*Template, error) {
	t := New("")

	var err error

	t.Template, err = t.Template.ParseFS(fs, patterns...)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func New(name string) *Template {
	return &Template{
		Template: template.New(name).Funcs(funcMap),
		states:   &proto.MapStore{},
	}
}
