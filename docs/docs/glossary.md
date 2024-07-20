# Glossary

## Webhooks
A webhook is a URL a client may specify along with the process request to Showdown.
Using this parameter, Showdown will execute the the process asynchronously and will
hit the webhook with the process results as a response.

Your application must set `WEBHOOK_SECRET` in [.env.creds](/creds-file) secrets file.
That will result into Showdown workers sending `Webhook-Secret` HTTP header in the webhook
requests, so that your own application can identified that the request is valid and is made
by an actual Showdown server. 

## Message queue
A message queue, which is RabbitMQ in our case, is used by the application to
queue process requests so that they can be executed later when the load is high
or when using Manager-Worker deployment.

## Showdown servers
A Showdown server is a http server, that allows the communication with the 
Showdown application instance.

### Worker
The only purpose a `worker` instance have, is to connect to the message queue,
and dequeue any available processes to execute it. A worker instance does not
communicate with the client directly, and will require an already existing
Showdown [manager](/glossary/#manager) instance to connect to.

There can be N number of worker instances connected to one manager instance,
hence providing scalability according to the requirements. The more load your
application gets the more number of workers can be deployed to handle the messages,
or just keep running it and let the workers work themselves to death, along 
with your application.

After processing the messages the results should be returned to the corresponding
clients. [Manager](/glossary/#manager) ensures that the client had sent a webhook
URL along with the request, which is used by the workers to send the response to.

### Manager
It listens on a port for process requests and just adds
the process into the message queue. It is responsible for authorization and
verifying correctness & format of the requests, after which a unique PID is assigned to the message
which is immediately returned in the response to the client and the message is
finally queued.

Manager itself does not do any processing of the requests, the message queue is consumed by
its N number of [worker](/glossary/#worker) instances. Every manager must have at least
one worker instance.

This Manager-Worker architecture is one way a Showdown application can be deployed.
Note that this architecture does not support immediate execution. A webhook is always required
and clients must implement the webhook requests to get the results.


### Standalone
A standalone instance of Showdown, as the name suggests, is standalone, and does
the work of both [worker](/glossary/#worker) and [manager](/glossary/#manager). A standalone
instance is best for small use cases. It is just one http server that can be deployed on one
machine. Clients can communicate to standalone instances and it will execute the requests
just in time and return the results, as it allows the client to not use `webhooks` and get
the results directly in the API response.

But note that as it does allow non webhook requests, the application will still act
upon the [MAX_ACTIVE_PROCESSES](/config-file) set in configuration file. Exceeding that limit,
a standalone instance will start to reject requests with error `max active processes limit reached`.
In that case the client must use webhook. The limit can be modified  in the [.config](/config-file)
according to your application needs.


