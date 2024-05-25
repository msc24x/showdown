sudo /usr/local/go/bin/go build -buildvcs=false -o ./build/showdown ./cmd/showdown
sudo build/showdown -start -m -p 8080 -config .wsl.config -creds .env.creds
