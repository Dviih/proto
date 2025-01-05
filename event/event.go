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
	if !event.Value().IsNull() {
		event.Value().Call("removeEventListener", name)
	}

	fn, ok := event.events.LoadAndDelete(name)
	if !ok {
		return
	}

	fn.(js.Func).Release()
}

func (event *Event) Running() bool {
	return event.running.Load()
}

func (event *Event) Value() js.Value {
	if event.value.IsUndefined() {
		event.forceValue()
	}

	return event.value
}

func (event *Event) forceValue() {
	event.value = proto.Document().Call("getElementById", event.id)
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
