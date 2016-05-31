package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/chzyer/godtrace"
)

func process() error {
	handle, err := godtrace.Open(0)
	if err != nil {
		return err
	}
	defer handle.Close()

	handle.SetBufSize("4m")

	prog, err := handle.Compile("syscall::read: { printf(\"a-%d-\\n\", uid) } tick-1sec { exit(0) }",
		godtrace.ProbeSpecName, godtrace.C_PSPEC, nil)
	if err != nil {
		return err
	}
	info, err := handle.Exec(prog)
	if err != nil {
		return err
	}
	println("matches:", info.Matches())

	pr, err := handle.ConsumePipe()
	if err != nil {
		return err
	}
	defer pr.Close()

	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			println("-", scanner.Text())
		}
	}()

	if err := handle.Go(); err != nil {
		return err
	}

	for {
		status, err := handle.Run()
		if err != nil {
			return fmt.Errorf("run error: %v", err)
		}
		if !status.IsOK() {
			break
		}
	}
	return nil
}

func main() {
	err := process()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
