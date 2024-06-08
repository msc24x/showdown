package judge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/msc24x/showdown/internal/api/urls"
	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/utils"
)

func AuthenticateInstance(instance_url string, expect_type string) (*InstanceState, error) {
	ping_url := fmt.Sprintf("%s%s", instance_url, urls.Url("status"))
	client := &http.Client{}
	req, err := http.NewRequest("GET", ping_url, nil)
	decline_public := config.ACCESS_TOKEN != ""

	if err != nil {
		return nil, err
	}

	if decline_public {
		req.Header.Set("Access-Token", config.ACCESS_TOKEN)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	state := InstanceState{}
	err = json.NewDecoder(res.Body).Decode(&state)

	if err != nil {
		return nil, err
	}

	if state.InstanceType != expect_type {
		return nil, fmt.Errorf("cannot connect to %s instance when %s is expected", state.InstanceType, expect_type)
	}

	if !state.Private && decline_public {
		return nil, errors.New("cannot connect to an open instance when access token is set")
	}

	return &state, err
}

func ConnectManager(url string) {
	failIf := func(err error, context string) {
		if err != nil {
			utils.LogWorker("%s\nFailed to %s the manager on %s\n", err.Error(), context, url)
			os.Exit(1)
		}
	}

	stats, err := AuthenticateInstance(url, config.T_MANAGER)
	failIf(err, "ping")

	config.MANAGER_INSTANCE_ID = stats.InstanceId
	utils.LogWorker("Ping successful to manager instance %d running on %s", stats.InstanceId, url)

	req_body_struct := WorkerRegistration{
		Address: fmt.Sprintf("%s://%s:%d", config.PROTOCOL, config.HOST, config.PORT),
	}

	req_body_bytes, err := json.Marshal(req_body_struct)
	utils.PanicIf(err)

	register_url := fmt.Sprintf("%s%s", url, urls.Url("workers-register"))
	client := &http.Client{}
	req, err := http.NewRequest("POST", register_url, bytes.NewBuffer(req_body_bytes))
	failIf(err, "connect")

	req.Header.Set("Access-Token", config.ACCESS_TOKEN)
	res, err := client.Do(req)
	failIf(err, "connect")

	if res.StatusCode != 200 {
		res_bytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		utils.LogWorker("%s : Failed to connect to the manager on %s\n", string(res_bytes), url)
		os.Exit(1)
	}

	res_obj := WorkerRegistrationResponse{}
	err = json.NewDecoder(res.Body).Decode(&res_obj)
	failIf(err, "connect")

	config.INSTANCE_ID = res_obj.AssignedInstanceId

	utils.LogWorker("Connection with manager instance %d running on %s successful", stats.InstanceId, url)
}
