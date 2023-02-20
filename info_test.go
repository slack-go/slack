package slack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Which tests do we want?
func TestJSONTime_UnmarshalJSONUpstream(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		wantTr  JSONTime
		wantErr bool
	}{
		{
			"acceptable int64 timestamp",
			args{[]byte(`1643435556`)},
			JSONTime(1643435556),
			false,
		},
		{
			"acceptable string timestamp",
			args{[]byte(`"1643435556"`)},
			JSONTime(1643435556),
			false,
		},
		{
			"null",
			args{[]byte(`null`)},
			JSONTime(0),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tr JSONTime
			if err := tr.UnmarshalJSON(tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("JSONTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
