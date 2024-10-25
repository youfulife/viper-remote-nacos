package main

import (
    "fmt"
    "log"
    "time"

    "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
    "github.com/spf13/viper"
    remote "github.com/youfulife/viper-remote-nacos"
)

func main() {

    // 设置 viper 使用 nacos 作为配置源
    v := viper.New()

    group := "example-remote"
    dataId := "example-remote"
    path := dataId + "/" + group

    err := v.AddRemoteProvider("nacos", "localhost:8848", path)
    if err != nil {
        log.Fatal(err)
    }

    v.SetConfigType("yaml") // 配置文件类型

    // 设置 nacos 客户端配置
    remote.SetNacosClientOption(
        constant.WithNamespaceId(""),
        constant.WithTimeoutMs(5000),
        constant.WithNotLoadCacheAtStart(true),
        constant.WithLogDir("tmp/nacos/log"),
        constant.WithCacheDir("tmp/nacos/cache"),
        constant.WithLogLevel("debug"),
    )

    // 从 nacos 中读取配置
    err = v.ReadRemoteConfig()
    if err != nil {
        log.Fatal(err)
    }

    // 监听配置变化
    err = v.WatchRemoteConfigOnChannel()
    if err != nil {
        log.Println("Failed to watch remote config:", err)
    }

    // 打印配置信息
    for {
        time.Sleep(5 * time.Second)
        fmt.Println("Database Host:", v.GetString("mysql.host"))
        fmt.Println("Database Port:", v.GetInt("mysql.port"))
    }

}
