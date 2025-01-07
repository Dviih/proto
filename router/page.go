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
	"fmt"
	"github.com/Dviih/Map"
	"github.com/Dviih/proto"
	"github.com/Dviih/proto/state"
	"log/slog"
	"reflect"
)

type Page struct {
	ctx      context.Context
	logger   *slog.Logger
	template string
	data     *Map.Map[string, interface{}]

	states *Map.Map[string, interface{}]
	c      chan string

	Arguments []string
	Query     map[string][]string
}

func (page *Page) Ctx() context.Context {
	return page.ctx
}

func (page *Page) Logger() *slog.Logger {
	return page.logger
}

func (page *Page) C() chan string {
	return page.c
}

func (page *Page) SetTemplate(template string) {
	page.template = template
}

func (page *Page) Set(key string, value interface{}) {
	page.data.Store(key, value)
}

func (page *Page) Get(key string) interface{} {
	v, err := page.data.Load(key)
	if err != nil {
		return nil
	}

	return v
}

func (page *Page) Join(m map[string]interface{}) {
	for k, v := range m {
		page.data.Store(k, v)
	}
}

