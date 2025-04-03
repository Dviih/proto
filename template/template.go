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

package template

import (
	"github.com/Dviih/proto"
	"github.com/Dviih/proto/state"
	"html/template"
	"io/fs"
	"reflect"
	"sync/atomic"
)

type Template struct {
	*template.Template
	states proto.Store
}

var funcMap = map[string]interface{}{
	"state": func(s string) string {
		return "<state state=\"" + s + "\"></state>"
	},
}

}

}

}



	if err != nil {
		return nil, err
	}

}

	return &Template{
	}
}
