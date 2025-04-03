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

package proto

import (
	"syscall/js"
	"time"
)

var AsyncDeadline = 3 * time.Second

func Promise(fn func(resolve, reject func(...interface{}) js.Value) js.Value) js.Value {
	var h js.Func
	h = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer h.Release()
		return fn(args[0].Invoke, args[1].Invoke)
	})

	return GPromise.Value().New(h)
}

func Await(promise js.Value) []js.Value {
	c := make(chan []js.Value, 1)

	var h js.Func

	go func() {
		h = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			defer h.Release()

			c <- args
			return nil
		})

		promise.Call("then", h)
	}()

	select {
	case data := <-c:
		return data
	case <-time.After(AsyncDeadline):
		return []js.Value{GError.Value().Invoke("deadline")}
	}
}
