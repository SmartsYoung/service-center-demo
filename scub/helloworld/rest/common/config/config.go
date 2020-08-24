package config

import (
    "fmt"
    "github.com/go-yaml/yaml"
    "io/ioutil"
    "log"
    "net"
    "os"
    "strconv"
)

var (
    Service   *ServiceConf
    Provider  *ServiceConf
    Registry  *RegistryConf
    Tenant    *TenantConf
    Instance  *InstanceConf
)

type MicroService struct {
    Service  *ServiceConf  `yaml:"service"`
    Instance *InstanceConf `yaml:"instance"`
    Provider *ServiceConf  `yaml:"provider"`
    Registry *RegistryConf `yaml:"registry"`
    Tenant   *TenantConf   `yaml:"tenant"`
}

type RegistryConf struct {
    Address string `yaml:"address"`
}

type ServiceConf struct {
    AppId       string `yaml:"appId"`
    ServiceName string `yaml:"serviceName"`
    Version     string `yaml:"version"`
}

// 实例配置
type InstanceConf struct {
    Hostname      string `yaml:"hostname"`
    Protocol      string `yaml:"protocol"`
    ListenAddress string `yaml:"listenAddress"`
}

type TenantConf struct {
    Domain  string `yaml:"domain"`
    Project string `yaml:"project"`
}

func LoadConfig(file string) error {

    config, err := ioutil.ReadFile(file)
    if err != nil {
        log.Fatal(err)
    }

    conf := MicroService{}
    err = yaml.Unmarshal(config, &conf)
    if err != nil {
        log.Fatalf("unmarshral yaml config eroors:", err)
    }

    if conf.Tenant == nil {
        conf.Tenant = &TenantConf{}
    }

    if conf.Tenant.Domain == "" {
        conf.Tenant.Domain = "default"
    }

    if conf.Instance != nil {
        if conf.Instance.Hostname == "" {
            conf.Instance.Hostname, _ = os.Hostname()
        }

        if conf.Instance.ListenAddress == "" {
            return fmt.Errorf("instance lister address is empty")
        }

        host, port, err := net.SplitHostPort(conf.Instance.ListenAddress)
        if err != nil {
            return fmt.Errorf("instance lister address is wrong: %s", err)
        }
        if host == "" {
            host = "127.0.0.1"
        }
        num, err := strconv.Atoi(port)
        if err != nil || num <= 0 {
            return fmt.Errorf("instance lister port %s is wrong: %s", port, err)
        }
        conf.Instance.ListenAddress = host + ":" + port
    }

    Service = conf.Service          //注意这里别写反了
    Instance = conf.Instance
    Registry = conf.Registry
    Provider = conf.Provider
    Tenant = conf.Tenant
    return nil
}
