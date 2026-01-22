package configs

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"runtime"
)

var UseConfig SysUseConfig

// 패키지 로드시 초기화
func init() {
	UseConfig = GetConfig()
}

// GetConfig :
func GetConfig() SysUseConfig {
	// 기본값 설정
	defaultConfig := DidConfig{
		SystemMode:         "LOCAL",
		SystemRunMode:      "N/A",
		SystemLogFlag:      "N/A",
		SystemLogMode:      "N/A",
		SystemLogPrintMode: "N/A",
		GenerationRule:     "hexdigit",
		RelService: struct {
			DidRegistry struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"did-registry"`
			DidRegistrar struct {
				Address           string   `yaml:"address"`
				Port              string   `yaml:"port"`
				AdoptedDriverList []string `yaml:"adopted_driver_list"`
			} `yaml:"did-sep"`
			ServiceEndpoint struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"service_endpoint"`
			RelyingParty struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"relying_party"`
			Issuer struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"issuer"`
			EthClient struct {
				RawUrl    string `yaml:"raw_url"`
				ScAddress string `yaml:"sc_address"`
			} `yaml:"eth_client"`
		}{
			DidRegistry: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50051",
				Port:    ":50051",
			},
			DidRegistrar: struct {
				Address           string   `yaml:"address"`
				Port              string   `yaml:"port"`
				AdoptedDriverList []string `yaml:"adopted_driver_list"`
			}{
				Address:           "localhost:50052",
				Port:              ":50052",
				AdoptedDriverList: []string{"byd50", "eth", "test"},
			},
			ServiceEndpoint: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50053",
				Port:    ":50053",
			},
			RelyingParty: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50054",
				Port:    ":50054",
			},
			Issuer: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50055",
				Port:    ":50055",
			},
			EthClient: struct {
				RawUrl    string `yaml:"raw_url"`
				ScAddress string `yaml:"sc_address"`
			}{
				RawUrl:    "https://data-seed-prebsc-1-s1.binance.org:8545",
				ScAddress: "0xEA40445e1C77071d62D8e25f11bD72314D266BD1",
			},
		},
		DevService: struct {
			DidRegistry struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"did-registry"`
			DidRegistrar struct {
				Address           string   `yaml:"address"`
				Port              string   `yaml:"port"`
				AdoptedDriverList []string `yaml:"adopted_driver_list"`
			} `yaml:"did-sep"`
			ServiceEndpoint struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"service_endpoint"`
			RelyingParty struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"relying_party"`
			Issuer struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"issuer"`
			EthClient struct {
				RawUrl    string `yaml:"raw_url"`
				ScAddress string `yaml:"sc_address"`
			} `yaml:"eth_client"`
		}{
			DidRegistry: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50051",
				Port:    ":50051",
			},
			DidRegistrar: struct {
				Address           string   `yaml:"address"`
				Port              string   `yaml:"port"`
				AdoptedDriverList []string `yaml:"adopted_driver_list"`
			}{
				Address:           "localhost:50052",
				Port:              ":50052",
				AdoptedDriverList: []string{"byd50", "eth", "test"},
			},
			ServiceEndpoint: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50053",
				Port:    ":50053",
			},
			RelyingParty: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50054",
				Port:    ":50054",
			},
			Issuer: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50055",
				Port:    ":50055",
			},
			EthClient: struct {
				RawUrl    string `yaml:"raw_url"`
				ScAddress string `yaml:"sc_address"`
			}{
				RawUrl:    "https://data-seed-prebsc-1-s1.binance.org:8545",
				ScAddress: "0xEA40445e1C77071d62D8e25f11bD72314D266BD1",
			},
		},
		LocalService: struct {
			DidRegistry struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"did-registry"`
			DidRegistrar struct {
				Address           string   `yaml:"address"`
				Port              string   `yaml:"port"`
				AdoptedDriverList []string `yaml:"adopted_driver_list"`
			} `yaml:"did-sep"`
			ServiceEndpoint struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"service_endpoint"`
			RelyingParty struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"relying_party"`
			Issuer struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			} `yaml:"issuer"`
			EthClient struct {
				RawUrl    string `yaml:"raw_url"`
				ScAddress string `yaml:"sc_address"`
			} `yaml:"eth_client"`
		}{
			DidRegistry: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50051",
				Port:    ":50051",
			},
			DidRegistrar: struct {
				Address           string   `yaml:"address"`
				Port              string   `yaml:"port"`
				AdoptedDriverList []string `yaml:"adopted_driver_list"`
			}{
				Address:           "localhost:50052",
				Port:              ":50052",
				AdoptedDriverList: []string{"byd50", "eth", "test"},
			},
			ServiceEndpoint: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50053",
				Port:    ":50053",
			},
			RelyingParty: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50054",
				Port:    ":50054",
			},
			Issuer: struct {
				Address string `yaml:"address"`
				Port    string `yaml:"port"`
			}{
				Address: "localhost:50055",
				Port:    ":50055",
			},
			EthClient: struct {
				RawUrl    string `yaml:"raw_url"`
				ScAddress string `yaml:"sc_address"`
			}{
				RawUrl:    "https://data-seed-prebsc-1-s1.binance.org:8545",
				ScAddress: "0xEA40445e1C77071d62D8e25f11bD72314D266BD1",
			},
		},
	}

	var config DidConfig
	var useConfig SysUseConfig

	// config.yml 파일 경로 설정
	yamlFilePath, _ := filepath.Abs("./configs.yml")
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		yamlFilePath, _ = filepath.Abs("../configs.yml")
		yamlFile, err = ioutil.ReadFile(yamlFilePath)
		if err != nil {
			yamlFilePath, _ = filepath.Abs("../../configs.yml")
			yamlFile, err = ioutil.ReadFile(yamlFilePath)
			if err != nil {
				yamlFilePath, _ = filepath.Abs("/sdcard/configs.yml")
				yamlFile, err = ioutil.ReadFile(yamlFilePath)
				if err != nil {
					log.Printf("Failed to read config file: %s. Using default configuration.", err)
					config = defaultConfig
				}
			}
		}
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Printf("Failed to parse config file: %s. Using default configuration.", err)
		config = defaultConfig
	}

	useConfig.SystemRunMode = config.SystemRunMode
	useConfig.SystemLogFlag = config.SystemLogFlag
	useConfig.SystemLogMode = config.SystemLogMode
	useConfig.SystemLogPrintMode = config.SystemLogPrintMode
	useConfig.GenerationRule = config.GenerationRule

	// 시스템 모드 별 설정
	switch config.SystemMode {
	case SystemModeRel:
		useConfig.DidRegistryAddress = config.RelService.DidRegistry.Address
		useConfig.DidRegistryPort = config.RelService.DidRegistry.Port
		useConfig.DidRegistrarAddress = config.RelService.DidRegistrar.Address
		useConfig.DidRegistrarPort = config.RelService.DidRegistrar.Port
		useConfig.AdoptedDriverList = config.RelService.DidRegistrar.AdoptedDriverList
		useConfig.ServiceEndpointAddress = config.RelService.ServiceEndpoint.Address
		useConfig.ServiceEndpointPort = config.RelService.ServiceEndpoint.Port
		useConfig.RelyingPartyAddress = config.RelService.RelyingParty.Address
		useConfig.RelyingPartyPort = config.RelService.RelyingParty.Port
		useConfig.IssuerAddress = config.RelService.Issuer.Address
		useConfig.IssuerPort = config.RelService.Issuer.Port
		useConfig.EthClientUrl = config.RelService.EthClient.RawUrl
		useConfig.EthClientScAddress = config.RelService.EthClient.ScAddress
	case SystemModeDev:
		useConfig.DidRegistryAddress = config.DevService.DidRegistry.Address
		useConfig.DidRegistryPort = config.DevService.DidRegistry.Port
		useConfig.DidRegistrarAddress = config.DevService.DidRegistrar.Address
		useConfig.DidRegistrarPort = config.DevService.DidRegistrar.Port
		useConfig.AdoptedDriverList = config.DevService.DidRegistrar.AdoptedDriverList
		useConfig.ServiceEndpointAddress = config.DevService.ServiceEndpoint.Address
		useConfig.ServiceEndpointPort = config.DevService.ServiceEndpoint.Port
		useConfig.RelyingPartyAddress = config.DevService.RelyingParty.Address
		useConfig.RelyingPartyPort = config.DevService.RelyingParty.Port
		useConfig.IssuerAddress = config.DevService.Issuer.Address
		useConfig.IssuerPort = config.DevService.Issuer.Port
		useConfig.EthClientUrl = config.DevService.EthClient.RawUrl
		useConfig.EthClientScAddress = config.DevService.EthClient.ScAddress
	case SystemModeLocal:
		useConfig.DidRegistryAddress = config.LocalService.DidRegistry.Address
		useConfig.DidRegistryPort = config.LocalService.DidRegistry.Port
		useConfig.DidRegistrarAddress = config.LocalService.DidRegistrar.Address
		useConfig.DidRegistrarPort = config.LocalService.DidRegistrar.Port
		useConfig.AdoptedDriverList = config.LocalService.DidRegistrar.AdoptedDriverList
		useConfig.ServiceEndpointAddress = config.LocalService.ServiceEndpoint.Address
		useConfig.ServiceEndpointPort = config.LocalService.ServiceEndpoint.Port
		useConfig.RelyingPartyAddress = config.LocalService.RelyingParty.Address
		useConfig.RelyingPartyPort = config.LocalService.RelyingParty.Port
		useConfig.IssuerAddress = config.LocalService.Issuer.Address
		useConfig.IssuerPort = config.LocalService.Issuer.Port
		useConfig.EthClientUrl = config.LocalService.EthClient.RawUrl
		useConfig.EthClientScAddress = config.LocalService.EthClient.ScAddress
	default:
	}
	return useConfig
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
