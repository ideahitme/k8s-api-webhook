package testutils

import (
	"encoding/csv"
	"io/ioutil"
	"os"
)

// GenerateTestData generates and writes a file, which is supposed to be removed later on
func GenerateTestData(lines [][]string) (tmpfile *os.File) {
	tmpfile, _ = ioutil.TempFile("", "authn-static-csv")
	csvWriter := csv.NewWriter(tmpfile)
	csvWriter.WriteAll(lines)
	return tmpfile
}
