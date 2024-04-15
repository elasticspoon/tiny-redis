package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	defer syscall.Close(fd)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating socket: %v\n", err)
		os.Exit(1)
	}

	sock_in := syscall.SockaddrInet4{
		Port: 1234,
		Addr: [4]byte{syscall.IN_LOOPBACKNET, 0, 0, 1},
	}

	err = syscall.Connect(fd, &sock_in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting address: %v\n", err)
		os.Exit(1)
	}

	if l, err := syscall.Write(fd, []byte("hello")); err != nil || l < 1 {
		fmt.Fprintf(os.Stderr, "error writing message: %v\n", err)
	}

	buf := make([]byte, 64)
	if l, err := syscall.Read(fd, buf); err != nil || l < 1 {
		fmt.Fprintf(os.Stderr, "error reading message: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("server says: %s\n", buf)
}
