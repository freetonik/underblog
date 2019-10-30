// +build linux darwin
// +build !windows

package internal

import (
	"log"
	"syscall"
)

const defaultWLimit = 100

// Change default system limits for unix like OS
func setRLimit(cur, max uint64) error {
	var rLimit syscall.Rlimit

	rLimit.Cur = cur
	rLimit.Max = max

	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	return err
}

// GetWorkersLimit returns workers limit according file descriptors (Rlimit / 2) for mac \ linux
// todo: cleanup required
func GetWorkersLimit(qSize int) int {

	limit := defaultWLimit

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[ERROR]: error getting Rlimit")
		return limit
	}

	limit = int(rLimit.Cur / 2)

	if qSize < limit {
		limit = qSize
	}

	return limit
}
