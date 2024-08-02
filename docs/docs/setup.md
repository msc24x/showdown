# Setup
## Perquisites
!!! warning
	After making sure you have your machine set up according to their requirements,
	you might still face some issues. Go through this guide, when you
	face the issues.

To run Showdown, you only need [Docker](https://www.docker.com/) on your machine.

Showdown uses ISOLATE to securely execute remote code. Which itself requires some
configuration for it to work properly. Although most of the requirements should
already be satisfied by any normal Linux distribution. Most of its configuration (based
on the ISOLATE's official [manual page](https://www.ucw.cz/moe/isolate.1.html)) is 
already done in Showdown's compiler docker image. There are some steps that
you might need to do on your servers to ensure it's properly working.


#### Not using --privileged
The ISOLATE uses root privileges of the host machine as explained [here.](https://www.ucw.cz/moe/isolate.1.html#:~:text=Isolate%20is%20designed%20to%20run%20setuid%20to%20root.%20The%20sub%2Dprocess%20inside%20the%20sandbox%20then%20switches%20to%20a%20non%2Dprivileged%20user%20ID%20(different%20for%20each%20%2D%2Dbox%2Did).%20The%20range%20of%20UIDs%20available%20and%20several%20filesystem%20paths%20are%20set%20in%20a%20configuration%20file%2C%20by%20default%20located%20in%20/usr/local/etc/isolate) 

Due to this, the docker container for any Showdown instance must be started as `--privileged`

#### Error while "Initializing isolate box"
This should be due to your control group's version. By default your machine might
have been using cg2, make sure to switch it to cg1.

To ensure that it is the issue, run `ls /sys/fs/cgroup/` and if you don't see
the directories (memory, cpuset, cpuacct), then it is the issue. To
fix it you need to switch your system to use control group version 1.
See the instructions [here](https://stackoverflow.com/a/76194598/11367677) to add
the following kernel parameters and restart the machine.

```
cgroup_enable=memory 
systemd.unified_cgroup_hierarchy=0
```

## Quick Start

!!! note
	This section only explains the most basic setup of only one mode of the Showdown.
	See the [installation](./setup.md) for depth and understanding.

### Configuration

You need to create two files in the same directory.

- .config ([configuration file](./config.md))
- .env.creds ([secrets file](./credentials.md))

Below are the example files that represent a valid configuration, but are 
not limited to the variables mentioned. Full details on configuration can be found in other
sections.

#### .config file example (manager/standalone)

	# Default paths for showdown to work (change only if you know what you are doing)
	C=/usr/bin/gcc
	CPP=/usr/bin/g++
	PY=/opt/python/3.12.0/bin/python3
	GO=/usr/local/go/bin/go
	JS=/usr/bin/node
	TS=/usr/bin/ts-node

	# Specify the path of creds file
	CREDS_FILE=env/.env.creds

	# HOST=showdown-manager # uncomment for manager
	# HOST=showdown-standalone # uncomment for standalone
	PORT=7070

#### .config file example (worker)

	# Default paths for showdown to work (change only if you know what you are doing)
	C=/usr/bin/gcc
	CPP=/usr/bin/g++
	PY=/opt/python/3.12.0/bin/python3
	GO=/usr/local/go/bin/go
	JS=/usr/bin/node
	TS=/usr/bin/ts-node

	# Specify the path of creds file
	CREDS_FILE=env/.env.creds

	MANAGER_INSTANCE_ADDRESS=http://showdown-manager:7070
	HOST=showdown-worker
	PORT=7071


#### .env.creds file example

	ACCESS_TOKEN=your-access-token
	WEBHOOK_SECRET=your-webhook-secret

	RABBIT_MQ_PORT=5672
	RABBIT_MQ_HOST=showdown-mq
	RABBIT_MQ_USER=<rabbitmq-user>
	RABBIT_MQ_PASSWORD=<rabbitmq-password>

	RABBITMQ_DEFAULT_USER=user
	RABBITMQ_DEFAULT_PASS=password


### Installation ([Standalone](/glossary/#standalone))

- Pull the latest docker images, and create a network
```bash
docker pull msc24x/showdown:latest-standalone
docker pull msc24x/showdown:latest-queue
docker network create showdown
```

- Run RabbitMQ
```bash
docker run -d -p 5672:5672 -p 15672:15672 \
	--name showdown-mq \
	--env-file <CREDS_FILE> \
	--network=showdown \
	msc24x/showdown:latest-queue
```

- Run Standalone Showdown
```bash
docker run -d -p 7070:7070 \
	--name showdown-standalone \
	-v <path-to-config>:/showdown/env \
	-v <path-to-data-dir>:/var/lib/showdown \
	--privileged \
	--network showdown \
	msc24x/showdown:latest-standalone
```

### Installation ([Manager-Worker](/glossary/#manager))
- Pull the latest docker images, and create a network
```bash
docker pull msc24x/showdown:latest-manager
docker pull msc24x/showdown:latest-worker
docker pull msc24x/showdown:latest-queue
docker network create showdown
```

- Run RabbitMQ
```bash
docker run -d -p 5672:5672 -p 15672:15672 \
	--name showdown-mq \
	--env-file <CREDS_FILE> \
	--network=showdown \
	msc24x/showdown:latest-queue
```

- Run Manager Instance
```bash
docker run -d -p 7070:7070 \
	--name showdown-manager \
	-v <path-to-config>:/showdown/env \
	-v <path-to-data-dir>:/var/lib/showdown \
	--network showdown \
	msc24x/showdown:latest-manager
```

- Run Worker instance
```bash
docker run -d -p 7071:7071 \
	--name showdown-worker \
	-v <path-to-config>:/showdown/env \
	-v <path-to-data-dir>:/var/lib/showdown \
	--privileged \
	--network showdown \
	msc24x/showdown:latest-worker
```
