package cmd

import (
	"errors"
	"fmt"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/checker"
	"kool-dev/kool/cmd/network"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"strings"
	"testing"
)

// FakeChannelCommand fake command not setting fake variables
// this way, it works inside go routines
type FakeChannelCommand struct {
	builder.FakeCommand
}

// Exec will send the command to shell execution.
func (f *FakeChannelCommand) Exec(args ...string) (outStr string, err error) {
	outStr = "output"
	return
}

func newFakeKoolStatus() *KoolStatus {
	return &KoolStatus{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&builder.FakeCommand{},
		&shell.FakeTableWriter{},
	}
}

func TestNewKoolStatus(t *testing.T) {
	k := NewKoolStatus()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Errorf("unexpected shell.OutputWriter on default KoolStatus instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolStatus instance")
	}

	if _, ok := k.DefaultKoolService.in.(*shell.DefaultInputReader); !ok {
		t.Errorf("unexpected shell.InputReader on default KoolStatus instance")
	}

	if _, ok := k.check.(*checker.DefaultChecker); !ok {
		t.Errorf("unexpected checker.Checker on default KoolStatus instance")
	}

	if _, ok := k.net.(*network.DefaultHandler); !ok {
		t.Errorf("unexpected network.Handler on default KoolStatus instance")
	}

	if _, ok := k.getServicesRunner.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Runner on default KoolStatus instance")
	}

	if _, ok := k.getServiceIDRunner.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Runner on default KoolStatus instance")
	}

	if _, ok := k.getServiceStatusPortRunner.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Runner on default KoolStatus instance")
	}

	if _, ok := k.table.(*shell.DefaultTableWriter); !ok {
		t.Errorf("unexpected shell.TableWriter on default KoolStatus instance")
	}
}

func TestStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesRunner.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDRunner.(*builder.FakeCommand).MockExecOut = "100"
	f.getServiceStatusPortRunner.(*builder.FakeCommand).MockExecOut = "Up About an hour|0.0.0.0:80->80/tcp, 9000/tcp"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := `Service | Running | Ports | State
app | Running | 0.0.0.0:80->80/tcp, 9000/tcp | Up About an hour`

	output := strings.TrimSpace(f.table.(*shell.FakeTableWriter).TableOut)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNotRunningStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesRunner.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDRunner.(*builder.FakeCommand).MockExecOut = "100"
	f.getServiceStatusPortRunner.(*builder.FakeCommand).MockExecOut = "Exited an hour ago"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := `Service | Running | Ports | State
app | Not running |  | Exited an hour ago`

	output := strings.TrimSpace(f.table.(*shell.FakeTableWriter).TableOut)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNoStatusPortStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesRunner.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDRunner.(*builder.FakeCommand).MockExecOut = "100"

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := `Service | Running | Ports | State
app | Not running |  |`

	output := strings.TrimSpace(f.table.(*shell.FakeTableWriter).TableOut)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestNoServicesStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()
	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := "No services found."

	output := fmt.Sprint(f.out.(*shell.FakeOutputWriter).WarningOutput...)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFailedGetServicesStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesRunner.(*builder.FakeCommand).MockError = errors.New("")

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := "No services found."

	output := fmt.Sprint(f.out.(*shell.FakeOutputWriter).WarningOutput...)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestFailedDependenciesStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()
	f.check.(*checker.FakeChecker).MockError = errors.New("")

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("expecting command to exit due to an error.")
	}
}

func TestFailedNetworkStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()
	f.net.(*network.FakeHandler).MockError = errors.New("")

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("expecting command to exit due to an error.")
	}
}

func TestFailedGetServiceIDStatusCommand(t *testing.T) {
	f := newFakeKoolStatus()

	f.getServicesRunner.(*builder.FakeCommand).MockExecOut = "app"
	f.getServiceIDRunner.(*builder.FakeCommand).MockError = errors.New("")

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("expecting command to exit due to an error.")
	}
}

func TestServicesOrderStatusCommand(t *testing.T) {
	f := &KoolStatus{
		*newFakeKoolService(),
		&checker.FakeChecker{},
		&network.FakeHandler{},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
		&FakeChannelCommand{},
		&FakeChannelCommand{},
		&shell.FakeTableWriter{},
	}

	f.getServicesRunner.(*builder.FakeCommand).MockExecOut = `cache
app`

	cmd := NewStatusCommand(f)

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing status command; error: %v", err)
	}

	expected := `Service | Running | Ports | State
app | Not running |  | output
cache | Not running |  | output`

	output := strings.TrimSpace(f.table.(*shell.FakeTableWriter).TableOut)

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}
