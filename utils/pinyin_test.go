package utils

import "testing"

func TestPinyin(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"test1", "你好", "nǐ hǎo"},
		{"test2", "", ""},
		{"test3", "干物妹", "gàn wù mèi"},
		{"test4", "你好 world", "nǐ hǎo"},
		{"test5", "world", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pinyin(tt.arg); got != tt.want {
				t.Errorf("Pinyin() = %v, want %v", got, tt.want)
			}
		})
	}
}
