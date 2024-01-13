package cmd

import (
	"fmt"
	"log/slog"
	"os"

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
var rootCmd = &cobra.Command{
	Use:   "letseat",
	Short: "Decide what to eat!",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		g, err = gout.NewWithCobraCmd(cmd, nil)
		cobra.CheckErr(err)
		g.SetWriter(os.Stdout)
	},
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.letseat.yaml)")
	bindRootArgs(rootCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func bindRootArgs(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("diary", "d", "diary.yaml", "diary file")
	cmd.PersistentFlags().StringP("format", "f", "yaml", "Format of the output")
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

func checkErr(err error) {
	if err != nil {
		slog.Error("fatal error occurred", "error", err)
		os.Exit(2)
	}
}
