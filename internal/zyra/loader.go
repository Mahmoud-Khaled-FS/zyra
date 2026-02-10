package zyra

import (
	"os"
	"path/filepath"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/utils"
)

type ZyraFile struct {
	File string
	Doc  *parser.Document
}

type ZyraDir struct {
	configPath string
	files      []ZyraFile
}

func loadConfig(path string) (*parser.Config, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parser.ParseConfig(string(data))
}

func loadDoc(path string) (*parser.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parser.ParseDocument(string(data))
}

func loadDir(path string) (*ZyraDir, error) {
	files, err := utils.ReadDirR(path)
	if err != nil {
		return nil, err
	}

	var zd ZyraDir

	for _, f := range files {
		if isConfigFile(path, f) {
			zd.configPath = f
			continue
		}

		doc, err := loadDoc(f)
		if err != nil {
			return nil, err
		}
		zd.files = append(zd.files, ZyraFile{
			File: f,
			Doc:  doc,
		})
	}

	return &zd, nil
}

func isConfigFile(root, path string) bool {
	return path == filepath.Join(root, configFileName)
}
