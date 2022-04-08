# Govader Backend Package

[![Go Reference](https://pkg.go.dev/badge/github.com/PIMPfiction/govader_backend.svg)](https://pkg.go.dev/github.com/PIMPfiction/govader_backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/PIMPfiction/govader_backend)](https://goreportcard.com/report/github.com/PIMPfiction/govader_backend)
[![codecov](https://codecov.io/gh/PIMPfiction/govader_backend/branch/master/graph/badge.svg?token=3KEBD30Q95)](https://codecov.io/gh/PIMPfiction/govader_backend)
![master](https://github.com/PIMPfiction/govader_backend/actions/workflows/tests.yml/badge.svg)

Govader-Backend is based on GoVader Package[https://github.com/jonreiter/govader](https://github.com/jonreiter/govader)


## Usage:

```sh
go get github.com/PIMPfiction/govader_backend
```

```go
package main

import (
	"github.com/PIMPfiction/govader_backend"
	echo "github.com/labstack/echo/v4"
	"fmt"
)

func main() {
	e := echo.New()
	err := Serve(e, "8080")
	if err != nil {
		panic(err)
	}
	fmt.Scanln()

}

```


### Sample Get Request

#### GET: http://localhost:8080?text=I%20am%20looking%20good

### Sample Post Request 

#### POST: http://localhost:8080/
### RequestBody: ```{"text": "I am looking good"}```
