package main

import "fmt"

func deploySuse(version string) (out int) {
	fmt.Println("Adding gpg key...")
	err := runCmd("sudo rpm --import https://download.newrelic.com/infrastructure_agent/gpg/newrelic-infra.gpg")
	if err != nil {
		fmt.Printf("add gpg Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Downloading agent...")
	err = runCmd(fmt.Sprintf("sudo curl -o /etc/zypp/repos.d/newrelic-infra.repo https://download.newrelic.com/infrastructure_agent/linux/zypp/sles/%s/x86_64/newrelic-infra.repo", version))
	if err != nil {
		fmt.Printf("download Error %s\n", err)
		out = 1
		return
	}

	fmt.Println("Referencing agent...")
	err = runCmd("sudo zypper -n ref -r newrelic-infra")
	if err != nil {
		fmt.Printf("reference Error: %s\n", err)
		out = 1
		return
	}

	if mode == "R" {
		out = suseRootInstall()
	} else {
		out = suseNonRootInstall()
	}
	return
}

func suseRootInstall() (out int) {
	fmt.Println("Installing New Relic agent...")
	if err := runCmd("sudo zypper -n install newrelic-infra"); err != nil {
		fmt.Printf("zypper install Error: %s\n", err)
		out = 1
		return
	}
	return
}

func suseNonRootInstall() (out int) {
	defer func() {
		_ = runCmd("unset NRIA_MODE")
	}()

	fmt.Printf("Running specific steps for %s...\n", modeName)

	if mode == "P" {
		fmt.Println("Installing libcap...")
		err := runCmd(`sudo zypper install libcap-progs`)
		if err != nil {
			fmt.Printf("zypper install Error: %s\n", err)
			out = 1
			return
		}
	}

	fmt.Println("Installing agent...")
	err := runCmd(fmt.Sprintf(`sudo %s zypper install newrelic-infra`, nria))
	if err != nil {
		fmt.Printf("zypper install Error: %s\n", err)
		out = 1
		return
	}

	fmt.Println("Doing zypper update...")
	cmd := fmt.Sprintf(`export %s;sudo -E zypper update newrelic-infra`, nria)
	if err = runCmd(cmd); err != nil {
		fmt.Printf("zypper update wasn't successful. Please try run the command `%s` manually\n", cmd)
	}

	return
}
