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
	"github.com/Dviih/sync"
	"log/slog"
	"reflect"
	"sync/atomic"
)

type Page struct {
	ctx      context.Context
	logger   *slog.Logger
	template atomic.Pointer[string]
	events   sync.Slice[*event.Event]
	post     sync.Slice[func()]

	states proto.Store
	c      chan string
	close  chan bool

	Store  proto.Store
	Router *Router

	Arguments []string
	Query     map[string][]string
}

func (page *Page) Logger() *slog.Logger {
	return page.logger
}

func (page *Page) SetTemplate(template string) {
	page.template.Store(&template)
}

func (page *Page) State(name string) interface{} {
	return page.states.Get(name)
}

func (page *Page) Event(value proto.Value) *event.Event {
	e := event.New(page.ctx, value)

	page.events.Append(e)
	return e
}

func (page *Page) Context() context.Context {
	return page.ctx
}

func (page *Page) Go(name string) {
	if h, _ := page.Router.match(name); h == nil {
		return
	}

	history.Default.Push(nil, name, &url.URL{Path: name})

	if err := page.Router.Handler(); err != nil {
		page.Logger().ErrorContext(page.Context(), "failed to go to other page", slog.Any("error", err))
	}
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

func (page *Page) Post(fn func()) {
	page.post.Append(fn)
}
