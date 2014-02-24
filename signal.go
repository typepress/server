package server

import (
	"os"
	"sync"

	"github.com/typepress/core"
)

func sigReceive(shutdownnow chan bool) func(sig os.Signal) bool {
	var onceDown sync.Once
	var onceStop sync.Once
	return func(sig os.Signal) bool {
		switch sig {
		case os.Kill, os.Interrupt:
			onceDown.Do(func() {
				core.Recover(func() {
					shutdownnow <- true
				})
			})
			return true
		}
		if sig.String() == stopString {
			onceStop.Do(door.Stop)
			return true
		}
		return true
	}
}
