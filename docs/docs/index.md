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


## Installation from source

### Standalone
Clone the Showdown repository on your machine
```
git clone https://github.com/msc24x/showdown 
cd showdown
```
Start the rabbitMQ instance 
```
docker compose up showdown-mq -d
```
Run the standalone instance of Showdown
```
docker compose up showdown-standalone -d
```