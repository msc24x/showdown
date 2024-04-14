sudo /usr/local/go/bin/go build -buildvcs=false -o ./build/showdown ./cmd/showdown
sudo build/showdown -start -w -c http://localhost:8080 -p 8081 -paths .wsl.paths -creds .env.creds