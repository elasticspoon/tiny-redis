package redis

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

var ErrUnexpectedEOR = errors.New("unexpected end of file")

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

func readFull(fd int, bytes int) error {
	buf := make([]byte, bytes)
	for bytes > 0 {
		n, err := syscall.Read(fd, buf)
		switch {
		case n <= 0:
			return ErrUnexpectedEOR
		case err != nil:
			return err
		}

		bytes -= n
	}
	return nil
}
