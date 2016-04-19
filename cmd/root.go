package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/buoyantio/namerctl/namer"
	"github.com/spf13/cobra"
	//jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var cfgFile string
var baseURLString string

func getBaseURL() (*url.URL, error) {
	if baseURLString == "" {
		baseURLString = viper.GetString("base-url")
	}
	if baseURLString == "" {
		return nil, errors.New("empty base URL")
	}
	u, err := url.Parse(baseURLString)
	if err != nil {
		return nil, err
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
	Short: "namerctl is a command-line client for namerd",
	Long: `namerd manages delegation tables for linkerd.

namerctl looks for a configuration file in the current working
directory or any of its parent directories. Configuration files are
named .namerctl.<ext> where <ext> is describes one of several formats
including yaml, json, toml, etc.  "base-url" is currently the only
supported configuration.  Furthermore, the base url may be specified
via the NAMERCTL_BASE_URL environment variable.

Find more information at https://linkerd.io`,
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
	//jww.SetStdoutThreshold(jww.LevelTrace)
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	RootCmd.PersistentFlags().StringVar(&baseURLString, "base-url", "",
		"namer location (e.g. http://namerd.example.com:4080)")
	viper.BindPFlag("base-url", RootCmd.PersistentFlags().Lookup("base-url"))
}

func addParentConfigPaths(dir string) {
	viper.AddConfigPath(dir + string(os.PathSeparator))
	if sep := strings.LastIndex(dir, string(os.PathSeparator)); sep != -1 {
		addParentConfigPaths(dir[:sep])
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // set on commandline
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigName(".namerctl")
	addParentConfigPaths(os.Getenv("PWD"))
	viper.SetEnvPrefix("namerctl")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
