package diff

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"

	localopt "github.com/yupsh/diff/opt"
)

// Flags represents the configuration options for the diff command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Diff creates a new diff command with the given parameters
func Diff(parameters ...any) yup.Command {
	cmd := command(opt.Args[string, Flags](parameters...))
	// Set defaults
	if cmd.Flags.UnifiedContext == 0 && bool(cmd.Flags.Unified) {
		cmd.Flags.UnifiedContext = 3
	}
	if cmd.Flags.ContextLines == 0 && bool(cmd.Flags.ContextDiff) {
		cmd.Flags.ContextLines = 3
	}
	return cmd
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if len(c.Positional) != 2 {
		fmt.Fprintln(stderr, "diff: need exactly 2 files")
		return fmt.Errorf("need exactly 2 files")
	}

	file1Name := c.Positional[0]
	file2Name := c.Positional[1]

	// Read files
	lines1, err := c.readFile(ctx, file1Name, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "diff: %s: %v\n", file1Name, err)
		return err
	}

	// Check for cancellation after reading first file
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	lines2, err := c.readFile(ctx, file2Name, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "diff: %s: %v\n", file2Name, err)
		return err
	}

	// Check for cancellation after reading second file
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Normalize lines if needed
	if bool(c.Flags.IgnoreCase) {
		lines1 = c.normalizeCase(ctx, lines1)
		lines2 = c.normalizeCase(ctx, lines2)
	}
	if bool(c.Flags.IgnoreWhitespace) {
		lines1 = c.normalizeWhitespace(ctx, lines1)
		lines2 = c.normalizeWhitespace(ctx, lines2)
	}

	// Check for cancellation after normalization
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Check if files are identical
	if c.linesEqual(ctx, lines1, lines2) {
		return nil // No output for identical files
	}

	// Brief mode: just report that files differ
	if bool(c.Flags.Brief) {
		fmt.Fprintf(stdout, "Files %s and %s differ\n", file1Name, file2Name)
		return nil
	}

	// Generate diff
	return c.generateDiff(ctx, lines1, lines2, file1Name, file2Name, stdout)
}

func (c command) readFile(ctx context.Context, filename string, stdin io.Reader) ([]string, error) {
	var reader io.Reader

	if filename == "-" {
		reader = stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}

	var lines []string
	scanner := bufio.NewScanner(reader)
	for yup.ScanWithContext(ctx, scanner) {
		lines = append(lines, scanner.Text())
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return lines, err
	}

	return lines, scanner.Err()
}

func (c command) normalizeCase(ctx context.Context, lines []string) []string {
	normalized := make([]string, len(lines))
	for i, line := range lines {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return normalized[:i] // Return partial results on cancellation
			}
		}
		normalized[i] = strings.ToLower(line)
	}
	return normalized
}

func (c command) normalizeWhitespace(ctx context.Context, lines []string) []string {
	normalized := make([]string, len(lines))
	for i, line := range lines {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return normalized[:i] // Return partial results on cancellation
			}
		}
		// Simple whitespace normalization: replace multiple spaces with single space
		fields := strings.Fields(line)
		normalized[i] = strings.Join(fields, " ")
	}
	return normalized
}

func (c command) linesEqual(ctx context.Context, lines1, lines2 []string) bool {
	if len(lines1) != len(lines2) {
		return false
	}
	for i := range lines1 {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return false // Return false on cancellation (assume files differ)
			}
		}
		if lines1[i] != lines2[i] {
			return false
		}
	}
	return true
}

func (c command) generateDiff(ctx context.Context, lines1, lines2 []string, file1Name, file2Name string, output io.Writer) error {
	if bool(c.Flags.Unified) || c.Flags.UnifiedContext > 0 {
		return c.generateUnifiedDiff(ctx, lines1, lines2, file1Name, file2Name, output)
	} else if bool(c.Flags.SideBySide) {
		return c.generateSideBySideDiff(ctx, lines1, lines2, file1Name, file2Name, output)
	} else {
		return c.generateNormalDiff(ctx, lines1, lines2, file1Name, file2Name, output)
	}
}

