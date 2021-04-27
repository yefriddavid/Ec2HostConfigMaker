package main

import "fmt"
import (
  "strconv"
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

var Version = ""
var VersionStr = ""
var GitCommit = ""
var GitShortCommit = ""
var Author = ""
var Homepage = ""
var ReleaseDate = ""
var SysConfigFile = ""

var (
	configFile = flag.String("configFile", "", "Path Configuration file")
	showPathConfigFile = flag.Bool("path", false, "Show configuration path file")
)

type Config struct {
	Mode                  string
	HostPrefix            string `mapstructure:"host-prefix"`
	AwsProfile            string `mapstructure:"aws-profile"`
	AwsRegion             string `mapstructure:"aws-region"`
	IdentityFileLocation  string `mapstructure:"identity-file-location"`
	TargetPathFile        string `mapstructure:"target-path-file"`
	Template              string
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

	if *configFile == "" && SysConfigFile != "" {
		*configFile = SysConfigFile
	}

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
	config, _,_ := loadSetting()

  //if true {
    //fmt.Println("aca")
    //fmt.Println(v)
    //return
  //}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(config.AwsRegion),
		},
		Profile: config.AwsProfile,
	})

	svc := ec2.New(sess)
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
		return
	}

	var instanceKeyName string
	f, err := os.Create(config.TargetPathFile)
	f.WriteString(config.Template)

	defer f.Close()
    var indexMachine int = 0
	for _, awsInstanceReservations := range awsInstances.Reservations {
      indexMachine++
		for _, instance := range awsInstanceReservations.Instances {
			instanceKeyName = GetArrayKeyValue(instance.Tags, "Name")
			//availabilityZone := strings.Split(*instance.Placement.AvailabilityZone, "-")

			if *instance.PublicDnsName != "" {

				check(err)
				//hostIdentifierName := instanceKeyName + "-" + availabilityZone[2]
				hostIdentifierName := instanceKeyName + "-" + strconv.Itoa(indexMachine)
				f.WriteString("Host " + config.HostPrefix + hostIdentifierName + "\n")
				f.WriteString("\tHostname " + *instance.PublicDnsName + "\n")
				f.WriteString("\tIdentityFile " + config.IdentityFileLocation + "/" + *instance.KeyName + ".pem\n")
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



