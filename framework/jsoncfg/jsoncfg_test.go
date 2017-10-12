package jsoncfg

import (
	"fmt"
	"os"
	"testing"
)

type Config struct {
	A int `json:"a"`
	B string
	C float64
}

func TestLoadConfig(t *testing.T) {
	cfg0 := Config{
		A: 1,
		B: "hello",
		C: 3.1415926,
	}

	if err := WriteToFile("cfg.json", &cfg0); err != nil {
		t.Errorf("write config failed: err[%v]", err)
		return
	}

	var cfg1 Config
	if err := LoadFromFile("cfg.json", &cfg1); err != nil {
		t.Errorf("load config failed: err[%v]", err)
		return
	}

	if cfg0 != cfg1 {
		t.Errorf("%v != %v", cfg0, cfg1)
	}

	os.Remove("cfg.json")
}

func ExampleLoadFromFile() {
	type Config struct {
		A int
		B string
		C float64
	}

	if err := WriteToFile("example.json", &Config{A: 1, B: "hello", C: 3.14}); err != nil {
		fmt.Printf("write to file: %v", err)
		return
	}

	var cfg Config
	if err := LoadFromFile("example.json", &cfg); err != nil {
		fmt.Printf("load from file: %v", err)
		return
	}

	fmt.Println(cfg)

	os.Remove("example.json")

	// output:
	// {1 hello 3.14}
}
