/*
Package cmd is the cli app
*/
package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	// special    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(highlight)
	docStyle   = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	/*
		infoStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(true).
				BorderForeground(subtle)
	*/

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	ratingRow  = lipgloss.NewStyle().Width(50)
	ratingKey  = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right).Width(20).PaddingRight(2)
	ratingItem = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right).Width(20)

	listItem      = lipgloss.NewStyle().PaddingLeft(2).Render
	listItemMajor = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Render
)

// analyzeCmd represents the analyze command
func newAnalyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "analyze",
		Short:   "Analyze the diary",
		Aliases: []string{"a"},
		RunE:    runAnalyze,
	}
	bindFilter(cmd)
	return cmd
}

func bindFilter(cmd *cobra.Command) {
	cmd.Flags().Bool("only-takeout", false, "Only include takeout meals")
	cmd.Flags().Bool("only-dinein", false, "Only include dine-in meals")
	cmd.Flags().StringP("earliest", "e", "90d", "Earliest date to include")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	diary := letseat.New(
		letseat.WithFilter(*mustNewEntryFilterWithCmd(cmd)),
		letseat.WithEntriesFile(mustGetCmd[string](*cmd, "diary")),
	)

	// Find best rated meals
	placesDetails := diary.PlaceDetails()

	// Print highest rated
	sort.Sort(placesDetails)

	ratings := []string{listHeader("\nHighest Rated")}
	for _, i := range placesDetails {
		ratings = append(ratings, ratingRow.Render(
			lipgloss.JoinHorizontal(lipgloss.Top, ratingKey.Render(i.Name), ratingItem.Render(letseat.Stars(i.AverageRating, "★"))),
		))
	}
	// Set up styling
	doc := strings.Builder{}
	doc.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(fmt.Sprintf("Most Popular: %v\n", diary.MostPopularPlace())),
		lipgloss.JoinVertical(lipgloss.Left, ratings...),
	))

	// Print least recent
	sort.Slice(placesDetails, func(i, j int) bool {
		return placesDetails[i].LastVisit.Before(*placesDetails[j].LastVisit)
	})

	lvisited := vistedStrings(placesDetails, *cmd)
	doc.WriteString(lipgloss.JoinVertical(lipgloss.Left, lvisited...))
	doc.WriteString("\n\n")

	entries := diary.Entries()
	lists := topList(entries.PeopleEnhanced())
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists...))

	fmt.Fprint(cmd.OutOrStdout(), docStyle.Render(doc.String()))
	return nil
}

func topList(people []letseat.Person) []string {
	lists := make([]string, len(people))
	for idx, person := range people {
		topn := person.FavoriteN(3)

		topx := topn[0:min(len(topn), 3)]
		topxI := make([]string, len(topx)+1)
		topxI[0] = listHeader(person.Name)
		for idx, topxitem := range topx {
			topxI[idx+1] = listItem(topxitem)
		}
		list := lipgloss.JoinVertical(
			lipgloss.Left,
			topxI...,
		)
		lists[idx] = list
	}
	return lists
}

func vistedStrings(pd letseat.PlaceDetails, cmd cobra.Command) []string {
	highlightTop := 3
	lvisited := []string{listHeader("\n\nLast Visited")}
	for i, v := range pd {
		lastD := int(getCurrentDate(&cmd).Sub(*v.LastVisit).Hours() / 24)
		var li string
		if i < highlightTop {
			li = listItemMajor(fmt.Sprintf("%20v %10v days ago", v.Name, lastD))
		} else {
			li = listItem(fmt.Sprintf("%20v %10v days ago", v.Name, lastD))
		}
		lvisited = append(lvisited, li)
	}
	return lvisited
}
