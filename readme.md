## RunnerGroup

[![GoDoc](https://pkg.go.dev/golang.cc/runnergroup?status.svg)](https://pkg.go.dev/golang.cc/runnergroup)

RunnerGroup is like [sync.WaitGroup](https://pkg.go.dev/sync?tab=doc#WaitGroup), the diffrence is if one task stops, all will be stopped.

### Install

    $ go get golang.cc/runnergroup

### Example

```
import (
	"context"
	"log"
	"net/http"
	"time"

	"golang.cc/runnergroup"
)

func Example() {
	g := runnergroup.New()

	s := &http.Server{
		Addr: ":9991",
	}
	g.Add(&runnergroup.Runner{
		Start: func() error {
			return s.ListenAndServe()
		},
		Stop: func() error {
			return s.Shutdown(context.Background())
		},
	})

	s1 := &http.Server{
		Addr: ":9992",
	}
	g.Add(&runnergroup.Runner{
		Start: func() error {
			return s1.ListenAndServe()
		},
		Stop: func() error {
			return s1.Shutdown(context.Background())
		},
	})

	go func() {
		time.Sleep(5 * time.Second)
		log.Println(g.Done())
	}()
	log.Println(g.Wait())
	// Output:
}

```
