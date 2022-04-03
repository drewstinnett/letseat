/*
Copyright © 2022 Drew Stinnett <drew@drewlink.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"sort"
	"time"

	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:     "analyze",
	Short:   "Analyze the diary",
	Aliases: []string{"a"},
	Run: func(cmd *cobra.Command, args []string) {
		diaryF, err := cmd.Flags().GetString("diary")
		CheckErr(err, "")

		onlyTakeout, err := cmd.Flags().GetBool("only-takeout")
		CheckErr(err, "")
		onlyDinein, err := cmd.Flags().GetBool("only-dinein")
		CheckErr(err, "")

		// Earliest
		earliestA, err := cmd.Flags().GetString("earliest")
		CheckErr(err, "")
		earliestD, err := letseat.ParseDuration(earliestA)
		CheckErr(err, "")
		earliest := time.Now().Add(-earliestD)

		diary, err := letseat.LoadDiaryWithFile(diaryF, &letseat.DiaryFilter{
			OnlyDineIn:  onlyDinein,
			OnlyTakeout: onlyTakeout,
			Earliest:    &earliest,
		})
		CheckErr(err, "")

		fmt.Printf("Most Popular: %v\n", diary.MostPopularPlace())

		// Find best rated meals
		places := diary.UniquePlaces()
		type kv struct {
			Key       string
			Rating    float64
			LastVisit *time.Time
		}
		var kvs []kv
		for _, place := range places {
			d, err := diary.PlaceDetails(place)
			CheckErr(err, "")
			kvs = append(kvs, kv{Key: place, Rating: d.AverageRating, LastVisit: d.LastVisit})
			// fmt.Printf("%20v %10.1f\n", place, d.AverageRating)
		}

		// Print highest rated
		fmt.Println("Highest Rated")
		sort.Slice(kvs, func(i, j int) bool {
			return kvs[i].Rating > kvs[j].Rating
		})
		for _, i := range kvs {
			fmt.Printf("%20v %10.1f\n", i.Key, i.Rating)
		}

		// Print least recent
		fmt.Println("Last Visited")
		sort.Slice(kvs, func(i, j int) bool {
			return kvs[i].LastVisit.Before(*kvs[j].LastVisit)
		})
		for _, i := range kvs {
			last := time.Now().Sub(*i.LastVisit)
			fmt.Printf("%20v %10v days ago\n", i.Key, int(last.Hours()/24))
		}

		people := diary.PeopleEnhanced()
		fmt.Println("Favorite 3 per person")
		for _, person := range people {
			fmt.Printf("%20v %20v\n", person.Name, person.FavoriteN(3))
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analyzeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// analyzeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	analyzeCmd.Flags().Bool("only-takeout", false, "Only include takeout meals")
	analyzeCmd.Flags().Bool("only-dinein", false, "Only include dine-in meals")
	analyzeCmd.Flags().StringP("earliest", "e", "90d", "Earliest date to include")
}