func (c command) generateNormalDiff(ctx context.Context, lines1, lines2 []string, file1Name, file2Name string, output io.Writer) error {
	// Simple line-by-line diff (ed-style)
	i, j := 0, 0
	lineCount := 0

	for i < len(lines1) || j < len(lines2) {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		lineCount++
		if lineCount%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return err
			}
		}

		if i >= len(lines1) {
			// Only lines2 remaining
			fmt.Fprintf(output, "%da%d\n", i, j+1)
			fmt.Fprintf(output, "> %s\n", lines2[j])
			j++
		} else if j >= len(lines2) {
			// Only lines1 remaining
			fmt.Fprintf(output, "%dd%d\n", i+1, j)
			fmt.Fprintf(output, "< %s\n", lines1[i])
			i++
		} else if lines1[i] == lines2[j] {
			// Lines match
			i++
			j++
		} else {
			// Lines differ
			fmt.Fprintf(output, "%dc%d\n", i+1, j+1)
			fmt.Fprintf(output, "< %s\n", lines1[i])
			fmt.Fprintln(output, "---")
			fmt.Fprintf(output, "> %s\n", lines2[j])
			i++
			j++
		}
	}

	return nil
}

func (c command) generateUnifiedDiff(ctx context.Context, lines1, lines2 []string, file1Name, file2Name string, output io.Writer) error {
	// Simplified unified diff format
	fmt.Fprintf(output, "--- %s\n", file1Name)
	fmt.Fprintf(output, "+++ %s\n", file2Name)

	contextLines := int(c.Flags.UnifiedContext)
	if contextLines == 0 {
		contextLines = 3
	}

	i, j := 0, 0
	lineCount := 0

	for i < len(lines1) || j < len(lines2) {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		lineCount++
		if lineCount%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return err
			}
		}

		if i >= len(lines1) {
			// Only lines2 remaining
			fmt.Fprintf(output, "@@ -%d,0 +%d,%d @@\n", i+1, j+1, len(lines2)-j)
			for ; j < len(lines2); j++ {
				// Check for cancellation in inner loop
				if j%100 == 0 {
					if err := yup.CheckContextCancellation(ctx); err != nil {
						return err
					}
				}
				fmt.Fprintf(output, "+%s\n", lines2[j])
			}
		} else if j >= len(lines2) {
			// Only lines1 remaining
			fmt.Fprintf(output, "@@ -%d,%d +%d,0 @@\n", i+1, len(lines1)-i, j+1)
			for ; i < len(lines1); i++ {
				// Check for cancellation in inner loop
				if i%100 == 0 {
					if err := yup.CheckContextCancellation(ctx); err != nil {
						return err
					}
				}
				fmt.Fprintf(output, "-%s\n", lines1[i])
			}
		} else if lines1[i] == lines2[j] {
			// Lines match - don't output unless in context
			i++
			j++
		} else {
			// Lines differ
			fmt.Fprintf(output, "@@ -%d,1 +%d,1 @@\n", i+1, j+1)
			fmt.Fprintf(output, "-%s\n", lines1[i])
			fmt.Fprintf(output, "+%s\n", lines2[j])
			i++
			j++
		}
	}

	return nil
}

func (c command) generateSideBySideDiff(ctx context.Context, lines1, lines2 []string, file1Name, file2Name string, output io.Writer) error {
	maxLen := len(lines1)
	if len(lines2) > maxLen {
		maxLen = len(lines2)
	}

	for i := 0; i < maxLen; i++ {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return err
			}
		}

		line1 := ""
		line2 := ""

		if i < len(lines1) {
			line1 = lines1[i]
		}
		if i < len(lines2) {
			line2 = lines2[i]
		}

		if line1 == line2 {
			fmt.Fprintf(output, "%-40s   %-40s\n", line1, line2)
		} else {
			fmt.Fprintf(output, "%-40s | %-40s\n", line1, line2)
		}
	}

	return nil
}

func (c command) String() string {
	return fmt.Sprintf("diff %v", c.Positional)
}
