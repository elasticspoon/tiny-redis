package redis

import (
	"errors"
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

const MAX_MESSAGE = 4096

func OneRequest(fd int) error {
	buf := make([]byte, 4+MAX_MESSAGE+1) // header size = 4

	if err := readFull(fd, &buf, 4); err != nil {
		return err
	}

	return nil
}

var ErrUnexpectedEOR = errors.New("unexpected end of file")

func readFull(fd int, buf *[]byte, exp_bytes int) error {
	for exp_bytes > 0 {
		bytes_read, err := syscall.Read(fd, *buf)
		switch {
		case bytes_read <= 0:
			return ErrUnexpectedEOR
		case err != nil:
			return err
		}

		assert(bytes_read <= exp_bytes)
		exp_bytes -= bytes_read
	}
	return nil
}

func writeAll(fd int, buf *[]byte, exp_bytes int) error {
	for exp_bytes > 0 {
		bytes_write, err := syscall.Write(fd, *buf)
		switch {
		case bytes_write <= 0:
			return ErrUnexpectedEOR
		case err != nil:
			return err
		}

		assert(bytes_write <= exp_bytes)
		exp_bytes -= bytes_write
	}
	return nil
}

func assert(v bool) {
	if !v {
		panic("failure")
	}
}
