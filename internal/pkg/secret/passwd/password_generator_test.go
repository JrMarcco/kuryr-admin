package passwd

import (
	"strings"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		length    int
		wantErr   bool
		checkFunc func(string) bool
	}{
		{
			name:      "default generator",
			generator: NewGenerator(),
			length:    16,
			wantErr:   false,
			checkFunc: func(pwd string) bool {
				return len(pwd) == 16 &&
					containsAny(pwd, LowerCase) &&
					containsAny(pwd, UpperCase) &&
					containsAny(pwd, Digits)
			},
		},
		{
			name:      "secure password generator",
			generator: NewGenerator(WithRequirements(true, true, true, true)),
			length:    20,
			wantErr:   false,
			checkFunc: func(pwd string) bool {
				return len(pwd) == 20 &&
					containsAny(pwd, LowerCase) &&
					containsAny(pwd, UpperCase) &&
					containsAny(pwd, Digits) &&
					containsAny(pwd, Symbols)
			},
		},
		{
			name:      "only digits password",
			generator: NewGenerator(WithCharset(Digits), WithRequirements(false, false, false, false)),
			length:    8,
			wantErr:   false,
			checkFunc: func(pwd string) bool {
				return len(pwd) == 8 && containsOnly(pwd, Digits)
			},
		},
		{
			name:      "password too short",
			generator: NewGenerator(),
			length:    4,
			wantErr:   true,
		},
		{
			name:      "custom minimum length",
			generator: NewGenerator(WithMinLength(12)),
			length:    10,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.generator.Generate(tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkFunc != nil && !tt.checkFunc(got) {
				t.Errorf("Generate() = %v, check failed", got)
			}
		})
	}
}

func TestSimpleGenerate(t *testing.T) {
	pwd, err := SimpleGenerate(16)
	if err != nil {
		t.Fatalf("SimpleGenerate() error = %v", err)
	}

	if len(pwd) != 16 {
		t.Errorf("SimpleGenerate() length = %d, want 16", len(pwd))
	}

	// 检查基本复杂度
	if !containsAny(pwd, LowerCase) || !containsAny(pwd, UpperCase) || !containsAny(pwd, Digits) {
		t.Errorf("SimpleGenerate() = %v, missing required character types", pwd)
	}
}

func TestSecureGenerate(t *testing.T) {
	pwd, err := SecureGenerate(20)
	if err != nil {
		t.Fatalf("SecureGenerate() error = %v", err)
	}

	if len(pwd) != 20 {
		t.Errorf("SecureGenerate() length = %d, want 20", len(pwd))
	}

	// 检查所有字符类型
	if !containsAny(pwd, LowerCase) || !containsAny(pwd, UpperCase) ||
		!containsAny(pwd, Digits) || !containsAny(pwd, Symbols) {
		t.Errorf("SecureGenerate() = %v, missing required character types", pwd)
	}
}

func TestGenerateWithPrefix(t *testing.T) {
	g := NewGenerator()
	prefix := "USER"

	pwd, err := g.GenerateWithPrefix(prefix, 20)
	if err != nil {
		t.Fatalf("GenerateWithPrefix() error = %v", err)
	}

	if !strings.HasPrefix(pwd, prefix+"_") {
		t.Errorf("GenerateWithPrefix() = %v, want prefix %s_", pwd, prefix)
	}

	if len(pwd) != 20 {
		t.Errorf("GenerateWithPrefix() length = %d, want 20", len(pwd))
	}
}

func TestUniqueness(t *testing.T) {
	g := NewGenerator()
	passwords := make(map[string]bool)

	// 生成100个密码，检查唯一性
	for i := 0; i < 100; i++ {
		pwd, err := g.Generate(16)
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}

		if passwords[pwd] {
			t.Errorf("Duplicate password generated: %v", pwd)
		}
		passwords[pwd] = true
	}
}

// 辅助函数
func containsAny(s, chars string) bool {
	for _, c := range chars {
		if strings.ContainsRune(s, c) {
			return true
		}
	}
	return false
}

func containsOnly(s, chars string) bool {
	for _, c := range s {
		if !strings.ContainsRune(chars, c) {
			return false
		}
	}
	return true
}
