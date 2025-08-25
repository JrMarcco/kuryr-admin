package passwd_test

import (
	"fmt"
	"log"

	"github.com/JrMarcco/kuryr-admin/internal/pkg/secret/passwd"
)

func ExampleSimpleGenerate() {
	// 简单生成16位密码
	pwd, err := passwd.SimpleGenerate(16)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated password length: %d\n", len(pwd))
	// Output: Generated password length: 16
}

func ExampleSecureGenerate() {
	// 生成包含所有字符类型的安全密码
	pwd, err := passwd.SecureGenerate(20)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated secure password length: %d\n", len(pwd))
	// Output: Generated secure password length: 20
}

func ExampleNewGenerator_basic() {
	// 创建默认生成器
	g := passwd.NewGenerator()

	pwd, err := g.Generate(12)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated password length: %d\n", len(pwd))
	// Output: Generated password length: 12
}

func ExampleNewGenerator_customCharset() {
	// 只使用字母和数字
	g := passwd.NewGenerator(
		passwd.WithCharset(passwd.LowerCase+passwd.UpperCase+passwd.Digits),
		passwd.WithRequirements(true, true, true, false), // 不要求特殊字符
	)

	pwd, err := g.Generate(16)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated alphanumeric password length: %d\n", len(pwd))
	// Output: Generated alphanumeric password length: 16
}

func ExampleNewGenerator_onlyDigits() {
	// 生成纯数字密码（如PIN码）
	g := passwd.NewGenerator(
		passwd.WithCharset(passwd.Digits),
		passwd.WithRequirements(false, false, false, false),
		passwd.WithMinLength(4),
	)

	pin, err := g.Generate(6)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated PIN length: %d\n", len(pin))
	// Output: Generated PIN length: 6
}

func ExampleGenerator_GenerateWithPrefix() {
	g := passwd.NewGenerator()

	// 生成带前缀的密码，适用于临时密码等场景
	pwd, err := g.GenerateWithPrefix("TEMP", 24)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Password starts with TEMP_: %v\n", len(pwd) >= 5 && pwd[:5] == "TEMP_")
	// Output: Password starts with TEMP_: true
}
