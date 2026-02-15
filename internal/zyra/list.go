package zyra

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/logger"
)

type ListZyraFilesOptions struct {
	Path        string
	ListCount   bool
	ListJSON    bool
	ListAbs     bool
	ListPattern string
}

func ListZyraFiles(options ListZyraFilesOptions) error {
	zDir, err := loadDir(options.Path)
	if err != nil {
		return err
	}

	if options.ListCount {
		fmt.Printf("Requests: %v\n", len(zDir.files))
		return nil
	}

	list := make([]logger.RequestMeta, 0, len(zDir.files))
	for _, f := range zDir.files {
		if options.ListPattern != "" && !(strings.Contains(f.File, options.ListPattern) || strings.Contains(f.Doc.Path, options.ListPattern)) {
			continue
		}
		list = append(list, logger.RequestMeta{
			FilePath:   f.File,
			Method:     f.Doc.Method,
			URL:        f.Doc.Path,
			Assertions: len(f.Doc.Assertions),
			HasHeaders: len(f.Doc.Headers) > 0,
			HasVars:    len(f.Doc.Vars) > 0,
			HasBody:    len(f.Doc.Body) > 0,
		})
	}

	if options.ListJSON {
		return printListJson(list)
	}

	logger.PrintList(list)
	return nil
}

func printListJson(files []logger.RequestMeta) error {
	out, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
