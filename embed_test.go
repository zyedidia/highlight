package highlight_test

import (
	"testing"

	// Namespace imports
	. "github.com/zyedidia/highlight"
)

func Test_Embed_001(t *testing.T) {
	defs, err := AllDefs()
	if err != nil {
		t.Fatal(err)
	}
	if len(defs) == 0 {
		t.Error("Expected non-zero defs array")
	}
	for n, def := range defs {
		t.Log(n, "=>", def)
	}
}

func Test_Embed_002(t *testing.T) {
	defs, err := AllDefs("yum")
	if err != nil {
		t.Fatal(err)
	}
	for n, def := range defs {
		if def.FileType != "yum" {
			t.Error("Expected FileType to be yum")
		}
		t.Log(n, "=>", def)
	}
}
