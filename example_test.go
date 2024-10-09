package runnergroup_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/txthinking/runnergroup"
)

func Example() {
	g := runnergroup.New()

	s1 := &http.Server{
		Addr: ":9991",
	}
	g.Add(&runnergroup.Runner{
		Start: func() error {
			log.Println("start1")
			return s1.ListenAndServe()
		},
		Stop: func() error {
			log.Println("stop1")
			return s1.Shutdown(context.Background())
		},
	})

	s2 := &http.Server{
		Addr: ":9992",
	}
	g.Add(&runnergroup.Runner{
		Start: func() error {
			log.Println("start2")
			return errors.New("fail")
			// return s2.ListenAndServe()
		},
		Stop: func() error {
			log.Println("stop2")
			return s2.Shutdown(context.Background())
		},
	})

	s3 := &http.Server{
		Addr: ":9993",
	}
	g.Add(&runnergroup.Runner{
		Start: func() error {
			log.Println("start3")
			return s3.ListenAndServe()
		},
		Stop: func() error {
			log.Println("stop3")
			time.Sleep(5 * time.Second)
			return s3.Shutdown(context.Background())
		},
	})

	go func() {
		time.Sleep(3 * time.Second)
		g.Done()
	}()
	log.Println(g.Wait())
	// Output:
}
