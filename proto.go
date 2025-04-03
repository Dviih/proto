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

func ToInterface(value Value) interface{} {
	switch value.Value().Type() {
	case js.TypeUndefined, js.TypeNull:
		return nil
	case js.TypeBoolean:
		return value.Value().Bool()
	case js.TypeNumber:
		return value.Value().Int()
	case js.TypeString:
		return value.Value().String()
	case js.TypeSymbol:
		return value.Value().Call("toString").String()
	case js.TypeObject:
		if value.Value().InstanceOf(GArray.Value()) {
			var m []interface{}

			for i := 0; i < value.Value().Length(); i++ {
				m = append(m, ToInterface(NewEmptyValue(value.Value().Index(i))))
			}

			return m
		}

		m := map[string]interface{}{}
		keys := GObject.Value().Call("keys", value)

		for i := 0; i < keys.Length(); i++ {
			key := keys.Index(i).String()
			m[key] = ToInterface(NewEmptyValue(value.Value().Get(key)))
		}

		return m
	case js.TypeFunction:
		fn := reflect.New(reflect.TypeFor[Func]()).Elem()

		fn.Set(reflect.ValueOf(Func(func(v ...interface{}) interface{} {
			var args []interface{}

			for _, v := range v {
				args = append(args, ToValue(v))
			}

			value.Value().Invoke(args...)
			return nil
		})))

		return fn.Interface()
	default:
		panic("invalid js.Constructor")
	}
}

}
