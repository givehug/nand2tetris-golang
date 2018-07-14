package utils

import (
	"testing"
)

func TestUtils(t *testing.T) {
	t.Run("ToXML", func(t *testing.T) {
		if ToXML("tag", "val", true) != "<tag> val </tag>" {
			t.Error("bad inline xml")
		}
		if ToXML("tag", "val", false) != "<tag>\n  val\n</tag>" {
			t.Error("bad multiline xml")
		}
	})
}
