package config

import (
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

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
