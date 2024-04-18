# Redaction

## About

Redaction is a library that helps protect accidental ex-filtration of sensitive/protected fields.
This package allows you to define stuct tags that will automatically remove any potentially sensitive data run the
object is passed through a filter.
This filter is intended to be applied against objects prior to

### Additional Design Notes

This library makes heavy use of reflection but provides a type-safe parametric wrapper for calling the library.
Reflection calls are cached using [github.com/weisbartb/rcache](https://github.com/weisbartb/rcache) and the
instruction.
Methods used by redact are memoized to prevent multiple parsing calls.

### Example

```go
package mypackage

import (
	"encoding/json"
	"github.com/weisbartb/redact"
	"net/http"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName" redact:"~[admin,csr]=remove(1)"`
	Password  string `redact:"all=zero"`
}

func MyHandler(r *http.Request, w http.ResponseWriter) {
	var u = User{
		FirstName: "fName",
		LastName:  "lName",
		Password:  "test",
	}

	err := redact.RedactRecord(&u)
	if err != nil {
		panic(err)
	}
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		panic(err)
	}
}
```

This will ensure that the password is zeroed out for all users
and that the last name is truncated to the first letter for anything that isn't an admin or csr.

## Built In Redaction Methods

### Zero

Zero will set a zero value for the field if the condition is met.
It takes no arguments and works with values and pointers.
**Note:** Pointers will be set to nil.

### Remove

Remove takes a single argument that is an integer
(positive or negative) and removes characters according to the direction.
If the argument is positive, it will remove characters after the first `n` characters.

`Example 555-555-5555 with remove(4) would be 555-`

If the argument is negative it will truncate to the last `n` characters in the string

`Example 555-555-5555 with remove(-4) would be 5555`

### Redact

Redact will redact non-matching characters (from argument 2) with argument 1. It takes two optional arguments,
the first argument allows you to specify the character to use for a redaction.
The second argument takes a list of characters (as a string) that are allowed.

`Example 555-555-5555 with redact("*","-") would be ***-***-****`

### Star

Star takes a single argument that is an integer
(positive or negative) and asterisks characters according to the direction.
If the argument is positive, it will asterisk characters after the first `n` characters.

`Example 555-555-5555 with remove(4) would be 555-********`

If the argument is negative it will asterisk to the last `n` characters in the string

`Example 555-555-5555 with remove(-4) would be ********5555`

## Adding new redaction methods

## Performance Notes

This library uses a copious amount of reflection, both in type introspection and re-assembly.
The goal is this library is to be "fast enough" and always "technically correct" from a data safety/integrity
standpoint.

## Safety Notes

This section has few notes around safety concerns and any caveats that are discovered that could lead to foot gun type
behavior.

- Objects are shallowly copied,
  having any objects who's underlying pointers are used by other objects may result in unexpected mutations.