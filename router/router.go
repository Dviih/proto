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
	"github.com/Dviih/Map"
)

type Handler func(*Page) error

type Router struct {
	pages    *Map.Map[string, Handler]
	_default string
}

func (router *Router) Add(route string, handler Handler) {
	router.pages.Store(route, handler)
}

func (router *Router) Remove(route string) {
	router.pages.Delete(route)
}

func (router *Router) Get(route string) (Handler, error) {
	handler := router.match(route)
	if handler != nil {
		return handler, nil
	}

	return nil, Map.KeyNotFound
}

func (router *Router) match(name string) Handler {
	var ret Handler

	router.pages.Range(func(s string, handler Handler) bool {
		route := split(s, '/')
		ns := split(name, '/')

		if len(ns) > len(route) {
			return false
		}

		for i, r := range route {
			if len(ns) >= i && len(ns[i]) == 0 {
				ret = nil
				return false
			}

			if len(r) == 0 || r[0] == ':' {
				if len(ns) < i+1 {
					ret = nil
				}
				continue
			}

			if len(ns) < i {
				return true
			}

			if ns[i] == r {
				ret = handler
			} else {
				return true
			}
		}

		return true
	})

	return ret
}

func split(s string, b byte) []string {
	var ret []string
	j := 1

	for i := 1; i < len(s); i++ {
		if s[i] == b {
			ret = append(ret, s[j:i])
			j = i + 1
		}
	}

	if j-len(s) == 0 {
		return ret
	}

	return append(ret, s[j:])
}

