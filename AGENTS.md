# AGENTS.md

## Purpose

This file defines rules and expectations for AI agents generating, modifying, or reviewing Go code in this repository.

All generated code must be idiomatic, maintainable, testable, and consistent with standard Go practices.

## General Principles

* Prefer simple, readable, idiomatic Go over clever abstractions.
* Write code that is easy to understand, test, debug, and maintain.
* Follow the standard Go project conventions unless the repository clearly uses a different style.
* Avoid unnecessary dependencies.
* Avoid premature abstraction.
* Preserve existing behavior unless explicitly asked to change it.
* Prefer explicit error handling over hidden control flow.
* Do not introduce global mutable state unless there is a strong reason.
* Do not silently ignore errors.
* Keep functions small and focused.
* Keep packages cohesive and minimal.

## Formatting

* All Go code must be formatted with `gofmt`.
* Use `goimports` when adding, removing, or reorganizing imports.
* Do not manually align code in ways that conflict with `gofmt`.
* Use tabs for indentation in Go files, as enforced by `gofmt`.
* Remove unused imports, variables, constants, and functions.
* Keep line length reasonable, but prioritize idiomatic Go formatting over arbitrary limits.

## Language Version

* Use the Go version declared by the project.
* Check `go.mod` before using newer language features.
* Do not introduce features unsupported by the module’s declared Go version.
* If no version is clear, prefer broadly compatible modern Go.

## Package Design

* Use short, meaningful package names.
* Avoid package names such as `utils`, `helpers`, or `common` unless already established in the project.
* Keep package APIs small and focused.
* Do not export identifiers unless they are required outside the package.
* Exported identifiers must have meaningful names and comments.
* Avoid circular dependencies.
* Prefer dependency injection over package-level globals.

## Naming

* Use idiomatic Go naming.
* Use `camelCase` for unexported identifiers.
* Use `PascalCase` for exported identifiers.
* Keep names concise but descriptive.
* Avoid unnecessary abbreviations.
* Use common Go initialisms consistently, such as:
    * `ID`
    * `URL`
    * `HTTP`
    * `JSON`
    * `SQL`
    * `API`
    * `UUID`
* Prefer names like `userID`, `httpClient`, and `jsonEncoder`.
* Do not use Hungarian notation.
* Avoid names that only describe the type, such as `stringValue` or `mapData`, unless useful for clarity.

## Error Handling

* Always handle returned errors.
* Do not discard errors with `_` unless there is a clear and documented reason.
* Prefer returning errors instead of logging and continuing silently.
* Wrap errors with context using `fmt.Errorf` and `%w`.
* Do not wrap errors if callers need direct comparison and no sentinel or helper is provided.
* Use `errors.Is` and `errors.As` for error inspection.
* Define sentinel errors only when callers need to compare against them.
* Prefer error messages that are concise, lowercase, and without trailing punctuation.
* Include useful context in errors.

Example:

```go
if err != nil {
	return fmt.Errorf("load config: %w", err)
}
```

Do not write:

```go
if err != nil {
	return fmt.Errorf("Failed to load config.")
}
```

## Logging

* Do not use logging as a substitute for error handling.
* Do not log sensitive data.
* Do not introduce a new logging library if the project already uses one.
* Use the existing logging style and logger.
* Library packages should generally return errors instead of logging directly.
* Application entry points may log errors before exiting.
* Prefer structured logging when the project already uses it.

## Context Usage

* Accept `context.Context` as the first argument when a function performs I/O, waits, calls external services, or may be canceled.
* Name the parameter `ctx`.
* Do not store contexts in structs unless there is a strong reason.
* Do not pass `nil` as a context.
* Use `context.Background()` only at top-level entry points.
* Use `context.TODO()` only as a temporary placeholder.
* Respect context cancellation and deadlines.

Example:

```go
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	// ...
}
```

## Concurrency

