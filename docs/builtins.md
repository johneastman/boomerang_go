# Builtin Functions
* `nArgs` as a argument name means the function takes any number of arguments
* `ANY` as a type means any data type is valid

## print

### Description
Output values to the console.

### Arguments
|Name|Type|Description|
|----|----|-----------|
|nArgs|LIST|list of object to be printed to the console|

### Returns:
* **Type:** LIST
* **Value:** `(false)`

### Examples
```
print <- (1, 2, 3)  # prints "1 2 3"
print <- ("hello, world",)  # prints '"hello, world"'
```

## input

### Description
Get input from the user.

### Arguments
|Name|Type|Description|
|----|----|-----------|
|prompt|STRING|Displayed to tbe user before user input. A colon and a space will be added after this string|

### Returns
* **Type:** STRING
* **Value:** user input as a string. If no input is given, the method returns an empty string

### Examples
```
# display: "input: "; 
# If user enters "hello, world", `user_input` will be a string with the value "hello, world".
user_input = input <- ("input",);
```

## unwrap

### Description
Get the actual return value from a monad. If no return value is found, return `defaultValue`.

### Arguments
|Name|Type|Description|
|----|----|-----------|
|monad|MONAD|A monad object|
|default_value|ANY|the value returned if `monad` contains a value|

### Returns
* **Type:** ANY
* **Value:** `VALUE` if `monad` contains a value; the default value of `monad` contains no value.

### Examples
```
list = (true, 5);
defaultValue = -1;
unwrap <- (list, defaultValue);  # returns "5" because the list contains "true"

list = (false);
unwrap <- (list, defaultValue);  # returns "-1" because the list contains "false" and no other values
```

## unwrap_all

### Description
Get a list of values from a list of monads.

### Arguments
|Name|Type|Description|
|----|----|-----------|
|list|LIST|a list of monads|
|default_value|ANY|the value used when any of the monads contain no values|

### Returns
* **Type:** ANY
* **Value:** a list of unwrapped monads

### Examples
```
list = (
  (true, 5),
  (false,),
  (true, 10),
  (true, 15),
);
defaultValue = -1;
unwrap_all <- (list, defaultValue); # return: (5, -1, 10, 15)
```

## len

### Description
Get the length of a list.

### Arguments
|Name|Type|Description|
|----|----|-----------|
|sequence|LIST or STRING|A LIST or STRING object|

### Returns:
* **Type:** NUMBER
* **Value:** the length of the given sequence

### Examples
```
len <- (1, 2, 3);  # 3

list = ("hello", "world");
len <- (list,);  # 2

len <- ((),);    # 0
len <- ((9,),);  # 1
len <- ("hello, world!");  # 13
```

## slice

### Description
Retrieve a sublist of `list` from `startPos` to `endPos` (inclusive). 

### Arguments
|Name|Type|Description|
|----|----|-----------|
|list|LIST or STRING|any LIST or STRING object|
|start_pos|NUMBER|index associated with the first element of the new list|
|end_pos|NUMBER|index associated with the last element of the new list|

### Returns
* **Type:** LIST
* **Value:** a new list from `list[start_pos:end_pos]`

### Examples
```
list = (1, 2, 3, 4, 5, 6, 7, 8, 9, 0);
slice <- (list, 1, 3);  # (2, 3, 4)
slice <- (list, 0, 0);  # (1,)
slice <- (list, 4, 2);  # error because the start index cannot be greater than the end index
```

## range

### Description
Return a list of numbers incrementing or decrementing from `start_value` to `end_value` (inclusive). If `start_value` is less than `end_value`, the numbers will increment; if `start_value` is greater than `end_value`, the numbers will decrement.

### Arguments
|Name|Type|Description|
|----|----|-----------|
|start_value|NUMBER|the number starting the sequence (the first element of the list)|
|end_value|NUMBER|the number ending the sequence (the last element in the list)|

### Returns
* **Type:** LIST
* **Value:** a list of incrementing numbers from `start_value` to `end_value`.

### Examples
```
range <- (0, 5)    # (0, 1, 2, 3, 4, 5)
range <- (5, 0)    # (5, 4, 3, 2, 1, 0)
range <- (10, 20)  # (10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
range <- (0, 0)    # (0)
range <- (0, -1)   # (0, -1)
range <- (5, -5)    # (5, 4, 3, 2, 1, 0, -1, -2, -3, -4, -5)
```

## random

### Description
Generate a random number between `min` and `max` (inclusive). The values of `min` and `max` can be negative or positive, but `min` must be less than `max`.

### Arguments
|Name|Type|Description|
|----|----|-----------|
|min|NUMBER|the minimum value the random number can be|
|max|NUMBER|the maximum value the random number can be|

### Returns
* **Type:** NUMBER
* **Value:** a random number between `min` and `max` (inclusive)

### Examples
```
random <- (0, 5)    # [0 to 5]
random <- (10, 20)  # [10 to 20]
random <- (0, 0)    # 0
random <- (-10, -5)  # [-10 to -5]
```


## is_success

### Description
Checks if a monad contains a value.

### Arguments
|Name|Type|Description|
|----|----|-----------|
|monad|MONAD|a monad object|

### Returns
* **Type:** BOOLEAN
* **Value:** `true` if the monad contains a value; `false` otherwise.

### Examples
```
add = func(a, b) {
  a + b;
};
sum = add <- (3, 4);
is_success <- (sum,)  # returns true

do_nothing = func() {};
did_nothing = do_nothing <- ();
is_success <- (did_nothing,)  # returns false
```
