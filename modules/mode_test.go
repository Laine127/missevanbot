package modules

import (
	"testing"
)

func TestInitMode(t *testing.T) {
	tests := []struct {
		name string
		rid  int
	}{
		{"init", TestRID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitMode(tt.rid)
		})
	}
}

func TestMode(t *testing.T) {
	type args struct {
		rid  int
		mode string
	}
	tests := []struct {
		name string
		args args
	}{
		{"pander", args{TestRID, ModePander}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mode(tt.args.rid, tt.args.mode)
			if err != nil {
				t.Errorf("Mode() error = %v", err)
				return
			}
			t.Logf("mode %s = %v", tt.args.mode, got)
		})
	}
}

func TestSetMode(t *testing.T) {
	type args struct {
		rid  int
		mode string
		val  bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"pander", args{TestRID, ModePander, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetMode(tt.args.rid, tt.args.mode, tt.args.val)
		})
	}
}
