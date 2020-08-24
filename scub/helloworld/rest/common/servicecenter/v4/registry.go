package v4

import (
    "encoding/json"
    "fmt"
    "github.com/apache/servicecomb-service-center/pkg/registry"
    "github.com/gorilla/websocket"
    "net/http"
    "net/url"

    "github.com/SmartsYoung/service-center-demo/scub/helloworld/rest/common/config"
    "github.com/SmartsYoung/service-center-demo/scub/helloworld/rest/common/restful"

)

var (
    domain  = "default"
    project = "default"

    existence        = "/v4/%s/registry/existence"
    MicroservicePath = "/v4/%s/registry/microservices"
    MicroserviceItem = "/v4/%s/registry/microservices/%s"
    InstancePath     = "/v4/%s/registry/microservices/%s/instances"
    Discovery        = "/v4/%s/registry/instances"
    Heartbeats       = "/v4/%s/registry/heartbeats"
    Watcher          = "/v4/%s/registry/microservices/%s/watcher"

    microserviceType string
)

type Client struct {
    RawUrl  string
    Domain  string
    Project string
}

func NewClient(addr string, domain string) *Client {
    return &Client{
        RawUrl: addr,
        Domain: domain,
    }
}

func (c *Client) existence(parm url.Values) (*registry.GetExistenceResponse, error) {
    reqUrl := c.RawUrl + fmt.Sprintf(existence, domain) + parm.Encode()
    req, err := restful.NewRequest(http.MethodGet, reqUrl, c.DefaultHeader(), nil)
    if err != nil {
        return nil, fmt.Errorf("Get Microservices failed, body is empty,  error: %s", err)
    }
    respData := &registry.GetExistenceResponse{}
    err = restful.DoRequest(req, respData)
    if err != nil {
        return nil, fmt.Errorf("Get Microservices failed, body is empty,  error: %s", err)
    }
    return respData, nil
}

// 获取微服务服务ID
func (c *Client) GetServiceID(svc *config.ServiceConf) (string, error) {
    val := url.Values{}
    val.Set("type", microserviceType)
    val.Set("appId", svc.AppId)
    val.Set("serviceName", svc.ServiceName)
    val.Set("version", svc.Version)
    respData, err := c.existence(val)
    if err == nil {
        return respData.ServiceId, nil
    }
    return "", fmt.Errorf("[GetServiceID]: %s", err)
}

// 注册微服务
func (c *Client) RegisterService(svc *config.ServiceConf) (string, error) {
    reqUrl := c.RawUrl + fmt.Sprintf(MicroservicePath, project)
    body := &registry.CreateServiceRequest{
        Service: &registry.MicroService{
            AppId:       svc.AppId,
            ServiceName: svc.ServiceName,
            Version:     svc.Version,
        },
    }
    req, err := restful.NewRequest(http.MethodPost, reqUrl, c.DefaultHeader(), body)
    if err != nil {
        return "", fmt.Errorf("[Post Microservices failed] %s", err)
    }
    respData := &registry.MicroService{}
    err = restful.DoRequest(req, respData)
    if err != nil {
        return "", fmt.Errorf("[Registry Microservices failed]: %s", err)
    }
    return respData.ServiceId, nil
}

// 注销微服务
func (c *Client) UnRegisterService(svcID string) error {
    reqUrl := c.RawUrl + fmt.Sprintf(MicroserviceItem, project, svcID)
    req, err := restful.NewRequest(http.MethodDelete, reqUrl, c.DefaultHeader(), nil)
    if err != nil {
        return fmt.Errorf("[Post Microservices failed] %s", err)
    }
    err = restful.DoRequest(req, nil)
    if err != nil {
        return fmt.Errorf("[Registry Microservices failed]: %s", err)
    }
    return nil
}

