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

package router

import (
	"context"
	"fmt"
	"github.com/Dviih/Map"
	"github.com/Dviih/proto"
	"github.com/Dviih/proto/state"
	"log/slog"
	"reflect"
)

type Page struct {
	ctx      context.Context
	logger   *slog.Logger
	template string
	data     *Map.Map[string, interface{}]

	states *Map.Map[string, interface{}]
	c      chan string
	close  chan bool

	Arguments []string
	Query     map[string][]string
}

func (page *Page) Ctx() context.Context {
	return page.ctx
}

func (page *Page) Logger() *slog.Logger {
	return page.logger
}

func (page *Page) C() chan string {
	return page.c
}

func (page *Page) SetTemplate(template string) {
	page.template = template
}

func (page *Page) Set(key string, value interface{}) {
	page.data.Store(key, value)
}

func (page *Page) Get(key string) interface{} {
	v, err := page.data.Load(key)
	if err != nil {
		return nil
	}

	return v
}

func (page *Page) Join(m map[string]interface{}) {
	for k, v := range m {
		page.data.Store(k, v)
	}
}

func (page *Page) State(name string, v interface{}) {
	page.states.Store(name, v)
}

func (page *Page) handle() {
	for {
		select {
		case <-page.close:
			return
		case name := <-page.c:
			s, err := page.states.Load(name)
			if err != nil {
				continue
			}

			out := reflect.ValueOf(s).MethodByName("Get").Call(nil)
			if out == nil || len(out) != 1 {
				continue
			}

			v := out[0].Interface()

			selectors := proto.GDocument.Call("querySelectorAll", "[state='"+name+"']")

			for i := 0; i < selectors.Length(); i++ {
				selectors.Index(i).Call("replaceChildren", proto.Create(fmt.Sprintf("%v", v))...)
			}
		}
	}
}

func Get[T interface{}](page *Page, name string) *state.State[T] {
	i, err := page.states.Load(name)
	if err != nil {
		return nil
	}

	s, ok := i.(*state.State[T])
	if !ok {
		return nil
	}

	return s
}
