package plugins

import (
	//"os"
	//"fmt"
	//"io/ioutil"
	"gopkg.in/yaml.v3"
)

type SessionSettings struct {
	Name           string   `yaml:"session_name"`
	StartDirectory string   `yaml:"start_directory"`
	Windows        []Window `yaml:"windows"`
}

// AnsibleInventory defines ansible inventory file struct
type Window struct {
	Name   string                 `yaml:"window_name"`
	Layout string                 `yaml:"layout"`
	Panes  map[string]interface{} `yaml:"panes"`
	//Panes     string     `yaml:"panes"`
}

func Make(config Config){

}

func GenerateYaml() string {
	panes := make(map[string]interface{})
	shellCommandsOptions := []string{"one command", "two commands"}
	panes["shell-command"] = shellCommandsOptions

	sessionSettings := &SessionSettings{
		Name:           "david",
		StartDirectory: "~",
		Windows: []Window{
			{"ss", "fff", panes},
		},
	}

	b, _ := yaml.Marshal(sessionSettings)
  return string(b)

	/*fmt.Println(err)

	    filename:= "/tmp/filename"

	  _, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	  err = ioutil.WriteFile(filename, b, 0644)*/
}

