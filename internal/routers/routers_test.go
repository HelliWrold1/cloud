package routers

import (
	"context"
	"testing"
	"time"

	"github.com/HelliWrold1/cloud/configs"
	"github.com/HelliWrold1/cloud/internal/config"

	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	err := config.Init(configs.Path("cloud.yml"))
	if err != nil {
		t.Fatal(err)
	}

	config.Get().App.EnableMetrics = false
	config.Get().App.EnableTrace = true
	config.Get().App.EnableHTTPProfile = true
	config.Get().App.EnableLimit = true
	config.Get().App.EnableCircuitBreaker = true

	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		gin.SetMode(gin.ReleaseMode)
		r := NewRouter()
		assert.NotNil(t, r)
		cancel()
	})
}

func TestNewRouter2(t *testing.T) {
	err := config.Init(configs.Path("cloud.yml"))
	if err != nil {
		t.Fatal(err)
	}

	config.Get().App.EnableMetrics = true

	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		gin.SetMode(gin.ReleaseMode)
		r := NewRouter()
		assert.NotNil(t, r)
		cancel()
	})
}

type mock struct{}

func (u mock) Create(c *gin.Context)      { return }
func (u mock) DeleteByID(c *gin.Context)  { return }
func (u mock) DeleteByIDs(c *gin.Context) { return }
func (u mock) UpdateByID(c *gin.Context)  { return }
func (u mock) GetByID(c *gin.Context)     { return }
func (u mock) ListByIDs(c *gin.Context)   { return }
func (u mock) List(c *gin.Context)        { return }

type mqttMock struct {
}

func (m mqttMock) Publish(ctx *gin.Context)   { return }
func (m mqttMock) Subscribe(ctx *gin.Context) { return }

func Test_frameRouter(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	frameRouter(r.Group("/"), &mock{})
}

func Test_downlinkRouter(f *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	downlinkRouter(r.Group("/"), &mock{})
}
func Test_MQTTRouter(f *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	MQTTRouter(r.Group("/"), &mqttMock{})
}
