package command

type ContextLines int
type UnifiedContext int

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

type flags struct {
	ContextLines     ContextLines
	UnifiedContext   UnifiedContext
	Unified          UnifiedFlag
	ContextDiff      ContextFlag
	Brief            BriefFlag
	IgnoreCase       IgnoreCaseFlag
	IgnoreWhitespace IgnoreWhitespaceFlag
	SideBySide       SideBySideFlag
	Recursive        RecursiveFlag
}

func (c ContextLines) Configure(flags *flags)         { flags.ContextLines = c }
func (u UnifiedContext) Configure(flags *flags)       { flags.UnifiedContext = u }
func (u UnifiedFlag) Configure(flags *flags)          { flags.Unified = u }
func (c ContextFlag) Configure(flags *flags)          { flags.ContextDiff = c }
func (b BriefFlag) Configure(flags *flags)            { flags.Brief = b }
func (i IgnoreCaseFlag) Configure(flags *flags)       { flags.IgnoreCase = i }
func (i IgnoreWhitespaceFlag) Configure(flags *flags) { flags.IgnoreWhitespace = i }
func (s SideBySideFlag) Configure(flags *flags)       { flags.SideBySide = s }
func (r RecursiveFlag) Configure(flags *flags)        { flags.Recursive = r }
