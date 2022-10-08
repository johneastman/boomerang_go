# Syntax
* TODO: add table of contents to each section
* TODO: add section on switch statements

## Comments
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

## Data Types
|Name|Examples|
|----|--------|
|NUMBER|`1`, `2`, `3.14159`, `100`, `1234567890`, `0.987654321`|
|BOOLEAN|`true`, `false`|
|STRING|`"hello, world!"`, `"1234567890"`, `"abcdefghijklmnopqrstuvwxyz"`, `"My number is {1 + 1}"`|
|LIST|`(1, 2)`, `(1, 2, 3)`, `(1, 2, 3 (6, 7, 8), 4, 5)`|

## Operators

### Binary (Infix) Operators
|Name|Literal|Left Valid Types|Right Valid Types|
|----|-------|----------------|-----------------|
|add|+|NUMBER|NUMBER|
|minus|-|NUMBER|NUMBER|
|multiply|*|NUMBER|NUMBER|
|divide|/|NUMBER|NUMBER|
|left pointer|<-|IDENTIFIER (of function)|LIST|
|at|@|LIST|NUMBER (must be an integer)|
|equal|==|Any Type|Any Type|

### Unary (Prefix) Operators
|Name|Literal|Valid Types|
|----|-------|-----------|
|minus|-|NUMBER|
|not/negate|`not`|BOOLEAN|

## Statements

### Variable Assignment
Syntax: `IDENTIFIER = EXPRESSION`


Examples:
```
number = 1;
number = 1 + (2 * 2) - 3;
number = -1 + 1;
```

### Print
Syntax: `print(EXPRESSION, EXPRESSION, ..., EXPRESSION)`


Examples:
```
print(1, 2, 3 + 4);

number = 3 + 4 / 2;
print(number, number * 2);

print(); # Does nothing
```

## Expressions

### Lists
Syntax: `(EXPRESSION, EXPRESSION, ..., EXPRESSION)`


To access an element in a list, use the `@` symbol, with the list on the left side and the index on the right side. The index for the first element is 0. Values outside the range 0 to (len <- (LIST)) - 1) will cause an index-out-of-range error.


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

### Functions
Syntax: `func(IDENTIFIER, IDENTIFIER, ..., IDENTIFIER) { STATEMENT; STATEMENT; ...; STATEMENT };`


The last statement in a function's body is returned. All custom functions (those defined in boomerang files) return a LIST object. If nothing is returned from the function or an error occurrs, `(false)` is returned. If the function does return successfully, `(true, <RETURN_VALUE>)` is returned, where `RETURN_VALUE` is what the function is expected to return.

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

### Boolean Operations
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

### If Expressions
Syntax: `if EXPRESSION { STATEMENT; STATEMENT; ...; STATEMENT; } else { STATEMENT; STATEMENT; ...; STATEMENT; };`


In Boomerang, if-else statements are actually expressions that return a value. They behave exactly like functions, returning `(true, <VALUE>)` if successful, and `(false)` if there are no statements/an error occurred. The only difference is that the condition after 'if` determines which block is executed.

To get the actual return value, use the builtin `unwrap` method, similar to function return values.

Below are some examples.


Examples:
```
# The condition after `if` returns `true`, so the value of `number` is `(true, 5)`
number = if true {
  5;
} else {
  10;
};
print(unwrap <- (number, 0));  # prints 5


# The condition after `if` returns `false`, so the value of `number` is `(true, 10)`
number = if false {
  5;
} else {
  10;
};
print(unwrap <- (number, 0));  # prints 10


# There are no statements in the `else` block, but the condition after `if` returns `false`, so the value of `number` is `(false)`
number = if false {
  5;
} else {};
print(unwrap <- (number, 0));  # prints 0, the default value passed to unwrap


# The condition after `if` is true, but the if-block has no statements, so the value of `number` is `(false)`
number = if true {} else {
  5;
};
print(unwrap <- (number, 0));  # prints 0, the default value passed to unwrap
```
