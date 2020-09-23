package main

import (
	"gopkg.in/yaml.v2"
	"ldy/forward/util"
	"log"
	"os"
)

func decodeYamlFile(path string, conf *util.Conf) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("读取配置文件:%s,失败", path)
	}
	defer file.Close()
	err = yaml.NewDecoder(file).Decode(conf)
	if err != nil {
		log.Fatalf("配置文件:%s，解析失败:%s", path, err)
	}
}

func main() {
	conf := new(util.Conf)
	decodeYamlFile("forward.yaml", conf)
	util.Run(conf)
}
