# Boomerang
A custom programming language written in Go.

**NOTE:** This project has been archived. The Boomerang project is being continued in [this repository](https://github.com/johneastman/boomerang).

Boomerang is:
* **Interpreted.** Code executes in the runtime of another language (in this case, Go).
* **Multi-paradigm** 
    * **Procedural.** Commands execute in the order they are defined.
    * **Functional.** Nothing is mutable, and constructs that return values (e.g., functions) return monads (there is no `nil` or `null` in the language).
    * **Imperative:** Explicitly define a series of commands to execute (what to do, as opposed to what to achieve).
* **Dynamically Typed.** Variable, function, parameter, etc. types are not declared explicitly; rather, they are interpreted from literal characters in the code (e.g., `1` and `3.14159` are numbers, `"hello world!"` is a string, `true` and `false` are booleans, etc.).
* **Strongly Typed.** The language has strict rules for how different types interact (e.g. `1 + 1` or `1 + 1.5` are valid, but `1 + "hello world!"` is invalid).

## Background
I originally started developing this project in [Python](https://github.com/johneastman/boomerang_old), but as the project grew, I ran into issues with Python's dynamic typing system because Python's runtime does not enforce type annotations (see the builtin [`typing`](https://docs.python.org/3/library/typing.html) module for type hints). I tried to resolve these issues with [mypy](https://github.com/python/mypy), a static code analysis tool that uses type hints, but I found myself regularly dealing with edge cases mypy could not handle, and refactoring code for the sake of mypy/the type checker. 

I realized if I wanted this project to grow, I would need to use a statically-typed language where type annotations are enforced during compile time or runtime. I settled on Go for it's balance of performance and a modern syntax (though I was also just interested in learning the language). Rewriting the project in Go has also allowed me to reflect on changes to Boomerang's syntax, as well as general implementation changes.

## Setup and Install
1. Setup and install [Go](https://go.dev/doc/install)
1. Clone/Download this repository
1. Open a terminal and `cd` into the downloaded repository's root directory
1. To run the main program, run `go run main.go`
1. To run the tests, run `go test -v ./tests`

## Language Specs
* [Grammar](docs/grammar.md)
* [Syntax](docs/syntax.md)
* [Builtin Functions](docs/builtins.md)

## License
This project is licensed under the [MIT License](LICENSE).

## Development Notes
Notes related to development. For example, notes on how to contribute, project design/structure, previous features and/or why they were removed, etc. Anything that a developer or contributer may find helpful to know.


This section is a work in progress.

### Language vs. Program Errors
There are two types of errors: Language Errors and Program Errors. Language errors are caused by users writing Boomerang code, such as syntax errors, evaluation errors, etc. These errors are created with `utils.CreateError` in `utils.go`.

Program errors are errors created during development or by the developer. These errors should never be raised by users writing Boomerang code and exist to inform developers when the code is broken in some way. To raise a program error, use `panic`.

### Notes on Previous Features

#### Removed `if-else` Expressions
In [this commit](https://github.com/johneastman/boomerang/commit/32397105ad307c3467f6936cee2a17b74b01b3f8), `if-else` expressions were removed. This is because `when` expressions can be used to perform the same functionality (see [When Expressions](docs/syntax.md#when-expressions))
