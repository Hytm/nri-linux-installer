package main

import (
	"fmt"
)

var tmpGPGName = "./nri-installer-linux-gpg"

func deployDeb(distribution, version string) (out int) {
	fmt.Println("Installing curl...")
	err := runCmd("sudo apt-get install curl -y")
	if err != nil {
		fmt.Printf("download curl Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Installing apt-transport-https...")
	err = runCmd("sudo apt-get install apt-transport-https -y")
	if err != nil {
		fmt.Printf("download apt-transport-https Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Retrieving gpg key...")
	err = runCmd(fmt.Sprintf("curl -o %s https://download.newrelic.com/infrastructure_agent/gpg/newrelic-infra.gpg", tmpGPGName))
	if err != nil {
		fmt.Printf("download gpg Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Adding gpg key...")
	err = runCmd(fmt.Sprintf("sudo apt-key add %s", tmpGPGName))
	if err != nil {
		fmt.Printf("add gpg Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Clean gpg temp file...")
	err = runCmd(fmt.Sprintf("rm %s", tmpGPGName))
	if err != nil {
		fmt.Printf("clean gpg Error: %s\n", err)
		out = 1
		return
	}

	fmt.Printf("Looking for distribution name based on %s %s...\n", distribution, version)
	n := getDebianDistribName(distribution, version)
	if n == "" {
		fmt.Println("No valid version Error")
		out = 1
		return
	}
	fmt.Printf("Adding to sources.list for %s...\n", n)
	err = writeToFile(fmt.Sprintf(`deb [arch=amd64] https://download.newrelic.com/infrastructure_agent/linux/apt %s main`, n), "/etc/apt/sources.list.d/newrelic-infra.list")
	if err != nil {
		fmt.Printf("source.list Error %s\n", err)
		out = 1
		return
	}

	fmt.Println("Doing apt update...")
	err = runCmd("sudo apt-get update")
	if err != nil {
		fmt.Printf("apt update Error: %s\n", err)
		out = 1
		return
	}

	if mode == "R" {
		out = debRootInstall()
	} else {
		out = debNonRootInstall()
	}
	return
}

func getDebianDistribName(distribution, version string) (v string) {
	switch distribution {
	case "debian":
		v = getDebianVersionName(version)
	case "ubuntu":
		v = getUbuntuVersionName(version)
	}
	return
}

func getDebianVersionName(version string) (v string) {
	switch version {
	case "7":
		v = "wheezy"
	case "8":
		v = "jessie"
	case "9":
		v = "stretch"
	}
	return
}

func getUbuntuVersionName(version string) (v string) {
	mv := version[:2]
	switch mv {
	case "12":
		v = "precise"
	case "14":
		v = "trusty"
	case "16":
		v = "xenial"
	case "18":
		v = "bionic"
	}
	return
}

func debRootInstall() (out int) {
	fmt.Println("Installing New Relic agent...")
	if err := runCmd("sudo apt-get install newrelic-infra -y"); err != nil {
		fmt.Printf("apt install Error: %s\n", err)
		out = 1
		return
	}
	return
}

func debNonRootInstall() (out int) {
	defer func() {
		_ = runCmd("unset NRIA_MODE")
	}()

	fmt.Printf("Running specific steps for %s using %s...\n", modeName, nria)

	if mode == "P" {
		fmt.Println("Installing libcap...")
		err := runCmd(`sudo apt-get install libcap2-bin -y`)
		if err != nil {
			fmt.Printf("apt install Error: %s\n", err)
			out = 1
			return
		}
	}

	fmt.Println("Installing agent...")
	err := runCmd(fmt.Sprintf(`sudo %s apt-get install newrelic-infra -y`, nria))
	if err != nil {
		fmt.Printf("apt install Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Doing dist-upgrade...")
	cmd := fmt.Sprintf(`export %s;sudo -E apt-get dist-upgrade -y`, nria)
	if err = runCmd(cmd); err != nil {
		fmt.Printf("apt dist-upgrade wasn't successful. Please try run the command `%s` manually\n", cmd)
	}

	return
}
