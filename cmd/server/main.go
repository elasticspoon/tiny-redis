package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/elasticspoon/tiny-redis"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating socket: %v\n", err)
		os.Exit(1)
	}

	val := 1
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, val)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting socket option: %v\n", err)
		os.Exit(1)
	}

	sock_in := syscall.SockaddrInet4{
		Port: 1234,
		Addr: [4]byte{0, 0, 0, 0},
	}

	err = syscall.Bind(fd, &sock_in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error binding address: %v\n", err)
		os.Exit(1)
	}

	err = syscall.Listen(fd, syscall.SOMAXCONN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listening on address: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Started listening on: %v:%v\n", sock_in.Addr, sock_in.Port)

	for {
		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			continue
		}

		redis.DoSomething(nfd)

		if err := syscall.Close(nfd); err != nil {
			fmt.Fprintf(os.Stderr, "error closing connection: %v\n", err)
		}
	}
}
