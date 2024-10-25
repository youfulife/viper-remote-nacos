package main

import (
    "fmt"
    "log"
    "time"

    "github.com/spf13/viper"

    _ "github.com/youfulife/viper-remote-nacos"
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

    // 从 nacos 中读取配置
    err = v.ReadRemoteConfig()
    if err != nil {
        log.Fatal(err)
    }

    // 打印配置信息
    fmt.Println("Database Host:", v.GetString("mysql.host"))
    fmt.Println("Database Port:", v.GetInt("mysql.port"))

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
