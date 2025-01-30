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
	"sync"
	"sync/atomic"
)

type Store interface {
	Set(string, interface{})
	Get(string) interface{}
	Range(func(string, interface{}) bool)
	Delete(string)
}

type MapStore struct {
	m sync.Map
}

func (store *MapStore) Set(key string, value interface{}) {
	store.m.Store(key, value)
}

func (store *MapStore) Get(key string) interface{} {
	v, ok := store.m.Load(key)
	if !ok {
		return nil
	}

	return v
}

func (store *MapStore) Range(fn func(string, interface{}) bool) {
	store.m.Range(func(key, value any) bool {
		s, ok := key.(string)
		if !ok {
			return false
		}

		return fn(s, value)
	})
}

func (store *MapStore) Delete(key string) {
	store.Delete(key)
}

type AnyStore struct {
	m atomic.Pointer[interface{}]
}

func (store *AnyStore) Set(_ string, v interface{}) {
	store.m.Store(&v)
}

func (store *AnyStore) Get(_ string) interface{} {
	v := store.m.Load()
	if v == nil {
		return nil
	}

	return *v
}

func (store *AnyStore) Range(fn func(string, interface{}) bool) {
	_ = fn("", store.Get(""))
}

func (store *AnyStore) Delete(_ string) {
	store.m.Store(nil)
}

func ToStore(v interface{}) Store {
	switch v := v.(type) {
	case *MapStore:
		return v
	case *AnyStore:
		return v
	case map[string]interface{}:
		return MapStoreFrom(v)
	default:
		store := &AnyStore{}

		store.Set("", v)
		return store
	}
}

