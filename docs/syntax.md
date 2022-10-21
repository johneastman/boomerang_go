# Syntax
* [Comments](#comments)
* [Data Types](#data-types)
* [Operators](#operators)
    * [Binary (Infix) Operators](#binary-infix-operators)
    * [Unary (Prefix) Operators](#unary-prefix-operators)
* [Statements](#statements)
    * [Variable Assignment](#variable-assignment)
    * [Print](#print)
    * [While Loop](#while-loop)
    * [Break](#break)
    * [Block Statements](#block-statements)
* [Expressions](#expressions)
    * [Lists](#lists)
    * [Functions](#functions)
    * [When Expressions](#when-expressions)
    * [For Loops](#for-loops)

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
when true {
  is i == 0 {
    i = i + 1; # inline comments appear after a single `#` and occupy a single line
  }
  else {
    i = i + 2;
  }
};
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

#### Addition
* `NUMBER + NUMBER`: add two numbers together

#### Subtraction
* `NUMBER - NUMBER`: subtract two numbers together

#### Multiplication
* `NUMBER * NUMBER`: multiply two numbers together
 
#### Division
* `NUMBER / NUMBER`: divide two numbers together. The right number cannot be zero

#### Send
* `FUNCTION <- LIST`: perform a function call, where the left side is a function and the right side is the arguments being passed to that function
* `LIST <- EXPRESSION`: append the right value to the list on the left
* `LIST <- LIST`: combine the two lists, adding the values in the list on the right to the end of the list on the left

#### At
* `LIST @ NUMBER`: get the element at a given position (right) in the list (left)

#### Equal
* `EXPRESSION == EXPRESSION`: Compare two values and return `true` if they are the same; `false` otherwise

#### And
* `BOOLEAN and BOOLEAN`: `true` if left and right are both `true`; `false` otherwise

#### Or
* `BOOLEAN or BOOLEAN`: `true` if left or right are `true`; `false` otherwise

#### In
* `EXPRESSION in LIST`: `true` if the left value is in the list on the right


### Unary (Prefix) Operators

#### Negative
* `-NUMBER`: negate a number. Positive numbers become negative and negative numbers become positive

#### Negate
* `not BOOLEAN`: flip a boolean value. `false` becomes `true` and `true` becomes `false`

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

### While Loop
Syntax: `while EXPRESSION { STATEMENT; STATEMENT; ...; STATEMENT; };`


Unlike for-loops, while-loops do not return any values. They simply execute the block statement until the condition evaluates to `false`.


Examples:
```
i = 0;
while i < 10 {
  i = i + 1;
};
print(i);  # "i" is 10
```

### Break
Syntax: `break;`


Break statements terminate a loop early. These statements are only allowed in loops--using them outside a loop will result in an error.


Examples:
```
i = 0;
while i < 10 {
  when {
    i == 5 {
      break;
    }
  };
  i = i + 1;
};
print(i);  # "i" is 5

# the loop will terminate when "i" is "5". The list returned from this loop would only be 5 elements long (0 - 4)
for i in range <- (0, 10) {
  when {
    i == 5 {
      break;
    }
  };
  i;  # Because "i" is returned for every element, the list returned by this loop will be `((true, 0), (true, 1), (true, 2), (true, 3), (true, 4))`.
};
```

### Block Statements
Block statements are multiple statements defined between `{` and `}`. These statements cannot be independently defined and appear as part of other constructs (functions, when expressions, etc.). Block statements return the result of the last statement/expression wrapped in a list.

If the block statement returns a value, the block statement will return `(true, <VALUE>)`, where `<VALUE>` is the returned value, and `true` indicates that a value was returned. For example, the below function utilizes a block statement that returns `a + b`, and because that statement/expression returns a value, the function will return `(true, a + b)`.
```
add = func(a, b) {
  a + b;
};

value = add <- (2, 3); # value: (true, 5)
```

However, block statements that return nothing simply return `(false)`, indicating that the block statement returned no value. For example, the below function takes a value and prints it to the output stream. But but because `print`, the last statement in the block statement, returns no value, the function returns `(false)`.
```
printVal = func(v) {
  print("The value is {v}");
};

# "The value is 2" is printed to output stream, and value equals (false)
value = printVal <- (2);
```

To extract the actual return value of a block statement, use the builtin `unwrap` method. See [builtin functions](../docs/builtins.md) for more information.

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


Functions return the result of their associated block statement (see [Block Statements](#block-statements) for more information).


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

### When Expressions
Syntax:
* **NOTE:** `NOTHING` is the absense of characters
* `|` means "or"
```
when [EXPRESSION | not | NOTHING] { 
  [is | NOTHING] EXPRESSION { 
    STATEMENT;
    STATEMENT;
    ...
    STATEMENT;
  }
  [is | NOTHING] EXPRESSION { 
    STATEMENT;
    STATEMENT;
    ...;
    STATEMENT
  }
  ...
  [else | NOTHING] {
    STATEMENT;
    STATEMENT;
    ...;
    STATEMENT;
  }
};
```
The block statement associated with the case matching the `when` expression is run and the value of the last statement or expression is returned (see [Block Statements](#block-statements) for more information).


`when` expressions act as both "if-'else if'-else" and switch statements, depending on how they are implemented (although in Boomerang, `when` is an expression and can return a value). When the implementation acts as a switch statement, an expression is provided after `when` and the `is` keyword comes before each case. For example:
```
num = 0;
when num {
  is 0 { ... }
  is 1 { ... }
  ...
  else { ... }
};
```

When then implementation acts as an "if-'else if'-else" statement, nothing is provided after `when` for `true`, and `not` is provided `false`. Below are some examples:
```
# The code block for "num == 0" is run because `num` does equal 1.
num = 0;
when {
  num == 0 { ... }
  num == 1 { ... }
  ...
  else { ... }
};

# The code block for "num == 1" is run because `num` does not equal 1.
num = 0;
when not {
  num == 0 { ... }
  num == 1 { ... }
  ...
  else { ... }
};
```

Additionally, the "else" block is not required:
```
when {
  true { ... }
  false { ... }
};
```
If none of the conditions are a match, the `when` expression will return `(false)`:
```
num = 0;
value = when num {
  is 1 { ... }
  is 2 { ... }
};
print(value); # value: (false)
```

Be aware that these slight syntactic differences are enforced by the language. The following examples will produce errors:
```
when {
  is true { ... }  # ERROR
};

num = 1;
when num {
  1 { ... }  # ERROR
};
```

### For Loops
Syntax: `for IDENTIFIER in LIST { STATEMENT, STATEMENT, ..., STATEMENT }`


For loops act similar to the `map` function in other languages. A new list is returned containing the result of the last expression or statement in the block statement. The resulting list will be a list of lists as a result of evaluating a block statement (see [Block Statements](#block-statements) for more information).


Examples:
```
# Use for-loop as a regular loop
list = (1, 2, 3, 4, 5);
for element in list {  # for loop returns ((false), (false), (false), (false), (false))
  print(element);
};

# use for-loop as map
squared = for element in list {
  element * element;
};
print(squared); # squared: ((true, 1), (true, 4), (true, 9), (true, 16), (true, 25))
```
