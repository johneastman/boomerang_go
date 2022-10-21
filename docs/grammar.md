# Grammar

## Statements
* Any action that does not return a value or produces side effects (e.g., printing or variable assignment)
* Any expression

## Expressions
* Any operation that returns a new value from existing values
* Any factor

## Factor
* Any value that returns itself when evaluated (numbers, strings, booleans, etc.)

```yaml
STATEMENT:
- ASSIGN
- PRINT
- WHILE_LOOP
- BREAK
- EXPRESSION
EXPRESSION:
- ADD('+')
- SUBTRACT('-')
- MULTIPLY('*')
- DIVIDE('/')
- SEND('<-')
- AT('@')
- NOT('not')
- IN('in')
- EQUAL('==')
- LESS_THAN('<')
- WHEN('when')
- FOR_LOOP('for')
- FACTOR
FACTOR:
- NUMBER('float64')
- STRING
- BOOLEAN('true' | 'false')
- LIST
- FUNCTION('func')
- IDENTIFIER  # variable, function calls
```
