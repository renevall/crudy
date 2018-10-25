// Copyright © 2018 René Vallecillo <reneval@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initiale", "initialise", "create"},
	Short:   "Initialize a CRUD aplication",
	Long: `Initialize (crudy init) will create a new application, with a licence
	and the appropiate structure for a CRUD application folling DDD practices and
	
  * If a name is provided, it will be created in the current directory;
  * If no name is provided, the current directory will be assumed;
  * If a relative path is provided, it will be created inside $GOPATH
    (e.g. github.com/spf13/cobra);
  * If an absolute path is provided, it will be created;
  * If the directory already exists but is empty, it will be used.
  
Init will not use an existing directory with contents.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		wd, err := os.Getwd()
		if err != nil {
			er(err)
		}

		var project *Project
		if len(args) == 0 {
			project = NewProjectFromPath(wd)
		} else if len(args) == 1 {
			arg := args[0]
			if arg[0] == '.' {
				arg = filepath.Join(wd, arg)
			}
			if filepath.IsAbs(arg) {
				project = NewProjectFromPath(arg)
			} else {
				project = NewProject(arg)
			}
		} else {
			er("please provide only one argument")
		}

		initializeProject(project)

		fmt.Fprintln(cmd.OutOrStdout(), `Your CRUD application is ready at
`+project.AbsPath()+`
Give it a try by going there and running `+"`go run main.go`."+`
Add a resource to it by running `+"`crudy generate userd`.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.
}

func initializeProject(pro *Project) {

}
