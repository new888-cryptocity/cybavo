package cybavo

import (
	"io/ioutil"
	"net/http"
	"strings"

	"666.com/gameserver/framework/config"
	"666.com/gameserver/framework/mlog"
)

var (

	//===== 正式服 預設值 (Server 啟動時會讀取 dat 中的資料覆蓋) =====
	Moke_Server_URL  = "http://127.0.0.1:18889"
	Contract_address = "0xa14g869asd7f4a5s6df456asdf456asdfasdf332" // 合約地址
	API_TOKEN        = "3rM4630assfNAEmUn"                          // TOKEN
	WALLET_ID        = "123456"                                     // 錢包ID
	ORDER_ID_PREFIX  = "123456"                                     // 提領單號前綴

	// API Name
	API_GET_RESULT = "getResult"
	API_SET_RESULT = "setResult"
)

func Init(Config *config.ConfigGroup) {
	Contract_address = Config.GetConfigValue("db", "Blockchain", "contract_address")
	API_TOKEN = Config.GetConfigValue("db", "Blockchain", "api_token")
	WALLET_ID = Config.GetConfigValue("db", "Blockchain", "wallet_id")
	ORDER_ID_PREFIX = Config.GetConfigValue("db", "Blockchain", "order_id_prefix")

	mlog.Info("Contract_address = %s", Contract_address)
	mlog.Info("API_TOKEN = %s", API_TOKEN)
	mlog.Info("WALLET_ID = %s", WALLET_ID)
	mlog.Info("ORDER_ID_PREFIX = %s", ORDER_ID_PREFIX)
}

// 將牌組資料上鏈
func AddingPokerBlocksToChain(orderID string, shoeURL string, pokerHex string) (string, error) {
	mlog.Debug("")
	client := &http.Client{}
	url := GenSetResultURL(WALLET_ID)
	method := "POST"
	order_id := orderID
	address := Contract_address
	amount := "0"
	contract_abi, _ := GenSetResultContractABIData(shoeURL, pokerHex)
	body := "{\"requests\":[{\"order_id\":\"" + order_id + "\",\"address\":\"" + address + "\",\"amount\":\"" + amount + "\",\"contract_abi\":\"" + contract_abi + "\"}]}"
	mlog.Info("Body = %s", body)
	reader := strings.NewReader(body)
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		mlog.Error("AddingPokerBlocksToChain Error %s", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-CODE", API_TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		mlog.Error("AddingPokerBlocksToChain Error! ", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mlog.Error("AddingPokerBlocksToChain Error %s", err)
		return "", err
	}

	// 後臺回應的內容
	mlog.Info(string(content))

	return string(content), nil
}

// 透過牌靴連結取得上鏈後的牌組資料
func GetBlocksResult(shoeURL string) (string, error) {

	client := &http.Client{}
	data := GenGetResoultData(shoeURL)
	url := GenGetResultURL(WALLET_ID, Contract_address, data)
	method := "GET"
	reader := strings.NewReader("")
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		mlog.Error("GetBlocksResult Error %s", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-CODE", API_TOKEN) // API Token

	resp, err := client.Do(req)
	if err != nil {
		mlog.Error("GetBlocksResult Error %s", err)
		return "", err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mlog.Error("GetBlocksResult Error %s", err)
		return "", err
	}

	// 後臺回應的內容
	mlog.Info(string(content))

	return string(content), nil
}
