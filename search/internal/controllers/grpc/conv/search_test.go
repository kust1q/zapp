package conv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kust1q/Zapp/search/internal/controllers/grpc/conv"
)

func TestConv_Search(t *testing.T) {
	ids := []int{1, 2}
	
	uResp := conv.ToSearchUserProtoResponse(ids)
	assert.Equal(t, int64(1), uResp.UserIds[0])
	assert.Equal(t, int64(2), uResp.UserIds[1])

	tResp := conv.ToSearchTweetProtoResponse(ids)
	assert.Equal(t, int64(1), tResp.TweetIds[0])
	assert.Equal(t, int64(2), tResp.TweetIds[1])
}
