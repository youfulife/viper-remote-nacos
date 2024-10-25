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
    opts   []constant.ClientOption
}

var defaultRemoteConfig = &remoteConfig{}

func SetNacosClientOption(opts ...constant.ClientOption) {
    defaultRemoteConfig.opts = append(defaultRemoteConfig.opts, opts...)
}

func (rc *remoteConfig) initClient(endpoint string, opts ...constant.ClientOption) error {
    machines := strings.Split(endpoint, ";")
    var serverConfigs []constant.ServerConfig
    for _, machine := range machines {
        ss := strings.Split(machine, ":")
        ipAddr := ss[0]
        port := cast.ToUint64(ss[1])
        serverConfigs = append(serverConfigs, *constant.NewServerConfig(ipAddr, port, constant.WithContextPath("/nacos")))
    }

    cc := *constant.NewClientConfig(opts...)
    client, err := clients.NewConfigClient(vo.NacosClientParam{ClientConfig: &cc, ServerConfigs: serverConfigs})
    if err != nil {
        return err
    }

    rc.client = client
    return nil
}

func (rc *remoteConfig) Get(rp viper.RemoteProvider) (io.Reader, error) {
    if rc.client == nil {
        err := rc.initClient(rp.Endpoint(), rc.opts...)
        if err != nil {
            return nil, err
        }
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

    quit := make(chan bool)
    resp := make(chan *viper.RemoteResponse)

    if rc.client == nil {
        err := rc.initClient(rp.Endpoint(), rc.opts...)
        if err != nil {
            return resp, quit
        }
    }

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
    viper.RemoteConfig = defaultRemoteConfig
}
