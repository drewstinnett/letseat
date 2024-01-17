package cmd

import (
	"fmt"
	"os"
	"time"

	// ``"github.com/apex/log"

	"github.com/spf13/cobra"

	"github.com/drewstinnett/go-output-format/v2/gout"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	g       *gout.Client
)

// rootCmd represents the base command when called without any subcommands
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "letseat",
		Short: "Decide what to eat!",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			g, err = gout.NewWithCobraCmd(cmd, nil)
			cobra.CheckErr(err)
			g.SetWriter(os.Stdout)
		},
		// Run: func(cmd *cobra.Command, args []string) { },
	}
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.letseat.yaml)")
	bindRootArgs(cmd)
	cmd.AddCommand(
		newAnalyzeCmd(),
		newLogCmd(),
		newRecommendCommand(),
	)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(newRootCmd().Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

func bindRootArgs(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("diary", "d", "diary.yaml", "diary file")
	cmd.PersistentFlags().StringP("format", "f", "yaml", "Format of the output")
	cmd.PersistentFlags().String("current-date", "", "Assume this as the current date, in the format YYYY-MM-DD")
}

func getCurrentDate(cmd *cobra.Command) time.Time {
	ds, err := cmd.Flags().GetString("current-date")
	if err != nil {
		panic(err)
	}
	if ds == "" {
		return time.Now()
	}
	t, err := time.Parse("2006-01-02", ds)
	if err != nil {
		panic(err)
	}
	return t
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".letseat")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
