package mflag

import (
	"github.com/nelhage/go.cli/completion"
	"strings"
)

func completeFlags(cl completion.CommandLine, flags *FlagSet) (completions []string, rest completion.CommandLine) {
	if len(cl) == 0 {
		return nil, cl
	}
	cl = cl[1:]
	var inFlag string
	for len(cl) > 1 {
		w := cl[0]
		if inFlag != "" {
			inFlag = ""
		} else if len(w) > 1 && w[0] == '-' && w != "--" {
			if !strings.Contains(w, "=") {
				var i int
				for i = 0; i < len(w) && w[i] == '-'; i++ {
				}
				inFlag = w[i:]
			}
			if flag := flags.Lookup(inFlag); flag != nil {
				if bf, ok := flag.Value.(boolFlag); ok && bf.IsBoolFlag() {
					inFlag = ""
				}
			}
		} else {
			if w == "--" {
				cl = cl[1:]
			}
			return nil, cl
		}
		cl = cl[1:]
	}

	if inFlag != "" {
		// Complete a flag value. No-op for now.
		return []string{}, nil
	} else if len(cl[0]) > 0 && cl[0][0] == '-' {
		// complete a flag name
		prefix := strings.TrimLeft(cl[0], "-")
		flags.VisitAll(func(f *Flag) {
			for _, name := range f.Names {
				if strings.HasPrefix(name, prefix) {
					completions = append(completions, "-"+name)
				}
			}
		})
		return completions, nil
	}

	if cl[0] == "" {
		flags.VisitAll(func(f *Flag) {
			for _, name := range f.Names {
				completions = append(completions, "-"+name)
			}
		})
	}
	return completions, cl
}

type flagCompleter struct {
	flags *FlagSet
	inner completion.Completer
}

func CompleterWithFlags(flags *FlagSet, completer completion.Completer) completion.Completer {
	return &flagCompleter{
		flags: flags,
		inner: completer,
	}
}

func (c *flagCompleter) Complete(cl completion.CommandLine) []string {
	completions, rest := completeFlags(cl, c.flags)
	if rest != nil {
		if extra := c.inner.Complete(rest); extra != nil {
			completions = append(completions, extra...)
		}
	}

	return completions
}
