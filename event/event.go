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

