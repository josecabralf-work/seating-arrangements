# Go guidelines

These guidelines MUST be followed by all Go code in this project

## Dependency management

Do not introduce new third-party dependencies unless explicitly told to do so.
Existing approved dependencies already present in the repository may be used
freely.

You are free to use anything from the standard library of any language.

## Documentation

Exported Go functions must have brief go docs explaining what they do and the
types of errors they may return. Input arguments are described naturally in the
flow of text. For example:
```go
// DoSomething does something with input. It may return ErrValidation if the
// input is invalid or ErrUnauthorized if the user identified in ctx does not
// have permissions to do something.
func DoSomething(ctx context.Context, input string) error {}
```

Non-exported functions that are not trivial may have brief go docs explaining
what they do and the types of errors they may return. Omit documentation for
trivial non-exported functions.

Comments in the middle of functions should NOT explain WHAT the code is doing,
but instead focus on WHY certain things are needed.

There may be comments that divide longer blocks of code and those may summarize
WHAT that longer block of code is doing. These must be used sparingly.

NEVER use hyphen-based delimiters for documentation.

## Code style

All code comments must end with a period in all languages. Prefer full
sentences, and if a comment seems too obvious, just omit it.

Log messages should always start with a lowercase letter unless referencing
something that is spelled with uppercase letters elsewhere in the codebase.

This project uses Go 1.26+, so it relies on new capabilities of the `new`
builtin to create a pointer to a literal value (for example, `new("foo")` or
`new(42)`) instead of defining helper functions like `strPtr` or `intPtr`.

Use `ID` for identifiers (not `Id`). Per [Go code review
comments](https://go.dev/wiki/CodeReviewComments#initialisms), initialisms
should be treated as words. When exposed as struct fields, variable names, or
method names, all letters of an initialism should be uppercase. For example:
`UserID`, `GetProfileID`, `orgEventID`.

## Error handling

Error messages should always start with a lowercase letter and should not end
with a period. Whenever reasonable, error messages should start with "cannot".

Always check errors with errors.Is or errors.As. Avoid returning bare
sentinels, preferring wrapped errors instead.

Prefer to use the `cannot do something` format for error messages.

Prefer using `%v` when returning errors. Only use `%w` with clear intent:
always document that the function is returning the type (or sentinel) of the
wrapped error.

## Testing

Use the `test` variable when iterating over subtests in a Go table test.

When declaring table tests, follow the style below:
- divide structs into multiple lines;
- keep braces compact;
- never redeclare types when they can be inferred from the slice type;
- slices with structs must be formatted in multiple lines;
- slices of primitive types should be in a single line.

Example of a well-formatted table test:
```go
var tests = []struct {
	desc   string

	thing  Thing
	things map[string][]Thing

	wantKv map[string]string
}{{
	desc: "...",
	thing: Thing{
		Name: "banana",
		Price: 2,
		Aliases: []string{"fruit", "yellow"},
	},
	things: map[string][]Thing{
		"thingset1": []Thing{{
			Name: "chair",
			Price: 100,
			Aliases: []string{"nicechair", "goodchair"},
		}, {
			Name: "table",
			Price: 200,
			Aliases: []string{"cooltable", "awesometable"},
		}},
		"thingset2": []Thing{{
			Name: "door",
			Price: 250,
		}},
	},
	wantKv: map[string]string{
		"a": "b",
		"c": "d",
	},
}, {
	desc: "...",
	// ...
}}
```

This project uses `github.com/frankban/quicktest`, imported as `qt` for all
test functions. Prefer to use `qt` utilities for making assertions, like
`qt.Equals`, `qt.DeepEquals`, `qt.Not`, `qt.IsNil`, etc.

Avoid writing test code that checks every field of a given struct. Instead,
prefer `qt.DeepEquals`.

If predictability is needed for time in tests, make sure to use the functions
from `internal/util/timeutil/timeutil.go`.

The description for each Go subtest should be named `desc` and it should always
start with a lowercase letter. Keep test descriptions short and direct, and
properly capitalize any acronyms.

For test struct definitions (but not declarations), put line breaks in between
the `desc` and input fields, and also between the input fields and the expected
output fields, for example:

```go
var tests = []struct {
	desc string

	name string

	wantName 			 string
	wantErr 			 error
	wantErrMatches string
}{...}
```

Prefer to group table test cases with similar characteristics together. Put all
the success cases first, followed by error cases. Within each group, order test
cases by increasing complexity.

When coming up with test URLs and email addresses, always use reserved domains
like `example.com` or `example.org`.

When writing tests, prefer attributes in test structs to define behavior,
rather than moving too much test logic to the test run function. Try to keep
the test run function as small as possible.

Never perform direct database manipulation with SQL in tests.

Always use `t.Parallel()` in the test functions and `c.Parallel()` inside table
test runner functions.

Since Go 1.22, loop variables are created per-iteration. NEVER use the
`test := test` pattern in test runner loops.

When using gomock and setting up more than one expected call in a mock, make
sure to use `gomock.InOrder`. It is fine to skip this recommendation when
there is a single function call being mocked.

NEVER patch global objects in tests. Instead, modify code being tested to
receive dependencies explicitly, and then mock dependencies in tests.

Tests should be in a separate `*_test` package. Only make exceptions when
explicitly asked by the developer, and exceptions have to be justified with a
brief comment at the beginning of the file.
