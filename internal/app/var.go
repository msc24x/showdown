// Maintains application state.
package app

import (
	"os"

	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/utils"
)

// Backups the instance state in a dump file in /var/lib/showdown.
func DumpInstanceState(c []byte) {
	os.Mkdir("/var/lib/showdown/", 0644)

	f, err := os.OpenFile(config.DUMP_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	utils.PanicIf(err)

	_, err = f.Write(c)
	utils.PanicIf(err)
	f.Close()
}

// Restores the instance state from the dump file.
func ReadInstanceState() ([]byte, error) {
	os.Mkdir("/var/lib/showdown/", 0644)

	return os.ReadFile(config.DUMP_FILE)
}
