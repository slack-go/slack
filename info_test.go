package slack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSONTime(t *testing.T) {
	{
		var res JSONTime
		err := (&res).UnmarshalJSON([]byte("null"))

		assert.NoError(t, err, fmt.Sprintf("Unexpected error %s", err))
		assert.Equal(t, JSONTime(0), res, "Expected 'null' to parse to an empty time")
	}

	{
		// parsing a number
		var res JSONTime
		err := (&res).UnmarshalJSON([]byte("1668593797"))

		assert.NoError(t, err, fmt.Sprintf("Unexpected error %s", err))
		assert.Equal(t, JSONTime(1668593797), res)
	}
}
