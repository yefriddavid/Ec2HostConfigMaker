package structs

import (
  "fmt"
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
  Identifiier string
  PublicDnsName string
}

type TmuxSessionsConfig struct {
	TargetPathFile  string `mapstructure:"target-path-file"`
	InstanceName    string `mapstructure:"instance-name"`
	PaneOptions     string `mapstructure:"pane-options"`
	Template        string
	HostPrefix      string `mapstructure:"host-prefix"`
}

type SshConfig struct {
	TargetPathFile       string `mapstructure:"target-path-file"`
	Template             string
	HostPrefix           string `mapstructure:"host-prefix"`
	IdentityFileLocation string `mapstructure:"identity-file-location"`
}

type Config struct {
	TmuxSessionsConfigs []TmuxSessionsConfig `mapstructure:"tmuxp-session-config"`
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

