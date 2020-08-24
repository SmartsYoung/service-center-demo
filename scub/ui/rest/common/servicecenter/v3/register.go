package v3

import (
	"encoding/json"
	"fmt"
	"github.com/SmartsYoung/service-center-demo/scub/ui/rest/common/config"
	"github.com/SmartsYoung/service-center-demo/scub/ui/rest/common/restful"
	"github.com/apache/servicecomb-service-center/pkg/registry"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

var (
	// 接口 API 定义
	microServices    = "/registry/v3/microservices"
	microServiceItem = "/registry/v3/microservices/%s"
	svcInstances     = "/registry/v3/microservices/%s/instances"
	svcInstanceItem  = "/registry/v3/microservices/%s/instances/%s"
	discovery        = "/registry/v3/instances"
	existence        = "/registry/v3/existence"
	heartbeats       = "/registry/v3/heartbeats"
	watcher          = "/registry/v3/microservices/%s/watcher"

	microServiceType sourceType = "microservice"

)

type sourceType string

type Client struct {
	rawURL string
	domain string
}


func NewClient(addr string, domain string) *Client {
	return &Client{rawURL: addr, domain: domain}
}

// 查询微服务是否存在
func (c *Client) existence(params url.Values) (*registry.GetExistenceResponse, error) {
	reqURL := c.rawURL + existence + "?" + params.Encode()
	req, err := restful.NewRequest(http.MethodGet, reqURL, c.DefaultHeaders(), nil)
	if err == nil {
		respData := &registry.GetExistenceResponse{}
		err = restful.DoRequest(req, respData)
		if err == nil {
			return respData, nil
		}
	}
	return nil, err
}

// 获取微服务服务ID
func (c *Client) GetServiceID(svc *config.ServiceConf) (string, error) {
	val := url.Values{}
	val.Set("type", string(microServiceType))
	val.Set("appId", svc.AppID)
	val.Set("serviceName", svc.Name)
	val.Set("version", svc.Version)

	fmt.Println(svc, "svc ")
	respData, err := c.existence(val)
	if err == nil {
		return respData.ServiceId, nil
	}
	return "", fmt.Errorf("[GetServiceID]: %s", err)
}

// 注册微服务
func (c *Client) RegisterService(svc *config.ServiceConf) (string, error) {
	log.Println("svc is :", svc)
	ms := &registry.CreateServiceRequest{
		Service: &registry.MicroService{
			AppId:       svc.AppID,
			ServiceName: svc.Name,
			Version:     svc.Version,
		},
	}

	reqURL := c.rawURL + microServices
	req, err := restful.NewRequest(http.MethodPost, reqURL, c.DefaultHeaders(), ms)
	if err == nil {
		respData := &registry.CreateServiceResponse{}
		err = restful.DoRequest(req, respData)
		if err == nil {
			return respData.ServiceId, nil
		}
	}

	return "", fmt.Errorf("[RegisterService]: %s", err)
}

// 注销微服务
func (c *Client) UnRegisterService(svcID string) error {
	reqURL := c.rawURL + fmt.Sprintf(microServiceItem, svcID)
	req, err := restful.NewRequest(http.MethodDelete, reqURL, c.DefaultHeaders(), nil)
	if err != nil {
		return fmt.Errorf("[UnRegisterService]: %s", err)
	}
	err = restful.DoRequest(req, nil)
	if err != nil {
		return fmt.Errorf("[UnRegisterService]: %s", err)
	}
	return nil
}

// 注册微服务实例
func (c *Client) RegisterInstance(svcID string, ins *config.InstanceConf) (string, error) {
	endpoint := ins.Protocol + "://" + ins.ListenAddress
	ms := &registry.RegisterInstanceRequest{
		Instance: &registry.MicroServiceInstance{
			Endpoints: []string{endpoint},
			HostName:  ins.Hostname,
		},
	}
	reqURL := c.rawURL + fmt.Sprintf(svcInstances, svcID)

	req, err := restful.NewRequest(http.MethodPost, reqURL, c.DefaultHeaders(), ms)
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
func (c *Client) UnRegisterInstance(svcID, insID string) error {
	reqURL := c.rawURL + fmt.Sprintf(svcInstanceItem, svcID, insID)
	req, err := restful.NewRequest(http.MethodDelete, reqURL, c.DefaultHeaders(), nil)
	if err != nil {
		return fmt.Errorf("[UnRegisterInstance]: %s", err)
	}
	err = restful.DoRequest(req, nil)
	if err != nil {
		return fmt.Errorf("[UnRegisterInstance]: %s", err)
	}
	return nil
}

// 心跳保活
func (c *Client) Heartbeat(svcID, insID string) error {
	hb := &registry.HeartbeatSetRequest{Instances: []*registry.HeartbeatSetElement{
		{ServiceId: svcID, InstanceId: insID},
	}}

	reqURL := c.rawURL + heartbeats

	req, err := restful.NewRequest(http.MethodPut, reqURL, c.DefaultHeaders(), hb)
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
	val.Set("appId", svc.AppID)
	val.Set("serviceName", svc.Name)
	val.Set("version", svc.Version)

	reqURL := c.rawURL + discovery + "?" + val.Encode()
	log.Println(reqURL)
	req, err := restful.NewRequest(http.MethodGet, reqURL, c.DefaultHeaders(), nil)
	if err == nil {
		req.Header.Set("x-consumerid", conID)
		fmt.Println("req",req, req.Body, req.URL, req.Header, req.RequestURI)
		respData := &registry.GetInstancesResponse{}
		err = restful.DoRequest(req, respData)
		if err == nil {
			return respData.Instances, nil
		}
	}
	return nil, fmt.Errorf("[Discovery]: %s", err)
}

func (c *Client) WatchService(svcID string, callback func(*registry.WatchInstanceResponse)) error {
	addr, err := url.Parse(c.rawURL + fmt.Sprintf(watcher, svcID))
	if err != nil {
		return fmt.Errorf("[WatchService]: parse repositry url faild: %s", err)
	}

	// 注： watch接口使用了 websocket 长连接
	addr.Scheme = "ws"
	conn, _, err := (&websocket.Dialer{}).Dial(addr.String(), c.DefaultHeaders())
	if err != nil {
		return fmt.Errorf("[WatchService]: start websocket faild: %s", err)
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if messageType == websocket.TextMessage {
			data := &registry.WatchInstanceResponse{}
			err := json.Unmarshal(message, data)
			if err != nil {
				log.Println(err)
				break
			}
			callback(data)
		}
	}
	return fmt.Errorf("[WatchService]: receive message faild: %s", err)

}

// 设置默认头部
func (c *Client) DefaultHeaders() http.Header {
	headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"X-Domain-Name": []string{"default"},
	}
	if c.domain != "" {
		headers.Set("X-Domain-Name", c.domain)
	}
	return headers
}
