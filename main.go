package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"strings"
	"os/exec"
	"log"
	"os"
	"path/filepath"
	"github.com/satori/go.uuid"
)

var (
	codes = kingpin.Arg("code", "").Required().Strings()
)

func main() {
	kingpin.Parse()
	for _, code := range *codes {
		fmt.Println(code)
	}

	tf, err := TempFile("", "igo-", ".go")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("tempfile: %s", tf.Name())

	header := `
package main
import (
	"fmt"
)
func p(a ...interface{}) {
	fmt.Println(a...)
}
func main() {
	`
	body := strings.Join(*codes, " ")
	footer := `
}
	`
	log.Printf("file content: %s", header + body + footer)

	if _, err := tf.WriteString(header + body + footer); err != nil {
		log.Fatal(err)
	}

	out, err := exec.Command("env", "go", "run", tf.Name()).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("output: %s", out)
}

// TempFile copy from io.TempFile, and enabled suffix option.
// suffix is not able to use current version.
// see also https://github.com/golang/go/issues/4896
func TempFile(dir, prefix, suffix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+uuid.NewV4().String()+suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			continue
		}
		break
	}
	return
}