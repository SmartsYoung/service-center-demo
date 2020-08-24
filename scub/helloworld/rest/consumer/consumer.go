package main

import (
    "fmt"
    "github.com/SmartsYoung/service-center-demo/scub/helloworld/rest/common/config"
    v4 "github.com/SmartsYoung/service-center-demo/scub/helloworld/rest/common/servicecenter/v4"
    "github.com/apache/servicecomb-service-center/pkg/registry"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "sync"
)

var caches = &sync.Map{}

func main() {
    // 配置文件加载
    err := config.LoadConfig("./conf/microservice.yaml")
    if err != nil {
        log.Fatalf("load config file faild: %s", err)
    }

    // 注册自身微服务
    svcID := registerService()

    // 服务发现 provider 实例信息
    discoveryProviderAndCache(svcID)

    // 与 provider 通讯
    log.Println(sayHello())

    // 提供对外服务，将请求转发到 helloServer 处理，验证 watch 功能
    sayHelloServer(svcID)
}

func discoveryProviderAndCache(svcId string) {
    cli := v4.NewClient(config.Registry.Address, config.Tenant.Domain)
    pris, err := cli.Discovery(svcId, config.Provider)
    if err != nil {
        log.Fatal(err)
    }
    if len(pris) == 0 {
        log.Fatalf("provider not found, serviceName: %s appId: %s, version: %s",
            config.Provider.ServiceName, config.Provider.AppId, config.Provider.Version)
    }
    if len(pris[0].Endpoints) == 0 {
        log.Fatalln("provider endpoints is empty")
    }

    caches.Store(config.Provider, pris)
}

// 提供对外服务，将请求转发到 helloServer 处理，验证 watch 功能
func sayHelloServer(svcId string)  {
    // 启动 provider 订阅
    go watch(svcId)

    // 启动 http 监听
    http.HandleFunc("/sayhello", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(sayHello()))
    })
    err := http.ListenAndServe(":8090", nil)
    log.Println(err)
}

// 注册服务
func registerService() string {
    cli := v4.NewClient(config.Registry.Address, config.Tenant.Domain)
    svcId, _ := cli.GetServiceID(config.Service)
    if svcId == "" {
        var err error
        svcId, err = cli.RegisterService(config.Service)
        if err != nil {
            log.Fatalln(err)
        }
    }
    return svcId
}

// 与 provider 通讯
func sayHello() string {
    addr, err := getProviderEndpoint()
    if err != nil {
        return err.Error()
    }
    req, err := http.NewRequest(http.MethodGet, addr+"/hello", nil)
    if err != nil {
        return fmt.Sprintf("create request faild: %s", err)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Sprintf("do request faild: %s", err)
    }

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Sprintf("read response body faild: %s, body: %s", err, string(data))
    }

    log.Printf("reply form provider: %s", string(data))
    return string(data)
}


func watch(svcId string) {
    cli := v4.NewClient(config.Registry.Address, config.Tenant.Domain)
    err := cli.WatchService(svcId, watchBack)
    if err != nil {
        log.Println(err)
    }
}

func watchBack(data *registry.WatchInstanceResponse){
    prisCache, ok := caches.Load(config.Provider)
    if !ok {
        log.Printf("provider \"%s\" not found", config.Provider.ServiceName)
        return
    }
    pris := prisCache.([]*registry.MicroServiceInstance)  // 注意这里是判断是否为Instance
    renew := false
    for i := 0; i < len(pris); i++ {
        if pris[i].InstanceId == data.Instance.InstanceId {
            pris[i] = data.Instance
            renew = true
            break
        }
    }
    if !renew {
        pris = append(pris, data.Instance)
    }
    caches.Store(config.Provider, pris)
}

// 获取在线的 provider endpoint
func getProviderEndpoint() (string, error) {
    prisCache, ok := caches.Load(config.Provider)
    if !ok {
        return "", fmt.Errorf("provider \"%s\" not found", config.Provider.ServiceName)
    }
    pris := prisCache.([]*registry.MicroServiceInstance)

    endpoint := ""

    for i := 0; i < len(pris); i++{
        if pris[i].Status == "UP" {
            endpoint = pris[i].Endpoints[0]
            break
        }
    }

    if endpoint != "" {
        addr, err := url.Parse(endpoint)
        if err!= nil{
            return "", fmt.Errorf("parse provider endpoint faild: %s", err)
        }
        if addr.Scheme == "rest" {
            addr.Scheme = "http"
        }
        return addr.String(), nil
    }
    return "", fmt.Errorf("provider \"%s\" endpoint not found", config.Provider.ServiceName)
}