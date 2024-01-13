/*
Package cmd is the cli app
*/
package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

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
var analyzeCmd = &cobra.Command{
	Use:     "analyze",
	Short:   "Analyze the diary",
	Aliases: []string{"a"},
	Run:     runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	bindFilter(analyzeCmd)
}

func bindFilter(cmd *cobra.Command) {
	cmd.Flags().Bool("only-takeout", false, "Only include takeout meals")
	cmd.Flags().Bool("only-dinein", false, "Only include dine-in meals")
	cmd.Flags().StringP("earliest", "e", "90d", "Earliest date to include")
}

func runAnalyze(cmd *cobra.Command, args []string) {
	diaryF, err := cmd.Flags().GetString("diary")
	checkErr(err)

	df, err := letseat.NewDiaryFilterWithCmd(cmd)
	checkErr(err)

	diary, err := letseat.LoadDiaryWithFile(diaryF, df)
	checkErr(err)

	// Set up styling
	doc := strings.Builder{}
	doc.WriteString(titleStyle.Render(fmt.Sprintf("Most Popular: %v\n", diary.MostPopularPlace())))

	// Find best rated meals
	places := diary.UniquePlaces()
	type kv struct {
		Key       string
		Rating    float64
		LastVisit *time.Time
	}
	kvs := make([]kv, len(places))
	for idx, place := range places {
		d, err := diary.PlaceDetails(place)
		checkErr(err)
		kvs[idx] = kv{
			Key:       place,
			Rating:    d.AverageRating,
			LastVisit: d.LastVisit,
		}

	}

	// Print highest rated
	ratings := []string{listHeader("\nHighest Rated")}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Rating > kvs[j].Rating
	})
	for _, i := range kvs {
		// ratings = append(ratings, listItem(fmt.Sprintf("%20v %10.1f", i.Key, i.Rating)))
		// ratings = append(ratings, listItem(fmt.Sprintf("%20v %v", i.Key, letseat.GetStars(i.Rating))))
		stars := letseat.GetStars(i.Rating, "★")
		row := lipgloss.JoinHorizontal(lipgloss.Top, ratingKey.Render(i.Key), ratingItem.Render(stars))
		ratings = append(ratings, ratingRow.Render(row))
	}
	doc.WriteString(lipgloss.JoinVertical(lipgloss.Left, ratings...))

	// Print least recent
	lvisited := []string{listHeader("\n\nLast Visited")}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].LastVisit.Before(*kvs[j].LastVisit)
	})

	highlightTop := 3
	for i, v := range kvs {
		lastD := int(time.Since(*v.LastVisit).Hours() / 24)
		var li string
		if i < highlightTop {
			li = listItemMajor(fmt.Sprintf("%20v %10v days ago", v.Key, lastD))
		} else {
			li = listItem(fmt.Sprintf("%20v %10v days ago", v.Key, lastD))
		}
		lvisited = append(lvisited, li)
	}
	doc.WriteString(lipgloss.JoinVertical(lipgloss.Left, lvisited...))
	doc.WriteString("\n\n")
	people := diary.PeopleEnhanced()
	lists := []string{}
	for _, person := range people {
		topn := person.FavoriteN(3)

		get := min(len(topn), 3)
		topx := topn[0:get]
		topxI := []string{
			listHeader(person.Name),
		}
		for _, topxitem := range topx {
			topxI = append(topxI, listItem(topxitem))
		}
		list := lipgloss.JoinVertical(
			lipgloss.Left,
			topxI...,
		)
		lists = append(lists, list)
	}
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists...))

	fmt.Println(docStyle.Render(doc.String()))
}
