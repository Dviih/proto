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
	"reflect"
	"sync/atomic"
)

type State struct {
	id  string
	ctx context.Context
	c   chan struct{}

	p reflect.Type
	m atomic.Value
}

func (state *State) Id() string {
	return state.id
}

func (state *State) Store(v interface{}) {
	if state.p != reflect.TypeOf(v) {
		return
	}

	state.m.Store(v)
	state.c <- struct{}{}
}

