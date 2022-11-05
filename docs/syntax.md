# Syntax
* [Comments](#comments)
* [Data Types](#data-types)
* [Operators](#operators)
    * [Binary (Infix) Operators](#binary-infix-operators)
    * [Unary (Prefix) Operators](#unary-prefix-operators)
* [Statements](#statements)
    * [While Loop](#while-loop)
    * [Break](#break)
    * [Continue](#continue)
    * [Block Statements](#block-statements)
* [Expressions](#expressions)
    * [Variable Assignment](#variable-assignment)
        * [Multiple Variable Assignment](#multiple-variable-assignments)
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
|MONAD|`Monad{}`, `Monad{5}`, `Monad{"hello, world"}`, `Monad{true}`, `Monad{false}`, `Monad{(1, 2, 3)}`|

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
* `LIST @ NUMBER`: get the element at the given position (right) in the list (left)
* `STRING @ NUMBER`: get the character at the given position (right) in the list (left). The character returned will be of type STRING.

#### Equal
* `EXPRESSION == EXPRESSION`: Compare two values and return `true` if they are the same; `false` otherwise

#### And
* `BOOLEAN and BOOLEAN`: `true` if left and right are both `true`; `false` otherwise

#### Or
* `BOOLEAN or BOOLEAN`: `true` if left or right are `true`; `false` otherwise

#### In
* `EXPRESSION in LIST`: `true` if the left value is in the list on the right; `false` otherwise


### Unary (Prefix) Operators

#### Negative
* `-NUMBER`: negate a number. Positive numbers become negative and negative numbers become positive

#### Negate
* `not BOOLEAN`: flip a boolean value. `false` becomes `true` and `true` becomes `false`

## Statements

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

### Continue
Syntax: `continue;`


Skip to the beginning of the next iteration in a loop. In for-loops, `contine` acts like filter functions in other languages; the last expression in the block statement will not be returned after a `continue` statement.


Examples:
```
i = 0;
while i < 10 {
  when {
    i % 2 == 0 {
      i = i + 1;
      continue;
    }
  };
  print <- (i,);  # only prints 1, 3, 5, 7, 9
  i = i + 1;
};

list = range <- (10, 0);
new_list = for e in list {
  when {
    e % 2 == 0 {
      e;
    } else {
      continue;
    }
  };
};
print <- (new_list,);  # new_list == (Monad{10}, Monad{8}, Monad{6}, Monad{4}, Monad{2}, Monad{0})
```

### Block Statements
Block statements are multiple statements defined between `{` and `}`. These statements cannot be independently defined and appear as part of other constructs (functions, when-expressions, for-loops, etc.). Block statements return the value of the last statement wrapped in a builtin Monad object.

Monad objects may or may not contain a value. Monads containing a value will be represented as `Monad{<VALUE>}`, whereas monads without a value are represented as `Monad{}`. Below are some examples:
```
# Monad from a function that does return a value
add = func(a, b) {
  a + b;
};

value = add <- (2, 3); # value: Monad{5}

# Monad from function that does not return a value
assign_val = func(v) {
  new_v = v;
};

value = assign_val <- (2);  # value: Monad{}
```

Monad objects cannot be independently instantiated.

To extract the actual return value of a Monad object, use the builtin `unwrap` method. See [builtin functions](../docs/builtins.md) for more information.

## Expressions

### Variable Assignment
Syntax: `IDENTIFIER = EXPRESSION`


Variable assignments are expression, so they return the value of `EXPRESSION`. This allows users to assign the same value to multiple variables.


Examples:
```
number = 1;
number = 1 + (2 * 2) - 3;
number = -1 + 1;

a = b = c = 20;  # a, b, and c are all assigned to the value "20"
```

#### Multiple Variable Assignments
Syntax: `(IDENTIFIER, IDENTIFIER, ..., IDENTIFIER) = (EXPRESSION, EXPRESSION, ..., EXPRESSION);`


Allows multiple variables to be assigned values on one line. Each identifier in the list on the left is assigned to the value at the corresponding position in the list on the right. For example:
Examples:
```
# a == 1
# b == 2
# c == 3
(a, b, c) = (1, 2, 3);
```

If the number of identifiers is greater than the number of values, any identifier without a corresponding value will be assigned an empty Monad object. For example:
```
# a == 1
# b == 2
# c == Monad{} because there is no third value in the list on the right
(a, b, c) = (1, 2);
```

If the number of identifiers is less than the number of values, the value corresponding with the last identifier, along with the remaining values, are stored in a list and assigned to the last identifier. For example:
```
# a == 1
# b == (2, 3) because there are 3 values but only two identifiers, so `2` and `3` are stored in a list and assigned to `b`.
(a, b) = (1, 2, 3);
```


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
Syntax: `func(IDENTIFIER|ASSIGN, IDENTIFIER|ASSIGN, ..., IDENTIFIER|ASSIGN) { STATEMENT; STATEMENT; ...; STATEMENT };`


Functions return monads (see [Block Statements](#block-statements) for more information).


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

Functions can also be created with default parameter values. These default values will be used if no values is provided in the function call, but providing a value for that parameter in the function call will override the default value. Additionally, default values must be declared after any non-default parameters. Below are some examples:
```
add = func(a, b = 2) {
  a + b;
};

sum = add <- (5,); # sum equals 7
sum = add <- (5, 10); # sum equals 15

add = func(a = 1, b) {
  a + b;
};

sum = add <- (5,); # this will cause an error because "5" will override "a", but "b" will have no value.
sum = add <- (5, 10); # "5" overrides "a", and "10" is passed for "b", so sum equals 15.
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
The block statement associated with the matching case is run and a monad is returned (see [Block Statements](#block-statements) for more information).


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


For loops act similar to `map` functions in other languages. For each value in the list, the block statement is run, and a monad is returned (see [Block Statements](#block-statements) for more information).


Examples:
```
# Use for-loop as a regular loop
list = (1, 2, 3, 4, 5);
new_list = for element in list {
  print <- (element,);
};
print <- (new_list,)  # new_list: (Monad{}, Monad{}, Monad{}, Monad{}, Monad{})

# use for-loop as map
squared = for element in list {
  element * element;
};
print <- (squared,); # squared: (Monad{1}, Monad{4}, Monad{9}, Monad{16}, Monad{25})
```
