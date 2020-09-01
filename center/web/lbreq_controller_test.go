package web

import (
	"bytes"
	"code.htres.cn/casicloud/alb/center/common"
	"code.htres.cn/casicloud/alb/pkg/model"
	"github.com/gin-gonic/gin/json"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var req = model.LBRequest{
	User:   model.NewADCString("12345"),
	Action: model.ActionAdd,
	Policy: model.LBPolicy{
		Record: model.LBRecord{
			IP:   model.NewADCString("192.168.100.200"),
			Port: 80,
			Type: model.TypeIP},
		Endpoints: []model.RealServer{
			{Name: "sever1", IP: "106.74.100.99", Port: 80},
			{Name: "sever2", IP: "106.74.100.98", Port: 80},
			{Name: "sever3", IP: "106.74.100.97", Port: 80},
		},
	},
}

func TestLBRequestController(t *testing.T) {

	common.GlobalConfig = common.Config{
		WorkDir:  ".",
		LogLevel: 0,

		AuditLogPath: ".adc",
		SysLogPath:   ".alb",

		AgentTimeout: 5 * time.Second,

		Port: 8080,

		DBArgs:        ":memory:",
		Dialect:       "sqlite3",
		EtcdEndpoints: []string{"localhost:1234"},
	}

	router, err := SetupRouter()
	assert.NoError(t, err)

	t.Run("test accept", func(t *testing.T) {
		w := httptest.NewRecorder()

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		req, _ := http.NewRequest("POST", "/1/lb", bytes.NewReader(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, `{"request_id":"0"}`, w.Body.String())
	})

	t.Run("test query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/1/lb?request_id=12345", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})
}
