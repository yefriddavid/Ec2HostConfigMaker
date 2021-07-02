package main

import "fmt"
import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
	plugins "github.com/yefriddavid/Ec2HostConfigMaker/src/plugins"
	structs "github.com/yefriddavid/Ec2HostConfigMaker/src/structs"
	// structs "Ec2HostConfigMaker/src/structs"
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
	debug               = flag.Bool("debug", false, "debug mode")
)

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

func apply(config structs.Config) {
	sess, _ := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(config.AwsRegion),
		},
		Profile: config.AwsProfile,
	})

	svc := ec2.New(sess)
	instances := getInstances(svc, config.SshConfig.HostPrefix)

	plugins.MakeSshHostsConfig(instances, config.SshConfig)
	plugins.MakeTmuxpSessionsFile(instances, config.TmuxpSessionConfigs)
	//config.makeConfig(instances)
	//config.makeTmuxSessions(instances)

}

func getInstances(svc *ec2.EC2, hostPrefix string) []structs.Host {
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
	// var hosts []structs.Host
	var hosts []structs.Host
	for _, awsInstanceReservations := range awsInstances.Reservations {

    if *debug == true {
      fmt.Println("-------------------------------------------")
    }
		for _, instance := range awsInstanceReservations.Instances {
			if currentInstance == GetArrayKeyValue(instance.Tags, "Name") {
				indexMachine++
			} else {
				indexMachine = 1
			}
			currentInstance = hostPrefix + GetArrayKeyValue(instance.Tags, "Name")
			if *instance.PublicDnsName != "" && instance.KeyName != nil {
				hostIdentifierName := currentInstance + "-" + strconv.Itoa(indexMachine)

        // currentInstance = strings.ReplaceAll(currentInstance, " ", "")
        // hostIdentifierName = strings.ReplaceAll(hostIdentifierName, " ", "")
        if *debug == true {
          fmt.Println(currentInstance)
          fmt.Println(hostIdentifierName)
          fmt.Println(*instance.PublicDnsName)
          fmt.Println("KeyName:", *instance.KeyName)
          fmt.Println(*instance.PrivateDnsName)
          fmt.Println(*instance.KeyName)
        }
				hosts = append(hosts, structs.Host{currentInstance, hostIdentifierName, *instance.PublicDnsName, *instance.KeyName,*instance.PrivateDnsName})

			}
		}

    if *debug == true {
      fmt.Println("-------------------------------------------")
	  }
	}

	return hosts
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

func loadSetting() (config structs.Config, v *viper.Viper, err error) {

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
