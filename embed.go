package highlight

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

//go:embed syntax_files
var syntaxFiles embed.FS

// AllDefs returns all the definitions in the embedded file system,
// filtering for named definitions if argument is provided
func AllDefs(fileType ...string) ([]*Def, error) {
	all, err := syntaxFiles.ReadDir("syntax_files")
	if err != nil {
		panic(err)
	}
	results := make([]*Def, 0, len(all))
	for _, file := range all {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".yaml" {
			continue
		}
		data, err := fs.ReadFile(syntaxFiles, filepath.Join("syntax_files", file.Name()))
		if err != nil {
			return nil, err
		}
		def, err := ParseDef(data)
		if err != nil {
			// Error handling
			return nil, fmt.Errorf("%v: %w", file.Name(), err)
		}
		if len(fileType) == 0 || sliceContains(fileType, def.FileType) {
			results = append(results, def)
		}
	}

	// Call ResolveIncludes
	ResolveIncludes(results)

	// Return all definitions
	return results, nil
}

func sliceContains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
