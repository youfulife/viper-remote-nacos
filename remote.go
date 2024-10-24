package viper_remote_nacos

import (
    "bytes"
    "io"
    "strings"

    "github.com/spf13/viper"
)

type remoteConfig struct{}

func (rc remoteConfig) Get(rp viper.RemoteProvider) (io.Reader, error) {
    client, err := NewClient(strings.Split(rp.Endpoint(), ";"))
    if err != nil {
        return nil, err
    }

    dataId := strings.Split(rp.Path(), "/")[0]
    group := strings.Split(rp.Path(), "/")[1]
    b, err := client.Get(dataId, group)
    if err != nil {
        return nil, err
    }

    return bytes.NewReader(b), nil
}

func (rc remoteConfig) Watch(rp viper.RemoteProvider) (io.Reader, error) {
    return rc.Get(rp)
}

func (rc remoteConfig) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
    client, err := NewClient(strings.Split(rp.Endpoint(), ";"))
    if err != nil {
        return nil, nil
    }

    quit := make(chan bool)
    resp := make(chan *viper.RemoteResponse)

    dataId := strings.Split(rp.Path(), "/")[0]
    group := strings.Split(rp.Path(), "/")[1]
    err = client.Watch(dataId, group, func(data string) {
        resp <- &viper.RemoteResponse{Value: []byte(data)}
    })

    return resp, quit
}

func init() {
    viper.SupportedRemoteProviders = append(viper.SupportedRemoteProviders, "nacos")
    viper.RemoteConfig = &remoteConfig{}
}
