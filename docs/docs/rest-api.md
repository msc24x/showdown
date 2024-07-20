# REST API

All communication to any of the [Showdown servers](/glossary/#showdown-servers)
is done via REST APIs. Currently, its interface is the simple as possible with limited
endpoints.

## Authorization
Showdown requests are authenticated with the HTTP header `Access-token`, which can be
specified by you in the [.env.creds](/creds-file) secrets file. By default, this is empty,
and hence unauthenticated requests will be able to communicate to your instance. Hence,
we recommend setting `ACCESS_TOKEN` in your secrets.

## Endpoints

### **GET** `/status`
Get the status of an Instance.
```bash title="curl"
curl http://localhost:8000/status \
	-H "Access-Token: abcdefghijklmnopqrstuvqxyz" \
```
```json title="example response (standalone instance)"
{
  "InstanceId": 1,
  "InstanceType": "standalone",
  "ManagerId": 0,
  "Private": true,
  "WorkerState": {
    "StartedSince": "2024-07-15T01:27:19.804548221Z",
    "TotalProcessed": 1,
    "ActiveProcesses": 0,
    "Processes": {}
  },
  "Workers": null
}
```

#### Understanding response using internal types
##### Status 200
```golang title="Top level schema"
// A general state struct for any showdown instance.
type InstanceState struct {
	InstanceId   uint
	InstanceType string
	// Provides manager id if instance type is worker.
	ManagerId uint
	// Specifies if instance uses Access-Token.
	Private bool

	// Not nil, if instance type is standalone/worker.
	WorkerState *WorkerState
	// Specifies connected workers if instance type is manager.
	Workers []*ShowdownWorker
}

```
```golang title="InstanceState.WorkerState schema"
// An in memory runtime information and statistics of the Showdown.
type WorkerState struct {
	StartedSince time.Time
	// Total number of requests denied since start.
	DeniedProcessed uint
	// Total number of requests processed since start.
	TotalProcessed uint
	// Number of active requests being processed.
	ActiveProcesses uint
	// Map of processes currently being processed.
	Processes map[string]bool
}

```
```golang title="InstanceState.Workers (an array of type ShowdownWorker)"
type ShowdownWorker struct {
	InstanceId uint
	Address    string
	Status     WorkerStatus
	// Number of retries a worker has pending before a manager marks it as dropped.
	Retries          uint8
	LastFetchedState *InstanceState

	CreatedSince time.Time
	// Time since a worker has stalled.
	InactiveSince time.Time
}
```
```golang title="Possible values of ShowdownWorker.Status"
type WorkerStatus int8

const (
	// Worker is active, valid and open to requests.
	SW_ACTIVE WorkerStatus = 1 << 0
	// Worker was active but is no longer responding after multiple retries.
	// Manager stops pinging such instances.
	SW_DROPPED WorkerStatus = 1 << 1
	// Worker was active but has stopped responding to unknown reasons.
	SW_STALLED WorkerStatus = 1 << 2
)
```

##### Status 401
When the HTTP header `Access-Token` is not identified by the application.

### **POST** `/judge`
Key endpoint to request Showdown for an execution of a process
```bash title="curl (example Hello World with C++)"
curl  http://localhost:8000/judge \
	-X POST \
	-H 'Access-Token: abcdefghijklmnopqrstuvqxyz' \
	-H 'Content-Type: application/json; charset=utf-8' \
		--data-binary @- << EOF
		{
			"judge_params": {
				"donotjudge": false,
				"webhook": ""
			},
			"exe": {
				"language": "cpp",
				"code": "#include <iostream>\n\nint main() {\n  std::cout << \"Hello, Wold!\";\n  return 0;\n}",
				"input": "10",
				"output": "1\n2\n3\n4\n5\n6\n7\n8\n\n9 \n 10\n \n "
			}
		}
		EOF
```

```json title="example response (standalone instance)"
{
  "pid": "d2aeaf2a-f2d1-4c17-998f-3a857b1810ee",
  "webhook": "",
  "success": false,
  "judged": true,
  "error": "",
  "output": "Hello, Wold!",
  "meta": "time:0.001\ntime-wall:0.001\nmax-rss:3576\ncsw-voluntary:3\ncsw-forced:1\ncg-mem:35648\nexitcode:0\n",
  "expected": "1\n2\n3\n4\n5\n6\n7\n8\n\n9 \n 10\n \n ",
  "server_fault": false
}
```

#### Understanding requests using internal types

```golang title="Top level schema"
// Struct to define http request to showdown
type JudgeRequest struct {
	JudgeParams judge.Params            `json:"judge_params"`
	Exe         engine.ExecutionRequest `json:"exe"`
}
```

```golang title="JudgeRequest.JudgeParams"
type Params struct {
	// Set the webhook, and showdown will send the response to that webhook
	// instead of an immediate response.
	Webhook string `json:"webhook"`

	// Set this to true and Showdown will only execute the code, not judge.
	DoNotJudge bool `json:"donotjudge"`
}
```

```golang title="JudgeRequest.Exe"
// Constructs of fields required for a execution request to be valid.
type ExecutionRequest struct {
	// Code submitted by the user.
	Code string `json:"code"`
	// Programming language of the code.
	Language string `json:"language"`
	// What should be streamed into the program.
	Input string `json:"input"`
	// What is the expected output for the program.
	Output string `json:"output"`
}
```

#### Understanding response using internal types

##### Status 200
```golang title="Top level schema"
// Struct to define the end results the users will receive.
type ExecutionResponse struct {
	PID         string `json:"pid"`
	Webhook     string `json:"webhook"`
	Success     bool   `json:"success"`
	Judged      bool   `json:"judged"`
	Error       string `json:"error"`
	Output      string `json:"output"`
	Meta        string `json:"meta"`
	Expected    string `json:"expected"`
	ServerFault bool   `json:"server_fault"`
}
```

##### Status 401
When the HTTP header `Access-Token` is not identified by the application.

##### Status 400
Following are the possible self-explanatory error strings.

- `Not allowed on worker instance`
- `max active processes limit reached`

##### Status 500
When something is wrong with the application servers.

