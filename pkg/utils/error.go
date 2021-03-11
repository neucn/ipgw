package utils

import (
	"net"
	"os"
	"syscall"
)

func IsNetworkError(err error) bool {
	netErr, ok := err.(net.Error)
	if !ok {
		return false
	}

	if netErr.Timeout() {
		return true
	}

	opErr, ok := netErr.(*net.OpError)
	if !ok {
		return false
	}

	switch t := opErr.Err.(type) {
	case *net.DNSError:
		return true
	case *os.SyscallError:
		if errno, ok := t.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ECONNREFUSED:
				return true
			case syscall.ETIMEDOUT:
				return true
			}
		}
	}

	return false
}
