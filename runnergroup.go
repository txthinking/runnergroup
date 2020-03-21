package runnergroup

import (
	"sync"
)

// RunnerGroup is like sync.WaitGroup,
// the diffrence is if one task stops, all will be stopped.
type RunnerGroup struct {
	runners []*Runner
	once    sync.Once
	done    chan byte
}

type Runner struct {
	// Start is a blocking func, until error or Stop called.
	Start func() error
	// After Stop called, the Start func must stop.
	Stop func() error
}

func New() *RunnerGroup {
	g := &RunnerGroup{}
	g.runners = make([]*Runner, 0)
	g.done = make(chan byte)
	return g
}

func (g *RunnerGroup) Add(r *Runner) {
	g.runners = append(g.runners, r)
}

// Call Wait after all task have been added,
// Return the first stopped Start's result, or return nil if stopped caused by Done.
func (g *RunnerGroup) Wait() error {
	errch := make(chan error)
	for _, v := range g.runners {
		go func(v *Runner) {
			err := v.Start()
			select {
			case <-g.done:
			case errch <- err:
				g.once.Do(func() {
					close(g.done)
					for _, v := range g.runners {
						_ = v.Stop()
					}
				})
			}
		}(v)
	}
	select {
	case <-g.done:
		return nil
	case err := <-errch:
		return err
	}
	return nil
}

// Call Done after Wait, if you want to stop all.
// According to the order of Add, return the first error that is not nil, or nil means all Stops have ended normally.
func (g *RunnerGroup) Done() error {
	var e error
	g.once.Do(func() {
		close(g.done)
		for _, v := range g.runners {
			if err := v.Stop(); err != nil {
				if e == nil {
					e = err
				}
			}
		}
	})
	return e
}
