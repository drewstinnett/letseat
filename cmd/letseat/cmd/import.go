package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Args:  cobra.ExactArgs(1),
		Short: "import entries from a flat yaml file",
		RunE:  runImport,
	}
	return cmd
}

func runImport(cmd *cobra.Command, args []string) error {
	diary := letseat.New(
		letseat.WithDBFilename(mustGetCmd[string](*cmd, "data")),
	)
	defer dclose(diary)
	eb, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	var entries letseat.Entries
	if yerr := yaml.Unmarshal(eb, &entries); yerr != nil {
		return yerr
	}
	zero := 0
	pb := pbar{
		progress: progress.New(progress.WithDefaultGradient()),
		diary:    diary,
		entries:  entries,
		current:  &zero,
	}
	// fn, _ := os.Open("/tmp/whatever.txt")
	_, rerr := tea.NewProgram(pb, tea.WithInput(os.Stdin)).Run()
	if rerr != nil {
		slog.Error("error running progressbar", "error", rerr)
	}
	return nil
}

type pbar struct {
	progress progress.Model
	diary    *letseat.Diary
	entries  letseat.Entries
	percent  float64
	current  *int
}

// Init satisfies the bubble interface
func (p pbar) Init() tea.Cmd {
	return p.log()
}

// View shows the current bar
func (p pbar) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + p.progress.View() + "\n\n" +
		// pad + fmt.Sprintf("%.2f percent complete", p.percent*100) + "\n\n" +
		pad + fmt.Sprintf("Importing %v of %v", *p.current-1, len(p.entries)) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func (p pbar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		panic("foo")
		return p, tea.Quit

	case tea.WindowSizeMsg:
		p.progress.Width = msg.Width - padding*2 - 4
		if p.progress.Width > maxWidth {
			p.progress.Width = maxWidth
		}
		return p, nil

	case logged:
		if p.progress.Percent() == 1.0 {
			return p, tea.Quit
		}

		// p.state = msg
		p.percent = float64(*p.current) / float64(len(p.entries))
		cmd := p.progress.SetPercent(p.percent)

		return p, tea.Batch(p.log(), cmd)
		// return p, cmd

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := p.progress.Update(msg)
		p.progress = progressModel.(progress.Model)
		return p, cmd

	default:
		return p, nil
	}
}

type logged bool

func (p *pbar) log() tea.Cmd {
	return func() tea.Msg {
		var res bool
		var c int
		if p == nil {
			return logged(false)
		}
		c = min(*p.current, len(p.entries)-1)
		if err := p.diary.Log(p.entries[c]); err != nil {
			slog.Error("error logging entriy")
		}
		*p.current++
		return logged(res)
	}
}
