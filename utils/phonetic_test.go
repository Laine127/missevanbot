package utils

import "testing"

func TestPinyin(t *testing.T) {
	tests := []struct {
		name  string
		arg   string
		wantA string
		wantB string
	}{
		{"test1", "你好", "nǐ hǎo", "[nǐ] [hǎo hào]"},
		{"test2", "", "", ""},
		{"test3", "干物妹", "gàn wù mèi", "[gàn gān] [wù] [mèi]"},
		{"test4", "你好 world", "nǐ hǎo", "[nǐ] [hǎo hào]"},
		{"test4", "宫商角徵羽", "gōng shāng jiǎo zhēng yǔ", "[gōng] [shāng] [jiǎo jué] [zhēng zhǐ] [yǔ hù]"},
		{"test5", "world", "", ""},
		{"test5", "world 你好", "nǐ hǎo", "[nǐ] [hǎo hào]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotA, gotB := Pinyin(tt.arg); gotA != tt.wantA || gotB != tt.wantB {
				t.Errorf("Pinyin() = %v, %v, want %v %v", gotA, gotB, tt.wantA, tt.wantB)
			}
		})
	}
}
