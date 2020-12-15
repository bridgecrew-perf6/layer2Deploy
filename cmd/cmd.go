package cmd

import (
	"encoding/json"
	"fmt"
	common2 "github.com/ontio/layer2deploy/common"
	"github.com/ontio/layer2deploy/layer2config"
	ontSdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/common/password"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

func SetOntologyConfig(ctx *cli.Context) error {
	cf := ctx.String(GetFlagName(ConfigfileFlag))
	if _, err := os.Stat(cf); os.IsNotExist(err) {
		// if there's no config file, use default config
		return err
	}

	file, err := os.Open(cf)
	if err != nil {
		return err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	cfg := &layer2config.Config{}
	err = json.Unmarshal(bs, cfg)
	if err != nil {
		return err
	}
	*layer2config.DefLayer2Config = *cfg

	if layer2config.DefLayer2Config.WalletName == "" || layer2config.DefLayer2Config.Layer2MainNetNode == "" || layer2config.DefLayer2Config.Layer2TestNetNode == "" {
		return fmt.Errorf("walletName/layer2MainNetAddress/layer2TestNetAddress  is nil")
	}

	wallet, err := ontSdk.OpenWallet(layer2config.DefLayer2Config.WalletName)
	if err != nil {
		return err
	}
	passwd, err := password.GetAccountPassword()
	if err != nil {
		return err
	}
	sagaAccount, err := wallet.GetDefaultAccount(passwd)
	if err != nil {
		return err
	}

	layer2Sdk := ontSdk.NewLayer2Sdk()
	switch layer2config.DefLayer2Config.NetWorkId {
	case layer2config.NETWORK_ID_MAIN_NET:
		log.Infof("currently Main net")
		layer2Sdk.NewRpcClient().SetAddress(layer2config.DefLayer2Config.Layer2MainNetNode)
	case layer2config.NETWORK_ID_POLARIS_NET:
		log.Infof("currently test net")
		layer2Sdk.NewRpcClient().SetAddress(layer2config.DefLayer2Config.Layer2TestNetNode)
	case layer2config.NETWORK_ID_SOLO_NET:
		log.Infof("currently solo net")
		// solo simulation with test net. but different contract and owner
		layer2Sdk.NewRpcClient().SetAddress(layer2config.DefLayer2Config.Layer2TestNetNode)
	default:
		return fmt.Errorf("error network id %d", layer2config.DefLayer2Config.NetWorkId)
	}

	layer2config.DefLayer2Config.Layer2Sdk = layer2Sdk
	layer2config.DefLayer2Config.AdminAccount = sagaAccount
	return CheckLayer2InitAddress()
}

func CheckLayer2InitAddress() error {
	layer2Sdk := layer2config.DefLayer2Config.Layer2Sdk
	if layer2Sdk == nil {
		return fmt.Errorf("layer2 sdk should not be nil")
	}

	if len(layer2config.DefLayer2Config.Layer2Contract) == 0 {
		return fmt.Errorf("layer2 contract address or contract not init")
	}

	log.Infof("layer2Contract %s", layer2config.DefLayer2Config.Layer2Contract)

	contractAddr, err := common.AddressFromHexString(layer2config.DefLayer2Config.Layer2Contract)
	// incase that is a file name.
	if err != nil || len(layer2config.DefLayer2Config.Layer2Contract) != common.ADDR_LEN*2 {
		code, err := ioutil.ReadFile(layer2config.DefLayer2Config.Layer2Contract)
		if err != nil {
			return fmt.Errorf("error in ReadFile: %s, %s\n", layer2config.DefLayer2Config.Layer2Contract, err)
		}

		codeHash := common.ToHexString(code)
		contractAddr, err = utils.GetContractAddress(codeHash)
		if err != nil {
			return fmt.Errorf("error get contract address %s", err)
		}

		payload, err := layer2Sdk.GetSmartContract(contractAddr.ToHexString())
		if payload == nil || err != nil {
			txHash, err := layer2Sdk.NeoVM.DeployNeoVMSmartContract(uint64(layer2config.DefLayer2Config.GasPrice), 200000000, layer2config.DefLayer2Config.AdminAccount, true, codeHash, "witness layer2 contract", "1.0", "witness", "email", "desc")
			if err != nil {
				return fmt.Errorf("deploy contract %s err: %s", layer2config.DefLayer2Config.Layer2Contract, err)
			}

			_, err = common2.GetLayer2EventByTxHash(txHash.ToHexString())
			if err != nil {
				return fmt.Errorf("deploy contract failed %s", err)
			}
			log.Infof("deploy concontract success")
		}

		log.Infof("the contractAddr hexstring is %s", contractAddr.ToHexString())
		layer2config.DefLayer2Config.Layer2Contract = contractAddr.ToHexString()
	}

	contractAddr, err = common.AddressFromHexString(layer2config.DefLayer2Config.Layer2Contract)
	if err != nil {
		return err
	}

	for {
		res, err := layer2Sdk.NeoVM.PreExecInvokeNeoVMContract(contractAddr, []interface{}{"init_status", 2})
		if err != nil {
			return fmt.Errorf("err get init_status %s", err)
		}

		if res.State == 0 {
			return fmt.Errorf("init statuc exec failed state is 0")
		}

		addrB, err := res.Result.ToByteArray()
		if err != nil {
			return fmt.Errorf("error init_status toByteArray %s", err)
		}

		sagaAddrBase58 := layer2config.DefLayer2Config.AdminAccount.Address.ToBase58()
		if len(addrB) != 0 {
			addrO, err := common.AddressParseFromBytes(addrB)
			if err != nil {
				return fmt.Errorf("AddressParseFromBytes err: %s", err)
			}

			log.Infof("layer2 address already init owner to addr %s", addrO.ToBase58())
			if addrO.ToBase58() != sagaAddrBase58 {
				return fmt.Errorf("contract addr not equal. owner is %s. but sagaAccount init to %s", addrO.ToBase58(), sagaAddrBase58)
			}
			break
		} else {
			//log.Infof("start init layer2 addr owner to address %s", sagaAddrBase58)
			txHash, err := layer2Sdk.NeoVM.InvokeNeoVMContract(uint64(layer2config.DefLayer2Config.GasPrice), 200000, nil, layer2config.DefLayer2Config.AdminAccount, contractAddr, []interface{}{"init", layer2config.DefLayer2Config.AdminAccount.Address})
			if err != nil {
				return fmt.Errorf("init layer2 owner err0 %s", err)
			}
			_, err = common2.GetLayer2EventByTxHash(txHash.ToHexString())
			if err != nil {
				return fmt.Errorf("init layer2 owner err1: %s", err)
			}
			log.Infof("init layer2 addr owner to address %s success.", sagaAddrBase58)
		}
	}

	txHash, err := layer2Sdk.NeoVM.InvokeNeoVMContract(uint64(layer2config.DefLayer2Config.GasPrice), 200000, nil, layer2config.DefLayer2Config.AdminAccount, contractAddr, []interface{}{"StoreHash", []interface{}{"6de9439834c9147569741d3c9c9fc010"}})
	if err != nil {
		return fmt.Errorf("StoreUsedNum test failed %s", err)
	}

	_, err = common2.GetLayer2EventByTxHash(txHash.ToHexString())
	if err != nil {
		return fmt.Errorf("init layer2 owner err1: %s", err)
	}
	log.Infof("test StoreUsedNum success ")
	return nil
}

func PrintErrorMsg(format string, a ...interface{}) {
	format = fmt.Sprintf("\033[31m[ERROR] %s\033[0m\n", format) //Print error msg with red color
	fmt.Printf(format, a...)
}
