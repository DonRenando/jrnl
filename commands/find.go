package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type FindCommand struct {
	options Configuration
	flags   *flag.FlagSet
}

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

// NewFindCommand creates a new command runner for finding entries
func NewFindCommand(config Configuration) *FindCommand {
	findCommand := FindCommand{
		options: config,
		flags:   flag.NewFlagSet("find", flag.ExitOnError),
	}
	return &findCommand
}

// Run the list-tags command
func (f *FindCommand) Run(ctx context.Context, subcommandArgs []string) {
	var tags arrayFlags
	f.flags.Var(&tags, "tag", "Find entries of a specific tag or tags.")
	if !f.flags.Parsed() {
		f.flags.Parse(subcommandArgs)
	}
	index, err := tagMap(f.options.JournalPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	seen := make(map[string]bool)
	for _, tag := range tags {
		if index[tag] != nil {
			for _, entry := range index[tag] {
				if !seen[entry] {
					seen[entry] = true
				}
			}
		}
	}
	output := make([]string, 0)
	for key := range seen {
		output = append(output, fmt.Sprintf("%s/entries/%s.md", f.options.JournalPath, key))
	}
	sort.Strings(output)
	fmt.Fprintln(os.Stdout, strings.Join(output, "\n"))
}
