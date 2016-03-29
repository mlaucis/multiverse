package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	statesPath   = "infrastructure/terraform/states"
	templatePath = "infrastructure/terraform/template"
)

func main() {
	var (
		env    = flag.String("env", "", "Environment name used for isolation")
		region = flag.String("region", "", "AWS region to deploy to")
	)
	flag.Parse()

	log.SetFlags(log.Lshortfile)

	if len(flag.Args()) != 1 {
		log.Fatal("provide command: setup, teardown, update")
	}

	if *env == "" {
		log.Fatal("provide env")
	}

	if *region == "" {
		log.Fatal("provide region")
	}

	if _, err := exec.LookPath("terraform"); err != nil {
		log.Fatal("terraform must be in your PATH")
	}

	var (
		namespace = fmt.Sprintf("%s-%s", *env, *region)
		planFile  = fmt.Sprintf("/tmp/%s.plan", namespace)
		statePath = filepath.Join(statesPath, namespace)
		stateFile = filepath.Join(statePath, fmt.Sprintf("%s.tfstate", namespace))
		varFile   = filepath.Join(statePath, fmt.Sprintf("%s.tfvars", namespace))
		environ   = append(
			os.Environ(),
			fmt.Sprintf("TF_VAR_env=%s", *env),
			fmt.Sprintf("TF_VAR_region=%s", *region),
		)
	)

	switch flag.Args()[0] {
	case "setup":
		if _, err := os.Stat(stateFile); err == nil {
			log.Fatalf("state file already exists: %s", stateFile)
		}

		if err := os.MkdirAll(statePath, 0755); err != nil {
			log.Fatalf("state dir creation failed: %s", err)
		}

		plan := prepareCmd(environ, "plan", []string{
			"-out", planFile,
			"-state", stateFile,
			"-var-file", varFile,
			templatePath,
		}...)

		err := plan.Run()
		if err != nil {
			os.Exit(1)
		}

		fmt.Println("Want to apply the plan? (type 'yes')")
		fmt.Print("(no) |> ")

		response := "no"
		fmt.Scanf("%s", &response)

		if response != "yes" {
			os.Exit(1)
		}

		apply := prepareCmd(environ, "apply", []string{
			"-state", stateFile,
			planFile,
		}...)

		err = apply.Run()
		if err != nil {
			os.Exit(1)
		}
	case "teardown":
		plan := prepareCmd(environ, "plan", []string{
			"-destroy",
			"-out", planFile,
			"-state", stateFile,
			"-var-file", varFile,
			templatePath,
		}...)

		err := plan.Run()
		if err != nil {
			os.Exit(1)
		}

		destroy := prepareCmd(environ, "destroy", []string{
			"-state", stateFile,
			"-var-file", varFile,
			templatePath,
		}...)

		err = destroy.Run()
		if err != nil {
			os.Exit(1)
		}
	case "update":
		if _, err := os.Stat(stateFile); err != nil {
			log.Fatalf("couldn't locate state file: %s", err)
		}

		plan := prepareCmd(environ, "plan", []string{
			"-out", planFile,
			"-state", stateFile,
			"-var-file", varFile,
			templatePath,
		}...)

		err := plan.Run()
		if err != nil {
			os.Exit(1)
		}

		fmt.Println("Want to apply the plan? (type 'yes')")
		fmt.Print("(no) |> ")

		response := "no"
		fmt.Scanf("%s", &response)

		if response != "yes" {
			os.Exit(1)
		}

		apply := prepareCmd(environ, "apply", []string{
			"-state", stateFile,
			planFile,
		}...)

		err = apply.Run()
		if err != nil {
			os.Exit(1)
		}
	default:
		log.Fatalf("'%s' not implemented", flag.Args()[0])
	}
}

func prepareCmd(environ []string, command string, args ...string) *exec.Cmd {
	args = append([]string{command}, args...)

	cmd := exec.Command("terraform", args...)
	cmd.Env = append(
		environ,
		"TF_LOG=TRACE",
		fmt.Sprintf("TF_LOG_PATH=/tmp/%s.log", command),
	)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd
}
