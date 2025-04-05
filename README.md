# Proto

### Dealing with bloated JavaScript frameworks is painful, let it Go and handle it by yourself with a tool designed for WebAssembly.

---

# Understand Proto

---

## Value

- A Proto Value is an enhanced `js.Value`.

## Value Interface
- `Name()` Returns the value name.
- `Value()` Returns the inner `js.Value`.

## How To Use
- `NewNamedValue(name, js.Value)` Returns `Value` based on name and `js.Value`.
- `NewEmptyValue(js.Value)` Returns `Value` with nil name.

---

## States
- A state is a variable that changes and notifies upon new store.

## State Interface
- `Id` Returns the State id, a string.
- `Store[T]` Store T.
- `Load` Loads T.
- `C` Returns a receive-only locked channel.

## Virtual Interface
- The Virtual state is a `State[interface{}]` by definition.
- A logic `virtual` struct exists, and it is what is used to store.

## Storer Interface
- `StoreState(id, Virtual)` Stores Virtual with id into Storer.
- `NewState(Type, string) Virtual` Creates a Virtual state with id, usage of the type is optional and might get used.

## Loader Interface
- `LoadState(id) interface{}` Loads a state into `interface{}` to be cast.

## How To Use
- `state.Store[T](Storer, Id) State[T]` Stores and returns a `State[T]` with the specified id.
- `state.Load[T](Loader, id) State[T]` Returns a `State[T]` with the specified id, if it doesn't exist it will return nil.

---

## Templates
- The template provided by Proto is built on top of `html/template` and supports states.

## Template Structure
##### Implement state.Storer and state.Loader.
- This structure is wrapped on `html/template` and includes direct access to it so any change is dealt there except for state implementations.

## Template State
##### Implement state.Virtual.
- This structure is a Virtual state implementing a template loader based on previously set value.
- The Store method behavior changes based on input. If it is a string, it will change the template name or if it is a `proto.MapStore` it will append to its `proto.Store`.

## How To Use
##### Import `github.com/Dviih/proto/template`.
- `template.New(name)` Returns a `*Template`.
- `template.ParseFS(FS, ...patterns)` Parses patterns inside FS and returns a `*Template`.

---

## The Router
- Proto includes a client side router responsible for the pages (or routes).
- The page has an exclusive state that evaluates into `textContent` in HTML.

## Router Structure
- This structure is responsible to hold context, logger, pages and the template.
- The router stores info such as default page and current page.
- `Add(route, Handler)` Adds a new route with specified handler.
- `Remove(route)` Removes a route from Router.
- `Template()` Returns the inner `*Template`.
- `Get(pattern)` Gets a route with matching pattern.
- `match(pattern)` Matches the pattern.
- `Handler()` Handles the page build and rendering; This function can be called as much as needed.
- `SetDefault(route)` Set default route.

## Page Structure
##### Implement state.Storer and state.Loader.
- A page holds context, logger, template name, events slice, post-rendering slice, states, store, router reference, arguments and query info.
- A post-rendering function is called after the template was executed, thus events should start from this step.
- `Logger` Returns `*slog.Logger`.
- `SetTemplate(name)` Sets the template name.
- `Event(proto.Value)` Returns a `*Event` based on `proto.Value` and managed by `*Page`.
- `Context()` Returns `context.Context`.
- `Go(route)` Changes URL path to route and calls Handler.
- `handle(state.Virtual)` Handles Virtual update events.
- `Post(func())` Add a function to post slice.

## Page State
##### Implement state.Virtual
- A page state exclusively sets `textContent` due to safety reasons.
- A page state is strict on its created type, any other set value call with mismatched type will return and not cause an update.

## How To Use
##### Import `github.com/Dviih/proto/router`.
- `router.New(Context, *Logger, *Template)` returns a `*Router`.

---

## Async Utility
- This utility is provided to handle promises of JavaScript.
- Since the structure of JavaScript requires things to just return, it causes a problem with Go because of I/O operations that can lock a thread for as long as it needs. However, in JavaScript it will cause the event loop to get stuck, thus the CPU will reach 100% really fast so we have `AsyncDeadline` global variable with a simple timeout.
- The previous issue can be solved if the function is fully in the background not requiring immediate return.

## How To Use
##### Import `github.com/Dviih/proto`.
- `proto.AsyncDeadline` The global deadline variable.
- `proto.Promise` Creates a promise, has `resolve` and `reject` parameters and immediately returns a `js.Value`.
- `proto.Await` Waits to resolve a `Promise` or else the deadline exceeds returning a deadline error.

---

## Store
- Store is a concept of unified storage across Proto.
- A key must be a string.

## Store Interface
- `Set(key, value)` Sets key to value.
- `Get(key) value` Returns value of the key.
- `Range(func(key, value) continue)` Ranges through store until `continue` is false or if all elements were sent.
- `Delete(key)` Deletes a key and its value.

## MapStore
- This structure implements `Store` on top of `sync.Map`.

## AnyStore
- This structure converts `interface{}` into `AnyStore` and simulates some methods to implement `Store`.
- Range will only run once.
- Delete will set value to nil.

## How To Use
##### Import `github.com/Dviih/proto`.
- `ToStore(interface{})` Returns Store based on input.
- `MapStoreFrom(interface{})` Creates a `MapStore` from input.
- `StoreToData(interface{})` Converts `proto.Store` back to `map[string]interface{}` or `interface{}`, currently used by `*Template`.

---

## Route and Template utilities
##### Import `github.com/Dviih/proto`.
- `URL` Returns `*url.URL` from current page.
- `Create` Creates an element from `innerHTML` bytes.

## Type Conversion utilities
- `ToInterface(proto.Value)` Returns `interface{}` from `proto.Value`.
- `ToValue(interface{})` Returns `js.Value` from `interface{}`.
- `BuildJSFunc(funcs)` Returns a crafted `js.Func` from multiple functions.
- `ToType[T](proto.Value)` Like `ToInterface` but converts to T or errors.
- `ToRType(Type, proto.Value)` Converts `proto.Value` into `reflect.Type`, returns `reflect.Value` of said type or error.
- `ToInt(proto.Value)` Converts multiple integers representation of JavaScript into Go's `int` type.

## Byte utilities
- `ToBytes` Returns a `[]byte` from `proto.Value`.
- `ToUint8Array` Returns a `Uint8Array` from `[]byte`.

---
#### Made for Gophers by Dviih
