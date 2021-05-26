package bin

import (
	"os"
	"testing"
)

func TestJsLoadFiles(t *testing.T) {
	f, err := os.Open("../../tests/binary/palm.js")
	defer f.Close()

	if err != nil {
		t.Error("error")
	}
	obj, err := ThreeJSObjFromJson(f)
	if err != nil {
		t.Error("error")
	}
	if obj != nil {
		t.Error("error")
	}
}

func TestBinLoadFiles(t *testing.T) {
	f, err := os.Open("../../tests/binary/palm.bin")
	defer f.Close()

	if err != nil {
		t.Error("error")
	}

	obj, err := Decode(f)
	if err != nil {
		t.Error("error")
	}
	if obj != nil {
		t.Error("error")
	}
}
