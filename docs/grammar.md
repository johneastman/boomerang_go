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
- IF-ELSE('if'/'else')
- SWITCH('when')
- FACTOR
FACTOR:
- NUMBER('float64')
- STRING
- BOOLEAN('true' | 'false')
- LIST
- FUNCTION('func')
- IDENTIFIER  # variable, function calls
```
