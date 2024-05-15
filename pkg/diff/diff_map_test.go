package diff

import (
	"encoding/json"
	"reflect"

	"testing"
)

func TestFindDifference(t *testing.T) {
	type args struct {
		input1 map[string]interface{}
		input2 map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantBefore map[string]interface{}
		wantAfter  map[string]interface{}
	}{
		// input1 address.street = "hawaii", address.postal = 8000
		// input2 address.street = "himalayan", address.postal = nil
		// result before & after will record address.street has changed
		{
			name: "Test 1",
			args: args{
				input1: map[string]interface{}{
					"address.street": "hawaii",
					"address.postal": 8000,
				},
				input2: map[string]interface{}{
					"address.street": "himalayan",
					"address.postal": nil,
				},
			},
			wantBefore: map[string]interface{}{
				"address.street": "hawaii",
			},
			wantAfter: map[string]interface{}{
				"address.street": "himalayan",
			},
		},
		{
			name: "Test 2",
			args: args{
				input1: map[string]interface{}{
					"address.street": "hawaii",
					"address.postal": 8000,
					"latitude":       8.0,
				},
				input2: map[string]interface{}{
					"address.street": "himalayan",
					"address.postal": nil,
					// latitude not available here
				},
			},
			wantBefore: map[string]interface{}{
				"address.street": "hawaii",
			},
			wantAfter: map[string]interface{}{
				"address.street": "himalayan",
			},
		},
		{
			name: "Test 3 Nested",
			args: args{
				input1: map[string]interface{}{ // nested input
					"address": map[string]interface{}{
						"street": "belitung",
						"postal": 8000,
					},
				},
				input2: map[string]interface{}{ // also nested input
					"address": map[string]interface{}{
						"street": "kalindo",
						"postal": 3000,
					},
				},
			},
			wantBefore: map[string]interface{}{
				"address.street": "belitung",
				"address.postal": 8000,
			},
			wantAfter: map[string]interface{}{
				"address.street": "kalindo",
				"address.postal": 3000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBefore, gotAfter := FindDifference(tt.args.input1, tt.args.input2)
			if !reflect.DeepEqual(gotBefore, tt.wantBefore) {
				t.Errorf("FindDifference() gotBefore = %v, want %v", gotBefore, tt.wantBefore)
			}
			if !reflect.DeepEqual(gotAfter, tt.wantAfter) {
				t.Errorf("FindDifference() gotAfter = %v, want %v", gotAfter, tt.wantAfter)
			}
		})
	}
}

type User struct {
	ID      int     `json:"id" mapstructure:"id"`
	Name    string  `json:"name" mapstructure:"name"`
	Age     int     `json:"age" mapstructure:"age"`
	Address Address `json:"address" mapstructure:"address"`
}

type Address struct {
	District int      `json:"district" mapstructure:"district"`
	Street   string   `json:"street" mapstructure:"street"`
	Home     Home     `json:"home" mapstructure:"home"`
	Tag      []string `json:"tag" mapstructure:"tag"`
}

type Home struct {
	RT string `json:"rt" mapstructure:"rt"`
}

type UserUpdate struct {
	ID      *int           `json:"id" mapstructure:"id,omitempty"`
	Name    *string        `json:"name" mapstructure:"name,omitempty"`
	Address *AddressUpdate `json:"address" mapstructure:"address,omitempty"`
}

type AddressUpdate struct {
	District *int       `json:"district" mapstructure:"district,omitempty"`
	Street   *string    `json:"street" mapstructure:"street,omitempty"`
	Street2  *string    `json:"street_2" mapstructure:"street_2,omitempty"`
	Home     HomeUpdate `json:"home" mapstructure:"home"`
	Tag      []string   `json:"tag" mapstructure:"tag"`
}

type HomeUpdate struct {
	RT *string `json:"rt" mapstructure:"rt,omitempty"`
}

func TestDiffJSON(t *testing.T) {

	dd := 6

	type args struct {
		existingData interface{}
		updatedData  interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantBefore map[string]interface{}
		wantAfter  map[string]interface{}
		wantErr    bool
	}{
		{
			name: "Test 1",
			args: args{
				existingData: User{
					ID:   8,
					Name: "Muchlis",
					Age:  18,
					Address: Address{District: 900,
						Street: "Belitung",
						Tag:    []string{"RT", "RW"},
					},
				},
				updatedData: UserUpdate{
					ID:   &dd,
					Name: nullString("Muchlisia"),
					Address: &AddressUpdate{
						District: nil,
						Street:   nullString("Belatang"),
						Street2:  nullString("Belatang2"),
						Home: HomeUpdate{
							RT: nullString("BBB"),
						},
						Tag: []string{"RT", "RX", "RCC"},
					},
				},
			},
			wantBefore: map[string]interface{}{
				"address.home.rt":  "",
				"address.street":   "Belitung",
				"address.street_2": "",
				"address.tag[1]":   "RW",
				"address.tag[2]":   "",
				"id":               8,
				"name":             "Muchlis",
			},

			wantAfter: map[string]interface{}{
				"address.home.rt":  "BBB",
				"address.street":   "Belatang",
				"address.street_2": "Belatang2",
				"address.tag[1]":   "RX",
				"address.tag[2]":   "RCC",
				"id":               6,
				"name":             "Muchlisia",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBefore, gotAfter, err := DiffJSON(tt.args.existingData, tt.args.updatedData)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiffJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			printedGotBefore, _ := mapToJSONString(gotBefore)
			printedWantBefore, _ := mapToJSONString(tt.wantBefore)

			if !reflect.DeepEqual(printedGotBefore, printedWantBefore) {
				t.Errorf("DiffJSON()\ngotBefore = %v,\nwantBefore = %v", printedGotBefore, printedWantBefore)
			}

			printedGotAfter, _ := mapToJSONString(gotAfter)
			printedWantAfter, _ := mapToJSONString(tt.wantAfter)

			if !reflect.DeepEqual(printedGotBefore, printedWantBefore) {
				t.Errorf("DiffJSON()\ngotAfter = %v,\nwantAfter = %v", printedGotAfter, printedWantAfter)
			}
		})
	}
}

func nullString(s string) *string {
	return &s
}

// comparing map[string]interface{} is dificult, i use string version to do that
// hopefully map is always sort correctly
func mapToJSONString(m map[string]interface{}) (string, error) {
	// Marshal the map to JSON
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	// Convert JSON bytes to string
	resultString := string(jsonBytes)

	return resultString, nil
}
