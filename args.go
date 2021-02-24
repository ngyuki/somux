package main

import (
	"flag"
	"fmt"
	"strings"
)

type setting struct {
	LocalForwards   []forwardDesc
	ReverseForwards []forwardDesc
	Command         []string
	Verbose         bool
}

type forwardDesc struct {
	Bind    string
	Connect string
}

type forwardDescFlags []forwardDesc

func (descs *forwardDescFlags) String() string {
	return "forwarding address"
}

func (descs *forwardDescFlags) Set(value string) error {
	desc := splitAddress(value)
	if desc == nil {
		return fmt.Errorf("invalid forwarding address specification \"%v\"", value)
	}
	*descs = append(*descs, *desc)
	return nil
}

func parseArgs() (*setting, error) {
	var locals forwardDescFlags
	var remotes forwardDescFlags
	var verbose bool

	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "Usage: %s [-L value]... [-R value]... [command args...]\n\n", flag.CommandLine.Name())
		flag.PrintDefaults()
	}

	flag.Var(&locals, "L", "local to remote forwarding (e.g. 8080:127.0.0.1:80)")
	flag.Var(&remotes, "R", "remote to local forwarding (e.g. 9000:127.0.0.1:9000)")
	flag.BoolVar(&verbose, "v", false, "verbose move")
	flag.Parse()

	setting := &setting{
		LocalForwards:   locals,
		ReverseForwards: remotes,
		Command:         flag.Args(),
		Verbose:         verbose,
	}

	return setting, nil
}

func splitAddress(addr string) *forwardDesc {
	arr := strings.Split(addr, ":")
	if len(arr) == 3 {
		arr = append([]string{"127.0.0.1"}, arr...)
	}
	if len(arr) != 4 {
		return nil
	}
	bind := arr[0] + ":" + arr[1]
	connect := arr[2] + ":" + arr[3]
	return &forwardDesc{
		Bind:    bind,
		Connect: connect,
	}
}
