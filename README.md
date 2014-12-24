tracerr
=======

Traceable errors in Go.

#### Example:

```
package main

import (
	"errors"
	"fmt"

	"github.com/st3v/tracerr"
)

func main() {
	foo := &foo{}

	err := nested(4, func() error {
		return foo.bar()
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}

func nested(depth int, fn func() error) error {
	if depth <= 1 {
		return fn()
	}
	return tracerr.Wrap(nested(depth-1, fn))
}

type foo struct{}

func (f *foo) bar() error {
	return tracerr.Wrap(errors.New("FooBarError"))
}
```

#### Output:

```
$ ./example
FooBarError
  at (*foo).bar (main/tracerr_example.go:36)
  at inner (main/tracerr_example.go:26)
  at nested (main/tracerr_example.go:19)
  at nested (main/tracerr_example.go:21)
  at nested (main/tracerr_example.go:21)
  at nested (main/tracerr_example.go:21)
  at main (main/tracerr_example.go:11)
  at main (runtime/proc.c:255)
  at goexit (runtime/proc.c:1445)
```
