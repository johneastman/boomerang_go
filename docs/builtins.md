# Builtin Functions

## print
* **Description:** output values to the console.
* **Arguments:**
    * args: any LIST object of any length.
* **Returns:** `(false)`.
* **Examples:**
  ```
  print <- (1, 2, 3)  # prints "1 2 3"
  print <- ("hello, world",)  # prints '"hello, world"'
  ```

## unwrap
* **Description:** get the actual return value from a block statement. If no return value is found, return `defaultValue`.
* **Arguments:**
    * list: the return value from a block statement (function call, if-else expression, etc.). This list will either be `(true, <VALUE>)` or `(false)`, depending on whether the block statement returned a value.
    * default_value: the value returned if `list` is `(false)`.
* **Returns:** `VALUE` if `list` is `(true, <VALUE>)`; default_value if `list` is `(false)`.
* **Examples:**
  ```
  list = (true, 5);
  defaultValue = -1;
  unwrap <- (list, defaultValue);  # returns "5" because the list contains "true"
  
  list = (false);
  unwrap <- (list, defaultValue);  # returns "-1" because the list contains "false" and no other values
  ```

## unwrap_all
* **Description:** get a list of values from a list of block statement return values.
* **Arguments:**
    * list: a list of values returned from a block statement (function call, if-else expression, for-loop, etc.). The values in this list will either be `(true, <VALUE>)` or `(false)`.
    * default_value: the value used if any of the values in `list` are `(false)`.
* **Returns:** a list of unwrapped values from block statements
* **Examples:**
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
* **Description:** return the length of a list.
* **Arguments:**
    * args: any LIST object of any length.
* **Returns:** the length of the given list.
* **Examples:**
  ```
  len <- (1, 2, 3);  # 3

  list = ("hello", "world");
  len <- list;  # 2
  len <- ();    # 0
  len <- (9,);  # 1
  ```

## slice
* **Description:** return a sublist of `list` from `startPos` to `endPos` (inclusive). 
* **Arguments:**
    * list: any LIST object.
    * start_pos: index associated with the first element of the new list.
    * end_pos: index associated with the last element of the new list.
* **Returns:** a new list from `list[start_pos:end_pos]`
* **Examples:**
  ```
  list = (1, 2, 3, 4, 5, 6, 7, 8, 9, 0);
  slice <- (list, 1, 3);  # (2, 3, 4)
  slice <- (list, 0, 0);  # (1,)
  slice <- (list, 4, 2);  # error because the start index cannot be greater than the end index
  ```

## range
* **Description:** return a list of numbers incrementing from the start value to the end value (inclusive).
* **Arguments:**
    * start_value: the number starting the sequence (the first element of the list).
    * end_value: the number ending the sequence (the last element in the list).
* **Returns:** a list of incrementing numbers from `start_value` to `end_value`.
* **Examples:**
  ```
  range <- (0, 5)    # (0, 1, 2, 3, 4, 5)
  range <- (10, 20)  # (10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
  range <- (0, 0)    # (0)
  range <- (0, -1)   # ()
  ```

## random
* **Description:** generate a random number between "min" and "max" (inclusive). The values of "min" and "max" can be negative or positive, but "min" must be less than "max".
* **Arguments:**
    * min: the minimum value the random number can be.
    * max: the maximum value the random number can be.
* **Returns:** a random number from `[min, max]`
* **Examples:**
  ```
  random <- (0, 5)    # [0 to 5]
  random <- (10, 20)  # [10 to 20]
  random <- (0, 0)    # 0
  random <- (-10, -5)  # [-10 to -5]
  ```
