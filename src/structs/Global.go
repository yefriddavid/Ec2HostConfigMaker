package structs

import (
  "fmt"
	// "io/ioutil"
	"gopkg.in/yaml.v3"
	//"flag"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/ec2"
	// "github.com/spf13/viper"
	"os"
	//"path"
	// "strconv"
	// "strings"
)

type Host struct {
  Name string
  Identifier string
  PublicDnsName string
  KeyName string
  PrivateDnsName string
}

type Window struct {
  Name string `yaml:"window_name"`
  Layout string `yaml:"layout" default:"tiled"`
  //Panes map[string]interface{} `yaml:"panes"`
  //Panes []map[string]interface{} `yaml:"panes"`
  //Panes map[string]interface{} `yaml:"panes"`
	Panes  []map[string]interface{} `yaml:"panes"`
  OptionsAfter interface{} `yaml:"options_after,omitempty"`


}
type StructureTemplate struct {
	SessionName     string `yaml:"session_name"`
	StartDirectory  string `yaml:"start_directory"`
  Windows         []Window `yaml:"windows"`
}

type TmuxpSessionConfig struct {
	TargetPathFile  string `mapstructure:"target-path-file"`
	InstanceName    string `mapstructure:"instance-name"`
	PaneOptions     string `mapstructure:"pane-options"`
	Template        string
	StructureTemplate StructureTemplate
	HostPrefix      string `mapstructure:"host-prefix"`
}

type SshConfig struct {
	TargetPathFile       string `mapstructure:"target-path-file"`
	Template             string
	HostPrefix           string `mapstructure:"host-prefix"`
	IdentityFileLocation string `mapstructure:"identity-file-location"`
}

type Config struct {
	TmuxpSessionConfigs []TmuxpSessionConfig `mapstructure:"tmuxp-session-config"`
	SshConfig           SshConfig            `mapstructure:"ssh-config"`
	Mode                string
	//HostPrefix            string `mapstructure:"host-prefix"`
	AwsProfile string `mapstructure:"aws-profile"`
	AwsRegion  string `mapstructure:"aws-region"`
	//IdentityFileLocation  string `mapstructure:"identity-file-location"`
	//TargetPathFile        string `mapstructure:"target-path-file"`
	//Template              string
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(err)
		return false
	}
	return !info.IsDir()
}

func (c *StructureTemplate) GetConf(template []byte) *StructureTemplate {

    /*yamlFile, err := ioutil.ReadFile("test.yaml")
    if err != nil {
        //log.Printf("yamlFile.Get err   #%v ", err)
        fmt.Println(err)

    }*/
    err := yaml.Unmarshal(template, c)
    // err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        //log.Fatalf("Unmarshal: %v", err)
        fmt.Println(err)
        return nil
    }
    //fmt.Println(c.SessionName)
    //fmt.Println(c)

    return c
}
