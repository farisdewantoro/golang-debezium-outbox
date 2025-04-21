package util

import (
	"testing"
)

type testStruct struct {
	ID   int
	Name string
}

func TestGetIndexWithFieldValue(t *testing.T) {
	tests := []struct {
		name        string
		arr         []testStruct
		fieldName   string
		targetValue interface{}
		wantIndex   int
		wantErr     bool
	}{
		{
			name: "found in slice",
			arr: []testStruct{
				{ID: 1, Name: "first"},
				{ID: 2, Name: "second"},
				{ID: 3, Name: "third"},
			},
			fieldName:   "ID",
			targetValue: 2,
			wantIndex:   1,
			wantErr:     false,
		},
		{
			name: "not found in slice",
			arr: []testStruct{
				{ID: 1, Name: "first"},
				{ID: 2, Name: "second"},
			},
			fieldName:   "ID",
			targetValue: 5,
			wantIndex:   -1,
			wantErr:     true,
		},
		{
			name: "string field search",
			arr: []testStruct{
				{ID: 1, Name: "first"},
				{ID: 2, Name: "second"},
			},
			fieldName:   "Name",
			targetValue: "second",
			wantIndex:   1,
			wantErr:     false,
		},
		{
			name:        "empty slice",
			arr:         []testStruct{},
			fieldName:   "ID",
			targetValue: 1,
			wantIndex:   -1,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, err := GetIndexWithFieldValue(tt.arr, tt.fieldName, tt.targetValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndexWithFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIndex != tt.wantIndex {
				t.Errorf("GetIndexWithFieldValue() = %v, want %v", gotIndex, tt.wantIndex)
			}
		})
	}
}

func TestRemoveAtIndex(t *testing.T) {
	tests := []struct {
		name    string
		slice   []int
		index   int
		want    []int
		wantErr bool
	}{
		{
			name:    "valid index middle",
			slice:   []int{1, 2, 3, 4, 5},
			index:   2,
			want:    []int{1, 2, 4, 5},
			wantErr: false,
		},
		{
			name:    "valid index first",
			slice:   []int{1, 2, 3},
			index:   0,
			want:    []int{2, 3},
			wantErr: false,
		},
		{
			name:    "valid index last",
			slice:   []int{1, 2, 3},
			index:   2,
			want:    []int{1, 2},
			wantErr: false,
		},
		{
			name:    "index out of range",
			slice:   []int{1, 2, 3},
			index:   5,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative index",
			slice:   []int{1, 2, 3},
			index:   -1,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "nil slice",
			slice:   nil,
			index:   0,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoveAtIndex(tt.slice, tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveAtIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("RemoveAtIndex() got length = %v, want length %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("RemoveAtIndex() got[%d] = %v, want[%d] %v", i, got[i], i, tt.want[i])
					}
				}
			}
		})
	}
}

func TestSliceHas(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		x     int
		want  bool
	}{
		{
			name:  "element exists",
			slice: []int{1, 2, 3, 4, 5},
			x:     3,
			want:  true,
		},
		{
			name:  "element does not exist",
			slice: []int{1, 2, 3, 4, 5},
			x:     6,
			want:  false,
		},
		{
			name:  "empty slice",
			slice: []int{},
			x:     1,
			want:  false,
		},
		{
			name:  "nil slice",
			slice: nil,
			x:     1,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceHas(tt.slice, tt.x); got != tt.want {
				t.Errorf("SliceHas() = %v, want %v", got, tt.want)
			}
		})
	}
}
