package viper_remote_nacos

import (
    "bytes"
    "io"
    "strings"

    "github.com/nacos-group/nacos-sdk-go/v2/clients"
    "github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
    "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
    "github.com/nacos-group/nacos-sdk-go/v2/vo"
    "github.com/spf13/cast"
    "github.com/spf13/viper"
)

type remoteConfig struct {
    client config_client.IConfigClient
}

func NewNacosClient(endpoint string, opts ...constant.ClientOption) (config_client.IConfigClient, error) {

    // create serverConfigs
    machines := strings.Split(endpoint, ";")
    var serverConfigs []constant.ServerConfig
    for _, machine := range machines {
        ss := strings.Split(machine, ":")
        ipAddr := ss[0]
        port := cast.ToUint64(ss[1])
        serverConfigs = append(serverConfigs, *constant.NewServerConfig(ipAddr, port, constant.WithContextPath("/nacos")))
    }

    //create ClientConfig
    cc := *constant.NewClientConfig(opts...)
    // create config client
    client, err := clients.NewConfigClient(
        vo.NacosClientParam{
            ClientConfig:  &cc,
            ServerConfigs: serverConfigs,
        },
    )
    if err != nil {
        return nil, err
    }

    return client, nil
}

func (rc *remoteConfig) initClient(endpoint string) {
    client, err := NewNacosClient(endpoint)
    if err != nil {
        panic(err)
    }
    rc.client = client
}

func (rc *remoteConfig) Get(rp viper.RemoteProvider) (io.Reader, error) {
    if rc.client == nil {
        rc.initClient(rp.Endpoint())
    }

    dataId := strings.Split(rp.Path(), "/")[0]
    group := strings.Split(rp.Path(), "/")[1]

    content, err := rc.client.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
    if err != nil {
        return nil, err
    }
    return bytes.NewReader([]byte(content)), nil

}

func (rc *remoteConfig) Watch(rp viper.RemoteProvider) (io.Reader, error) {
    return rc.Get(rp)
}

func (rc *remoteConfig) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
    if rc.client == nil {
        rc.initClient(rp.Endpoint())
    }

    quit := make(chan bool)
    resp := make(chan *viper.RemoteResponse)

    dataId := strings.Split(rp.Path(), "/")[0]
    group := strings.Split(rp.Path(), "/")[1]
    _ = rc.client.ListenConfig(vo.ConfigParam{
        DataId: dataId,
        Group:  group,
        OnChange: func(namespace, group, dataId, data string) {
            resp <- &viper.RemoteResponse{Value: []byte(data)}
        },
    })

    return resp, quit
}

func init() {
    viper.SupportedRemoteProviders = append(viper.SupportedRemoteProviders, "nacos")
    viper.RemoteConfig = &remoteConfig{}
}
