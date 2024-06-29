package config

import (
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

// After config file path has been set, call this to import the config from file.
func ImportConfig() bool {
	paths, err := godotenv.Read(CONFIG_FILE)

	if err != nil {
		log.Printf("Could not read file %s, continuing with default configuration\n", CONFIG_FILE)
		return false
	}

	var (
		conf_key string
		sbuff    string
		ibuff    int
	)

	conf_key = "MAX_WORKER_RETRIES"
	ibuff, err = strconv.Atoi(paths[conf_key])
	if err == nil {
		MAX_WORKER_RETRIES = uint8(ibuff)
	}

	conf_key = "MAX_ACTIVE_PROCESSES"
	ibuff, err = strconv.Atoi(paths[conf_key])
	if err == nil {
		MAX_ACTIVE_PROCESSES = uint(ibuff)
	}

	conf_key = "ACTIVE_POLLING_RATE"
	ibuff, err = strconv.Atoi(paths[conf_key])
	if err == nil {
		ACTIVE_POLLING_RATE = ibuff
	}

	conf_key = "REVIVAL_POLLING_RATE"
	ibuff, err = strconv.Atoi(paths[conf_key])
	if err == nil {
		REVIVAL_POLLING_RATE = ibuff
	}

	conf_key = "PORT"
	ibuff, err = strconv.Atoi(paths[conf_key])
	if err == nil {
		PORT = ibuff
	}

	conf_key = "HOST"
	sbuff = paths[conf_key]
	if sbuff != "" {
		HOST = sbuff
	}

	conf_key = "PROTOCOL"
	sbuff = paths[conf_key]
	if sbuff != "" {
		PROTOCOL = sbuff
	}

	conf_key = "CREDS_FILE"
	sbuff = paths[conf_key]
	if sbuff != "" {
		CREDS_FILE = sbuff
	}

	conf_key = "INSTANCE_TYPE"
	sbuff = paths[conf_key]
	if sbuff != "" {
		INSTANCE_TYPE = sbuff
	}

	conf_key = "MANAGER_INSTANCE_ADDRESS"
	sbuff = paths[conf_key]
	if sbuff != "" {
		MANAGER_INSTANCE_ADDRESS = sbuff
	}

	return true
}
