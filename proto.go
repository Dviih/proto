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
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"syscall/js"
	"unsafe"
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

func ToValue(v interface{}) js.Value {
	value := reflect.ValueOf(v)

	for value.Kind() == reflect.Pointer {
		if value.IsZero() {
			return js.Undefined()
		}

		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Invalid:
		return js.Null()
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String, reflect.Uintptr, reflect.UnsafePointer:
		return js.ValueOf(value.Interface())
	case reflect.Complex64, reflect.Complex128:
		c := value.Complex()

		return js.ValueOf(map[string]interface{}{
			"real": ToValue(real(c)),
			"imag": ToValue(imag(c)),
		})
	case reflect.Array, reflect.Slice:
		tmp := reflect.MakeSlice(reflect.SliceOf(reflect.TypeFor[interface{}]()), value.Len(), value.Cap())

		for i := 0; i < value.Len(); i++ {
			tmp.Index(i).Set(value.Index(i))
		}

		return js.ValueOf(tmp.Interface())
	case reflect.Chan:
		// Channel cannot block since JS will make the CPU go to 100%.
		// The operations for channels must be `TrySend` and `TryRecv`
		// since they are non-blocking operations.
		return js.ValueOf(map[string]interface{}{
			"send": js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
				if len(args) < 1 {
					switch value.Type().Elem() {
					case reflect.TypeFor[interface{}](), reflect.TypeFor[struct{}]():
						value.TrySend(reflect.ValueOf(struct{}{}))
						return nil
					default:
						return GError.Value().Invoke("must send at least one")
					}
				}

				for _, arg := range args {
					value.TrySend(convert(value.Type().Elem(), reflect.ValueOf(ToInterface(NewEmptyValue(arg)))))
				}

				return nil
			}),
			"receive": js.FuncOf(func(js.Value, []js.Value) interface{} {
				x, ok := value.TryRecv()
				if !ok {
					return js.Null()
				}

				return ToValue(x.Interface())
			}),
			"close": js.FuncOf(func(js.Value, []js.Value) interface{} {
				value.Close()
				return nil
			}),
			"len": js.FuncOf(func(js.Value, []js.Value) interface{} {
				return value.Len()
			}),
			"cap": js.FuncOf(func(js.Value, []js.Value) interface{} {
				return value.Cap()
			}),
			"type": js.FuncOf(func(js.Value, []js.Value) any {
				return nil
			}),
			"p": value.Pointer(),
		})
	case reflect.Func:
		return BuildJSFunc(value).Value
	case reflect.Interface:
		if value.Type().Implements(reflect.TypeFor[error]()) {
			return js.ValueOf(value.Call([]reflect.Value{reflect.ValueOf("Error")})[0].String())
		}

		panic("proto.ToValue: interface")
	case reflect.Map:
		tmp := reflect.MakeMapWithSize(reflect.MapOf(reflect.TypeFor[string](), reflect.TypeFor[interface{}]()), value.Len())

		m := value.MapRange()
		for m.Next() {
			tmp.SetMapIndex(reflect.ValueOf(fmt.Sprintf("%v", m.Key().Interface())), m.Value())
		}

		return js.ValueOf(tmp.Interface())
	case reflect.Pointer:
		panic("proto.ToValue: pointer")
	case reflect.Struct:
		m := map[string]interface{}{}

		vt := value.Type()

		for i := 0; i < value.NumField(); i++ {
			ft := vt.Field(i)

			if ft.PkgPath != "" {
				continue
			}

			m[ft.Name] = value.Field(i).Interface()
		}

		for i := 0; i < value.NumMethod(); i++ {
			mt := vt.Method(i)

			if mt.PkgPath != "" {
				continue
			}

			m[mt.Name] = BuildJSFunc(value.Method(i))
		}

		pv := reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr()))
		pv.Elem().Set(value)

		for i := 0; i < pv.NumMethod(); i++ {
			pvt := pv.Type().Method(i)

			if pvt.PkgPath != "" {
				continue
			}

			m[pvt.Name] = BuildJSFunc(pv.Method(i))
		}

		return js.ValueOf(m)
	}

	return js.Value{}
}

