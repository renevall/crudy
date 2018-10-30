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
	createConfigModel(pro)
	createEnvModel(pro)
}

func createMainFile(pro *Project) {
	tpl := `{{ comment .copyright }}
	{{if .license}}{{ comment .license }}{{end}}
		package main

		import (
			"log"
			{{if .domain}} "{{ .domain }}" {{end}}
			{{if .router}} "{{ .router }}" {{end}}


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
	data["domain"] = filepath.Join(pro.Name(), "/model")
	data["router"] = filepath.Join(pro.Name(), "/router")

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

	tpl := `package main

	import (
		"fmt"
		"log"
		"time"
		
		{{if .domain}} "{{ .domain }}" {{end}}
		"github.com/jinzhu/gorm"
		_ "github.com/jinzhu/gorm/dialects/postgres"
	)
	
	// InitDB starts the DB
	func InitDB(config *model.Config) (*gorm.DB, error) {
		log.Println("Connecting to database")
		cnx := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	
		db, err := gorm.Open("postgres", cnx)
		if err != nil {
			fmt.Println(err)
			return nil, err	
		}

		db.LogMode(true)
		
		// Ping until connection comes alive, docker.
		var dbError error
		maxAttempts := 5
		for attempts := 1; attempts <= maxAttempts; attempts++ {
			dbError = db.DB().Ping()
			if dbError == nil {
				break
			}
			log.Println(dbError)
			time.Sleep(time.Duration(attempts) * time.Second)
		}

		if dbError != nil {
			log.Fatal(dbError)
		}
	
		db.AutoMigrate()
		return db, nil
	}`

	data := make(map[string]interface{})
	data["domain"] = filepath.Join(pro.Name(), "/model")

	mainScript, err := executeTemplate(tpl, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(pro.AbsPath(), "db.go"), mainScript)
	if err != nil {
		er(err)
	}
}

func createConfigFile(pro *Project) {
	tpl := `package main

	import (
		"os"
		"path/filepath"

		{{if .domain}} "{{ .domain }}" {{end}}
		"github.com/spf13/viper"
	)

	//InitConfig reads configuration files
	func InitConfig() (*model.Config, error) {

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
		viper.SetDefault("PREFIX_SECRET", "generatecode")
		viper.SetDefault("PREFIX_DBHOST", "localhost")
		viper.SetDefault("PREFIX_DBUSER", "user")
		viper.SetDefault("PREFIX_DBPASSWORD", "password")
		viper.SetDefault("PREFIX_DBNAME", "sample")
		viper.SetDefault("PREFIX_DBPORT", 5432)

		secret := viper.GetString("PREFIX_SECRET")
		dbhost := viper.GetString("PREFIX_DBHOST")
		dbuser := viper.GetString("PREFIX_DBUSER")
		dbpassword := viper.GetString("PREFIX_DBPASSWORD")
		dbname := viper.GetString("PREFIX_DBNAME")
		dbport := viper.GetInt("PREFIX_DBPORT")

		return NewConfig(secret,dbhost,dbuser,dbpassword,dbname,dbport), nil
	}

	//NewConfig Returns a new configuration object
	func NewConfig(vals ...interface{}) *model.Config {
		return &model.Config{
			Secret:     vals[0].(string),
			DBHost:     vals[1].(string),
			DBUser:     vals[2].(string),
			DBPassword: vals[3].(string),
			DBName:     vals[4].(string),
			DBPort:     vals[5].(int),
		}
	}`

	data := make(map[string]interface{})
	data["domain"] = filepath.Join(pro.Name(), "/model")

	mainScript, err := executeTemplate(tpl, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(pro.AbsPath(), "config.go"), mainScript)
	if err != nil {
		er(err)
	}
}

func createRouterFile(pro *Project) {
	tpl := `package router

 import (
	 "fmt"
	 "net/http"
 
	 {{if .domain}} "{{ .domain }}" {{end}}
	 "github.com/gin-gonic/gin"
 )
 
 // We initialize router and set the basic routes.
 func InitRouter(config *model.Config, env *model.Env) *gin.Engine {
	 router := gin.Default()
	 router.Use(CORSMiddleware())
 
	 // Sample User CRUD
	//  user := router.Group("/user")
	//  {
	// 	 user.POST("/", UserCreateHandler(config, env.User)) //Create User
	// 	 user.GET("/", UserListHandler(config, env.User))    //Read All Users
	// 	 user.GET("/:id", UserFindHandler(config, env.User))
	// 	 user.PATCH("/:id", UserUpdateHandler(config, env.User)) //Update User
	// 	 user.DELETE("/:id", NotImplementedHandler())            //Delete User44
 
	//  }
 
	 return router
 
 }
 
 // NotImplementedHandler is returned when the handler is not done
 func NotImplementedHandler() gin.HandlerFunc {
	 return func(c *gin.Context) {
		 c.JSON(http.StatusNotFound, gin.H{"status": "Fail", "message": "Handler not implemented"})
	 }
 }
 
 // CORSMiddleware set the CORS headers
 func CORSMiddleware() gin.HandlerFunc {
	 return func(c *gin.Context) {
		 c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		 c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		 c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding")
		 c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		 c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
 
		 if c.Request.Method == "OPTIONS" {
			 fmt.Println("OPTIONS")
			 c.AbortWithStatus(200)
		 } else {
			 c.Next()
		 }
	 }
 }
 `
	data := make(map[string]interface{})
	data["domain"] = filepath.Join(pro.Name(), "/model")

	mainScript, err := executeTemplate(tpl, data)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(pro.AbsPath(), "/router", "router.go"), mainScript)
	if err != nil {
		er(err)
	}
}

func createConfigModel(pro *Project) {
	tpl := `package model
	//Config struct, must be injected where needed.
	type Config struct {
		Secret string
		DBHost string
		DBUser string
		DBPassword string
		DBName string
		DBPort int
	}`

	mainScript, err := executeTemplate(tpl, nil)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(pro.AbsPath(), "/model", "config.go"), mainScript)
	if err != nil {
		er(err)
	}
}

func createEnvModel(pro *Project) {
	tpl := `package model
	//Env struct, used to help wiring up depencencies sent to the router
	type Env struct {
		// SampleStore InterfaceName
	}`

	mainScript, err := executeTemplate(tpl, nil)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(filepath.Join(pro.AbsPath(), "/model", "env.go"), mainScript)
	if err != nil {
		er(err)
	}
}
