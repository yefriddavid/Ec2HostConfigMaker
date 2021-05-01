package main

import "fmt"
import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
	structs "github.com/yefriddavid/Ec2HostConfigMaker/src/structs"
	"os"
	"path"
	"strconv"
	"strings"
)

var Version = "No Provided"
var GitCommit = "No Provided"
var GitShortCommit = "No Provided"
var Date = "No Provided"
var VersionStr = ""

var Author = ""
var Homepage = ""
var ReleaseDate = ""

// var SysConfigFile = ""

var (
	configFile         = flag.String("configFile", "/etc/ConfigRefreshEc2HostMaker.yml", "Path Configuration file")
	showPathConfigFile = flag.Bool("path", false, "Show configuration path file")
	showVersion        = flag.Bool("version", false, "Show version")
)


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

func init() {

	flag.Parse()

	/*if *configFile == "" && SysConfigFile != "" {
		*configFile = SysConfigFile
	}*/

	if fileExists(*configFile) {
		// fmt.Println("File exist")
		return
	}

	// if !fileExists("config.yml") {
	// 	*configFile = "config.yml"
	// } else {
	// 	fmt.Println("File not exist")
	// 	//os.Exit(2)
	// }
	//if !fileExists("/etc/Ec2HostMakerConfig.yml") {
	if !fileExists(*configFile) {
		*configFile = "/etc/ConfigRefreshEc2HostMaker.yml"
	} else {
		fmt.Println("File not exist")
		//os.Exit(2)
	}

}

func main() {
	flag.Parse()

	if *showVersion == true {

		fmt.Println("Commit:", GitCommit)
		fmt.Println("Version:", Version)
		fmt.Println("Date:", Date)
		fmt.Println("Author:", Author)
		return
	}

	config, _, _ := loadSetting()

	apply(config)
}

func apply(config Config) {
	sess, _ := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(config.AwsRegion),
		},
		Profile: config.AwsProfile,
	})

	svc := ec2.New(sess)
	instances := getInstancesV2(svc)
  fmt.Println(instances)

	//config.makeConfig(instances)
	//config.makeTmuxSessions(instances)

}

func getInstances(svc *ec2.EC2) *ec2.DescribeInstancesOutput {
	input := &ec2.DescribeInstancesInput{
	/*InstanceIds: []*string{
	    //aws.String("i-1234567890abcdef0"),
	},*/
	}
	awsInstances, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil
	}

	return awsInstances
}

