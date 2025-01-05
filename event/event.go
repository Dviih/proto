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

package event

import (
	"context"
	"github.com/Dviih/Map"
	"github.com/Dviih/proto"
	"sync/atomic"
	"syscall/js"
)

type Event struct {
	ctx     context.Context
	value   proto.Value
	events  *Map.Map[string, js.Func]
	running atomic.Bool
}

var (
	Global = Attached(context.Background(), &global{})
)

func (event *Event) Id() string {
	return event.value.Name()
}

func (event *Event) Running() bool {
	return event.running.Load()
}

func (event *Event) Subscribe(name string, fn func(js.Value, []js.Value) interface{}) {
	_, err := event.events.LoadOrStore(name, js.FuncOf(fn))
	if err != nil {
		event.Unsubscribe(name)
	}

	event.value.Value().Call("addEventListener", name, fn)
}

func (event *Event) Unsubscribe(name string) {
	event.value.Value().Call("removeEventListener", name, nil)

	fn, err := event.events.LoadAndDelete(name)
	if err != nil {
		return
	}

	fn.Release()
}

	}

	return event.value
}

	event.running.Store(true)
	return event
}

func New(id string, c chan bool) *Event {
	return &Event{
		id:         id,
		conditions: sync.Map{},
		events:     sync.Map{},
		running:    atomic.Bool{},
		c:          c,
	}
}

func Attached(value js.Value) *Event {
	event := &Event{
		value:    value,
		events:   sync.Map{},
		running:  atomic.Bool{},
		attached: true,
	}

	event.running.Store(true)
	return event
}
