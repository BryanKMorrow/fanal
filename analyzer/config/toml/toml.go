package toml

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/open-policy-agent/conftest/parser/toml"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/config/scanner"
	"github.com/aquasecurity/fanal/types"
)

const version = 1

var requiredExts = []string{".toml"}

type ConfigScanner struct {
	parser *toml.Parser
	scanner.Scanner
}

func NewConfigScanner(filePattern *regexp.Regexp, namespaces, policyPaths, dataPaths []string) (ConfigScanner, error) {
	s, err := scanner.New(filePattern, namespaces, policyPaths, dataPaths)
	if err != nil {
		return ConfigScanner{}, xerrors.Errorf("unable to initialize config scanner: %w", err)
	}

	return ConfigScanner{
		parser:  &toml.Parser{},
		Scanner: s,
	}, nil
}

func (s ConfigScanner) Analyze(target analyzer.AnalysisTarget) (*analyzer.AnalysisResult, error) {
	var parsed interface{}
	if err := s.parser.Unmarshal(target.Content, &parsed); err != nil {
		return nil, xerrors.Errorf("unable to parse TOML (%s): %w", target.FilePath, err)
	}

	return &analyzer.AnalysisResult{
		Configs: []types.Config{
			{
				Type:     types.TOML,
				FilePath: target.FilePath,
				Content:  parsed,
			},
		},
	}, nil
}

func (s ConfigScanner) Required(filePath string, _ os.FileInfo) bool {
	if s.Match(filePath) {
		return true
	}

	ext := filepath.Ext(filePath)
	for _, required := range requiredExts {
		if ext == required {
			return true
		}
	}
	return false
}

func (s ConfigScanner) Type() analyzer.Type {
	return analyzer.TypeTOML
}

func (s ConfigScanner) Version() int {
	return version
}
