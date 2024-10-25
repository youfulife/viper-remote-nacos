# viper-remote-nacos

viper-remote-nacos for Go Viper library allows you to use nacos as a remote configuration center.

# Installation
Use go get to install SDK：

```shell
$ go get -u github.com/youfulife/viper-remote-nacos
```

# Quick Examples

see example package main.go

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/spf13/viper"
    _ "github.com/youfulife/viper-remote-nacos"
)

func main() {
    v := viper.New()
    // path is group/dataId
    group := "example-remote"
    dataId := "example-remote"
    path := dataId + "/" + group
    err := v.AddRemoteProvider("nacos", "localhost:8848", path)
    if err != nil {
        log.Fatal(err)
    }
    // set config type
    v.SetConfigType("yaml")
    
    // read remote config
    err = v.ReadRemoteConfig()
    if err != nil {
        log.Fatal(err)
    }

    // listen remote config change
    err = v.WatchRemoteConfigOnChannel()
    if err != nil {
        log.Println("Failed to watch remote config:", err)
    }

    // print config every 5 seconds
    for {
        time.Sleep(5 * time.Second)
        fmt.Println("Database Host:", v.GetString("mysql.host"))
        fmt.Println("Database Port:", v.GetInt("mysql.port"))
    }
}

```

if you need set nacos client option, use like this.

```go

import (
    remote "github.com/youfulife/viper-remote-nacos"
)

// you must set nacos client option before read remote config or watch remote config on channel
remote.SetNacosClientOption(
    constant.WithNamespaceId(""),
    constant.WithTimeoutMs(5000),
    constant.WithNotLoadCacheAtStart(true),
    constant.WithLogDir("tmp/nacos/log"),
    constant.WithCacheDir("tmp/nacos/cache"),
    constant.WithLogLevel("debug"),
)
```

