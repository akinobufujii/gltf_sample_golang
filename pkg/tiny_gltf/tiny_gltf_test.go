package tiny_gltf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewInstanceFromFile(t *testing.T) {

	filenamelist := []string{
		"testdata/minimal.gltf",
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for _, filename := range filenamelist {
		loadfilename := filepath.Join(wd, filename)

		instance, err := NewInstanceFromFile(loadfilename)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%s\n", filename)
		t.Logf("%+v\n", instance)
		t.Logf("\n")
	}
}
