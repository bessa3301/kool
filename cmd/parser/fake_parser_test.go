package parser

import (
	"errors"
	"kool-dev/kool/cmd/builder"
	"testing"
)

func TestFakeParser(t *testing.T) {
	f := &FakeParser{
		MockParsedCommands: []builder.Command{&builder.FakeCommand{}},
	}

	_ = f.AddLookupPath("path")

	if !f.CalledAddLookupPath || len(f.TargetFiles) != 1 {
		t.Error("failed to use mocked AddLookupPath function on FakeParser")
	}

	_ = f.AddLookupPath("path")

	if len(f.TargetFiles) != 2 {
		t.Error("failed to use mocked AddLookupPath function more then once on FakeParser")
	}

	commands, _ := f.Parse("script")

	if !f.CalledParse || len(commands) != 1 {
		t.Error("failed to use mocked Parse function on FakeParser")
	}

	f.MockScripts = []string{"script"}

	scripts, _ := f.GetScripts()

	if !f.CalledGetScripts || len(scripts) != 1 || scripts[0] != "script" {
		t.Error("failed to use mocked GetScripts function on FakeParser")
	}
}

func TestFakeFailedParser(t *testing.T) {
	f := &FakeParser{
		MockParseError:      errors.New("parser error"),
		MockGetScriptsError: errors.New("get scripts error"),
	}

	_, err := f.Parse("script")

	if !f.CalledParse || err == nil {
		t.Error("failed to use mocked failing Parse function on FakeParser")
	}

	_, err = f.GetScripts()

	if !f.CalledGetScripts || err == nil {
		t.Error("failed to use mocked failing GetScripts function on FakeParser")
	}
}
