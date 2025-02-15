package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mpppk/twitter/internal/option"
	"github.com/spf13/afero"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func NewRootCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "twitter",
		Short: "CLI for collect tweets",
	}

	configFlag := &option.StringFlag{
		Flag: &option.Flag{
			Name:         "config",
			IsPersistent: true,
			Usage:        "config file (default is $HOME/.config/.twitter.yaml)",
		},
	}

	dbPathFlag := &option.StringFlag{
		Flag: &option.Flag{
			Name:         "db-path",
			ViperName:    "DBPath",
			IsPersistent: true,
			Usage:        "DB file path",
		},
	}

	if err := option.RegisterStringFlag(cmd, configFlag); err != nil {
		return nil, err
	}
	if err := option.RegisterStringFlag(cmd, dbPathFlag); err != nil {
		return nil, err
	}

	var subCmds []*cobra.Command
	for _, cmdGen := range cmdGenerators {
		subCmd, err := cmdGen(fs)
		if err != nil {
			return nil, err
		}
		subCmds = append(subCmds, subCmd)
	}
	cmd.AddCommand(subCmds...)
	cmd.SetOut(os.Stdout)

	return cmd, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd, err := NewRootCmd(afero.NewOsFs())
	if err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".twitter" (without extension).
		viper.AddConfigPath(path.Join(home, ".config"))
		viper.SetConfigName(".twitter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
