package main

import (
	"context"
	"os/exec"
	"time"
)

var (
	binary = "wg"
)

type Dumper func() ([]byte, error)

func WGDump() (out []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	return exec.CommandContext(ctx, binary, "show", "all", "dump").Output()
}
