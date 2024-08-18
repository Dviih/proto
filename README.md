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

#### For examples, it should be coming and will be available in `examples/`.

---
#### Made for Gophers by Dviih