* Use goroutines deliberately.
* Avoid goroutine leaks.
* Ensure every started goroutine can terminate.
* Use channels when they clarify ownership or synchronization.
* Prefer `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, and `errgroup` where appropriate.
* Protect shared mutable state.
* Avoid data races.
* Use `go test -race` when modifying concurrent code.
* Do not close channels from the receiving side.
* Do not close a channel unless the sender owns it.
* Prefer context cancellation for stopping goroutines.

## Interfaces

* Accept interfaces, return concrete types when practical.
* Define interfaces close to where they are consumed.
* Keep interfaces small.
* Avoid unnecessary interfaces with only one implementation unless they improve testability or decoupling.
* Do not create large “service” interfaces without a clear need.
* Prefer standard interfaces such as:
    * `io.Reader`
    * `io.Writer`
    * `io.Closer`
    * `fmt.Stringer`
    * `http.Handler`

## Structs

* Keep structs focused on a single responsibility.
* Prefer explicit fields over map-based unstructured data when the schema is known.
* Use field tags consistently.
* Do not export struct fields unless callers need direct access.
* Validate struct inputs at boundaries.
* Avoid embedding types unless it clearly improves the API.
* Do not use embedding only to save typing.

## Functions and Methods

* Keep functions short and focused.
* Prefer returning early to reduce nesting.
* Avoid deeply nested control flow.
* Group related parameters into structs when parameter lists become long.
* Avoid boolean parameters that make call sites unclear.
* Prefer clear method receivers.
* Use pointer receivers when methods mutate the receiver or copying would be expensive.
* Use value receivers for small immutable values when appropriate.
* Be consistent with receiver names.

Example:

```go
func (s *Store) Save(ctx context.Context, user User) error {
	// ...
}
```

## Generics

* Use generics only when they reduce duplication and improve clarity.
* Do not use generics for simple cases where interfaces or concrete types are clearer.
* Keep type constraints simple.
* Prefer standard constraints and comparable types where appropriate.
* Avoid overly abstract generic utilities.

## Comments and Documentation

* Write comments for exported identifiers.
* Exported comments must begin with the identifier name.
* Explain why, not merely what.
* Avoid redundant comments.
* Document non-obvious behavior, edge cases, and assumptions.
* Keep comments up to date when modifying code.

Example:

```go
// UserStore persists and retrieves users.
type UserStore struct {
	// ...
}
```

## Testing

* Add or update tests for behavior changes.
* Prefer table-driven tests for multiple cases.
* Use subtests with `t.Run`.
* Prefer `testify` for assertions and requirements.
    * Use `require` when the test cannot continue after a failure.
    * Use `assert` when the test can continue and report additional failures.
* Keep tests deterministic.
* Avoid relying on test execution order.
* Avoid sleeping in tests when synchronization is possible.
* Use temporary directories with `t.TempDir()`.
* Use `t.Cleanup()` for cleanup.
* Use `httptest` for HTTP handlers and clients.
* Use `testing/iotest` where useful for I/O edge cases.
* Test public behavior, not implementation details.
* Include tests for error paths.
* Include tests for boundary cases.
* Do not weaken or remove tests unless they are obsolete and the reason is clear.
* Write tests the way they are written in the example below.

Example:

```go
func TestParseUserID(t *testing.T) {
    type given struct {
        id string
    }

    type want struct {
        id  string
        err bool
    }

    tests := map[string]struct {
        given given
        want  want
    }{
        "valid id": {
            given: given{
                id: "user-123",
            },
            want: want{
                id: "user-123",
            },
        },
        "empty id": {
            given: given{
                id: "",
            },
            want: want{
                err: true,
            },
        },
    }

    for name, test := range tests {
        t.Run(name, func(t *testing.T) {
            got, err := ParseUserID(test.given.id)

            if test.want.err {
                require.Error(t, err)
                return
            }

            require.NoError(t, err)
            require.Equal(t, test.want.id, got)
        })
    }
}
```

## Benchmarks

* Add benchmarks only when performance is relevant.
* Use `testing.B`.
* Reset timers where setup should not be measured.
* Avoid misleading microbenchmarks.
* Report allocations when useful with `b.ReportAllocs()`.

## Fuzzing

* Use fuzz tests for parsers, decoders, validators, and boundary-heavy code.
* Keep fuzz targets deterministic.
* Add seed corpus values for known edge cases.
* Do not rely on external services in fuzz tests.

## HTTP Code

* Always set appropriate timeouts on HTTP clients and servers.
* Do not use `http.DefaultClient` for production code unless acceptable in the project.
* Always close response bodies.
* Check HTTP status codes explicitly.
* Use `http.NewRequestWithContext` when creating outbound requests.
* Avoid reading unbounded response bodies.
* Use `http.MaxBytesReader` for request bodies where appropriate.
* Validate and sanitize inputs.

Example:

```go
request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
if err != nil {
	return fmt.Errorf("create request: %w", err)
}

