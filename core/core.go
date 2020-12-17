package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ontio/layer2deploy/http"
	"github.com/ontio/layer2deploy/layer2config"
	ontSdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
)

type RequestLog struct {
	Timestamp       time.Time     `json:"time"`
	ResponseStatus  int64         `json:"status"`
	ProcessDuration time.Duration `json:"process_duration"`
	HTTPMethod      string        `json:"http_method"`
	Path            string        `json:"path"`
	RequestParam    string        `json:"request_param"`
	ResponseParam   string        `json:"response_param"`
}

type SendService struct {
	Wg      *sync.WaitGroup
	QuitS   chan bool
	Cfg     *layer2config.Config
	Enabled uint32
}

var DefSendService *SendService

func NewSendService(cfg *layer2config.Config) *SendService {
	return &SendService{
		Wg:      new(sync.WaitGroup),
		QuitS:   make(chan bool),
		Cfg:     cfg,
		Enabled: 1, //default enabled
	}
}

func (self *SendService) RepeantSendLogToChain() {
	self.Wg.Add(1)
	defer self.Wg.Done()

	client := http.NewClient()
	layer2Sdk := self.Cfg.Layer2Sdk
	contractAddr, err := common.AddressFromHexString(self.Cfg.Layer2Contract)
	if err != nil {
		log.Errorf("RepeantSendLogToChain N.0 %s", err)
		return
	}

	log.Infof("RepeantSendLogToChain Y.0. start to send log to chain.")
	var count uint64
	for {
		atomic.LoadUint32(&self.Enabled)

		select {
		case <-self.QuitS:
			log.Infof("RepeantSendLogToChain get QuitS signal")
			return
		default:
			if self.Enabled == 0 {
				log.Infof("RepeantSendLogToChain Disabled")
				time.Sleep(time.Second * time.Duration(self.Cfg.Layer2RecordInterval*10))
				continue
			}

			data, err := client.Get(self.Cfg.RequestLogServer)
			if err != nil {
				log.Errorf("RepeantSendLogToChain N.1 %s", err)
				continue
			}

			Rl := []RequestLog{}
			err = json.Unmarshal(data, &Rl)
			if err != nil {
				log.Errorf("RepeantSendLogToChain N.2 %s", err)
				continue
			}

			for _, r := range Rl {
				rlog, err := json.Marshal(r)
				if err != nil {
					log.Errorf("RepeantSendLogToChain N.2 %s", err)
					continue
				}

				sum := sha256.Sum256(rlog)
				_, err = layer2Sdk.NeoVM.InvokeNeoVMContract(uint64(self.Cfg.GasPrice), 200000, nil, self.Cfg.AdminAccount, contractAddr, []interface{}{"StoreHash", []interface{}{fmt.Sprintf("%x", sum)}})
				if err != nil {
					log.Errorf("RepeantSendLogToChain N.2 %s. %s: %x", err, string(rlog), sum)
					continue
				}
				count++
				if count%uint64(self.Cfg.Layer2RecordBatchC) == 0 {
					time.Sleep(time.Second * time.Duration(self.Cfg.Layer2RecordInterval))
				}

				if count%uint64(self.Cfg.Layer2RecordBatchC*20) == 0 {
					log.Infof("RepeantSendLogToChainY.3 %x : %s", sum, string(rlog))
				}
			}

			time.Sleep(time.Second * time.Duration(self.Cfg.Layer2RecordInterval))
		}
	}
}

func GetCommitedLayer2Height(ontSdk *ontSdk.OntologySdk, contract common.Address) (uint32, error) {
	tx, err := ontSdk.NeoVM.NewNeoVMInvokeTransaction(0, 0, contract, []interface{}{"getCurrentHeight", []interface{}{}})
	if err != nil {
		return 0, err
	}
	result, err := ontSdk.PreExecTransaction(tx)
	if err != nil {
		fmt.Printf("PreExecTransaction failed! err: %s", err.Error())
		return 0, err
	}
	if result == nil {
		fmt.Printf("can not find the result")
		return 0, fmt.Errorf("can not find current height!")
	}
	height, err := result.Result.ToInteger()
	if err != nil {
		return 0, fmt.Errorf("current height is not right!")
	}
	return uint32(height.Uint64()), nil
}

