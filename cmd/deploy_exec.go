package cmd

import (
	"fmt"
	"io/ioutil"
	"kool-dev/kool/api"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/environment"
	"os"

	"github.com/spf13/cobra"
)

// KoolDeployExec holds handlers and functions for using Deploy API
type KoolDeployExec struct {
	DefaultKoolService

	kubectl, kool builder.Command

	env     environment.EnvStorage
	apiExec api.ExecCall
}

// NewDeployExecCommand initializes new kool deploy Cobra command
func NewDeployExecCommand(deployExec *KoolDeployExec) *cobra.Command {
	return &cobra.Command{
		Use:   "exec [service]",
		Short: "Executes a command in a service from your deployed application on Kool cloud",
		Args:  cobra.MinimumNArgs(1),
		Run:   DefaultCommandRunFunction(deployExec),
	}
}

// NewKoolDeployExec creates a new pointer with default KoolDeployExec service dependencies
func NewKoolDeployExec() *KoolDeployExec {
	return &KoolDeployExec{
		*newDefaultKoolService(),
		builder.NewCommand("kubectl"),
		builder.NewCommand("kool"),
		environment.NewEnvStorage(),
		api.NewDefaultExecCall(),
	}
}

// Execute runs the deploy exec logic - integrating with Deploy API
func (e *KoolDeployExec) Execute(args []string) (err error) {
	var (
		domain  string
		service string = args[0]
		resp    *api.ExecResponse
	)

	args = args[1:]

	e.Println("kool deploy exec - start")

	if domain = e.env.Get("KOOL_DEPLOY_DOMAIN"); domain == "" {
		err = fmt.Errorf("missing deploy domain (env KOOL_DEPLOY_DOMAIN)")
		return
	}

	e.apiExec.Body().Set("domain", domain)
	e.apiExec.Body().Set("service", service)

	if resp, err = e.apiExec.Call(); err != nil {
		return
	}

	if resp.Token == "" {
		err = fmt.Errorf("failed to generate access credentials to cloud deploy")
		return
	}

	CAPath := fmt.Sprintf("%s/.kool-cluster-CA", os.TempDir())
	if err = ioutil.WriteFile(CAPath, []byte(resp.CA), os.ModePerm); err != nil {
		return
	}

	e.kubectl.AppendArgs("--server", resp.Server)
	e.kubectl.AppendArgs("--token", resp.Token)
	e.kubectl.AppendArgs("--namespace", resp.Namespace)
	e.kubectl.AppendArgs("--certificate-authority", CAPath)
	e.kubectl.AppendArgs("exec", "-i")
	if e.IsTerminal() {
		e.kubectl.AppendArgs("-t")
	}
	e.kubectl.AppendArgs(resp.Path, "--")
	if len(args) == 0 {
		args = []string{"bash"}
	}
	e.kubectl.AppendArgs(args...)

	if e.LookPath(e.kubectl) == nil {
		// the command is available on current PATH, so let's use it
		err = e.Interactive(e.kubectl)
		return
	}

	// we do not have 'kubectl' on current path... let's use a container!
	e.kool.AppendArgs(
		"docker", "--",
		"-v", fmt.Sprintf("%s:%s", CAPath, CAPath),
		"kooldev/toolkit:full",
		e.kubectl.Cmd(),
	)
	e.kool.AppendArgs(e.kubectl.Args()...)

	err = e.Interactive(e.kool)
	return
}
