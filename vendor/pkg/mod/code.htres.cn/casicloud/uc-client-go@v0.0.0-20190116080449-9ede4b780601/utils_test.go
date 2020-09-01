package uclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_utils_signStringQuery(t *testing.T) {
	query := "abc=1&bcd=2&ts=10000"
	signed, err := signStringQuery(query, "")
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, signed, "abc=1&bcd=2&ts=10000&sign=69d8f9b51ddf80391f2704c413405c15")
}
