package command_test

import (
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/diff"
)

func TestDiff_Basic(t *testing.T) {
	result := run.Quick(command.Diff("testdata/a.txt", "testdata/b.txt"))
	assertion.NoError(t, result.Err)
	// Should show differences between the files (diff produces multiple lines)
}

func TestDiff_Unified(t *testing.T) {
	result := run.Quick(command.Diff("testdata/a.txt", "testdata/b.txt", command.Unified))
	assertion.NoError(t, result.Err)
	// Unified diff format
}

func TestDiff_ContextDiff(t *testing.T) {
	result := run.Quick(command.Diff("testdata/a.txt", "testdata/b.txt", command.ContextDiff))
	assertion.NoError(t, result.Err)
	// Context diff format
}

func TestDiff_Brief(t *testing.T) {
	result := run.Quick(command.Diff("testdata/a.txt", "testdata/b.txt", command.Brief))
	assertion.NoError(t, result.Err)
	// Brief output
}

func TestDiff_Identical(t *testing.T) {
	result := run.Quick(command.Diff("testdata/a.txt", "testdata/a.txt"))
	assertion.NoError(t, result.Err)
	// Same file should have no output
}

func TestDiff_MissingFile(t *testing.T) {
	result := run.Quick(command.Diff("nonexistent.txt", "testdata/a.txt"))
	assertion.Error(t, result.Err)
}

