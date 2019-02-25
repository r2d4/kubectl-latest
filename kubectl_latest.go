package main // import "matt-rickard.com/kubectl-latest"

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	binaryName  = "kubectl-latest"
	defaultVerb = "get"
	defaultNoun = "all"
	longUsage   = `Usage: ` + binaryName + `
	Returns the most recently created resource of a particular type.

* All resources are supported.
* Arbitrary flags can be passed to the underlying commands.
* "get" and "describe" are the only kubectl output subcommands supported currently.

Trigger this output:
	kubectl-latest help

Defaults:
	Verb: ` + defaultVerb + `
	Noun: ` + defaultNoun + `

Examples:
	# Return the "get" output of the most recent resource (across all types)
	kubectl-latest get 

	# Return the "get" output of the most recent pod, using with the pod short syntax "po"
	kubectl-latest po
	
	# or equivalently
	kubectl-latest get po

	# Return the logs of the most recently pod
    kubectl-latest logs

	# Returns the "get" output in yaml format of the most recent deployment. 
	# kubectl-latest will pass on arbitrary flags to kubectl
	kubectl-latest deployment -o yaml

	# Return the "describe" output of the most recent service.
	kubectl-latest describe svc`
)

var rootCmd = &cobra.Command{
	Use:                binaryName,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && args[0] == "help" {
			cmd.Usage()
			return nil
		}

		v := defaultVerb
		n := defaultNoun
		vargs := args

		if len(args) == 0 {
			vargs = []string{defaultNoun}
		}

		if vargs[0] == "logs" {
			return runLogs(vargs[1:])
		}

		// If a verb (get or describe) is provided, use that instead of the default
		if vargs[0] == "get" || vargs[0] == "describe" {
			v = args[0]
			vargs = args[1:]
		}

		// If no args (or resource type is provided), default to latest across all resources
		if len(vargs) > 0 {
			n = vargs[0]
			vargs = vargs[1:]
		}

		return run(v, n, vargs)
	},
}

func runLogs(args []string) error {
	n, _, err := latest("pod")
	if err != nil {
		return errors.Wrap(err, "getting latest pod")
	}
	kargs := append([]string{"logs", n}, args...)
	c := exec.Command("kubectl", kargs...)
	out, err := runCmd(c)
	if err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Println(strings.TrimSpace(string(out)))
	return nil
}

func run(v, n string, args []string) error {
	name, k, err := latest(n)
	if err != nil {
		return errors.Wrapf(err, "getting latest %s", n)
	}
	kargs := append([]string{v, fmt.Sprintf("%s/%s", k, name)}, args...)
	c := exec.Command("kubectl", kargs...)
	out, err := runCmd(c)
	if err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Println(strings.TrimSpace(string(out)))
	return nil
}

func latest(noun string) (string, string, error) {
	out, err := runCmd(exec.Command("kubectl", latestArgs(noun)...))
	if err != nil {
		return "", "", errors.Wrapf(err, "getting latest %s", noun)
	}
	nk := strings.Split(strings.TrimSpace(string(out)), " ")
	return nk[0], nk[1], nil
}

func latestArgs(noun string) []string {
	var args []string
	args = append(args, "get")
	args = append(args, noun)
	args = append(args, "--sort-by={.metadata.creationTimestamp}")
	args = append(args, "-o=go-template")
	args = append(args, `--template={{$noun := "" }}{{range .items}}{{$noun = (printf "%s %s" .metadata.name .kind)}}{{end}}{{printf "%s" $noun}}`)
	return args
}

// util

func main() {
	rootCmd.SetUsageFunc(func(*cobra.Command) error { fmt.Println(longUsage); return nil })
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func runCmd(cmd *exec.Cmd) ([]byte, error) {
	logrus.Debugf("command: %s", cmd.Args)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrapf(err, "starting command %v", cmd)
	}

	stdout, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		return nil, err
	}

	stderr, err := ioutil.ReadAll(stderrPipe)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return stdout, errors.Wrapf(err, "running %s: stdout %s, stderr: %s, err: %v", cmd.Args, stdout, stderr, err)
	}

	if len(stderr) > 0 {
		logrus.Debugf("Command output: [%s], stderr: %s", stdout, stderr)
	} else {
		logrus.Debugf("Command output: [%s]", stdout)
	}

	return stdout, nil
}
