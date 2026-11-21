//go:build integration

package dvls

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Server(t *testing.T) {
	info, err := testClient.GetServerInfo()
	require.NoError(t, err)
	assert.NotEmpty(t, info.Version)
}
