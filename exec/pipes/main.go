package main

import (
	"bytes"
	"context"
	"flag"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	host := flag.String("host", "", "")
	remoteDir := flag.String("dir", "/", "")

	flag.Parse()

	if *host == "" {
		panic("-host cannot be blank")
	}

	r1, w1, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	r2, w2, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sshLs := exec.CommandContext(ctx, "ssh", *host, "ls "+*remoteDir)
	sshLs.Stdout = w1

	shuf := exec.Command("shuf")
	shuf.Stdin = r1
	shuf.Stdout = w2

	var buf bytes.Buffer
	head1 := exec.Command("head", "-n", "1")
	head1.Stdin = r2
	head1.Stdout = &buf

	if err := sshLs.Start(); err != nil {
		panic(err)
	}

	if err := shuf.Start(); err != nil {
		panic(err)
	}

	if err := head1.Start(); err != nil {
		panic(err)
	}

	if err := sshLs.Wait(); err != nil {
		panic(err)
	}

	w1.Close()

	if err := shuf.Wait(); err != nil {
		panic(err)
	}

	w2.Close()

	if err := head1.Wait(); err != nil {
		panic(err)
	}

	if _, err := io.Copy(os.Stdout, &buf); err != nil {
		panic(err)
	}
}