// 注册微服务实例
func (c *Client) RegisterInstance(scvId string, ins *config.InstanceConf) (string, error) {
    reqUrl := c.RawUrl + fmt.Sprintf(InstancePath, project, scvId)
    endpoint := ins.Protocol + "://" + ins.ListenAddress
    ms := &registry.RegisterInstanceRequest{
        Instance: &registry.MicroServiceInstance{
            HostName:  ins.Hostname,
            Endpoints: []string{endpoint},
        },
    }
    req, err := restful.NewRequest(http.MethodPost, reqUrl, c.DefaultHeader(), ms)
    if err == nil {
        respData := &registry.RegisterInstanceResponse{}
        err = restful.DoRequest(req, respData)
        if err == nil {
            return respData.InstanceId, nil
        }
    }
    return "", fmt.Errorf("[RegisterInstance]: %s", err)
}

// 注销微服务实例
func (c *Client) UnRegisterInstance(svcID , insID string) (error) {
    reqURL := c.RawUrl + fmt.Sprintf(InstancePath, svcID, insID)
    req, err := restful.NewRequest(http.MethodDelete, reqURL, c.DefaultHeader(), nil)
    if err == nil {
        err = restful.DoRequest(req, nil)
        if err == nil {
            return nil
        }
    }
    return fmt.Errorf("[UNRegisterInstance]: %s", err)
}


// 心跳保活
func (c *Client) Heartbeat(svcID, insID string) error {
    hb := &registry.HeartbeatSetRequest{
        Instances: []*registry.HeartbeatSetElement{
            {ServiceId: svcID, InstanceId: insID},
        },
    }
    reqURL := c.RawUrl + fmt.Sprintf(Heartbeats, project)
    req, err := restful.NewRequest(http.MethodPut, reqURL, c.DefaultHeader(), hb)
    if err != nil {
        return fmt.Errorf("[Heartbeat]: %s", err)
    }
    err = restful.DoRequest(req, nil)
    if err != nil {
        return fmt.Errorf("[Heartbeat]: %s", err)
    }
    return nil
}

// 服务发现
func (c *Client) Discovery(conID string, svc *config.ServiceConf) ([]*registry.MicroServiceInstance, error) {
    val := url.Values{}
    val.Set("appId",svc.AppId)
    val.Set("serviceName", svc.ServiceName)
    val.Set("version", svc.Version)

    reqURL := c.RawUrl + fmt.Sprintf(Discovery, project) + "?" + val.Encode()
    req, err := restful.NewRequest(http.MethodGet, reqURL, c.DefaultHeader(), nil)
    if err == nil {
        req.Header.Set("x-consumerid", conID)
        fmt.Println(req)
        respData := &registry.GetInstancesResponse{}
        err = restful.DoRequest(req, respData)
        if err == nil {
            return respData.Instances, nil
        }
    }
    return nil, fmt.Errorf("[Discovery]: %s", err)
}

// 服务订阅
func (c *Client) WatchService(svcID string, callback func(*registry.WatchInstanceResponse)) error {
    addr, err:= url.Parse(c.RawUrl + fmt.Sprintf(Watcher, project, svcID))
    if err != nil {
        return fmt.Errorf("[WatchService]: parse repositry url faild: %s", err)
    }

    // 注： watch接口使用了 websocket 长连接
    addr.Scheme = "ws"
    conn, _, err := (&websocket.Dialer{}).Dial(addr.String(),c.DefaultHeader())
    if err != nil {
        return fmt.Errorf("[WatchService]: start websocket faild: %s", err)
    }
    for  {
        messageType, message, err := conn.ReadMessage()
        if err != nil{
            fmt.Errorf("[conn]: %s", err)
        }

        if messageType == websocket.TextMessage {
            var resp  *registry.WatchInstanceResponse
            err := json.Unmarshal(message, &resp)
            if err != nil {
                break
            }
            callback(resp)
        }
    }
    return fmt.Errorf("[WatchService]: receive message faild: %s", err)
}


// 设置默认头部
func (c *Client) DefaultHeader() http.Header {
    headers := http.Header{
        "Content-Type":  []string{"application/json"},
        "X-Domain-Name": []string{"default"},
    }
    if c.Domain != "" {
        headers.Set("X-Domain-Name", c.Domain)
    }
    return headers
}

