package env

import (
	"os"
	"reflect"
	"testing"
)

func TestGetEnv(t *testing.T) {
	type args struct {
		key string
		def any
	}
	tests := []struct {
		name string
		set  string
		args args
		want any
	}{
		{
			name: "integer",
			set:  "10",
			args: args{
				key: "number",
				def: 1,
			},
			want: 10,
		},
		{
			name: "integer empty",
			set:  "",
			args: args{
				key: "number",
				def: 1,
			},
			want: 1,
		},
		{
			name: "boolean",
			set:  "true",
			args: args{
				key: "boolean",
				def: false,
			},
			want: true,
		},
		{
			name: "boolean empty",
			set:  "",
			args: args{
				key: "boolean",
				def: false,
			},
			want: false,
		},
		{
			name: "string",
			set:  "http://muchlis.dev",
			args: args{
				key: "url",
				def: "http://muchlis.dev/test",
			},
			want: "http://muchlis.dev",
		},
		{
			name: "string empty",
			set:  "",
			args: args{
				key: "url",
				def: "http://muchlis.dev/test",
			},
			want: "http://muchlis.dev/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.args.key, tt.set)

			switch x := tt.args.def.(type) {
			case int:
				if got := Get(tt.args.key, x); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			case bool:
				if got := Get(tt.args.key, x); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			case string:
				if got := Get(tt.args.key, x); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}
