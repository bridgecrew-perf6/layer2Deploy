package layer2config

import (
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
)

var Version = ""

var (
	DEFAULT_LOG_LEVEL = log.InfoLog
	DEFAULT_REST_PORT = uint(8080)
)

type Config struct {
	NetWorkId            int    `json:"netWorkId"`
	WalletName           string `json:"walletName"`
	GasPrice             int    `json:"gasPrice"`
	Layer2Contract       string `json:"layer2Contract"`
	Layer2MainNetNode    string `json:"layer2MainNetNode"`
	Layer2TestNetNode    string `json:"layer2TestNetNode"`
	Layer2RecordInterval int    `json:"layer2RecordInterval"`
	Layer2Sdk            *ontology_go_sdk.Layer2Sdk
	AdminAccount         *ontology_go_sdk.Account
}

var DefLayer2Config = &Config{}
