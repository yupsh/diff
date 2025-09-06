package opt

// Custom types for parameters
type ContextLines int
type UnifiedContext int

// Boolean flag types with constants
type UnifiedFlag bool

const (
	Unified   UnifiedFlag = true
	NoUnified UnifiedFlag = false
)

type ContextFlag bool

const (
	ContextDiff   ContextFlag = true
	NoContextDiff ContextFlag = false
)

type BriefFlag bool

const (
	Brief   BriefFlag = true
	NoBrief BriefFlag = false
)

type IgnoreCaseFlag bool

const (
	IgnoreCase    IgnoreCaseFlag = true
	CaseSensitive IgnoreCaseFlag = false
)

type IgnoreWhitespaceFlag bool

const (
	IgnoreWhitespace   IgnoreWhitespaceFlag = true
	NoIgnoreWhitespace IgnoreWhitespaceFlag = false
)

type SideBySideFlag bool

const (
	SideBySide   SideBySideFlag = true
	NoSideBySide SideBySideFlag = false
)

type RecursiveFlag bool

const (
	Recursive   RecursiveFlag = true
	NoRecursive RecursiveFlag = false
)

// Flags represents the configuration options for the diff command
type Flags struct {
	ContextLines     ContextLines         // Context lines for context diff (-C)
	UnifiedContext   UnifiedContext       // Context lines for unified diff (-U)
	Unified          UnifiedFlag          // Unified diff format (-u)
	ContextDiff      ContextFlag          // Context diff format (-c)
	Brief            BriefFlag            // Report only whether files differ (-q)
	IgnoreCase       IgnoreCaseFlag       // Ignore case differences (-i)
	IgnoreWhitespace IgnoreWhitespaceFlag // Ignore whitespace differences (-w)
	SideBySide       SideBySideFlag       // Side-by-side format (-y)
	Recursive        RecursiveFlag        // Recursively compare directories (-r)
}

// Configure methods for the opt system
func (c ContextLines) Configure(flags *Flags)         { flags.ContextLines = c }
func (u UnifiedContext) Configure(flags *Flags)       { flags.UnifiedContext = u }
func (u UnifiedFlag) Configure(flags *Flags)          { flags.Unified = u }
func (c ContextFlag) Configure(flags *Flags)          { flags.ContextDiff = c }
func (b BriefFlag) Configure(flags *Flags)            { flags.Brief = b }
func (i IgnoreCaseFlag) Configure(flags *Flags)       { flags.IgnoreCase = i }
func (i IgnoreWhitespaceFlag) Configure(flags *Flags) { flags.IgnoreWhitespace = i }
func (s SideBySideFlag) Configure(flags *Flags)       { flags.SideBySide = s }
func (r RecursiveFlag) Configure(flags *Flags)        { flags.Recursive = r }
