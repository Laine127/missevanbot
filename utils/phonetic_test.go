package utils

import "testing"

func TestPinyin(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"test1", "你好", "nǐ hǎo\n[nǐ] [hǎo hào]"},
		{"test2", "", ""},
		{"test3", "干物妹", "gàn wù mèi\n[gàn gān] [wù] [mèi]"},
		{"test4", "你好 world", "nǐ hǎo\n[nǐ] [hǎo hào]"},
		{"test4", "宫商角徵羽", "gōng shāng jiǎo zhēng yǔ\n[gōng] [shāng] [jiǎo jué] [zhēng zhǐ] [yǔ hù]"},
		{"test5", "world", ""},
		{"test5", "world 你好", "nǐ hǎo\n[nǐ] [hǎo hào]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pinyin(tt.arg); got != tt.want {
				t.Errorf("Pinyin() = %v, want %v", got, tt.want)
			}
		})
	}
}
