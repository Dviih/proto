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

package state

import "reflect"

type State[T interface{}] interface {
	Id() string
	Store(T)
	Load() T
	C() <-chan struct{}
}

type Storer interface {
	StoreState(string, interface{})
	NewState(reflect.Type, string) Virtual
}

func Store[T interface{}](storer Storer, id string) State[T] {
	state := storer.NewState(reflect.TypeFor[T](), id)

	virtual := &virtual[T]{
		id:   state.Id,
		set:  state.Store,
		load: state.Load,
	}

	storer.StoreState(id, virtual)
	return virtual
}

}

func (state *State[T]) Set(t T) {
	if state.ctx.Err() != nil {
		return
	}

	state.m.Store(&t)
	state.c <- state.id
}
