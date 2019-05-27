package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	tmpFileName = "./nri-linux-installer-tmp"
	m           = flag.String("m", "R", "Installation mode for the Agent: (R)oot, (P)rivileged, (U)nprivileged. Default set to root")
	msg         = "Starting installation as"
	mode        string
	modeName    string
	nria        string
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("bad number of argument")
		usageAndExit()
	}

	mode := *m
	if mode != "R" && mode != "P" && mode != "U" {
		fmt.Printf("invalid mode: %s\n", mode)
		usageAndExit()
	}

	switch mode {
	case "R":
		modeName = "root"
	case "P":
		modeName = "PRIVILEGED"
		nria = `NRIA_MODE="PRIVILEGED"`
	case "U":
		modeName = "UNPRIVILEGED"
		nria = `NRIA_MODE="UNPRIVILEGED"`
	}
	fmt.Printf("%s %s user\n", msg, modeName)

	lk := flag.Args()[0]
	err := writeToFile(fmt.Sprintf("license_key: %s", lk), "/etc/newrelic-infra.yml")
	if err != nil {
		os.Exit(1)
	}
	out, err := exec.Command("/bin/sh", "-c", "cat /etc/os-release").Output()
	if err != nil {
		fmt.Printf("ReleaseError: %s\n", err)
		os.Exit(1)
	}

	d, v := parseOSRelease(string(out))

	var o int
	switch d {
	case "centos":
		o = deployRHEL(d, v)
	case "sles":
		o = deploySuse(v)
	case "ubuntu":
		o = deployDeb(d, v)
	case "debian":
		o = deployDeb(d, v)
	case "rhel":
		o = deployRHEL(d, v)
	case "amzn":
		o = deployRHEL(d, v)
	}
	if o == 1 {
		fmt.Printf("Unable to install New Relic Infrastructure agent on %s %s\n", d, v)
	} else {
		fmt.Printf("New Relic Infrastructure agent deployed successfully on %s %s\n", d, v)
	}
}

func usageAndExit() {
	usage := `Usage: nri-linux-installer [options...] <license key>
Options:	
   -mode  Installation mode (default as root). Use R for root, P for Privileged and U for Unprivileged.
`
	fmt.Fprintf(os.Stderr, usage)
	os.Exit(1)
}
