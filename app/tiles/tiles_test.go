package tiles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PrepareTiles(t *testing.T) {
	err := PrepareTiles("../loader/images", "50x50", 50)

	require.NoError(t, err)
}
