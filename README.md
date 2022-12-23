# Fork 自 github.com/philchia/agollo 的 Apollo 配置中心客户端

## 功能

* 多 namespace 支持
* 容错，本地缓存
* 实时更新通知
* 支持Umarshal

## 依赖

**go 1.9** 或更新

## 安装

```sh
    go get -u github.com/qiwen698/agollo
```

## 使用（仅介绍Fork后的主要改动）

### 1, 启动

#### 启动客户端, 该方法将按顺序依次从`环境变量`、`CLI命令行参数`、执行文件所处目录的`agollo.json`文件中读取apollo的连接配置。

```golang
    agollo.Start()
```

#### 使用自定义结构启动

```golang
    agollo.StartAndUnmarshal(v interface{})
```

#### 使用自定义结构启动, 并在成功后自动进行配置热更新
```golang
    agollo.StartAndUnmarshalOnChange(v interface{}, run func(*ChangeEvent, error))
```

### 2, 配置说明

#### 环境变量
```shell
agollo_appid
# (string) Apollo的AppID

agollo_cluster
# (string) Apollo的群组配置

agollo_namespaces
# (string) Apollo的命名空间配置, 多个用半角逗号","分隔

agollo_ip
# (string) Apollo服务地址

agollo_tagname
# (string) 配置中struct的标签名, 默认为: config

agollo_onenamespacemode
# (string) 是否为单命名空间模式, yes 或 1 时开启
```

#### 命令行参数
```shell
-agollo_appid
# (string) Apollo的AppID

-agollo_cluster
# (string) Apollo的群组配置

-agollo_namespaces
# (string) Apollo的命名空间配置, 多个用半角逗号","分隔

-agollo_ip
# (string) Apollo服务地址

-agollo_tagname
# (string) 配置中struct的标签名, 默认为: config

-agollo_onenamespacemode
# (bool) 是否为单命名空间模式, false 或 true
```

#### agollo.json
```json
{
  "appId": "shareapi",
  "cluster" : "default",
  "namespaceNames" : ["app.yaml"],
  "ip" : "192.168.67.117:8080",
  "oneNamespaceMode" : true,
  "tagname" : "config"
}
```

## 代码示例：
```golang
package srv

import (
	"fmt"
	"log"

	"github.com/qiwen698/agollo"
)

// defDBConfig 数据库配置结构
type defDBConfig struct {
	Host string `config:"host"` //数据库地址
	Port int `config:"port"` //数据库端口
	Dbname string `config:"dbname"` //数据库库名
	User string `config:"user"` //数据库用户名
	Pass string `config:"pass"` //数据库密码
	Charset string `config:"charset"` //数据库字符集
	MaxIdle int `config:"max_idle"` //最大闲置连接数
	MaxConnection int `config:"max_conncetion"` //数据库最大连接数
}

// defRedisConfig Redis配置
type defRedisConfig struct {
	Addr string `config:"addr"` // Redis连接地址
	Password string `config:"password"` // Redis认证密码密码
	DBIndex int `config:"dbindex"` // Redis数据库索引
}

// defConfig 全局配置
type defConfig struct {
	Port int `config:"port"`
	WebDB *defDBConfig `config:"web_db"` // Web库配置
	Redis *defRedisConfig `config:"redis"`
}

var Config = &defConfig{}

func init() {
	err := agollo.StartAndUnmarshalOnChange(Config, func(e *agollo.ChangeEvent, err error) {
		if err != nil {
			log.Panicf("配置重载失败: %s", err.Error())
			return
		}

		log.Printf("========== 配置已重载 ==========\nNamespace: %s\n", e.Namespace)
		for k, v := range e.Changes {
			log.Printf(
				"[%s] %s\n----- old -----\n%s\n----- new -----\n%s\n",
				v.ChangeType.String(),
				k,
				v.OldValue,
				v.NewValue,
			)
		}
		log.Printf("最新的配置为: %v\n", Config)
	})
	if err != nil {
		panic(fmt.Sprintf("failed to init apollo: %s", err.Error()))
	}

	log.Printf("载入配置: %v\n", Config)
}
```
