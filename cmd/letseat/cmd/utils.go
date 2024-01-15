package cmd

import (
	"fmt"
	"reflect"
	"time"

	letseat "github.com/drewstinnett/letseat/pkg"
	"github.com/spf13/cobra"
)

// mustGetCmd uses generics to get a given flag with the appropriate Type from a cobra.Command
func mustGetCmd[T []int | []string | int | string | bool | time.Duration](cmd cobra.Command, s string) T {
	switch any(new(T)).(type) {
	case *int:
		item, err := cmd.Flags().GetInt(s)
		panicIfErr(err)
		return any(item).(T)
	case *string:
		item, err := cmd.Flags().GetString(s)
		panicIfErr(err)
		return any(item).(T)
	case *bool:
		item, err := cmd.Flags().GetBool(s)
		panicIfErr(err)
		return any(item).(T)
	case *[]int:
		item, err := cmd.Flags().GetIntSlice(s)
		panicIfErr(err)
		return any(item).(T)
	case *[]string:
		item, err := cmd.Flags().GetStringSlice(s)
		panicIfErr(err)
		return any(item).(T)
	case *time.Duration:
		item, err := cmd.Flags().GetDuration(s)
		panicIfErr(err)
		return any(item).(T)
	default:
		panic(fmt.Sprintf("unexpected use of mustGetCmd: %v", reflect.TypeOf(s)))
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func mustNewEntryFilterWithCmd(cmd *cobra.Command) *letseat.EntryFilter {
	got, err := newEntryFilterWithCmd(cmd)
	if err != nil {
		panic(err)
	}
	return got
}

func newEntryFilterWithCmd(cmd *cobra.Command) (*letseat.EntryFilter, error) {
	earliestD, err := letseat.ParseDuration(mustGetCmd[string](*cmd, "earliest"))
	if err != nil {
		return nil, err
	}
	return &letseat.EntryFilter{
		OnlyTakeout: mustGetCmd[bool](*cmd, "only-takeout"),
		OnlyDineIn:  mustGetCmd[bool](*cmd, "only-dinein"),
		Earliest:    toPTR(getCurrentDate(cmd).Add(-earliestD)),
	}, nil
}

func toPTR[V any](v V) *V {
	return &v
}
