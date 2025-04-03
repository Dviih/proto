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

package template

import (
	"github.com/Dviih/bin/buffer"
	"github.com/Dviih/proto"
	"sync/atomic"
)

type State struct {
	id       string
	template *Template
	current  atomic.Pointer[string]
	store    proto.Store
	c        chan struct{}
}

func (state *State) Id() string {
	return state.id
}

func (state *State) Load() interface{} {
	v := state.current.Load()
	if v == nil {
		return ""
	}

	b := buffer.New()
	if err := state.template.ExecuteTemplate(b, *v, proto.StoreToData(state.store)); err != nil {
		return err.Error()
	}

	return string(b.Data())
}

func (state *State) C() <-chan struct{} {
	return state.c
}
