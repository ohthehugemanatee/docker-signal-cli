package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestCollectMessages(t *testing.T) {
	cmd := exec.Command("cat", "./testoutput.txt")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Errorf("Could not get stdOut from cat command.")
	}
	buffer := bytes.Buffer{}
	FilterMessages(stdout, &buffer)

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}

	got := buffer.String()
	expected := "Send attachment id 6190387717992389667"
	if got != expected {
		t.Errorf("Got %q but expected %q", got, expected)
	}
}
