package convert

import (
	"testing"
	"time"
)

func TestStringEpochToTime(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "Valid epoch time",
			args: args{
				str: "1738138980",
			},
			want:    time.Date(2025, 1, 29, 16, 23, 0, 0, time.FixedZone("WITA", 8*60*60)), // 2025-01-29 16:23:00 +0800 WITA
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringEpochToTime(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringEpochToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("StringEpochToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
