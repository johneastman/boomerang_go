# Boomerang
A custom interpreted programming language written in Go.

## Setup and Install
1. Setup and install [Go](https://go.dev/doc/install).
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

### Data Types
|Name|Examples|
|----|--------|
|NUMBER|`1`, `2`, `3.14159`, `100`, `1234567890`, `0.987654321`|
|BOOLEAN|`true`, `false`|
|STRING|`"hello, world!"`, `"1234567890"`, `"abcdefghijklmnopqrstuvwxyz"`, `"My number is {1 + 1}"`|
|LIST|`(1, 2)`, `(1, 2, 3)`, `(1, 2, 3 (6, 7, 8), 4, 5)`|

### Math Operators

#### Binary (Infix) Operators
|Name|Literal|Valid Types|
|----|-------|-----------|
|add|+|NUMBER|
|minus|-|NUMBER|
|multiply|*|NUMBER|
|divide|/|NUMBER|
|left pointer|<-|left expression: FUNCTION, BUILTIN_FUNCTION, right expression: LIST|

#### Unary (Prefix) Operators
|Name|Literal|Valid Types|
|----|-------|-----------|
|minus|-|NUMBER|

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

### Expressions

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
