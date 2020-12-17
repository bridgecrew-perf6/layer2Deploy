package layer2config

import (
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
)

var (
	DEFAULT_LOG_LEVEL = log.InfoLog
	DEFAULT_REST_PORT = uint(18082)
)

type Config struct {
	NetWorkId            int    `json:"netWorkId"`
	WalletName           string `json:"walletName"`
	GasPrice             int    `json:"gasPrice"`
	Layer2Contract       string `json:"layer2Contract"`
	Layer2WitContract    string `json:"layer2WitContract"`
	Layer2MainNetNode    string `json:"layer2MainNetNode"`
	Layer2TestNetNode    string `json:"layer2TestNetNode"`
	Layer2RecordInterval int    `json:"layer2RecordInterval"`
	Layer2RecordBatchC   uint32 `json:"layer2RecordBatchC"`
	Layer2RetryCount     int    `json:"layer2RetryCount"`
	RequestLogServer     string `json:"requestLogServer"`
	RestPort             string `json:"restPort"`
	EnableSendService    bool   `json:"enableSendService"`
	Layer2Sdk            *ontology_go_sdk.Layer2Sdk
	OntSdk               *ontology_go_sdk.OntologySdk
	AdminAccount         *ontology_go_sdk.Account
}
