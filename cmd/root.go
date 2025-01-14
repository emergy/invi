/*
Copyright © 2024 alex.emergy@gmail.com

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
	"os"

	"github.com/emergy/invi/internal/config"
	"github.com/emergy/invi/internal/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "invi",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	commands := viper.Get("commands")
	// 	spew.Dump("commands in rootcmd:", commands)
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	LoadCommands()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// LoadCommands loads all commands from the configuration files
func LoadCommands() {
	for _, command := range viper.Get("cfg").(config.Config).Commands {
		cmd := &cobra.Command{
			Use:   command.Use,
			Short: command.Description,
			Run: func(cmd *cobra.Command, args []string) {
				ctx := map[string]interface{}{
					"args":  args,
					"flags": cmd.Flags(),
					"tasks": command.Tasks,
				}
				err := tasks.RunTasks(ctx)
				if err != nil {
					fmt.Printf("error running tasks: %v\n", err)
					os.Exit(1)
				}
			},
			// Annotations: map[string]string{},
		}

		for flagName, flagAttrs := range command.Flags {
			switch flagAttrs.Type {
			case "string":
				cmd.Flags().StringP(flagName, flagAttrs.Short, flagAttrs.Value, flagAttrs.Description)
			case "bool":
				cmd.Flags().BoolP(flagName, flagAttrs.Short, false, flagAttrs.Description)
			case "int":
				cmd.Flags().IntP(flagName, flagAttrs.Short, 0, flagAttrs.Description)
			default:
				fmt.Printf("unknown flag type: %s\n", flagAttrs.Type)
				os.Exit(1)
			}
		}

		rootCmd.AddCommand(cmd)
	}
}
