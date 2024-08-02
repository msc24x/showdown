# Credentials file (.env.creds)

| Configuration | Description                        | Defaults     | Remarks  |
| ------------- | ---------------------------------- | ------------ | -------- |
| ACCESS_TOKEN | Set it to make your instance private, i.e. HTTP HEADER 'Access-Token' would be required to make requests | nil | By default anyone can make requests |
| WEBHOOK_SECRET | Shared token between Showdown and its clients, so that webhook requests can be identified if its really from Showdown's servers | nil | Should be set to ensure webhook requests are trusted |
| RABBIT_MQ_HOST | RabbitMQ hostname | localhost | - |
| RABBIT_MQ_PORT | RabbitMQ port | 5672 | - |
| RABBIT_MQ_USER | RabbitMQ user name | guest | - |
| RABBIT_MQ_PASSWORD | RabbitMQ user password | guest | - |

