# Grammar

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
- AT('@')
- FACTOR
- IF
FACTOR:
- NUMBER('float64')
- STRING
- BOOLEAN('true' | 'false')
- MINUS('-')  # unary operator
- OPEN_PAREN('(')
- FUNCTION('func')
- IDENTIFIER  # variable, function calls
- LIST
```
