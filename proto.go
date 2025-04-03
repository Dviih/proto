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
	"errors"
	"net/url"
	"syscall/js"
)

type Func func(...interface{}) interface{}

var (
	GGlobal = NewNamedValue("global", js.Global())

	GDocument    = NewNamedValue("document", GGlobal.Value().Get("document"))
	GElement     = NewNamedValue("element", GGlobal.Value().Get("Element"))
	GWindow      = NewNamedValue("window", GGlobal.Value().Get("window"))
	GUint8Array  = NewNamedValue("Uint8Array", GGlobal.Value().Get("Uint8Array"))
	GArrayBuffer = NewNamedValue("ArrayBuffer", GGlobal.Value().Get("ArrayBuffer"))
	GObject      = NewNamedValue("Object", GGlobal.Value().Get("Object"))
	GError       = NewNamedValue("Error", GGlobal.Value().Get("Error"))
	GArray       = NewNamedValue("Array", GGlobal.Value().Get("Array"))
	GPromise     = NewNamedValue("Promise", GGlobal.Value().Get("Promise"))


	UnsupportedType = errors.New("unsupported type")
	OutOfRange      = errors.New("number out of range")
)

func URL() *url.URL {
	u, err := url.Parse(GDocument.Value().Get("URL").String())
	if err != nil {
		panic(err)
	}

	return u
}

func Create[T []byte | string](data T) []interface{} {
	create := GDocument.Value().Call("createElement", "create")
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
