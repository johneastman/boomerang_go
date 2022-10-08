# Builtin Functions

## unwrap
* **Example:** `unwrap <- (list, defaultValue)`
* **Description:** get the actual return value from a function. If no return value is found, return `defaultValue`
* **Arguments:**
    * list: the return value from a block statement (function call, if-else expression, etc.). This list will either be `(true, <VALUE>)` or `(false)`, depending on whether a value was returned from the block statement.
    * defaultValue: this value is returned if `list` is `(false)`, meaning no value was returned from the block statement. 
