package main

import "fmt"

func deployRHEL(distribution, version string) (out int) {
	fmt.Println("Installing curl...")
	err := runCmd("sudo yum install curl -y")
	if err != nil {
		fmt.Printf("download curl Error: %s\n", err)
		out = 1
		return
	}

	fmt.Printf("Looking for distribution version based on %s %s...\n", distribution, version)
	n := getRHELDistribName(distribution, version)
	if n == "" {
		fmt.Println("No valid version Error")
		out = 1
		return
	}

	fmt.Println("Adding repo...")
	err = runCmd(fmt.Sprintf("sudo curl -o /etc/yum.repos.d/newrelic-infra.repo https://download.newrelic.com/infrastructure_agent/linux/yum/el/%s/x86_64/newrelic-infra.repo", n))
	if err != nil {
		fmt.Printf("add repo Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Enabling repo...")
	err = runCmd("sudo yum -q makecache -y --disablerepo='*' --enablerepo='newrelic-infra'")
	if err != nil {
		fmt.Printf("makecache Error %s\n", err)
		out = 1
		return
	}

	if mode == "R" {
		out = rhelRootInstall()
	} else {
		out = rhelNonRootInstall()
	}
	return
}

func getRHELDistribName(distribution, version string) (v string) {
	v = version
	if distribution == "amzn" {
		switch version {
		case "1":
			v = "6"
		case "2":
			v = "7"
		}
	}
	return
}

func rhelRootInstall() (out int) {
	fmt.Println("Installing New Relic agent...")
	if err := runCmd("sudo yum install newrelic-infra -y"); err != nil {
		fmt.Printf("yum install Error: %s\n", err)
		out = 1
		return
	}
	return
}

func rhelNonRootInstall() (out int) {
	defer func() {
		_ = runCmd("unset NRIA_MODE")
	}()

	fmt.Printf("Running specific steps for %s using %s...\n", modeName, nria)

	if mode == "P" {
		fmt.Println("Installing libcap...")
		err := runCmd(`sudo yum install libcap`)
		if err != nil {
			fmt.Printf("yum install Error: %s\n", err)
			out = 1
			return
		}
	}

	fmt.Println("Installing agent...")
	err := runCmd(fmt.Sprintf(`sudo %s yum install newrelic-infra -y`, nria))
	if err != nil {
		fmt.Printf("yum install Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Doing yum update...")
	cmd := fmt.Sprintf(`export %s;sudo -E yum update -y`, nria)
	if err = runCmd(cmd); err != nil {
		fmt.Printf("yum update wasn't successful. Please try run the command `%s` manually\n", cmd)
	}

	return
}
