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

var retry = time.Second / 20

func Retry() time.Duration {
	return retry
}

func SetRetry(duration time.Duration) {
	retry = duration
}

func Wait(fn func() js.Value) js.Value {
	for {
		if v := fn(); !v.IsNull() {
			return v
		}

		time.Sleep(retry)
	}
}
