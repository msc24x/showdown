# Overview

## What is Showdown?
Showdown is a portable application, written in [GoLang](https://go.dev/), to execute and judge code for multiple languages.
It comes with an already built docker image with compilers and runners for selected
programming languages and configuration to use required binaries to execute the code
based on the language provided.
	
It acts as a server that listens to requests for code
execution and processes the requests either in asynchronous mode or immediate mode,
depending upon the configuration. Its customizable in terms of execution and the configuration
gives you enough control over your application specific setup. It uses [RabbitMQ](https://www.rabbitmq.com/) for
queue that allows the Showdown to consist of multiple instances on different servers
and act together as a distributed network workers.

## How it works?
There are 4 key parts of its design.

1. Application (Standalone/Manager/Worker)
2. Isolate ([Sandbox for securely executing untrusted programs](https://github.com/ioi/isolate))
3. Compilers (Docker image with collection of compilers/runners)
4. Message queue (RabbitMQ)

The application, which is Showdown, is a mediator between the client/user and the compilers. It
is responsible for listening to user requests for any code execution. It parses the request, extracts
the language, on the basis of which it selects the appropriate compiler or the runner for the language
and at last executes it.

For execution, it utilizes Isolate, which is a sandbox that allows Showdown to isolate the code
submitted by the user and execute it safely. Using it we can set up CPU time or memory limits as well.

Now the application can be started on your local machine. You can specify the paths of the required
compilers in its config file. Otherwise, an ideal condition (recommended) to start it would be through
the docker images provided. The images are based upon our Compilers image which is already a package
of supported languages and Isolate installation which is required for Showdown to work.

At last, the application can work in 3 modes. _Worker_, whose only job is to consume the available
execution requests lying in message queue. _Manager_, which is only responsible for listening to user
requests and just adding the request into the message queue. _Standalone_, which is quite self explanatory;
starting a _standalone_ instance would be sufficient itself for small or most use cases, as you won't
be required to run _manager_ or _worker_ instance separately.


## Perquisites
!!! warning
	After making sure you have your machine setup according to their requirements,
	it is possible you might still face some issues. Go through this guide, only when you
	face the issues.

Showdown uses ISOLATE to securely execute remote code. Which itself requires some
configuration for it to work properly. Although most of the requirements should
already be satisfied by any normal linux distribution. There are some steps that
you will need to do on your servers to ensure its proper working.
The ISOLATE's official manual page is [here.](https://www.ucw.cz/moe/isolate.1.html)

#### Not using --privileged
The ISOLATE uses root privileges of the host machine as explained [here.](https://www.ucw.cz/moe/isolate.1.html#:~:text=Isolate%20is%20designed%20to%20run%20setuid%20to%20root.%20The%20sub%2Dprocess%20inside%20the%20sandbox%20then%20switches%20to%20a%20non%2Dprivileged%20user%20ID%20(different%20for%20each%20%2D%2Dbox%2Did).%20The%20range%20of%20UIDs%20available%20and%20several%20filesystem%20paths%20are%20set%20in%20a%20configuration%20file%2C%20by%20default%20located%20in%20/usr/local/etc/isolate) 

Due to which the docker container for any Showdown instance must be started as `--privileged`

#### Error while "Initializing isolate box"
This should be due to your control groups version. By default your machine might
have been using cg2, make sure to switch it to cg1.

To ensure that it is the issue, run `ls /sys/fs/cgroup/` and if you don't see
the directories (memory, cpuset, cpuacct), then it is definitely the issue. To
fix it you need to switch your system to use control group version 1.
See instructions [here](https://stackoverflow.com/a/76194598/11367677)
to add following kernel parameters and restart the machine.

```
cgroup_enable=memory 
systemd.unified_cgroup_hierarchy=0
```


## Quick Start (Standalone)

!!! note
	This section only explains the most basic setup of only one mode of the Showdown.
	See [installation](./setup.md) for depths and understanding.

### Configuration

You need to create two files in same directory.

- .config (configuration file)
- .env.creds (secrets file)

Below are the example files that represents a valid configuration, but not limited
to the variables mentioned. Full details on configuration can be found in other
sections.

#### .config file example

	# Default paths for showdown to work (change only if you know what you are doing)
	C=/usr/bin/gcc
	CPP=/usr/bin/g++
	PY=/opt/python/3.12.0/bin/python3
	GO=/usr/local/go/bin/go
	JS=/usr/bin/node
	TS=/usr/bin/ts-node

	# Specify instance type
	INSTANCE_TYPE=standalone

	# Specify the path of creds file
	CREDS_FILE=env/.env.creds


#### .env.creds file example

	ACCESS_TOKEN=your-access-token
	WEBHOOK_SECRET=your-webhook-secret

	RABBIT_MQ_PORT=5672
	RABBIT_MQ_HOST=<rabbitmq-host>
	RABBIT_MQ_USER=<rabbitmq-user>
	RABBIT_MQ_PASSWORD=<rabbitmq-password>

	RABBITMQ_DEFAULT_USER=user
	RABBITMQ_DEFAULT_PASS=password


### Installation (docker cli)

- Pull the latest docker images, and create network

```
docker pull msc24x/showdown:latest-standalone
docker pull msc24x/showdown:latest-queue
docker network create showdown
```

- Run RabbitMQ

```
docker run -d -p 5672:5672 -p 15672:15672 \
	--name showdown-mq \
	--env-file <CREDS_FILE> \
	--network=showdown \
	msc24x/showdown:latest-queue
```

- Run Standalone Showdown

```
docker run -d -p 7070:7070 \
	--name showdown-standalone \
	-v <path-to-config>:/showdown/env \
	-v <path-to-data-dir>:/var/lib/showdown \
	--privileged \
	--network showdown \
	msc24x/showdown:latest-standalone

```

### Installation (from source)

Clone the Showdown repository on your machine
```
git clone https://github.com/msc24x/showdown 
cd showdown
```
Start the rabbitMQ instance 
```
docker compose \
	-f docker-compose.yml \
	--env-file <.env.creds> \
	up showdown-mq \
	-d --build
```
Run the standalone instance of Showdown
```
docker compose \
	-f docker-compose.yml \
	--env-file <.config> \
	up showdown-standalone \
	-d --build
```