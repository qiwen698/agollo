package agollo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

type DataConfig struct {
	Driver  string `yaml:"driver" json:"driver"`
	Host    string `yaml:"host" json:"host"`
	Port    int    `yaml:"port" json:"port"`
	Timeout string `yaml:"timeout" json:"timeout"`
}
type RedisConfig struct {
	Addr     string `yaml:"addr" json:"addr"`
	Password string `yaml:"password" json:"password"`
	DBIndex  int    `yaml:"db_index" json:"db_index"`
}
type DemoYaml struct {
	AppSalt  string      `yaml:"appsalt" json:"appsalt"`
	Redis    RedisConfig `yaml:"redis" json:"redis"`
	Database DataConfig  `yaml:"database" json:"database"`
}

var conf DemoYaml
var mu sync.RWMutex

func TestStartAndUnmarshalOnChange(t *testing.T) {
	newConf := &DemoYaml{}
	// 初始化配置
	err := StartAndUnmarshalOnChange(newConf, func(e *ChangeEvent, err error) {
		if err != nil {
			fmt.Printf("failed to reload config: %s", err.Error())
			return
		}
		if err := reload(*newConf); err != nil {
			fmt.Printf("failed to reconnection thirdparty service: %s", err.Error())
		}
	})
	if err != nil {
		panic(fmt.Sprintf("failed to init apollo: %s", err.Error()))
	}
	err = reload(*newConf)
	if err != nil {
		fmt.Printf("failed to reload config: %s", err)
	}
	http.HandleFunc("/config", IndexHandle)
	err = http.ListenAndServe(":8083", nil)
	if err != nil {
		fmt.Println(err)
	}
}
func IndexHandle(w http.ResponseWriter, r *http.Request) {
	marshal, _ := json.Marshal(conf)
	_, err := w.Write(marshal)
	if err != nil {
		fmt.Println(err)
	}
}

func reload(NewConf DemoYaml) error {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(NewConf)
	conf = NewConf
	fmt.Println("load new config")
	fmt.Println(conf)
	return nil
}
