# Proto

### Dealing with bloated JavaScript frameworks is painful, let it Go and handle it by yourself with a zero-dependency tool designed for WebAssembly.

---

## Install: `go get -u github.com/Dviih/proto`

## Utilities
- `String` - This is an interface that requires `String() string` to be implemented.
- `IsValue(interface{})` - Check if a value can be used as `js.Value`.

## Usage
- `Subscription.Channel()` - Returns a sender channel.
- `Subscription.Handler(Handler)` - However data should be handled here; the Handler argument is a function where it takes data and `js.Value`.
- `New()` - Returns a `*Subscription` instance and guaranteed to exist, else it panics.

In HTML the element must include a field `id`.

## Example

#### The example is split into two parts, HTML shows the required HTML code, and Go shows the required WebAssembly Go code.

```html
<html lang="en">
    <head>
	    <title>Proto</title>
	    <meta charset="utf-8"/>
	    <script src="wasm_exec.js"></script>
	    <script>
		    const go = new Go();
		
		    WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject).then((result) => {
			    go.run(result.instance);
		    });
	    </script>
    </head>
    <body>
        <span id="example">Loading...</span>
    </body>
</html>
```

```go
package main

import (
	"github.com/Dviih/proto"
	"syscall/js"
	"time"
)

func main() {
	example := proto.New("example")

	go func() {
		for {
			example.Channel() <- "Example"
			time.Sleep(time.Nanosecond) // Go is quite fast.
		}
	}()

	go example.Handler(func(i interface{}, value js.Value) {
		value.Set("innerHTML", i)
	})
}
```

#### More examples are available at `examples/`.

---
#### Made for Gophers by Dviih