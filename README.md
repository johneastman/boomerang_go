# Boomerang
A custom interpreted programming language written in Go.

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

## Wiki
Additional documentation and notes are on the [wiki](https://github.com/johneastman/boomerang/wiki).