// BuildJSFunc v is expected to be an array of functions,
// these functions can return errors and any value in a
// parameter of a future function.
func BuildJSFunc(v ...interface{}) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		for _, fn := range v {
			rv := reflect.ValueOf(fn)

			if rv.Kind() != reflect.Func || rv.IsNil() {
				continue
			}

			rt := rv.Type()

			delta := 0

			switch rt {
			case reflect.TypeFor[func(js.Value)]():
				rv.Call([]reflect.Value{reflect.ValueOf(args[delta])})
			case reflect.TypeFor[func(js.Value) interface{}]():
				rv.Call([]reflect.Value{reflect.ValueOf(args[delta])})
			case reflect.TypeFor[func(js.Value, []js.Value)]():
				in := []reflect.Value{reflect.ValueOf(this)}

				for _, arg := range args {
					in = append(in, reflect.ValueOf(arg))
				}

				rv.Call(in)
			case reflect.TypeFor[func(js.Value, []js.Value) interface{}]():
				in := []reflect.Value{reflect.ValueOf(this)}

				for _, arg := range args {
					in = append(in, reflect.ValueOf(arg))
				}

				rv.Call(in)
			default:

			}

			rt.NumIn()
		}

		return nil
	})
}

func checkError(out []reflect.Value) ([]reflect.Value, error) {
	var (
		ret []reflect.Value
		err error
	)

	for _, o := range out {
		if !o.Type().Implements(reflect.TypeFor[error]()) {
			ret = append(ret, o)
			continue
		}

		if o.IsNil() {
			continue
		}

		err = errors.Join(err, o.Interface().(error))
	}

	return ret, err
}

func convert(t reflect.Type, value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	switch t.Kind() {
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Struct:
		panic("no conversion")
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String, reflect.Pointer, reflect.UnsafePointer:
		if value.CanConvert(t) {
			return value.Convert(t)
		}

		return reflect.Value{}
	case reflect.Array:
		if value.Type().Len() != t.Len() {
			// Size must be equal
			return reflect.Value{}
		}

		if value.Type().Elem().Kind() == reflect.Interface {
			array := reflect.New(t).Elem()

			for i := 0; i < value.Len(); i++ {
				array.Index(i).Set(convert(t.Elem(), value.Index(i)))
			}

			return array
		}

		return reflect.Value{}
	case reflect.Interface:
		for value.Kind() == reflect.Interface {
			value.Elem()
		}

		return value
	case reflect.Map:
		kt := t.Key()
		vt := t.Elem()

		m := reflect.MakeMapWithSize(t, value.Len())

		r := value.MapRange()
		for r.Next() {
			rk := convert(kt, r.Key())
			if !rk.IsValid() {
				continue
			}

			rv := convert(vt, r.Value())
			if !rv.IsValid() {
				continue
			}

			m.SetMapIndex(rk, rv)
		}

		return m
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Interface {
			slice := reflect.MakeSlice(t, value.Len(), value.Cap())

			for i := 0; i < slice.Len(); i++ {
				slice.Index(i).Set(convert(t.Elem(), value.Index(i)))
			}

			return slice
		}

		return reflect.Value{}
	}

	return value
}

func constructor(t reflect.Type) js.Func {
	n := t.NumField()
	var fields []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.PkgPath != "" {
			n--
			continue
		}

		fields = append(fields, field.Name)
	}

	return js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		if len(args) != n {
			return GError.Value().Invoke("invalid args length")
		}

		ptr := reflect.New(t)

		for i, field := range fields {
			ptr.Elem().FieldByName(field).Set(reflect.ValueOf(ToInterface(NewEmptyValue(args[i]))))
		}

		return ToValue(ptr.Interface())
	})
}

