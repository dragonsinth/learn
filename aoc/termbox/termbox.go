package termbox

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
	"sync"
	"syscall"
)

func New(enable bool) *terminal {
	c := &terminal{}
	if enable {
		c.Start()
	}
	return c
}

type terminal struct {
	mu      sync.RWMutex
	enabled bool
	done    chan struct{}
}

func (c *terminal) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.enabled {
		return
	}

	if err := termbox.Init(); err != nil {
		fmt.Println(err)
		return
	}

	c.enabled = true
	c.done = make(chan struct{})
	go func() {
		defer close(c.done)
		for {
			ev := termbox.PollEvent()
			switch ev.Type {
			case termbox.EventKey:
				if ev.Ch == 0 {
					switch ev.Key {
					case termbox.KeyCtrlC:
						c.stop(false)
						// resend the signal
						_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
						return
					case termbox.KeyEsc:
						c.stop(false)
						return
					}
				}
			case termbox.EventError:
				panic(ev.Err)

			case termbox.EventInterrupt:
				return
			}

		}
	}()
}

func (c *terminal) Stop() {
	c.stop(true)
}

func (c *terminal) stop(interrupt bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.enabled {
		return
	}

	termbox.Close()
	if interrupt {
		termbox.Interrupt()
		<-c.done
	}
	c.done = nil
	c.enabled = false
}

func (c *terminal) Render(data [][]byte, ifDisabled io.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.enabled {
		if ifDisabled != nil {
			for _, line := range data {
				_, _ = fmt.Fprintln(ifDisabled, string(line))
			}
		}
		return
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()
	for y, line := range data {
		for x, c := range line {
			termbox.SetCell(x, y, rune(c), termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func (c *terminal) Enabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.enabled
}
