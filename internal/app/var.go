package app

import (
	"msc24x/showdown/config"
	"msc24x/showdown/internal/utils"
	"os"
)

func DumpInstanceState(c []byte) {
	os.Mkdir("/var/lib/showdown/", 0644)

	f, err := os.OpenFile(config.DUMP_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	utils.PanicIf(err)

	_, err = f.Write(c)
	utils.PanicIf(err)
	f.Close()
}

func ReadInstanceState() ([]byte, error) {
	os.Mkdir("/var/lib/showdown/", 0644)

	return os.ReadFile(config.DUMP_FILE)
}
