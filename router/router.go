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
	"errors"
	"github.com/Dviih/bin/buffer"
	"github.com/Dviih/proto"
	"github.com/Dviih/proto/pkg/js/history"
	"github.com/Dviih/proto/template"
	"github.com/Dviih/sync"
	"io"
	"log/slog"
)

type Handler func(*Page) error

type Router struct {
	ctx    context.Context
	logger *slog.Logger

	pages    *sync.Map[string, Handler]
	template *template.Template
	current  *Page
	_default string
}

var (
	RouteNotFound          = errors.New("route not found")
	InvalidTemplateHandler = errors.New("invalid template handler")
	TemplateIsNil          = errors.New("template is nil")
)

func (router *Router) Add(route string, handler Handler) {
	router.pages.Store(route, handler)
}

func (router *Router) Remove(route string) {
	router.pages.Delete(route)
}

func (router *Router) Get(route string) (Handler, error) {
	handler, _ := router.match(route)
	if handler != nil {
		return handler, nil
	}

	return nil, RouteNotFound
}

func (router *Router) match(name string) (Handler, []string) {
	var ret Handler
	var args []string

	router.pages.Range(func(s string, handler Handler) bool {
		route := split(s, '/')
		ns := split(name, '/')

		if len(ns) != len(route) {
			return true
		}

		for i, r := range route {
			if len(ns) > i && len(ns[i]) == 0 {
				ret = nil
				return false
			}

			if len(r) == 0 || r[0] == ':' {
				if len(ns) <= i {
					ret = nil
					return false
				}

				args = append(args, ns[i])
				ret = handler
				continue
			}

			if len(ns) <= i {
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

	return ret, args
}

func (router *Router) Handler() error {
	url := proto.URL()

	if url.Path == "/" {
		url.Path = router._default
		history.Default.Push(nil, router._default, url)
	}

	handler, args := router.match(url.Path)
	if handler == nil {
		return RouteNotFound
	}

	if router.current == nil {
		router.current = &Page{
			ctx:       router.ctx,
			logger:    router.logger.With("route", url.Path),
			states:    &proto.MapStore{},
			Store:     &proto.MapStore{},
			Router:    router,
			Arguments: args,
			Query:     url.Query(),
		}
	} else {
		router.current = &Page{
			ctx:       router.ctx,
			logger:    router.logger.With("route", url.Path),
			states:    router.current.states,
			Store:     &proto.MapStore{},
			Router:    router,
			Arguments: args,
			Query:     url.Query(),
		}
	}

	if err := handler(router.current); err != nil {
		return err
	}

	currentTemplate := router.current.template.Load()
	if currentTemplate == nil {
		return TemplateIsNil
	}

	b := buffer.New()

	if err := router.template.ExecuteTemplate(b, *currentTemplate, proto.StoreToData(router.current.Store)); err != nil {
		return err
	}

	create := proto.Create(b.Data())

	root := proto.GDocument.Value().Call("getElementById", "root")
	root.Call("replaceChildren", create...)

	router.logger.Debug("render done")

	router.current.post.Range(func(_ int, post func()) bool {
		go post()
		return true
	})

	router.logger.Debug("post done")
	return nil
}

func (router *Router) SetDefault(_default string) {
	router._default = _default
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

func New(ctx context.Context, logger *slog.Logger, template *template.Template) *Router {
	if logger == nil {
		logger = slog.Default()
	}

	return &Router{
		ctx:      ctx,
		logger:   logger.WithGroup("router"),
		pages:    &sync.Map[string, Handler]{},
		template: template,
	}
}
