package toolstester

import (
	"os"

	"gopkg.in/yaml.v2"
)

var Conf = conf{}

type conf struct {
	Database dbConfig
}

func InitConfig(confPath string) {
	dir, err := os.Getwd()
	if err != nil {
		panic("get current dir err: " + err.Error())
	}
	yamlFile, err := os.ReadFile(dir + confPath)
	if err != nil {
		panic("config file read err: " + err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		panic("yaml config parse error:" + err.Error())
	}

}
