package redis

import (
	"fmt"
	"os"
	"syscall"
)

func DoSomething(fd int) error {
	buf := make([]byte, 64)
	if l, err := syscall.Read(fd, buf); err != nil || l < 1 {
		fmt.Fprintf(os.Stderr, "error reading message: %v\n", err)
		return err
	}

	fmt.Printf("client says: %s\n", buf)

	if l, err := syscall.Write(fd, []byte("world")); err != nil || l < 1 {
		fmt.Fprintf(os.Stderr, "error writing message: %v\n", err)
		return err
	}

	return nil
}
