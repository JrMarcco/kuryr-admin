package passwd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret"
)

const (
	// 字符集定义
	LowerCase    = "abcdefghijklmnopqrstuvwxyz"
	UpperCase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits       = "0123456789"
	Symbols      = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	DefaultChars = LowerCase + UpperCase + Digits + Symbols
)

var _ secret.Generator = (*Generator)(nil)

type Generator struct {
	// 可选配置
	charset    string
	minLength  int
	mustLower  bool
	mustUpper  bool
	mustDigit  bool
	mustSymbol bool
}

type Option func(*Generator)

// WithCharset 自定义字符集
func WithCharset(charset string) Option {
	return func(g *Generator) {
		g.charset = charset
	}
}

// WithMinLength 最小长度
func WithMinLength(minLength int) Option {
	return func(g *Generator) {
		g.minLength = minLength
	}
}

// WithRequirements 设置密码复杂度要求
func WithRequirements(lower, upper, digit, symbol bool) Option {
	return func(g *Generator) {
		g.mustLower = lower
		g.mustUpper = upper
		g.mustDigit = digit
		g.mustSymbol = symbol
	}
}

// NewGenerator 创建密码生成器
func NewGenerator(opts ...Option) *Generator {
	g := &Generator{
		charset:    DefaultChars,
		minLength:  8,
		mustLower:  true,
		mustUpper:  true,
		mustDigit:  true,
		mustSymbol: false,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Generate 生成随机密码
func (g *Generator) Generate(length int) (string, error) {
	if length < g.minLength {
		return "", fmt.Errorf("password length must be at least %d", g.minLength)
	}

	var password strings.Builder
	password.Grow(length)

	// 如果有复杂度要求，先确保满足每种字符类型
	guaranteed := 0
	if g.mustLower {
		char, err := randomChar(LowerCase)
		if err != nil {
			return "", err
		}
		password.WriteString(char)
		guaranteed++
	}

	if g.mustUpper {
		char, err := randomChar(UpperCase)
		if err != nil {
			return "", err
		}
		password.WriteString(char)
		guaranteed++
	}

	if g.mustDigit {
		char, err := randomChar(Digits)
		if err != nil {
			return "", err
		}
		password.WriteString(char)
		guaranteed++
	}

	if g.mustSymbol {
		char, err := randomChar(Symbols)
		if err != nil {
			return "", err
		}
		password.WriteString(char)
		guaranteed++
	}

	// 填充剩余长度
	for i := guaranteed; i < length; i++ {
		char, err := randomChar(g.charset)
		if err != nil {
			return "", err
		}
		password.WriteString(char)
	}

	// 打乱密码字符顺序
	return shuffle(password.String())
}

// GenerateWithPrefix 带前缀的密码生成
func (g *Generator) GenerateWithPrefix(prefix string, length int) (string, error) {
	prefixLen := len(prefix)
	if prefixLen >= length {
		return "", fmt.Errorf("prefix length must be less than total length")
	}

	password, err := g.Generate(length - prefixLen - 1) // -1 for separator
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s_%s", prefix, password), nil
}

// randomChar 从字符集中随机选择一个字符
func randomChar(charset string) (string, error) {
	if len(charset) == 0 {
		return "", fmt.Errorf("charset is empty")
	}

	max := big.NewInt(int64(len(charset)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	return string(charset[n.Int64()]), nil
}

// shuffle 打乱字符串顺序
func shuffle(s string) (string, error) {
	if len(s) <= 1 {
		return s, nil
	}

	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		max := big.NewInt(int64(i + 1))
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("failed to shuffle: %w", err)
		}
		j := n.Int64()
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}

// SimpleGenerate 简单快速生成密码的函数
func SimpleGenerate(length int) (string, error) {
	g := NewGenerator()
	return g.Generate(length)
}

// SecureGenerate 生成安全密码（包含所有字符类型）
func SecureGenerate(length int) (string, error) {
	g := NewGenerator(WithRequirements(true, true, true, true))
	return g.Generate(length)
}
