package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/buoyantio/namerctl/namer"
	"github.com/spf13/cobra"
)

var (
	dtabCmd = &cobra.Command{
		Use:   "dtab",
		Short: "Control namer's dtab interface",
		Long:  ``,
	}

	dtabListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List delegation table names",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				ctl, err := getController()
				if err != nil {
					return err
				}
				names, err := ctl.List()
				if err != nil {
					return err
				}
				fmt.Println(strings.Join(names, "\n"))
				return nil

			default:
				return errors.New("list does not take arguments")
			}
		},
	}

	dtabGetPretty = true

	dtabGetCmd = &cobra.Command{
		Use:     "get [name]",
		Aliases: []string{"cat"},
		Short:   "Get a delegation table by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 1:
				ctl, err := getController()
				if err != nil {
					return err
				}
				name := args[0]
				vd, err := ctl.Get(name)
				if err != nil {
					return err
				}
				if dtabGetPretty {
					if vd.Version != namer.Version("") {
						fmt.Printf("# version %s\n", vd.Version)
					}
					fmt.Print(vd.Dtab.Pretty())
				} else {
					fmt.Println(vd.Dtab.String())
				}

				return nil

			default:
				return errors.New("get requires a name argument")
			}
		},
	}

	dtabCreateCmd = &cobra.Command{
		Use:     "create [name] [file]",
		Aliases: []string{"new"},
		Short:   "Create a new delegation table.",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 2:
				ctl, err := getController()
				if err != nil {
					return err
				}
				name := args[0]
				dtabstr, err := readDtabPath(args[1])
				if err != nil {
					return err
				}
				_, err = ctl.Create(name, dtabstr)
				if err != nil {
					return err
				}
				fmt.Printf("Created %s\n", name)
				return nil

			default:
				return errors.New("create requires a name and file path")
			}
		},
	}

	dtabUpdateCmd = &cobra.Command{
		Use:     "update [name] [file]",
		Aliases: []string{"up"},
		Short:   "Update a delegation table.",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 2:
				ctl, err := getController()
				if err != nil {
					return err
				}
				name := args[0]
				dtabstr, err := readDtabPath(args[1])
				if err != nil {
					return err
				}
				versioned, err := ctl.Get(name)
				if err != nil {
					return err
				}
				_, err = ctl.Update(name, dtabstr, versioned.Version)
				if err != nil {
					return err
				}
				fmt.Printf("Updated %s\n", name)
				return nil

			default:
				return errors.New("update requires a name and file path")
			}
		},
	}

	dtabDeleteCmd = &cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"del", "rm"},
		Short:   "Delete a delegation by name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 1:
				ctl, err := getController()
				if err != nil {
					return err
				}
				name := args[0]
				if err = ctl.Delete(name); err != nil {
					return err
				}
				fmt.Printf("Deleted %s\n", name)
				return nil

			default:
				return errors.New("delete requires one argument")
			}
		},
	}
)

func init() {
	dtabCmd.AddCommand(dtabListCmd)

	dtabGetCmd.PersistentFlags().BoolVar(&dtabGetPretty, "pretty", true, "pretty-print dtabs")
	dtabCmd.AddCommand(dtabGetCmd)

	dtabCmd.AddCommand(dtabCreateCmd)

	dtabCmd.AddCommand(dtabUpdateCmd)

	dtabCmd.AddCommand(dtabDeleteCmd)

	RootCmd.AddCommand(dtabCmd)
}

func readDtabPath(path string) (string, error) {
	var file io.Reader
	var err error
	switch path {
	case "":
		return "", errors.New("empty dtab path")
	case "-":
		file = os.Stdin
	default:
		file, err = os.Open(path)
		if err != nil {
			return "", err
		}
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
