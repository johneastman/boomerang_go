# Boomerang
A simple recursive descent parser written in Go.

## Setup and Install
1. Setup and install [Go](https://go.dev/doc/install).
1. Clone/Download this repository
1. Open a terminal in the root directory
1. To run the main program, run `go run main.go`
1. To run the tests, run `go test -v ./tests`

## Language Specs

### Grammar
```yaml
STATEMENT:
- ASSIGN
- PRINT
- EXPRESSION
EXPRESSION:
- ADD('+')
- SUBTRACT('-')
- MULTIPLY('*')
- DIVIDE('/')
- LEFT_POINTER('<-')
- RIGHT_POINTER('->')
- FACTOR
FACTOR:
- NUMBER('float64')
- STRING
- MINUS('-')  # unary operator
- OPEN_PAREN('(')
- FUNCTION('func')
- IDENTIFIER  # variable
- PARAMETER
```

### Data Types
|Name|Examples|
|----|--------|
|NUMBER|`1`, `2`, `3.14159`, `100`, `1234567890`, `0.987654321`|
|STRING|`"hello, world!"`, `"1234567890"`, `"abcdefghijklmnopqrstuvwxyz"`, `"My number is {1 + 1}"`|
|PARAMETER|`(1, 2)`, `(1, 2, 3)`, `(1, 2, 3 (6, 7, 8), 4, 5)`|

### Math Operators

#### Infix Operators
|Name|Literal|Valid Types|
|----|-------|-----------|
|add|+|NUMBER|
|minus|-|NUMBER|
|multiply|*|NUMBER|
|divide|/|NUMBER|
|left pointer|<-|left expression: FUNCTION, right expression: PARAMETER|
|right pointer|->|right expression: PARAMETER, left expression: FUNCTION|

#### Prefix Operators
|Name|Literal|Valid Types|
|----|-------|-----------|
|minus|-|NUMBER|

### Statements

#### Variable Assignment
Syntax: `IDENTIFIER = EXPRESSION`
<br/>
Examples:
```
number = 1;
number = 1 + (2 * 2) - 3;
number = -1 + 1;
```

#### Print
Syntax: `print(EXPRESSION, EXPRESSION, ..., EXPRESSION)`
<br/>
Examples:
```
print(1, 2, 3 + 4);

number = 3 + 4 / 2;
print(number, number * 2);

print(); # Does nothing
```

#### Functions
Syntax: `func(IDENTIFIER, IDENTIFIER, ..., IDENTIFIER) { STATEMENT; STATEMENT; ...; STATEMENT };`
<br/>
The last statement in a function's body is returned. Currently, functions with no statements in their body are allowed.
<br/>
Examples:
```
add = func(a, b) {
  a + b;
};
add <- (1, 2);
(3, 4) -> add;

sum = func(c, d) {
  c + d;
} <- (1, 2);

value = () -> func() {
  number = 1 + 1;
  (number + 2) * 6;
};
```