func getInstancesV2(svc *ec2.EC2) []Host {
	input := &ec2.DescribeInstancesInput{}

	awsInstances, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil
	}

  var currentInstance string
  var indexMachine int
  var hosts []structs.Host
  for _, awsInstanceReservations := range awsInstances.Reservations {
    for _, instance := range awsInstanceReservations.Instances {
        if currentInstance == GetArrayKeyValue(instance.Tags, "Name") {
          indexMachine++
        } else {
          indexMachine = 1
        }
        currentInstance = GetArrayKeyValue(instance.Tags, "Name")
        if *instance.PublicDnsName != "" {
          hostIdentifierName := currentInstance + "-" + strconv.Itoa(indexMachine)
          hosts = append(hosts, Host{currentInstance, hostIdentifierName, *instance.PublicDnsName})

        }
    }
  }

	return hosts
}
func (config Config) makeTmuxSessions(awsInstances *ec2.DescribeInstancesOutput) {

	for _, itemConfig := range config.TmuxSessionsConfigs {
					f, _ := os.Create(itemConfig.TargetPathFile)
					f.WriteString(itemConfig.Template)
		var instanceName string
		var indexMachine int = 0
		for _, awsInstanceReservations := range awsInstances.Reservations {
			for _, instance := range awsInstanceReservations.Instances {
				if GetArrayKeyValue(instance.Tags, "Name") == itemConfig.InstanceName {
					if instanceName == GetArrayKeyValue(instance.Tags, "Name") {
						indexMachine++
					} else {
						indexMachine = 1
					}
					instanceName = GetArrayKeyValue(instance.Tags, "Name")
					if *instance.PublicDnsName != "" {
						hostIdentifierName := instanceName + "-" + strconv.Itoa(indexMachine)
						f.WriteString("      - shell_command:\n")
						f.WriteString("          - ssh " + hostIdentifierName + "\n")
            //f.WriteString("          - autossh " + hostIdentifierName + "\n")
						f.WriteString("          - sudo tail -f /var/log/web.stdout.log\n")
					}
				}
			}
		}
		f.Close()
	}
	if true {
		return
	}

	/*
		  var instanceName string
		  var indexMachine int = 0
			for _, awsInstanceReservations := range awsInstances.Reservations {
				for _, instance := range awsInstanceReservations.Instances {
		      if instanceName == GetArrayKeyValue(instance.Tags, "Name") {
		        indexMachine++
		      } else {
		        indexMachine = 1
		      }
		      instanceName = GetArrayKeyValue(instance.Tags, "Name")
					//availabilityZone := strings.Split(*instance.Placement.AvailabilityZone, "-")

					if *instance.PublicDnsName != "" {

						check(err)
						//hostIdentifierName := instanceKeyName + "-" + availabilityZone[2]
						hostIdentifierName := instanceName + "-" + strconv.Itoa(indexMachine)
		        //fmt.Println(instanceName)
		        fmt.Println(hostIdentifierName)

		        f.WriteString("Host " + config.SshConfig.HostPrefix + hostIdentifierName + "\n")
						f.WriteString("\tHostname " + *instance.PublicDnsName + "\n")
						f.WriteString("\tIdentityFile " + config.SshConfig.IdentityFileLocation + "/" + *instance.KeyName + ".pem\n")
						f.WriteString("\n")

					}

				}

			}*/
}

func (config Config) makeConfig(awsInstances *ec2.DescribeInstancesOutput) {

	f, err := os.Create(config.SshConfig.TargetPathFile)
	f.WriteString(config.SshConfig.Template)

	defer f.Close()

	var instanceName string
	var indexMachine int = 0
	for _, awsInstanceReservations := range awsInstances.Reservations {
		for _, instance := range awsInstanceReservations.Instances {
			if instanceName == GetArrayKeyValue(instance.Tags, "Name") {
				indexMachine++
			} else {
				indexMachine = 1
			}
			instanceName = GetArrayKeyValue(instance.Tags, "Name")
			//availabilityZone := strings.Split(*instance.Placement.AvailabilityZone, "-")

			if *instance.PublicDnsName != "" {

				check(err)
				//hostIdentifierName := instanceKeyName + "-" + availabilityZone[2]
				hostIdentifierName := instanceName + "-" + strconv.Itoa(indexMachine)
				//fmt.Println(instanceName)
				//fmt.Println(hostIdentifierName)

				f.WriteString("Host " + config.SshConfig.HostPrefix + hostIdentifierName + "\n")
				f.WriteString("\tHostname " + *instance.PublicDnsName + "\n")
				f.WriteString("\tIdentityFile " + config.SshConfig.IdentityFileLocation + "/" + *instance.KeyName + ".pem\n")
				f.WriteString("\n")

			}

		}
	}
}

func GetArrayKeyValue(values []*ec2.Tag, keySearch string) string {
	for _, currentValue := range values {
		if *currentValue.Key == keySearch {
			return *currentValue.Value
		}
	}
	return ""
}

type TagsType struct {
	Value string
	Key   string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadSetting() (config Config, v *viper.Viper, err error) {

	v = viper.New()
	if *configFile == "" {
		v.AddConfigPath("./")
	} else {
		dir, file := path.Split(*configFile)
		ext := path.Ext(file)
		var absoluteFileName string
		if ext == "" {
			absoluteFileName = file
		} else {
			absoluteFileName = strings.TrimRight(file, ext)
		}
		v.AddConfigPath(dir)
		v.SetConfigName(absoluteFileName)

	}
	err = v.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	v.Unmarshal(&config)
	return
}
