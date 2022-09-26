package clip

import (
	"fmt"
	"os"
)

type Flag struct {
	Command        *Command
	Word           string
	Aliases        []string // i and interactive
	Description    string
	Default        string
	Type           string // string, int, bool, float
	MultiValued    bool
	SpaceSeparated bool // only true if MultiValued is true
	Changed        bool

	Process func(flag *Flag, word string, args []string) error
}

type Command struct {
	Word            string
	Aliases         []string
	Short           string
	Description     string
	Use             []string
	Indent          int // First line indent
	Margin          int // 2+ line indent
	FlagIndent      int
	OptionHeader    bool
	SupportHelp     bool
	SupportQuestion bool

	flags []Flag
}

func (cmd *Command) Docs() string {
	str := fmt.Sprintf("NAME\n")
	return str
}

func (cmd *Command) Usage() {
}

func (cmd *Command) Execute() error {
	return cmd.ExecuteWithArgs(os.Args[1:])
}

func (cmd *Command) ExecuteWithArgs(args []string) error {
	return nil
}
