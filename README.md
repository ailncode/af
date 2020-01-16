What is AF
-------------------------------------------------------------------------------------------------
*The name of AF comes from air refueling.*

*AF is a package of Go.*

*You can use the AF write a graceful reload HTTP server easily.*

Support
-------------------------------------------------------------------------------------------------
* Linux
* OS X
* Windows (only build and run)

Features
-------------------------------------------------------------------------------------------------

* graceful reload
* graceful stop
* custom signal handler


Install
-------------------------------------------------------------------------------------------------

```	go
go get -u github.com/ailncode/af
#or
go mod edit -require=github.com/ailncode/af@latest
```

Usage
-------------------------------------------------------------------------------------------------

1. Import the AF

```go
import "github.com/ailncode/af"
```

2. Simple use

```go
package main

import(
	"fmt"
	"net/http"
	"github.com/ailncode/af"
)

func main(){
    http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
    refuel := af.Default()
    //refuel.Addr = ":8080"
    //:8080 is the default Addr
    //refuel.Handler = http.DefaultServeMux
    //http.DefaultServeMux is the default HTTP handler
    //refuel.ShutdownTimeOut = time.Second * 10
    //time.Second * 10 is the default ShutdownTimeOut
    refuel.Run()
    //You can check error in here.
}
```
3. Reload & Stop your server

```shell
#Reload
kill -USR2 <pid>
#Stop
kill -INT <pid> #or kill -TERM <pid>
```

4. Use your custom Listen Address Handler ShutdownTimeOut

```go
package main

import(
    "net/http"
    "time"
    "github.com/ailncode/af"
    "github.com/gin-gonic/gin"
)

func main(){
    //Use handler like *gin.Engine
    router :=gin.Default()
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK,"ok")
	})
    refuel := af.Default()
    refuel.Handler = router
    //You can use the below code also.
    //af := af.NewWithHandler(r)
    refuel.Addr = ":80"
    refuel.ShutdownTimeOut = time.Second * 10
    refuel.Run()
}
```

5. Use your custom system signal handler

```go
package main

import (
	"github.com/ailncode/af"
	"syscall"
	"time"
)

func main() {
	refuel := af.Default()
	refuel.ShutdownTimeOut = time.Second * 10
	refuel.HandleSignal(func(r *af.AF) {
		r.Stop()
	}, syscall.SIGINT, syscall.SIGTERM)
	refuel.HandleSignal(func(r *af.AF) {
		r.Reload()
		//You can check error in here
	}, syscall.SIGUSR2)
	refuel.Run()
}
```

6. Use your custom server

```go
package main

import (
	"af"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:"8080",
		Handler:http.DefaultServeMux,
		//...
	}
	refuel := af.NewWithServer(server)
	refuel.Run()
}
```

