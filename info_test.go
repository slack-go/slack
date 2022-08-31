package slack

import (
	"testing"
)

func TestJSONTime_UnmarshalJSON(t *testing.T) {
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
