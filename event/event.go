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
	"errors"
	"github.com/Dviih/proto"
	"html"
	"sync"
	"sync/atomic"
	"syscall/js"
)

type Event struct {
	id    string
	value js.Value

	conditions sync.Map
	events     sync.Map

	running atomic.Bool
	c       chan bool

	attached bool
}

var isAttached = errors.New("event is attached")

func (event *Event) Match() bool {
	matched := true

	event.conditions.Range(func(condition, expected interface{}) bool {
		matched = !proto.Document().Call("querySelector", "["+condition.(string)+"="+expected.(string)+"]").IsNull()
		return matched
	})

	return matched
}

func (event *Event) Run() {
	if !event.attached {
		event.forceValue()

		if !event.Match() {
			event.running.Store(false)
			return
		}
	}

	event.events.Range(func(name, fn any) bool {
		event.Value().Call("addEventListener", name, fn)
		return true
	})

	event.running.Store(true)
}

func (event *Event) Add(attribute string, value js.Value) {
	if event.attached {
		panic(isAttached)
	}

	event.events.Store(attribute, value)
}

func (event *Event) Remove(attribute string) {
	if event.attached {
		panic(isAttached)
	}

	event.events.Delete(attribute)
}

func (event *Event) Condition(condition, expected string) {
	if event.attached {
		panic(isAttached)
	}

	event.conditions.Store(html.EscapeString(condition), html.EscapeString(expected))
}

func (event *Event) Subscribe(name string, fn func(js.Value, []js.Value) interface{}) {
	event.events.Store(name, js.FuncOf(fn))

	if event.c != nil {
		event.c <- true
	}

	if event.attached {
		event.Run()
	}
}

func (event *Event) Unsubscribe(name string) {
	event.events.Delete(name)
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

