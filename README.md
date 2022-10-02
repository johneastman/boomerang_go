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

### Grammar
```yaml
STATEMENT:
- ASSIGN
- PRINT
- EXPRESSION
- RETURN('return')
EXPRESSION:
- ADD('+')
- SUBTRACT('-')
- MULTIPLY('*')
- DIVIDE('/')
- LEFT_POINTER('<-')
- AT('@')
- FACTOR
FACTOR:
- NUMBER('float64')
- STRING
- BOOLEAN('true' | 'false')
- MINUS('-')  # unary operator
- OPEN_PAREN('(')
- FUNCTION('func')
- IDENTIFIER  # variable
- LIST
```

### Comments
There are two types of comments:
* Inline
* Block

```
##
block comments appear between a pair of double `#`s
and can occupy
multiple
lines
##
i = 0;
if true {
  i = i + 1; # inline comments appear after a single # and occupy a single line
}
print(i)
```

### Data Types
|Name|Examples|
|----|--------|
|NUMBER|`1`, `2`, `3.14159`, `100`, `1234567890`, `0.987654321`|
|BOOLEAN|`true`, `false`|
|STRING|`"hello, world!"`, `"1234567890"`, `"abcdefghijklmnopqrstuvwxyz"`, `"My number is {1 + 1}"`|
|LIST|`(1, 2)`, `(1, 2, 3)`, `(1, 2, 3 (6, 7, 8), 4, 5)`|

### Operators

#### Binary (Infix) Operators
|Name|Literal|Left Valid Types|Right Valid Types|
|----|-------|----------------|-----------------|
|add|+|NUMBER|NUMBER|
|minus|-|NUMBER|NUMBER|
|multiply|*|NUMBER|NUMBER|
|divide|/|NUMBER|NUMBER|
|left pointer|<-|IDENTIFIER (of function)|LIST|
|at|@|LIST|NUMBER (must be an integer)|
|equal|==|Any Type|Any Type|

#### Unary (Prefix) Operators
|Name|Literal|Valid Types|
|----|-------|-----------|
|minus|-|NUMBER|
|not/negate|`not`|BOOLEAN|

### Statements

#### Variable Assignment
Syntax: `IDENTIFIER = EXPRESSION`


Examples:
```
number = 1;
number = 1 + (2 * 2) - 3;
number = -1 + 1;
```

#### Print
Syntax: `print(EXPRESSION, EXPRESSION, ..., EXPRESSION)`


Examples:
```
print(1, 2, 3 + 4);

number = 3 + 4 / 2;
print(number, number * 2);

print(); # Does nothing
```

#### Return
Syntax: `return EXPRESSION`


Examples:
```
return 1;
return 1 + 1;
return "hello, world!";
return (1, 2 + 3, 5);
```

#### If Statements
Syntax: `if EXPRESSION { STATEMENT; STATEMENT; ...; STATEMENT; };`


Examples:
```
number = 1;
if true {
  number = 2;
};
print(number);  # number: 2

number = 1;
if false {
  number = 2;
}
print(number)  # number: 1
```

### Expressions

#### Lists
Syntax: `(EXPRESSION, EXPRESSION, ..., EXPRESSION)`


To access an element in a list, use the `@` symbol, with the list on the left side and the index on the right side. The index for the first element is 0. Values outside the range 0 to (len(LIST) - 1) will cause an index-out-of-range error.


Examples:
```
numbers = (5, 10, 15, 20);
value = numbers @ 0;  # value: 5
value = numbers @ 1;  # value: 10
value = numbers @ 2;  # value: 15
value = numbers @ 3;  # value: 20
```

To append values to a list, use the `<-` operator. On the right side of that operator, a single value can be passed, which will add that value to the end of the list, or a LIST can be passed, which combines the two lists. Be aware that this operation creates a new list, and the original list is not modified.
```
names = ("John", "Joe", "Jerry");

names = names <- "James"; # names: ("John", "Joe", "Jerry", "James")

names = names <- ("Jimmy", "Jack", "Jacob"); # names: ("John", "Joe", "Jerry", "James", "Jimmy", "Jack", "Jacob")
```

#### Functions
Syntax: `func(IDENTIFIER, IDENTIFIER, ..., IDENTIFIER) { STATEMENT; STATEMENT; ...; STATEMENT };`


If a `return` statement is not present, the last statement in a function's body is returned. All custom functions (those defined in boomerang files) return a LIST object. If nothing is returned from the function or an error occurrs, `(false)` is returned. If the function does return successfully, `(true, <RETURN_VALUE>)` is returned, where `RETURN_VALUE` is what the function is expected to return.

To get the actual return value from a function, call `unwrap` on the function's return value. That method takes a LIST object and a default value to return. If the function successfully returns a value, the actual value will be returned. Otherwise, the provided default value is returned.


Examples:
```
add = func(a, b) {
  a + b;
};
sum = add <- (1, 2); # sum: (true, 3)
value = unwrap <- (sum, 0) # value: 3

sum = func(c, d) {
  c + d;
} <- (1, 2); # sum = (true, 3)
value = unwrap <- (sum, 0) # value: 3

value = func() { # value: (true, 24)
  number = 1 + 1;
  (number + 2) * 6;
} <- ();
value = unwrap <- (value, 0) # value: 24

value = func() {} <- ();  # value: (false)

result = unwrap <- (value, 2) # result: 2
```

#### Boolean Operations
* Negate a boolean expression
  ```
  not true;  # false
  not false; # true
  ```
* Compare two values. The values being compared do not have to be compatible or the same type
  ```
  1 == 1; # true
  1 == 2; # false
  true == "hello, world!"; # false
  ```