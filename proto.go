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
	"net/url"
	"syscall/js"
)

type Func func(...interface{}) interface{}

var (
	GGlobal = js.Global()

	GDocument    = GGlobal.Get("document")
	GWindow      = GGlobal.Get("window")
	GUint8Array  = GGlobal.Get("Uint8Array")
	GArrayBuffer = GGlobal.Get("ArrayBuffer")
	GObject      = GGlobal.Get("Object")
	GError       = GGlobal.Get("Error")
	GArray       = GGlobal.Get("Array")
	GPromise     = GGlobal.Get("Promise")

)

func URL() *url.URL {
	u, err := url.Parse(GDocument.Get("URL").String())
	if err != nil {
		panic(err)
	}

	return u
}

func Create[T []byte | string](data T) []interface{} {
	create := GDocument.Call("createElement", "create")
	create.Set("innerHTML", string(data))

	var children []interface{}

	tmp := create.Get("children")

	if tmp.Length() == 0 {
		return []interface{}{create}
	}

	for i := 0; i < tmp.Length(); i++ {
		children = append(children, tmp.Index(i))
	}

	return children
}

}
