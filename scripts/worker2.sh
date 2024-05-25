sudo /usr/local/go/bin/go build -buildvcs=false -o ./build/showdown ./cmd/showdown
sudo build/showdown -start -w -c http://localhost:8080 -p 8082 -config .wsl.config -creds .env.creds