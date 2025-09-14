# pstr - A Regex Engine in Golang

## Reliability Issues to Fix (Made by Cursor)

- [X] Fix processRepeat() - Array bounds crash when ctx.Tokens is empty. Line 269: lastToken := ctx.Tokens[len(ctx.Tokens)-1] will panic if no tokens exist. Need to check if len(ctx.Tokens) > 0 before accessing last element.
- [X] Fix findNextSymbol() - Array bounds issue and logic error. Line 125: 'if currPos > len(regex)' should be 'if currPos >= len(regex)'. Also need to handle case where symbol is never found to prevent infinite loop.
- [X] Fix getMinMaxRange() - Missing error handling for invalid range syntax. Currently returns utils.Infinite for both values on error instead of returning an error. Need to validate range format and return proper errors for malformed {m,n} syntax.
- [X] Fix processBrackets() - Array bounds crash in line 211-212. insideRegex[1] and insideRegex[len(insideRegex)-2] will panic if insideRegex has less than 2 elements. Need to validate bracket content length before accessing elements.
- [X] Fix processGroup() - Insufficient validation. Line 162: newPos == 1 check is too simplistic. Need to validate that closing ) exists and that group content is not empty. Also need to handle nested groups properly.
- [X] Fix Token.ToNFA() - Unsafe type assertions that can panic. Lines 52, 63, 72, 83, 140: Using single-value type assertions like token.Value.([]Token) will panic if type is wrong. Need to use two-value form with ok check.
- [X] Fix ToNFA() - Array bounds crash when no tokens exist. Line 10: ctx.Tokens[0] will panic if ctx.Tokens is empty. Need to validate that tokens exist before processing.
- [X] Add input validation to Parse() - No validation for empty or nil regex strings. Should return error for empty input and validate basic regex structure before processing.
- [X] Fix processOr() - Incomplete implementation and potential bounds issues. Function processes RHS but doesn't handle edge cases properly. Need to validate that | operator has valid operands on both sides.
- [X] Add comprehensive error handling throughout parser - Many functions don't return errors when they should. Need to establish consistent error handling pattern and propagate errors properly up the call stack.

## References
-  [How to build a regex engine from scratch - rhaeguard](https://rhaeguard.github.io/posts/regex/)
-  [Go by Example](https://gobyexample.com)
-  [Tests from CPython Regex](https://github.com/python/cpython/blob/main/Lib/test/re_tests.py)
-  [regex101](https://regex101.com/)
-  [ASCII Table - Wikimedia Commons](https://commons.wikimedia.org/wiki/File:ASCII-Table-wide.svg)

