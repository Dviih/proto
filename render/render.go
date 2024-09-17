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

package render

import (
	"bytes"
	"embed"
	"github.com/Dviih/proto"
	"github.com/Dviih/proto/event"
	"html/template"
	"sync"
	"syscall/js"
)

type Render struct {
	root     *Element
	template *template.Template

	m    sync.Mutex
	data map[string]interface{}

	events sync.Map
	c      chan bool

}

func (render *Render) Add(key string, value interface{}) {
	defer render.m.Unlock()
	render.m.Lock()

	render.data[key] = value
}

func (render *Render) Remove(key string) {
	defer render.m.Unlock()
	render.m.Lock()

	delete(render.data, key)
}

func (render *Render) Element(id string) *Element {
	defer render.m.Unlock()
	render.m.Lock()

	value := proto.Wait(func() js.Value {
		return proto.Document().Call("getElementById", id)
	})

	return &Element{
		m:     sync.Mutex{},
		value: value,
	}
}

func (render *Render) Root() *Element {
	if render.root == nil {
		render.root = render.Element("root")
	}

	return render.root
}

func (render *Render) Event(id string) *event.Event {
	e := event.New(id, render.c)

	render.events.Store(id, e)
	return e
}

func (render *Render) Execute(name string) error {
	return render.template.ExecuteTemplate(render.Root(), name, render.data)
}

func (render *Render) Create(name string) (*Element, error) {
	element := &Element{
		m:     sync.Mutex{},
		value: proto.Document().Call("createElement", "create"),
	}

	if err := render.template.ExecuteTemplate(element, name, render.data); err != nil {
		return nil, err
	}

	element.value = element.Value().Get("children").Index(0)
	return element, nil
}

func (render *Render) hook() {
	for {
		select {
		case <-render.c:
			render.events.Range(func(_, e interface{}) bool {
				e.(*event.Event).Run()
				return true
			})
		}
	}
}

func New(fs embed.FS, patterns ...string) (*Render, error) {
	t, err := template.ParseFS(fs, patterns...)
	if err != nil {
		return nil, err
	}

	render := &Render{
		template: t,
		m:        sync.Mutex{},
		data:     make(map[string]interface{}),
		events:   sync.Map{},
		c:        make(chan bool),
	}

	go render.hook()
	return render, nil
}
