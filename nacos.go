package viper_remote_nacos

import (
    "strings"

    "github.com/nacos-group/nacos-sdk-go/v2/clients"
    "github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
    "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
    "github.com/nacos-group/nacos-sdk-go/v2/vo"
    "github.com/spf13/cast"
)

type Client struct {
    client config_client.IConfigClient
}

func NewClient(machines []string) (*Client, error) {

    var serverConfigs []constant.ServerConfig
    for _, machine := range machines {
        ss := strings.Split(machine, ":")
        ipAddr := ss[0]
        port := cast.ToUint64(ss[1])
        serverConfigs = append(serverConfigs, *constant.NewServerConfig(ipAddr, port, constant.WithContextPath("/nacos")))
    }

    //create ClientConfig
    cc := *constant.NewClientConfig(
        constant.WithNotLoadCacheAtStart(true),
        constant.WithNamespaceId(""), // default public
    )
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

    return &Client{client: client}, nil
}

func (c *Client) Get(dataId, group string) ([]byte, error) {
    content, err := c.client.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
    if err != nil {
        return nil, err
    }
    return []byte(content), nil
}

func (c *Client) Watch(dataId, group string, f func(string)) error {
    err := c.client.ListenConfig(vo.ConfigParam{
        DataId: dataId,
        Group:  group,
        OnChange: func(namespace, group, dataId, data string) {
            f(data)
        },
    })
    return err
}
