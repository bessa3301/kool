package cmd

import (
	"bytes"
	"errors"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/shell"
	"kool-dev/kool/environment"
	"testing"
)

func newFakeKoolExec() *KoolExec {
	return &KoolExec{
		*newFakeKoolService(),
		&KoolExecFlags{[]string{}, false},
		&shell.FakeTerminalChecker{MockIsTerminal: true},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{},
	}
}

func newFailedFakeKoolExec() *KoolExec {
	return &KoolExec{
		*newFakeKoolService(),
		&KoolExecFlags{[]string{}, false},
		&shell.FakeTerminalChecker{MockIsTerminal: true},
		environment.NewFakeEnvStorage(),
		&builder.FakeCommand{MockError: errors.New("error exec")},
	}
}

func TestNewKoolExec(t *testing.T) {
	k := NewKoolExec()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Errorf("unexpected shell.OutputWriter on default KoolExec instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolExec instance")
	}

	if _, ok := k.DefaultKoolService.in.(*shell.DefaultInputReader); !ok {
		t.Errorf("unexpected shell.InputReader on default KoolExec instance")
	}

	if k.Flags == nil {
		t.Errorf("Flags not initialized on default KoolExec instance")
	} else {
		if len(k.Flags.EnvVariables) > 0 {
			t.Errorf("bad default value for EnvVariables flag on default KoolExec instance")
		}

		if k.Flags.Detach {
			t.Errorf("bad default value for Detach flag on default KoolExec instance")
		}
	}

	if _, ok := k.composeExec.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Command on default KoolExec instance")
	}
}

func TestNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !f.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Error("did not call SetWriter")
	}

	if !f.composeExec.(*builder.FakeCommand).CalledInteractive {
		t.Error("did not call Interactive on KoolExec.composeExec Command")
	}

	interactiveArgs := f.composeExec.(*builder.FakeCommand).ArgsInteractive

	if len(interactiveArgs) != 2 || interactiveArgs[0] != "service" || interactiveArgs[1] != "command" {
		t.Error("bad arguments to Interactive on KoolExec.composeExec Command")
	}
}

func TestNoArgsNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetOut(bytes.NewBufferString(""))

	if err := cmd.Execute(); err == nil {
		t.Error("expecting no arguments error executing exec command")
	}
}

func TestKoolUserEnvNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	f.envStorage.(*environment.FakeEnvStorage).Envs["KOOL_ASUSER"] = "user_testing"

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !f.composeExec.(*builder.FakeCommand).CalledAppendArgs {
		t.Error("did not call AppendArgs on KoolExec.composeExec Command")
	}

	argsAppend := f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 2 || argsAppend[0] != "--user" || argsAppend[1] != "user_testing" {
		t.Error("bad arguments to KoolExec.composeExec Command with KOOL_USER environment variable")
	}
}

func TestEnvFlagNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"--env=VAR_TEST=1", "service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !f.composeExec.(*builder.FakeCommand).CalledAppendArgs {
		t.Error("did not call AppendArgs on KoolExec.composeExec Command")
	}

	argsAppend := f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 2 || argsAppend[0] != "--env" || argsAppend[1] != "VAR_TEST=1" {
		t.Errorf("bad arguments to KoolExec.composeExec Command with EnvVariables flag")
	}
}

func TestDetachFlagNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"--detach", "service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !f.composeExec.(*builder.FakeCommand).CalledAppendArgs {
		t.Error("did not call AppendArgs on KoolExec.composeExec Command")
	}

	argsAppend := f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 1 || argsAppend[0] != "--detach" {
		t.Errorf("bad arguments to KoolExec.composeExec Command with Detach flag")
	}
}

func TestFailingNewExecCommand(t *testing.T) {
	f := newFailedFakeKoolExec()
	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("expecting command to exit due to an error.")
	}

	if err := f.out.(*shell.FakeOutputWriter).Err; err.Error() != "error exec" {
		t.Errorf("expecting error 'error exec', got '%s'", err.Error())
	}
}

func TestNonTerminalNewExecCommand(t *testing.T) {
	f := newFakeKoolExec()
	f.terminal.(*shell.FakeTerminalChecker).MockIsTerminal = false

	cmd := NewExecCommand(f)

	cmd.SetArgs([]string{"service", "command"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing exec command; error: %v", err)
	}

	argsAppend := f.composeExec.(*builder.FakeCommand).ArgsAppend

	if len(argsAppend) != 1 || argsAppend[0] != "-T" {
		t.Errorf("bad arguments to KoolExec.composeExec Command on non terminal environment")
	}
}
