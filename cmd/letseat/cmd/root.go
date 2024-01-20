package cmd

import (
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/charmbracelet/log"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"

	"github.com/drewstinnett/go-output-format/v2/gout"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	g       *gout.Client
	config  configPaths
	version string = "dev"
	verbose bool
)

type configPaths struct {
	ConfigPath string
	ConfigFile string
	DataPath   string
	DataFile   string
}

// rootCmd represents the base command when called without any subcommands
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "letseat",
		Short:         "Decide what to eat!",
		Version:       version,
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
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "use a custom configuration")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "use verbose logging")
	bindRootArgs(cmd)
	cmd.AddCommand(
		newAnalyzeCmd(),
		newLogCmd(),
		newRecommendCommand(),
		newConfigCmd(),
	)

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := newRootCmd().Execute()
	if err != nil {
		slog.Error("exiting", "error", err)
		os.Exit(2)
	}
}

func init() {
	data := path.Join(xdg.DataHome, "letseat")
	if !exists(data) {
		if err := os.MkdirAll(data, 0o700); err != nil {
			panic(err)
		}
	}
	config = configPaths{
		ConfigPath: path.Join(xdg.ConfigHome, "letseat"),
		DataPath:   data,
		DataFile:   path.Join(data, "diary.yaml"),
	}
	cobra.OnInitialize(initConfig)
}

func bindRootArgs(cmd *cobra.Command) {
	// cmd.PersistentFlags().StringP("diary", "d", config.DataFile, "diary file")
	cmd.PersistentFlags().StringP("diary", "d", config.DataFile, "diary file")
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
	// Find home directory.
	// home, err := homedir.Dir()
	// cobra.CheckErr(err)
	opts := log.Options{}
	if verbose {
		opts.Level = log.DebugLevel
	}
	logger := slog.New(
		log.NewWithOptions(os.Stderr, opts),
	)
	slog.SetDefault(logger)

	// Search config in home directory with name ".cli" (without extension).
	// viper.AddConfigPath(home)
	viper.SetConfigName("letseat")
	viper.AddConfigPath(config.ConfigPath)

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		config.ConfigFile = viper.ConfigFileUsed()
		slog.Debug("using config file", "file", config.ConfigFile)
	}
}
