---
title: Syntax, Variables & Types
sidebar_label: Syntax, Variables & Types
description: Go's shape on the page — declarations, the type system, constants, and zero values.
---

# Syntax, Variables & Types

## Learning objectives

- Read any Go file top to bottom: package clause, imports, declarations.
- Declare variables all four ways and know when to use which.
- Explain zero values and why Go code leans on them instead of null checks.
- Use basic types, type conversions (never implicit), and typed constants with `iota`.

## Prerequisites

- [Module 0 — Environment Setup](../module-0-setup/environment)

## Time estimate

**3 hours**

## Concepts

### The shape of a Go file

```go
// Package greet says hello. (Every package gets a doc comment.)
package greet

import (
	"fmt"
	"strings"
)

// MaxRetries is exported (capitalized); maxDelay is not.
const MaxRetries = 3

func Greet(name string) string {
	return fmt.Sprintf("Hello, %s", strings.TrimSpace(name))
}
```

Three things carry meaning that other languages put in keywords:

- **Capitalization is visibility.** `Greet` is exported (public), `greet` is package-private. There is no `public`/`private` keyword.
- **The import block is managed by tooling** (`goimports`), grouped stdlib-first. You never hand-maintain it.
- **Doc comments** (`// Greet ...` directly above the declaration, starting with the name) are the API documentation — `pkg.go.dev` renders them.

### Declaring variables

```go
var count int            // declaration; count == 0 (zero value)
var name = "dx"          // type inferred
limit := 100             // short form — only inside functions; declare-and-assign
var a, b = 1, "two"      // multiple, mixed types
```

House style (matches the platform's linters): use `:=` inside functions when you're initializing with a meaningful value; use `var x T` when the zero value is what you want. Declare variables **as close as possible to first use**, not at the top of the function.

### Zero values

Every type has a usable zero value: `0`, `""`, `false`, `nil` for pointers/slices/maps/interfaces. Well-designed Go types are useful without initialization:

```go
var sb strings.Builder   // ready to use, no constructor
sb.WriteString("dx")

var mu sync.Mutex        // a valid, unlocked mutex
```

This is a design principle you'll meet again when defining your own structs: *make the zero value meaningful*, and when it can't be, provide a `NewX(...)` constructor.

### Basic types and conversions

`int`, `int64`, `float64`, `bool`, `string`, `byte` (= `uint8`), `rune` (= `int32`, a Unicode code point). Go **never converts numeric types implicitly**:

```go
var n int = 42
var f float64 = float64(n) // explicit, always
var u uint = uint(f)
```

This looks pedantic and prevents an entire class of silent-truncation bugs. Strings are immutable byte slices under the hood; indexing gives bytes, ranging gives runes — details in [Collections](collections).

### Constants and iota

```go
type Status int

const (
	StatusUnknown Status = iota // 0 — zero value means "not set"
	StatusPending               // 1
	StatusActive                // 2
	StatusRevoked               // 3
)
```

`iota` auto-increments within a `const` block. Platform convention (from the Uber style guide): **start enums at 1** unless the zero value is genuinely meaningful — here we made zero explicitly "unknown", which is the other acceptable pattern. Constants are untyped until used, which is why `const timeout = 5 * time.Second` works so cleanly.

:::info[Platform connection]
Open any DX service and the first file you read will use all of this: e.g. config structs full of zero-value-friendly fields, `iota` enums for domain states, and exported-vs-unexported as the entire access-control story. Note there are **no classes and no inheritance anywhere** — by the end of this module you'll see what Go uses instead (composition + interfaces).
:::

## Exercises

1. Write a program that declares the same variable four different ways, prints each with `fmt.Printf("%v %T\n", x, x)`, and comments on when you'd use each form.
2. Define a `Role` enum (`RoleConsumer`, `RoleProvider`, `RoleAdmin`, …) with `iota`, and a `String() string` method for it using a `switch`. (You'll formalize methods in [Structs & Methods](structs-methods).)
3. Predict the output before running: what is `var s []string; fmt.Println(s == nil, len(s))`?

## Check yourself

- What makes an identifier exported?
- Why does `var f float64 = 3; var n int = f` fail to compile, and why is that good?
- What's the zero value of a string? A pointer? A slice?
- When should an enum start at 1?

## References

- [A Tour of Go — Basics](https://go.dev/tour/basics/1)
- [Go by Example: Variables](https://gobyexample.com/variables), [Constants](https://gobyexample.com/constants)
- [Effective Go — Names, Constants](https://go.dev/doc/effective_go#names)
- [Google Go Style Guide](https://google.github.io/styleguide/go/) — the platform's primary style reference