func ToRType(t reflect.Type, value Value) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Invalid:
		return reflect.ValueOf(nil), UnsupportedType
	case reflect.Bool:
		return reflect.ValueOf(value.Value().Bool()), nil
	case reflect.Int:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(i), nil
	case reflect.Int8:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		if i < -128 || i > 127 {
			return reflect.ValueOf(nil), OutOfRange
		}

		return reflect.ValueOf(int8(i)), nil
	case reflect.Int16:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		if i < -32768 || i > 32767 {
			return reflect.ValueOf(nil), OutOfRange
		}

		return reflect.ValueOf(int16(i)), nil
	case reflect.Int32:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		if i < -2147483648 || i > 2147483647 {
			return reflect.ValueOf(nil), OutOfRange
		}

		return reflect.ValueOf(int32(i)), nil
	case reflect.Int64:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(int64(i)), nil
	case reflect.Uint:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(uint(i)), nil
	case reflect.Uint8:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		u := uint(i)

		if u > 255 {
			return reflect.ValueOf(nil), OutOfRange
		}

		return reflect.ValueOf(uint8(u)), nil
	case reflect.Uint16:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		u := uint(i)

		if u > 65535 {
			return reflect.ValueOf(nil), OutOfRange
		}

		return reflect.ValueOf(uint16(u)), nil
	case reflect.Uint32:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		u := uint(i)

		if u > 4294967295 {
			return reflect.ValueOf(nil), OutOfRange
		}

		return reflect.ValueOf(uint32(u)), nil
	case reflect.Uint64:
		i, err := ToInt(value)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(uint64(i)), nil
	case reflect.Uintptr:
		return reflect.ValueOf(nil), UnsupportedType
	case reflect.Float32:
		f := value.Value().Float()

		if f > math.MaxFloat32 {
			return reflect.ValueOf(nil), OutOfRange
		}

		return reflect.ValueOf(float32(f)), nil
	case reflect.Float64:
		return reflect.ValueOf(value.Value().Float()), nil
	case reflect.Complex64:
		if value.Value().Type() != js.TypeObject {
			return reflect.ValueOf(nil), UnsupportedType
		}

		r := value.Value().Get("real").Float()

		if r > math.MaxFloat32 {
			return reflect.ValueOf(nil), OutOfRange
		}

		i := value.Value().Get("imag").Float()
		if i > math.MaxFloat32 {
			return reflect.ValueOf(nil), OutOfRange
		}

		return reflect.ValueOf(complex(float32(r), float32(i))), nil
	case reflect.Complex128:
		if value.Value().Type() != js.TypeObject {
			return reflect.ValueOf(nil), UnsupportedType
		}

		r := value.Value().Get("real").Float()
		i := value.Value().Get("imag").Float()

		return reflect.ValueOf(complex(r, i)), nil
	case reflect.Array:
		array := reflect.New(t)

		for i := 0; i < array.Len(); i++ {
			v, err := ToRType(t.Elem(), NewEmptyValue(value.Value().Index(i)))
			if err != nil {
				return reflect.ValueOf(nil), err
			}

			array.Index(i).Set(v)
		}

		return array, nil
	case reflect.Chan:
		if value.Value().Type() != js.TypeObject {
			return reflect.ValueOf(nil), UnsupportedType
		}

		return reflect.ValueOf(nil), UnsupportedType
	case reflect.Func:
		return reflect.ValueOf(nil), errors.New("to be func")
	case reflect.Interface:
		return reflect.ValueOf(ToInterface(value)), nil
	case reflect.Map:
		m := reflect.MakeMapWithSize(t, value.Value().Length())

		keys := GObject.Value().Call("keys", value)
		values := GObject.Value().Call("values", value)

		for i := 0; i < keys.Length(); i++ {
			key, err := ToRType(t.Key(), NewEmptyValue(keys.Index(i)))
			if err != nil {
				return reflect.ValueOf(nil), err
			}

			value, err := ToRType(t.Elem(), NewEmptyValue(values.Index(i)))
			if err != nil {
				return reflect.ValueOf(nil), err
			}

			m.SetMapIndex(key, value)
		}

		return m, nil
	case reflect.Pointer:
		return ToRType(t.Elem(), value)
	case reflect.Slice:
		slice := reflect.MakeSlice(t, value.Value().Length(), value.Value().Length())

		for i := 0; i < slice.Len(); i++ {
			v, err := ToRType(t.Elem(), NewEmptyValue(value.Value().Index(i)))
			if err != nil {
				return reflect.ValueOf(nil), err
			}

			slice.Index(i).Set(v)
		}

		return slice, nil
	case reflect.String:
		return reflect.ValueOf(value.Value().String()), nil
	case reflect.Struct:
		return reflect.ValueOf(nil), errors.New("to be struct")
	case reflect.UnsafePointer:
		return reflect.ValueOf(nil), UnsupportedType
	}
}

func ToInt(value Value) (int, error) {
	switch value.Value().Type() {
	case js.TypeString:
		return strconv.Atoi(value.Value().String())
	case js.TypeNumber:
		return value.Value().Int(), nil
	case js.TypeUndefined, js.TypeNull:
		return 0, nil
	case js.TypeBoolean:
		if value.Value().Bool() {
			return 1, nil
		}

		return 0, nil
	default:
		return 0, UnsupportedType
	}
}

}
