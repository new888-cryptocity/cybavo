package cybavo

import (
	"fmt"
	"strings"

	"666.com/gameserver/framework/mlog"
	ethABI "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// 建立OroderID
func GenOrderID(gameID int32, groupID int32, tableNo int32, shoeNo int32) string {
	return fmt.Sprintf("%s_%d_%d_%d_%d", ORDER_ID_PREFIX, gameID, groupID, tableNo, shoeNo)
}

// 取得 setResult 的網址
func GenSetResultURL(wallet_id string) string {
	url := Moke_Server_URL + "/v1/mock/wallets/" + wallet_id + "/sender/transactions"
	return url
}

// 產出 setResult 中需要帶入的參數 contract_abi 資料
func GenSetResultContractABIData(url string, hex string) (string, error) {
	mlog.Debug("GenSetResultContractABIData:\nurl=%s\nhex=%s\n", url, hex)
	output := ""
	st, err := ethABI.NewType("string", "", nil)
	if err != nil {
		mlog.Error(err.Error())
		return output, err
	}

	i256, err := ethABI.NewType("uint256", "", nil)
	if err != nil {
		mlog.Error(err.Error())
		return output, err
	}

	as := ethABI.Arguments{
		{Name: "_link", Type: st},
		{Name: "_hash", Type: i256},
	}
	hexbytes := common.Hex2Bytes(hex)
	bi := common.BytesToHash(hexbytes).Big()
	args := make([]interface{}, 2)
	args[0] = url
	args[1] = bi
	bs, err := as.PackValues(args)

	//convert form bytes to hex string
	output = fmt.Sprintf("setResult:0x%x", bs)

	return output, err
}

// 取得 getResult 的網址
func GenGetResultURL(wallet_id string, contract_address string, data string) string {
	url := Moke_Server_URL + "/v1/mock/wallets/" + wallet_id + "/contract/read?contract=" + contract_address + "&data=" + data
	return url
}

// 產出 getResult 中所需要的 data 資料
func GenGetResoultData(url string) string {
	output := ""
	definition := `[{"inputs":[{"name":"_link","type":"string"}],"name":"getResult","type":"function"}]`
	abi, err := ethABI.JSON(strings.NewReader(definition))
	if err != nil {
		mlog.Error(err.Error())
		return output
	}

	bs, err := abi.Pack(API_GET_RESULT, url)
	if err != nil {
		mlog.Error(err.Error())
		return output
	}

	//convert form bytes to hex string
	output = fmt.Sprintf("%x", bs)
	return output
}

func ParserSetReoultResponse(content string) {

}
