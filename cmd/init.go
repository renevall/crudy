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
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
Add a resource to it by running `+"`crudy generate user`.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.
}

func initializeProject(pro *Project) {
	if !exists(pro.AbsPath()) { // If path doesn't yet exist, create it
		err := os.MkdirAll(pro.AbsPath(), os.ModePerm)
		if err != nil {
			er(err)
		}
	} else if !isEmpty(pro.AbsPath()) { // If path exists and is not empty don't use it
		er("Crudy will not create a new project in a non empty directory: " + pro.AbsPath())
	}

	createMainFile(pro)
	createDBFile(pro)
	createConfigFile(pro)
	createRouterFile(pro)
}

func createMainFile(pro *Project) {
	tpl := `{{ comment .copyright }}
	{{if .license}}{{ comment .license }}{{end}}
		package main

		import (
			"log"

		)

		func main() {
			config, err := InitConfig()
			if err != nil {
				log.Fatal("Could not load config")
			}

			db, err := InitDB(config)
			if err != nil {
				log.Fatal("Could not connect to the db")
			}
			defer db.Close()

			env := &model.Env{

			}

			router := router.InitRouter(config, env)
			router.Run(":2323")
    }`

	data := make(map[string]interface{})
	data["copyright"] = copyrightLine()
	data["viper"] = viper.GetBool("useViper")
	data["license"] = pro.License().Header
	data["appName"] = path.Base(pro.Name())

	mainScript, err := executeTemplate(tpl, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(pro.AbsPath(), "main.go"), mainScript)
	if err != nil {
		er(err)
	}
}

func createDBFile(pro *Project) {

}

func createConfigFile(pro *Project) {
	tpl := `package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

//Config struct, must be injected where needed.
type Config struct {
	Secret string
}

//InitConfig reads configuration files
func InitConfig() (*Config, error) {

	viper.SetEnvPrefix("prefix")
	if os.Getenv("Enviroment") == "dev" {
		viper.SetConfigName(".conf")
		viper.SetConfigType("toml")
		viper.AddConfigPath(filepath.Dir(""))
		viper.ReadInConfig()
	} else {
		viper.AutomaticEnv()
	}

	//defaults
	viper.SetDefault("PREFIX_SECRET", "Random string")
	secret := viper.GetString("PREFIX_SECRET")

	return NewConfig(secret), nil
}

//NewConfig Returns a new configuration object
func NewConfig(vals ...interface{}) *Config {
	return &Config{
		Secret: vals[0].(string),
	}
}`

	mainScript, err := executeTemplate(tpl, nil)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(pro.AbsPath(), "config.go"), mainScript)
	if err != nil {
		er(err)
	}
}

func createRouterFile(pro *Project) {

}
