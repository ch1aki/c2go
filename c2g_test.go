package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRun_versionFlag(t *testing.T) {
	inStream, outStream, errStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{inStream: inStream, outStream: outStream, errStream: errStream}
	args := strings.Split("c2g -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("c2g version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("Output=%q, want %q", errStream.String(), expected)
	}
}

func TestRun_parseError(t *testing.T) {
	inStream, outStream, errStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{inStream: inStream, outStream: outStream, errStream: errStream}
	args := strings.Split("c2g -dr.peppar", " ")

	status := cli.Run(args)
	if status != ExitCodeParseFlagError {
		t.Errorf("expected %d to eq %d", status, ExitCodeParseFlagError)
	}
}

func TestRun_printGraph(t *testing.T) {
	stdin := bytes.NewBufferString(" 452 64.242.88.10\n" +
		" 270 10.0.0.153\n" +
		"  51 h24-71-236-129.ca.shawcable.net\n" +
		"  44 cr020r01-3.sac.overture.com\n" +
		"  32 h24-70-69-74.ca.shawcable.net\n")
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{inStream: stdin, outStream: outStream, errStream: errStream}
	args := strings.Split("c2g", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

	expected := []byte(" 452 [||||||||||||||||||||||||||||||] 64.242.88.10\n" +
		" 270 [|||||||||||||||||             ] 10.0.0.153\n" +
		"  51 [|||                           ] h24-71-236-129.ca.shawcable.net\n" +
		"  44 [||                            ] cr020r01-3.sac.overture.com\n" +
		"  32 [||                            ] h24-70-69-74.ca.shawcable.net\n")
	if bytes.Compare(expected, outStream.Bytes()) != 0 {
		t.Errorf("Output=%q, want %q", errStream.String(), expected)
	}
}