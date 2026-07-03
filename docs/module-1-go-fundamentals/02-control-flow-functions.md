---
title: Control Flow & Functions
sidebar_label: Control Flow & Functions
description: if, for, switch, defer — and functions with multiple returns, variadics, and closures.
---

# Control Flow & Functions

## Learning objectives

- Use Go's single loop (`for`), `if` with init statements, and `switch` (including type-free switches).
- Write functions with multiple return values — the foundation of Go error handling.
- Use `defer` for cleanup and understand its evaluation rules.
- Write closures and variadic functions; recognize the early-return style the platform enforces.

## Prerequisites

- [Syntax, Variables & Types](syntax-variables-types)

## Time estimate

**3 hours**

## Concepts

### One loop to rule them all

```go
for i := 0; i < 10; i++ {}     // classic
for count > 0 {}               // while
for {}                         // forever (until break/return)
for i, v := range items {}     // over slices, maps, strings, channels
for range time.Tick(interval) {} // you'll meet this in workers
```

### if with init; switch without pain

```go
if err := validate(req); err != nil {
	return err
}
// err is scoped to the if — it doesn't leak
```

`switch` doesn't fall through by default (no `break` needed), cases can be expressions, and the condition is optional:

```go
switch {
case n < 0:
	return "negative"
case n == 0:
	return "zero"
default:
	return "positive"
}
```

### Early returns — the platform's readability rule

Go code handles the error/edge case first and returns, keeping the happy path at minimal indentation. This is an explicit rule in the DX style guide:

```go
// Good — happy path flows down the left margin.
func process(order *Order) error {
	if order == nil {
		return errors.New("nil order")
	}
	if !order.Paid {
		return errors.New("unpaid order")
	}
	ship(order)
	return nil
}
```

No `else` after a branch that returns. If you find yourself nesting three levels deep, restructure.

### Multiple return values

```go
func parsePort(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: %w", s, err)
	}
	return n, nil
}
```

`(T, error)` is *the* Go function signature. There are no exceptions to catch — errors are values you return and check. This gets a full page ([Error Handling](error-handling)); for now internalize the shape.

### defer

`defer` schedules a call to run when the function returns — however it returns. It's Go's cleanup mechanism:

```go
f, err := os.Open(path)
if err != nil {
	return err
}
defer f.Close() // runs on every exit path below this line
```

Rules that bite people: deferred calls run **LIFO**; arguments are **evaluated at defer time**, not at run time; and defers run on panic too (which is how HTTP middleware recovers — you'll see `Recoverer` in Module 3).

```go
func main() {
	x := 1
	defer fmt.Println("deferred x =", x) // captures 1 now
	x = 2
	fmt.Println("final x =", x)
}
// prints: final x = 2, then: deferred x = 1
```

### Closures and variadics

```go
func counter() func() int {
	n := 0
	return func() int { n++; return n } // captures n by reference
}

func sum(nums ...int) int { // variadic
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
```

:::info[Platform connection]
`defer` is everywhere in DX services: `defer f.Close()`, `defer rows.Close()`, and — critically — `defer tx.Rollback(ctx)` right after starting a database transaction (rollback becomes a no-op if the transaction commits first — Module 3's [Transactions](../module-3-advanced/transactions) page builds on exactly this). The early-return style is enforced in code review across every service.
:::

## Exercises

1. Write `func classify(n int) string` using an expression-less `switch`, then rewrite it with if/else and compare readability.
2. Write a function that opens two files and copies one to the other, using `defer` for both closes. What order do they close in?
3. Predict the output: a loop `for i := 0; i < 3; i++ { defer fmt.Println(i) }` — then run it.
4. Refactor this into early-return style:
   ```go
   func handle(u *User) error {
   	if u != nil {
   		if u.Active {
   			doWork(u)
   			return nil
   		} else {
   			return errors.New("inactive")
   		}
   	} else {
   		return errors.New("nil user")
   	}
   }
   ```

## Check yourself

- Why does Go not need `break` in switch cases?
- When are a deferred call's arguments evaluated?
- What does "keep the happy path left-aligned" mean, and why does the platform enforce it?

## References

- [A Tour of Go — Flow control](https://go.dev/tour/flowcontrol/1)
- [Go by Example: If/Else](https://gobyexample.com/if-else), [Switch](https://gobyexample.com/switch), [Defer](https://gobyexample.com/defer), [Closures](https://gobyexample.com/closures)
- [Effective Go — Control structures](https://go.dev/doc/effective_go#control-structures), [Defer](https://go.dev/doc/effective_go#defer)
