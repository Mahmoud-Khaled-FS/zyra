package zyra

import (
	"fmt"
	"os"
	"sync"

	"github.com/Mahmoud-Khaled-FS/zyra/internal/assert/builtin"
	"github.com/Mahmoud-Khaled-FS/zyra/internal/parser"
)

const configFileName = "zyra.config"

type RunOption struct {
	Path       string
	ConfigPath string
	NoTest     bool
}

func Run(options RunOption) error {
	stat, err := os.Stat(options.Path)
	if err != nil {
		return err
	}

	builtin.InitBuiltin()

	if stat.IsDir() {
		return RunDir(options)
	}
	return RunFile(options)

}

func RunFile(options RunOption) error {
	var config *parser.Config = nil

	if options.ConfigPath != "" {
		var err error
		config, err = loadConfig(options.ConfigPath)
		if err != nil {
			return err
		}
	}

	doc, err := loadDoc(options.Path)
	if err != nil {
		return err
	}

	z := NewZyra(config, options.NoTest)
	r, err := z.Process(ZyraFile{
		File: options.Path,
		Doc:  doc,
	})
	if err != nil {
		return err
	}

	results := make([]ZyraResult, 1)
	results[0] = r

	BeautyLogger(results)

	return err
}

func RunDir(options RunOption) error {
	zDir, err := loadDir(options.Path)
	if err != nil {
		return err
	}

	if options.ConfigPath != "" {
		zDir.configPath = options.ConfigPath
	}

	var config *parser.Config = nil

	if zDir.configPath != "" {
		config, err = loadConfig(zDir.configPath)
		if err != nil {
			return err
		}
	}

	results, err := runDirSync(zDir, config, options.NoTest)
	if err != nil {
		return err
	}

	BeautyLogger(results)
	return nil
}

func runDirSync(zd *ZyraDir, config *parser.Config, noTest bool) ([]ZyraResult, error) {

	z := NewZyra(config, noTest)
	results := make([]ZyraResult, len(zd.files))

	for i, f := range zd.files {
		r, err := z.Process(f)
		if err != nil {
			return results, fmt.Errorf("file %s: %w", f.File, err)
		}
		results[i] = r
	}
	return results, nil
}

func RunDirConcurrent(zd *ZyraDir, config *parser.Config, noTest bool) ([]ZyraResult, error) {
	z := NewZyra(config, noTest)
	results := make([]ZyraResult, len(zd.files))

	var wg sync.WaitGroup
	errCh := make(chan error, len(zd.files))
	resCh := make(chan struct {
		Index  int
		Result ZyraResult
	}, len(zd.files))

	for i, f := range zd.files {
		wg.Add(1)
		go func(idx int, file ZyraFile) {
			defer wg.Done()
			r, err := z.Process(file)
			if err != nil {
				errCh <- fmt.Errorf("file %s: %w", file.File, err)
				return
			}
			resCh <- struct {
				Index  int
				Result ZyraResult
			}{Index: idx, Result: r}
		}(i, f)
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(resCh)
	}()

	if len(errCh) > 0 {
		return results, <-errCh
	}

	for r := range resCh {
		results[r.Index] = r.Result
	}
	return results, nil
}