func GetCommitedLayer2StateByHeight(ontSdk *ontSdk.OntologySdk, contract common.Address, height uint32) ([]byte, uint32, error) {
	tx, err := ontSdk.NeoVM.NewNeoVMInvokeTransaction(0, 0, contract, []interface{}{"getStateRootByHeight", []interface{}{height}})
	if err != nil {
		fmt.Printf("new transaction failed!")

	}
	result, err := ontSdk.PreExecTransaction(tx)
	if err != nil {
		fmt.Printf("PreExecTransaction failed! err: %s", err.Error())
		return nil, 0, err

	}
	if result == nil {
		fmt.Printf("can not find the result")
		return nil, 0, fmt.Errorf("can not find state of heigh: %d", height)

	}
	tt, _ := result.Result.ToArray()
	if len(tt) != 3 {
		fmt.Printf("result is not right")
		return nil, 0, fmt.Errorf("result is not right, height: %d", height)

	}
	item0, _ := tt[0].ToString()
	item1, _ := tt[1].ToInteger()
	item2, _ := tt[2].ToInteger()
	fmt.Printf("item0: %s, item1: %d, item2: %d\n", item0, item1, item2)
	stateRoot, err := common.Uint256FromHexString(item0)
	if err != nil {
		return nil, 0, fmt.Errorf("state hash is not right, height: %d", height)

	}
	return stateRoot.ToArray(), uint32(item1.Uint64()), nil
}

type VerifyService struct {
	Cfg *layer2config.Config
}

var DefVerifyService VerifyService

type VerifyResult struct {
	Exist            bool   `json:"exist"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Proof            string `json:"proof"`
	Layer2Height     uint32 `json:"layer2Height"`
	WitnessStateRoot string `json:"witnessStateRoot"`
	WitnessContract  string `json:"witnessContract"`
}

// verify the store
func (self *VerifyService) VerifyHashCore(hash string) (*VerifyResult, error) {
	layer2Sdk := self.Cfg.Layer2Sdk
	ontSdk := self.Cfg.OntSdk

	// 1. get the store key
	//    get the store data, store proof by the key
	key, _ := layer2Sdk.GetLayer2StoreKey(self.Cfg.Layer2Contract, []byte(hash))
	store, err := layer2Sdk.GetLayer2StoreProof(key)
	if err != nil {
		log.Errorf("VerifyHashCore N.0 key: %s. %s", hash, err)
		return nil, err
	}

	log.Infof("verify key: %s , value: %s, proof: %s, layer2height: %d", hash, store.Value, store.Proof, store.Height)

	// 2. ensure the state root of the store is commited to ontology
	contractAddress, err := common.AddressFromHexString(self.Cfg.Layer2WitContract)
	if err != nil {
		log.Errorf("VerifyHashCore N.1 key: %s. %s", hash, err)
		return nil, err
	}

	var count uint32
	for {
		if count > uint32(self.Cfg.Layer2RetryCount) {
			return nil, fmt.Errorf("VerifyHashCore Retry over times")
		}

		count++
		curHeight, err := GetCommitedLayer2Height(ontSdk, contractAddress)
		if err != nil {
			log.Errorf("VerifyHashCore N.3 key: %s. %s", hash, err)
			return nil, err
		}

		if curHeight < store.Height {
			log.Infof("VerifyHashCore N.3.0 : %s.  wait layer2 relayer commit layer2 block to height %d. currHeight: %d", hash, store.Height, curHeight)
			time.Sleep(time.Second * 1)
			continue
		}
		break
	}

	// 3. get the state root which is commited to ontology
	stateRoot, height, err := GetCommitedLayer2StateByHeight(ontSdk, contractAddress, store.Height)
	if err != nil {
		log.Errorf("VerifyHashCore N.4 key: %s. %s", hash, err)
		return nil, err
	}
	log.Infof("state root: %s, height: %d\n", hex.EncodeToString(stateRoot), height)

	// 4. verify the data is stored through the store proof and state root
	proof_byte, err := hex.DecodeString(store.Proof)
	if err != nil {
		log.Errorf("VerifyHashCore N.5 key: %s. %s", hash, err)
		return nil, err
	}
	value_bytes, err := hex.DecodeString(store.Value)
	if err != nil {
		log.Errorf("VerifyHashCore N.6 key: %s. %s", hash, err)
		return nil, err
	}
	result, err := layer2Sdk.VerifyLayer2StoreProof(key, value_bytes, proof_byte, stateRoot)
	if err != nil {
		log.Errorf("VerifyHashCore N.7 key: %s. %s", hash, err)
		return nil, err
	}

	return &VerifyResult{
		Exist:            result,
		Key:              hash,
		Value:            store.Value,
		Proof:            store.Proof,
		Layer2Height:     store.Height,
		WitnessStateRoot: fmt.Sprintf("%x", stateRoot),
		WitnessContract:  self.Cfg.Layer2WitContract,
	}, nil
}
