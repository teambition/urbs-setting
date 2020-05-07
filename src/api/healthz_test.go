package api

import (
	"fmt"
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
)

func TestHealthzAPIs(t *testing.T) {
	tt, cleanup := SetUpTestTools()
	defer cleanup()

	t.Run(`"GET /healthz" should work`, func(t *testing.T) {
		assert := assert.New(t)

		res, err := request.Get(fmt.Sprintf("%s/healthz", tt.Host)).End()
		assert.Nil(err)
		assert.Equal(200, res.StatusCode)

		json := map[string]interface{}{}
		res.JSON(&json)
		assert.True(json["dbConnect"].(bool))
	})
}
