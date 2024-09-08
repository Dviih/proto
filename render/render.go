/*
 *     An easy way to provide conectivity across multiple servers.
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

package render

import (
	"embed"
	"github.com/Dviih/proto"
	"html/template"
	"sync"
	"syscall/js"
)

type Render struct {
	root     *Element
	template *template.Template

	m    sync.Mutex
	data map[string]interface{}
}

func (render *Render) Add(key string, value interface{}) {
	defer render.m.Unlock()
	render.m.Lock()

	render.data[key] = value
}

