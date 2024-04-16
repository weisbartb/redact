# Redaction

## About
Redaction is a library that helps protect accidental ex-filtration of sensitive/protected fields.
This package allows you
to define stuct tags
that will automatically remove any potentially sensitive data run the object is passed through a filter. 
This filter is intended to be applied against objects prior to

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

