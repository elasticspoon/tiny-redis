package redis

import (
	"encoding/binary"
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

	len := binary.LittleEndian.Uint32(buf[:4])
	if len > MAX_MESSAGE {
		return errors.New("message too long")
	}

	headlessArr := buf[4:]
	if err := readFull(fd, &headlessArr, int(len)); err != nil {
		return err
	}

	// NOTE: should we be appending a null '\0' bytes at
	// the end of this array or not?
	fmt.Printf("client says: %s", headlessArr)

	return nil
}

var ErrUnexpectedEOR = errors.New("unexpected end of file")

// readFull takes in a socket num, a pointer to a buffer and an expected number of bytes
// it keeps reading in from the socket to the buffer until exp_bytes are read
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

// writeAll takes in a socket num, a pointer to a buffer and an expected number of bytes
// it keeps writing from the buffer to the socket until exp_bytes are read
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
