---
title: Reflection — and When Not to Use It
sidebar_label: Reflection (and when not)
description: How reflect works, why struct tags depend on it, and why application code should almost never touch it.
---

# Reflection — and When Not to Use It

## Learning objectives

- Understand what `reflect` can do: inspect types, read struct tags, walk values at runtime.
- Explain how the libraries you use daily (JSON, viper, pgx row mapping) are built on reflection.
- Apply the platform's stance: reflection belongs **inside libraries**, not in service code.

## Prerequisites

- [Structs & Methods](structs-methods) (struct tags), [Generics](generics)

## Time estimate

**2 hours** — deliberately short. This is a "know it exists" topic.

## Concepts

### What reflection is

`reflect` lets a program inspect its own types and values at runtime:

```go
func describe(v any) {
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("%s %s tag:%q\n", f.Name, f.Type, f.Tag.Get("json"))
	}
}

describe(ServerConfig{})
// Port int tag:"port"
// Host string tag:"host,omitempty"
```

You give up compile-time safety: mistakes surface as runtime panics (`reflect: call of ... on zero Value`) instead of build failures. That trade is worth it exactly once per problem — inside a well-tested library.

### Where reflection is already working for you

Every one of these platform staples is reflection under the hood:

| You write | The library reflects |
|---|---|
| `json.Unmarshal(body, &req)` | walks `req`'s fields, matches `json:` tags |
| `dxconfig.LoadService[Config]()` | viper + mapstructure map YAML/env keys onto `mapstructure:` tags |
| `pgx.RowToStructByName[T](row)` | matches SQL column names to `db:` tags on `T` |

Notice the division of labor: **you** declare intent with struct tags (declarative, checked by linters and tests); **the library** does the dynamic work. Your service code stays fully statically typed.

### The platform rule

> Application code in DX services does not import `reflect`. If you think you need it, you usually want generics, an interface, or a type switch — in that order.

Reasons: reflected code is slower, panics instead of failing to compile, defeats `go vet` and the linters, and reads like a puzzle. The rare legitimate uses (a validation library, a test helper comparing arbitrary structs) belong in shared libraries where the complexity is paid once and covered by tests.

A related honesty note: `reflect.DeepEqual` appears in older test code, but prefer `cmp.Diff` (`google/go-cmp`) in new tests — better output, configurable comparisons.

:::info[Platform connection]
Struct tags are the visible half of the platform's reflection use — the invisible half lives in viper, `encoding/json`, and pgx. Search any DX service for `import "reflect"`: you should come up empty, and that's the point. When Module 3 shows `pgx.RowToStructByName` mapping rows into your domain structs "magically", you now know exactly what the magic is.
:::

## Exercises

1. Write `printTags(v any)` that prints every exported field's name, type, and `json` tag. Test it on a config struct.
2. Break it: pass an `int` instead of a struct, watch the panic, then guard with `reflect.ValueOf(v).Kind() != reflect.Struct` — and reflect (pun intended) on how much checking the compiler normally does for you.
3. Write the same "copy matching fields between two structs" helper twice: once with reflection, once as plain explicit code for the two concrete types. Compare length, readability, and what happens when a field is renamed.

## Check yourself

- Name three libraries in the platform stack built on reflection.
- Why does the platform ban `reflect` in service code but rely on it in libraries?
- What should you reach for instead of reflection, in order?

## References

- [Go Blog: The Laws of Reflection](https://go.dev/blog/laws-of-reflection)
- [reflect package docs](https://pkg.go.dev/reflect)
- [go-cmp](https://pkg.go.dev/github.com/google/go-cmp/cmp) — the better DeepEqual for tests
