# Configuration file (.config)

| Configuration | Description                        | Defaults     | Remarks  |
| ------------- | ---------------------------------- | ------------ | -------- |
| ACTIVE_POLLING_RATE | After how many seconds should a manager ping workers to synchronize states and check if they are live or not | 10 | For Manager instances only |
| REVIVAL_POLLING_RATE | After how many seconds should a manager retry to ping dropped workers to check if they are live again and can be revived | 1800 | For Manager instances only |
| PORT | Specifies the port on which the Showdown instance should run | 7070 | - |
| HOST | Specifies the hostname/ip address on which the Showdown instance should run | 0.0.0.0, string | - |
| CREDS_FILE | Specifies the path of the directory where the secrets file .env.creds resides | env/.env.creds | - |
| INSTANCE_TYPE | Specifies the type of the Showdown Instance that should be run | standalone | - |
| MANAGER_INSTANCE_ADDRESS | Specifies the full address of the manager instance to which current worker instance must connect | nil | For Worker instances only, required |
| C | gcc compiler path | /usr/bin/gcc | - |
| CPP | g++ compiler path | /usr/bin/g++ | - |
| PY | python interpreter path | /opt/python/3.12.0/bin/python3 | - |
| GO | go compiler path | /usr/local/go/bin/go | - |
| JS | node path | /usr/bin/node | - |
| TS | ts-node path | /usr/bin/ts-node | - |

