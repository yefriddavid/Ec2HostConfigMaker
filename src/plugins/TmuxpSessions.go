package plugins

import (
	"os"
	"fmt"
	"io/ioutil"
	structs "github.com/yefriddavid/Ec2HostConfigMaker/src/structs"
	"gopkg.in/yaml.v2"
)

type SessionSettings struct {
	Name           string   `yaml:"session_name"`
	StartDirectory string   `yaml:"start_directory"`
	Windows        []Window `yaml:"windows"`
}

type WorkspaceItem struct {
	Hosts         []structs.Host
	SessionConfig structs.TmuxpSessionConfig
}

// AnsibleInventory defines ansible inventory file struct
type Window struct {
	Name   string                   `yaml:"window_name"`
	Layout string                   `yaml:"layout"`
	Panes  []map[string]interface{} `yaml:"panes"`
	//Panes     string     `yaml:"panes"`
}

func MakeTmuxpSessionsFile(instances []structs.Host, configs []structs.TmuxpSessionConfig)  {

	// Make workspaces

	workspaces := map[string]WorkspaceItem{}
	for _, config := range configs {

		workspaceItem := WorkspaceItem{SessionConfig: config}
		for _, instance := range instances {
			if instance.Name == config.InstanceName {
				workspaceItem.Hosts = append(workspaceItem.Hosts, instance)
			}
		}

		workspaces[config.InstanceName] = workspaceItem
	}

	buildTemplate := func(workspace WorkspaceItem) []byte{
		var panes []map[string]interface{}
		var c structs.StructureTemplate

		hosts := workspace.Hosts
		config := workspace.SessionConfig

    // Create Panes
		for _, host := range hosts {
			pane := make(map[string]interface{})
			shellCommandsOptions := []string{"autossh " + host.Identifier, "sudo tail -f /var/log/web.stdout.log"}
			pane["shell_command"] = shellCommandsOptions
			panes = append(panes, pane)
		}

    c.GetConf([]byte(config.Template))

		if len(c.Windows[0].Panes) > 0 {
			c.SessionName = config.InstanceName
			c.StartDirectory = "~"
			c.Windows = []structs.Window{{Name: "Default", Panes: panes}}
		} else {
			c.Windows[0].Panes = panes

		}
		b, _ := yaml.Marshal(c)
		return b //fmt.Println(string(b))

	}

	for _, workspace := range workspaces {
    finalTemplate := buildTemplate(workspace)
		//fmt.Println(string(finalTemplate))

    filename := workspace.SessionConfig.TargetPathFile

    os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	  ioutil.WriteFile(filename, finalTemplate, 0644)
    fmt.Println("Created:", filename)

		//break
	}

	/*for _, config := range configs {
		for _, instance := range instances {
			if instance.Name == config.InstanceName {
				//filteredInstances[instance.Name] = append(filteredInstances[instance.Name], instance)

				createFile(instance)
				break

			}
		}
	}*/

	/*
		  windows := []Window{{"ss", "tiled", panes}}
			sessionSettings := &SessionSettings{
				Name:           "david",
				StartDirectory: "~",
				Windows: windows,
		  }

			b, _ := yaml.Marshal(sessionSettings)
		  fmt.Println(string(b))*/
	//return string(b)

}

func GenerateYaml() string {
	panes := make(map[string]interface{})
	shellCommandsOptions := []string{"ssh command", "two commands"}
	panes["shell_command"] = shellCommandsOptions

	sessionSettings := &SessionSettings{
		Name:           "david",
		StartDirectory: "~",
		Windows:        []Window{
		//{"ss", "fff", panes},
		//{"ss", "fff", panes},
		},
	}

	b, _ := yaml.Marshal(sessionSettings)
	return string(b)

	/*fmt.Println(err)

	    filename:= "/tmp/filename"

	  _, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	  err = ioutil.WriteFile(filename, b, 0644)
    */
}

/*func getConf(c structs.TmuxpSessionConfig) {

    yamlFile, err := ioutil.ReadFile("test.yaml")
    if err != nil {
        //log.Printf("yamlFile.Get err   #%v ", err)
        fmt.Println(err)
    }
    var obj SessionSettings
    //err := yaml.Unmarshal([]byte(c.Template), obj)
    err = yaml.Unmarshal(yamlFile, obj)
    if err != nil {
        //log.Fatalf("Unmarshal: %v", err)
        //fmt.Println(err)
    }
    //fmt.Println(c.Template)

    fmt.Println(obj)
    // fmt.Println(obj)

}*/
