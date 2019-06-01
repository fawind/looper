package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

// GetDumpFileName returns the default filename for a test dump file
func GetDumpFileName(service string, testCmd string) string {
	return fmt.Sprint(service, "-", getFingerprint(testCmd), ".mitmdump")
}

// DumpExists checks if the dump file exists
func DumpExists(dumpDir string, outFile string) bool {
	_, err := os.Stat(path.Join(dumpDir, outFile))
	return !os.IsNotExist(err)
}

func getFingerprint(testCmd string) string {
	hash := sha1.New()
	hash.Write([]byte(testCmd))
	return hex.EncodeToString(hash.Sum(nil))
}