response, err := client.Do(request)
if err != nil {
	return fmt.Errorf("send request: %w", err)
}
defer response.Body.Close()

if response.StatusCode != http.StatusOK {
	return fmt.Errorf("unexpected status: %s", response.Status)
}
```

## Database Code

* Use context-aware database methods such as `QueryContext`, `QueryRowContext`, and `ExecContext`.
* Always close rows.
* Always check `rows.Err()`.
* Use transactions where multiple related changes must be atomic.
* Roll back transactions on error.
* Avoid building SQL by string concatenation with untrusted input.
* Use parameterized queries.
* Handle `sql.ErrNoRows` explicitly when relevant.
* Keep database logic testable.

Example:

```go
rows, err := db.QueryContext(ctx, query, arg)
if err != nil {
	return fmt.Errorf("query users: %w", err)
}

defer rows.Close()

for rows.Next() {
	// ...
}

err := rows.Err()
if err != nil {
	return fmt.Errorf("iterate users: %w", err)
}
```

## File and I/O Code

* Always close files and other resources.
* Check errors from file operations.
* Use `os.ReadFile` and `os.WriteFile` for simple whole-file operations.
* Use streaming APIs for large files.
* Avoid loading unbounded data into memory.
* Use `io.Reader` and `io.Writer` to improve testability.
* Use `fs.FS` when working with abstract filesystems.
* Preserve file permissions intentionally.

## JSON and Encoding

* Handle encoding and decoding errors.
* Use struct tags intentionally.
* Avoid `map[string]interface{}` when a concrete type is known.
* Use `json.Decoder` for streams.
* Consider `DisallowUnknownFields` for strict API inputs.
* Avoid silently ignoring unknown or invalid fields unless intended.
* Do not expose internal struct fields through JSON accidentally.

## Security

* Never hardcode secrets, tokens, passwords, or private keys.
* Do not log sensitive information.
* Validate all external input.
* Use parameterized SQL queries.
* Avoid command injection by not passing untrusted strings to shells.
* Prefer `exec.CommandContext` over shell execution.
* Use secure random values from `crypto/rand` for security-sensitive randomness.
* Do not use `math/rand` for tokens or secrets.
* Set appropriate timeouts for network operations.
* Avoid path traversal vulnerabilities when handling file paths.
* Use `filepath.Clean` carefully and verify paths remain within allowed directories.
* Do not disable TLS verification unless explicitly required for tests.

## Dependencies

* Prefer the Go standard library when it is sufficient.
* Do not add dependencies without a clear benefit.
* Use existing project dependencies where possible.
* Keep dependencies maintained and reputable.
* Avoid introducing large frameworks for small tasks.
* Update `go.mod` and `go.sum` consistently.
* Run `go mod tidy` after dependency changes.
* Do not vendor dependencies unless the project already does so.

## CLI Code

* Return clear error messages.
* Use non-zero exit codes for failures.
* Keep command parsing separate from business logic.
* Do not call `os.Exit` deep inside library code.
* Write user-facing output to stdout.
* Write errors and diagnostics to stderr.
* Respect context cancellation where possible.

## Configuration

* Validate configuration at startup.
* Fail fast on invalid required configuration.
* Avoid hidden defaults for critical settings.
* Document environment variables and config fields.
* Do not read environment variables throughout the codebase; centralize configuration loading.
* Keep configuration structs explicit.

## Time Handling

* Use `time.Time` and `time.Duration`.
* Avoid representing durations as raw integers unless required by an external format.
* Use UTC for storage and comparisons unless local time is explicitly required.
* Inject clocks or time providers in tests when behavior depends on current time.
* Avoid flaky tests based on real wall-clock timing.

## Numeric Code

* Choose integer types intentionally.
* Use `int` for general in-memory counts and indexes.
* Use fixed-size integers when required by protocols, storage formats, or external APIs.
* Check for overflow when relevant.
* Avoid floating-point for money.
* Use decimal or integer minor units for monetary values.

## API Compatibility

* Preserve public APIs unless a breaking change is requested.
* Avoid changing exported names without need.
* Keep serialized formats backward compatible when possible.
* Add fields in a backward-compatible way.
* Do not remove or rename JSON fields casually.
* Document intentional breaking changes.

## Repository Consistency

* Follow existing project layout and patterns.
* Match existing naming, error handling, logging, and testing style.
* Reuse existing helpers where appropriate.
* Do not introduce a competing architecture.
* Prefer minimal diffs.
* Avoid unrelated refactoring.
* Keep generated changes scoped to the request.

## Code Generation Rules

When generating Go code:

* Include all required imports.
* Ensure code compiles.
* Avoid placeholder code unless explicitly requested.
* Do not leave `TODO` comments unless the missing work is intentional and clearly explained.
* Do not invent APIs that are not present in the repository.
* Do not change behavior outside the requested scope.
* Preserve existing comments unless they are inaccurate.
* Prefer small, incremental changes.
* Add tests for new behavior.
* Update documentation when public behavior changes.
* Ensure examples are valid Go code.
* Always use `status.Errorf`, `status.Error`, etc. for wrapping Go errors to gRPC status codes.

## Review Checklist

Before finalizing changes, verify:

* Code is formatted with `gofmt`.
* Imports are organized.
* Tests are added or updated.
* Errors are handled.
* Error messages include useful context.
* Public identifiers are documented.
* No sensitive data is logged or committed.
* No unnecessary dependencies were added.
* Concurrency is safe.
* Resources are closed.
* Contexts are used where appropriate.
* Existing behavior is preserved unless intentionally changed.
* The code compiles.
* The relevant tests pass.

## Commands

Use these commands when applicable:

```sh
go fmt ./...
go test ./...
go test -race ./...
go vet ./...
go mod tidy
```

If the repository uses `make`, `task`, `just`, or another build tool, prefer the existing project commands.

## Do Not Do

* Do not ignore errors.
* Do not introduce data races.
* Do not use panics for normal error handling.
* Do not use global mutable state unnecessarily.
* Do not add dependencies without justification.
* Do not write overly abstract code.
* Do not change unrelated files.
* Do not remove tests to make changes pass.
* Do not hardcode credentials or environment-specific paths.
* Do not log secrets or personal data.
* Do not use `interface{}` or `any` when a concrete type is clear.
* Do not create large catch-all utility packages.
* Do not bypass context cancellation.
* Do not silently swallow failures.

## Preferred Style Summary

* Simple is better than clever.
* Explicit is better than implicit.
* Small interfaces are better than large interfaces.
* Composition is better than inheritance-like embedding.
* Clear errors are better than vague errors.
* Tests should describe behavior.
* The standard library is preferred unless a dependency is justified.
* Maintainability matters more than novelty.
