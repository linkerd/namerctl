package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/buoyantio/namerctl/namer"
	"github.com/spf13/cobra"
)

var cfgFile string
var baseURLString string

func getBaseURL() (*url.URL, error) {
	u, err := url.Parse(baseURLString)
	if err != nil {
		return nil, err
	}
	if baseURLString == "" {
		return nil, errors.New("empty base URL")
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("invalid base URL: " + baseURLString)
	}
	if u.Path != "" {
		return nil, errors.New("base URL may not have a path: " + baseURLString)
	}
	return u, nil
}

func getController() (namer.Controller, error) {
	baseURL, err := getBaseURL()
	if err != nil {
		return nil, err
	}
	ctl := namer.NewHttpController(baseURL, &http.Client{})
	return ctl, nil
}

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "namerctl",
	Short: "namerctl controls the namer delegation management service",
	Long:  `Find more information at https://linkerd.io`,
}

// Execute adds all child commands to the root command sets flags
// appropriately.  This is called by main.main(). It only needs to
// happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.namerctl.yaml)")

	RootCmd.PersistentFlags().StringVar(&baseURLString, "base-url",
		os.Getenv("NAMERCTL_BASE_URL"),
		"namer location (e.g. http://namerd.example.com:4080)")
}

//TODO
// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" { // enable ability to specify config file via flag
// 		viper.SetConfigFile(cfgFile)
// 	}
// 	viper.SetConfigName(".namerctl") // name of config file (without extension)
// 	viper.AddConfigPath("$PWD")     // adding current directory as first search path
// 	viper.AddConfigPath("$HOME")    // adding home directory as second search path
// 	viper.AutomaticEnv()            // read in environment variables that match
// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
