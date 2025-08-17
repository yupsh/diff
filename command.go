package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[string, flags]

func Diff(parameters ...any) yup.Command {
	cmd := command(yup.Initialize[string, flags](parameters...))
	if cmd.Flags.UnifiedContext == 0 && bool(cmd.Flags.Unified) {
		cmd.Flags.UnifiedContext = 3
	}
	if cmd.Flags.ContextLines == 0 && bool(cmd.Flags.ContextDiff) {
		cmd.Flags.ContextLines = 3
	}
	return cmd
}

func (p command) Executor() yup.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Need two file paths to compare
		if len(p.Positional) < 2 {
			_, _ = fmt.Fprintf(stderr, "diff: missing operand after '%s'\n", strings.Join(p.Positional, " "))
			return fmt.Errorf("diff requires two files to compare")
		}

		file1Path := p.Positional[0]
		file2Path := p.Positional[1]

		// Read both files
		lines1, err := readFileLines(file1Path)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "diff: %s: %v\n", file1Path, err)
			return err
		}

		lines2, err := readFileLines(file2Path)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "diff: %s: %v\n", file2Path, err)
			return err
		}

		// Check if files are identical
		if areIdentical(lines1, lines2, bool(p.Flags.IgnoreCase), bool(p.Flags.IgnoreWhitespace)) {
			// Files are identical, no output
			return nil
		}

		// Brief mode - just report that files differ
		if bool(p.Flags.Brief) {
			_, _ = fmt.Fprintf(stdout, "Files %s and %s differ\n", file1Path, file2Path)
			return nil
		}

		// Perform diff and output
		if bool(p.Flags.Unified) {
			outputUnifiedDiff(stdout, file1Path, file2Path, lines1, lines2, int(p.Flags.UnifiedContext))
		} else if bool(p.Flags.ContextDiff) {
			outputContextDiff(stdout, file1Path, file2Path, lines1, lines2, int(p.Flags.ContextLines))
		} else {
			outputNormalDiff(stdout, lines1, lines2)
		}

		return nil
	}
}

// readFileLines reads all lines from a file
func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// areIdentical checks if two sets of lines are identical
func areIdentical(lines1, lines2 []string, ignoreCase, ignoreWhitespace bool) bool {
	if len(lines1) != len(lines2) {
		return false
	}

	for i := range lines1 {
		l1, l2 := lines1[i], lines2[i]

		if ignoreWhitespace {
			l1 = strings.TrimSpace(l1)
			l2 = strings.TrimSpace(l2)
		}

		if ignoreCase {
			l1 = strings.ToLower(l1)
			l2 = strings.ToLower(l2)
		}

		if l1 != l2 {
			return false
		}
	}

	return true
}

// outputNormalDiff outputs in normal diff format
func outputNormalDiff(w io.Writer, lines1, lines2 []string) {
	// Simple line-by-line comparison for normal format
	maxLen := len(lines1)
	if len(lines2) > maxLen {
		maxLen = len(lines2)
	}

	for i := 0; i < maxLen; i++ {
		if i >= len(lines1) {
			fmt.Fprintf(w, "%da%d\n", len(lines1), i+1)
			fmt.Fprintf(w, "> %s\n", lines2[i])
		} else if i >= len(lines2) {
			fmt.Fprintf(w, "%dd%d\n", i+1, len(lines2))
			fmt.Fprintf(w, "< %s\n", lines1[i])
		} else if lines1[i] != lines2[i] {
			fmt.Fprintf(w, "%dc%d\n", i+1, i+1)
			fmt.Fprintf(w, "< %s\n", lines1[i])
			fmt.Fprintf(w, "---\n")
			fmt.Fprintf(w, "> %s\n", lines2[i])
		}
	}
}

// outputUnifiedDiff outputs in unified diff format
func outputUnifiedDiff(w io.Writer, file1, file2 string, lines1, lines2 []string, context int) {
	fmt.Fprintf(w, "--- %s\n", file1)
	fmt.Fprintf(w, "+++ %s\n", file2)

	// Simple unified diff implementation
	for i := 0; i < len(lines1) || i < len(lines2); i++ {
		if i >= len(lines1) {
			fmt.Fprintf(w, "+%s\n", lines2[i])
		} else if i >= len(lines2) {
			fmt.Fprintf(w, "-%s\n", lines1[i])
		} else if lines1[i] != lines2[i] {
			fmt.Fprintf(w, "-%s\n", lines1[i])
			fmt.Fprintf(w, "+%s\n", lines2[i])
		} else {
			fmt.Fprintf(w, " %s\n", lines1[i])
		}
	}
}

// outputContextDiff outputs in context diff format
func outputContextDiff(w io.Writer, file1, file2 string, lines1, lines2 []string, context int) {
	fmt.Fprintf(w, "*** %s\n", file1)
	fmt.Fprintf(w, "--- %s\n", file2)

	// Simple context diff implementation
	for i := 0; i < len(lines1) || i < len(lines2); i++ {
		if i >= len(lines1) {
			fmt.Fprintf(w, "+ %s\n", lines2[i])
		} else if i >= len(lines2) {
			fmt.Fprintf(w, "- %s\n", lines1[i])
		} else if lines1[i] != lines2[i] {
			fmt.Fprintf(w, "! %s\n", lines1[i])
			fmt.Fprintf(w, "! %s\n", lines2[i])
		} else {
			fmt.Fprintf(w, "  %s\n", lines1[i])
		}
	}
}
