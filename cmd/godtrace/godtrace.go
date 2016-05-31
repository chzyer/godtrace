package main

import (
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

	prog, err := handle.Compile("syscall::read: { printf(\"a-%d\\n\", uid) } tick-1sec { exit(0) }",
		godtrace.ProbeSpecName, godtrace.C_PSPEC, nil)
	if err != nil {
		return err
	}
	info, err := handle.Exec(prog)
	if err != nil {
		return err
	}
	println("matches:", info.Matches())

	if err := handle.Go(); err != nil {
		return err
	}

	for {
		handle.Sleep()
		status := handle.Work()
		if status != godtrace.WS_OKAY {
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
