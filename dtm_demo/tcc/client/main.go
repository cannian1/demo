package main

import (
	"demo/dtm_demo/tcc/dto"
	"fmt"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const (
	// DtmServer       = "http://dtm.shop-cluster.svc.cluster.local:36789/api/dtmsvr"
	DtmServer       = "http://127.0.0.1:36789/api/dtmsvr"
	BusinessServer1 = "http://localhost:8080"
	BusinessServer2 = "http://localhost:8081"
)

func main() {
	// TccGlobalTransaction 会开启一个全局事务
	gid, _ := uuid.NewUUID()
	err := dtmcli.TccGlobalTransaction(DtmServer, gid.String(), func(tcc *dtmcli.Tcc) (*resty.Response, error) {
		resp, err := tcc.CallBranch(&dto.ConsumeReq{Amount: 3000}, BusinessServer1+"/try", BusinessServer1+"/confirm", BusinessServer1+"/cancel")
		if err != nil {
			return nil, err
		}

		fmt.Println("resp:", resp)
		return nil, nil
	})
	if err != nil {
		fmt.Println("err:", err)
		return
	}
}
