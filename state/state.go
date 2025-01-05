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

import (
	"context"
	"sync/atomic"
)

type State[T interface{}] struct {
	ctx context.Context
	c   chan string

	id string
	m  atomic.Pointer[T]
}

func (state *State[T]) Get() T {
	if state.ctx.Err() != nil {
		var zero T
		return zero
	}

	return *state.m.Load()
}

func (state *State[T]) Set(t T) {
	if state.ctx.Err() != nil {
		return
	}

	state.m.Store(&t)
	state.c <- state.id
}
