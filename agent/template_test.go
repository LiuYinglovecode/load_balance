package agent

import (
	"bytes"
	"os"
	"testing"

	"code.htres.cn/casicloud/alb/pkg/model"
)

func TestWriteHaproxyCfg(t *testing.T) {
	test := model.LBPolicy{
		Record: model.LBRecord{
			IP:   model.NewADCString("192.168.100.200"),
			Port: 80,
			Type: model.TypeIP},
		Endpoints: []model.RealServer{
			{Name: "sever1", IP: "106.74.100.99", Port: 80},
			{Name: "sever2", IP: "106.74.100.98", Port: 80},
			{Name: "sever3", IP: "106.74.100.97", Port: 80},
		},
	}
	WriteHaproxyCfg(os.Stdout, test)
}

func TestWriteKeepalivedCfg(t *testing.T) {
	config := KeepalivedConfig{
		INet:            "127.0.0.1",
		VirutalRouterID: "51",
		State:           "MASTER",
		Priority:        100,
		UnicastSrcIP:    "127.0.0.1",
		UnicastPeer: []string{
			"127.0.0.2",
			"127.0.0.3",
		},
		VirtualIP: "8.8.8.8",
	}

	var buf bytes.Buffer
	WriteKeepalivedCfg(&buf, config)
	t.Logf("cfg is: %s", buf.String())
}
