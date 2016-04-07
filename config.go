//配置文件解析

package main

import (
	"fmt"
	"log"
	"strings"

	ini "github.com/widuu/goini"
)

var Conf *ini.Config

//加载配置文件
func init() {
	log.Println("配置文件加载...")
	Conf = ini.SetConfig("conf.ini")
}

func GetValue(node, key string) string {
	if Conf != nil {
		return Conf.GetValue(node, key)
	}
	return ""
}

func DeleteValue(node, key string) bool {
	if Conf != nil {
		return Conf.DeleteValue(node, key)
	}
	return false
}

func SetValue(node, key, v string) bool {
	if Conf != nil {
		return Conf.SetValue(node, key, v)
	}
	return false
}

func GetValues() []map[string]map[string]string {
	if Conf != nil {
		return Conf.ReadList()
	}
	return nil
}

func GetTcpAddr() string {
	return fmt.Sprintf("%s:%s", GetValue("tcp", "bind"), GetValue("tcp", "port"))
}

func GetSioAddr() string {
	return fmt.Sprintf("%s:%s", GetValue("sio", "bind"), GetValue("sio", "port"))
}

func GetZks() []string {
	zks := GetValue("zks", "zkServers")
	return strings.Split(zks, ",")
}

func GetBroker() string {
	return GetValue("broker", "id")
}

func GetNode() string {
	return GetValue("snode", "path")
}
