package configs

// DidConfig : API 서버 환경 설정
type DidConfig struct {
	SystemMode         string `yaml:"system_mode"`
	SystemRunMode      string `yaml:"system_run_mode"`
	SystemLogFlag      string `yaml:"system_log_flag"`
	SystemLogMode      string `yaml:"system_log_mode"`
	SystemLogPrintMode string `yaml:"system_log_print_mode"`
	GenerationRule     string `yaml:"generation_rule"`

	RelService struct {
		DidRegistry struct {
			Address string `yaml:"address"`
			Port    string `yaml:"port"`
		} `yaml:"did-registry"`
		DidRegistrar struct {
			Address           string   `yaml:"address"`
			Port              string   `yaml:"port"`
			AdoptedDriverList []string `yaml:"adopted_driver_list"`
		} `yaml:"did-registrar"`
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
	} `yaml:"rel_service"`
	DevService struct {
		DidRegistry struct {
			Address string `yaml:"address"`
			Port    string `yaml:"port"`
		} `yaml:"did-registry"`
		DidRegistrar struct {
			Address           string   `yaml:"address"`
			Port              string   `yaml:"port"`
			AdoptedDriverList []string `yaml:"adopted_driver_list"`
		} `yaml:"did-registrar"`
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
	} `yaml:"dev_service"`
	LocalService struct {
		DidRegistry struct {
			Address string `yaml:"address"`
			Port    string `yaml:"port"`
		} `yaml:"did-registry"`
		DidRegistrar struct {
			Address           string   `yaml:"address"`
			Port              string   `yaml:"port"`
			AdoptedDriverList []string `yaml:"adopted_driver_list"`
		} `yaml:"did-registrar"`
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
	} `yaml:"local_service"`
}

// SysUseConfig : API 서버 환경 설정
type SysUseConfig struct {
	SystemRunMode          string
	SystemLogFlag          string
	SystemLogMode          string
	SystemLogPrintMode     string
	DidRegistryAddress     string
	DidRegistryPort        string
	DidRegistrarAddress    string
	DidRegistrarPort       string
	AdoptedDriverList      []string
	ServiceEndpointAddress string
	ServiceEndpointPort    string
	RelyingPartyAddress    string
	RelyingPartyPort       string
	IssuerAddress          string
	IssuerPort             string
	GenerationRule         string
	EthClientUrl           string
	EthClientScAddress     string
}
