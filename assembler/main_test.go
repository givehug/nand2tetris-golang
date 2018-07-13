package assembler

import (
	"fmt"
	"io/ioutil"
	"nand2tetris-golang/common/utils"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// For each file in /asm dir compile hack file
// and comapre with corresponding one in /compare dir
func TestOuput(t *testing.T) {
	var wg sync.WaitGroup

	asmFiles, err := filepath.Glob("test/asm/**/*")
	if err != nil {
		t.Error(err)
	}
	compareFiles, err := filepath.Glob("test/compare/**/*")
	if err != nil {
		t.Error(err)
	}

	for i, f := range asmFiles {
		wg.Add(1)
		go func(i int, f string) {
			defer wg.Done()
			outFile := fmt.Sprintf("test/out_test_%d.hack", i)
			defer os.Remove(outFile)

			Compile(f, outFile)

			a, err := ioutil.ReadFile(compareFiles[i])
			if err != nil {
				t.Error(err)
			}
			b, err := ioutil.ReadFile(outFile)
			if err != nil {
				t.Error(err)
			}

			if utils.FilterNewLines(string(a)) != utils.FilterNewLines(string(b)) {
				t.Error("Error processing file", f)
			}
		}(i, f)
	}

	wg.Wait()
}
