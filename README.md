# Govader Backend Package



Govader-Backend is based on GoVader Package[https://github.com/jonreiter/govader](https://github.com/jonreiter/govader)


## Usage:

```sh
go get github.com/PIMPfiction/govader_backend
```

```go
package main

import (
	"github.com/PIMPfiction/govader_backend"
)

func main() {
	govader_backend.Serve()
}

```


### Sample Get Request

#### GET: http://localhost:8080?text=I%20am%20looking%20good

### Sample Post Request 

#### POST: http://localhost:8080/
### RequestBody: ```{"text": "I am looking good"}```