package zyra

import (
	"os"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/assert/builtin"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
)

type RunRequestFileOption struct {
	FilePath   string
	ConfigPath string
}

func RunRequestFile(options RunRequestFileOption) error {
	builtin.InitBuiltin()

	var config *parser.Config = nil

	if options.ConfigPath != "" {
		bytesConfig, err := os.ReadFile(options.ConfigPath)
		if err != nil {
			return err
		}
		config, err = parser.ParseConfig(string(bytesConfig))
		if err != nil {
			return err
		}
	}

	bytes, err := os.ReadFile(options.FilePath)
	if err != nil {
		return err
	}

	doc, err := parser.ParseDocument(string(bytes))
	if err != nil {
		return err
	}

	z := NewZyra(config)
	return z.Process(doc)
}
