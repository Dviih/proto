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

package main

import (
	"embed"
	"flag"
	"net/http"
)

var (
	address     string
	certificate string
	key         string
)

//go:embed public
var embedded embed.FS

func main() {
	flag.StringVar(&address, "address", ":3000", "The address to listen on")
	flag.StringVar(&certificate, "certificate", "", "Path to a TLS certificate file")
	flag.StringVar(&key, "key", "", "Path to a TLS key file")

	flag.Parse()

	handler := http.FileServer(http.FS(embedded))

	if len(certificate) != 0 && len(key) != 0 {
		panic(http.ListenAndServeTLS(address, certificate, key, handler))
	}

	panic(http.ListenAndServe(address, handler))
}
