package plugins

import (
	"fmt"
	structs "github.com/yefriddavid/Ec2HostConfigMaker/src/structs"
	"os"
	//"io/ioutil"
	//"gopkg.in/yaml.v3"
)

func MakeSshHostsConfig(instances []structs.Host, sshConfig structs.SshConfig) {
	//targetPath string, identityFileLocation string, template string
	f, _ := os.Create(sshConfig.TargetPathFile)
	f.WriteString(sshConfig.Template)

	defer f.Close()

	//fmt.Println(instances)
	for _, instance := range instances {

        if false {
		    fmt.Println(instance.Identifier)

        }
		//f.WriteString("Host " + instance.Identifier + "\n")
		f.WriteString("Host " + instance.Identifier + " " + instance.PrivateDnsName + "\n")
		f.WriteString("\tHostname " + instance.PublicDnsName + "\n")
		f.WriteString("\tIdentityFile " + sshConfig.IdentityFileLocation + "/" + instance.KeyName + ".pem\n")
		f.WriteString("\n")
	}

}
