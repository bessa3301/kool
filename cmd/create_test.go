package cmd

import (
	"bytes"
	"kool-dev/kool/cmd/builder"
	"kool-dev/kool/cmd/presets"
	"kool-dev/kool/cmd/shell"
	"testing"
)

func newFakeKoolCreate() *KoolCreate {
	return &KoolCreate{
		*newFakeKoolService(),
		&presets.FakeParser{},
		&builder.FakeCommand{},
		*newFakeKoolPreset(),
	}
}

func TestNewKoolCreate(t *testing.T) {
	k := NewKoolCreate()

	if _, ok := k.DefaultKoolService.out.(*shell.DefaultOutputWriter); !ok {
		t.Errorf("unexpected shell.OutputWriter on default KoolCreate instance")
	}

	if _, ok := k.DefaultKoolService.exiter.(*shell.DefaultExiter); !ok {
		t.Errorf("unexpected shell.Exiter on default KoolCreate instance")
	}

	if _, ok := k.DefaultKoolService.in.(*shell.DefaultInputReader); !ok {
		t.Errorf("unexpected shell.InputReader on default KoolCreate instance")
	}

	if _, ok := k.createCommand.(*builder.DefaultCommand); !ok {
		t.Errorf("unexpected builder.Command on default KoolCreate instance")
	}

	if _, ok := k.parser.(*presets.DefaultParser); !ok {
		t.Errorf("unexpected presets.Parser on default KoolCreate instance")
	}
}

func TestNewKoolCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	f.parser.(*presets.FakeParser).MockExists = true
	f.KoolPreset.presetsParser.(*presets.FakeParser).MockExists = true
	f.parser.(*presets.FakeParser).MockCreateCommand = "kool docker create command"

	cmd := NewCreateCommand(f)
	cmd.SetArgs([]string{"laravel", "my-app"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing create command; error: %v", err)
	}

	if !f.parser.(*presets.FakeParser).CalledLoadPresets {
		t.Error("did not call parser.LoadPresets")
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.parser.(*presets.FakeParser).CalledGetCreateCommand {
		t.Error("did not call parser.GetCreateCommand")
	}

	if !f.createCommand.(*builder.FakeCommand).CalledParseCommand {
		t.Error("did not call Parse on KoolCreate.createCommand Command")
	}

	if !f.createCommand.(*builder.FakeCommand).CalledInteractive {
		t.Error("did not call Interactive on KoolCreate.createCommand Command")
	}

	if !f.out.(*shell.FakeOutputWriter).CalledSetWriter {
		t.Error("did not call SetWriter")
	}
}

func TestInvalidPresetCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()
	cmd := NewCreateCommand(f)

	cmd.SetArgs([]string{"invalid", "my-app"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("unexpected error executing preset command; error: %v", err)
	}

	if !f.parser.(*presets.FakeParser).CalledLoadPresets {
		t.Error("did not call parser.LoadPresets")
	}

	if !f.parser.(*presets.FakeParser).CalledExists {
		t.Error("did not call parser.Exists")
	}

	if !f.out.(*shell.FakeOutputWriter).CalledError {
		t.Error("did not call Error")
	}

	expected := "Unknown preset invalid"
	output := f.out.(*shell.FakeOutputWriter).Err.Error()

	if expected != output {
		t.Errorf("expecting error '%s', got '%s'", expected, output)
	}

	if !f.exiter.(*shell.FakeExiter).Exited() {
		t.Error("did not call Exit")
	}
}

func TestNoArgsNewCreateCommand(t *testing.T) {
	f := newFakeKoolCreate()

	cmd := NewCreateCommand(f)
	cmd.SetOut(bytes.NewBufferString(""))

	if err := cmd.Execute(); err == nil {
		t.Error("expecting no arguments error executing create command")
	}
}
