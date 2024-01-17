package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
)

func newRecommendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "recommend",
		Aliases: []string{"rec", "r"},
		Short:   "Recommend some places to eat!",
		RunE:    runRecommend,
	}
	bindFilter(cmd)
	cmd.PersistentFlags().Int("top", 3, "return N number recommendations")
	return cmd
}

func runRecommend(cmd *cobra.Command, args []string) error {
	diary := letseat.New(
		letseat.WithFilter(*mustNewEntryFilterWithCmd(cmd)),
		letseat.WithEntriesFile(mustGetCmd[string](*cmd, "diary")),
	)
	topN := mustGetCmd[int](*cmd, "top")
	placesDetails := diary.PlaceDetails()
	sort.Slice(placesDetails, func(i, j int) bool {
		return placesDetails[i].LastVisit.Before(*placesDetails[j].LastVisit)
	})

	lvisited := placesDetails[0:min(topN, len(placesDetails))]
	doc := strings.Builder{}
	doc.WriteString("# Recommendations Based on Time\n")
	for _, item := range lvisited {
		doc.WriteString(fmt.Sprintf("* %v (%v days ago)\n", item.Name, int(time.Since(*item.LastVisit).Hours()/24)))
	}
	doc.WriteString("\n")

	out, err := getRenderer().Render(doc.String())
	if err != nil {
		return err
	}
	fmt.Print(out)

	return nil
}

func getRenderer() *glamour.TermRenderer {
	r, err := glamour.NewTermRenderer(
		// glamour.WithAutoStyle(),
		glamour.WithStandardStyle("light"),
		glamour.WithWordWrap(80),
	)
	panicIfErr(err)
	return r
}
