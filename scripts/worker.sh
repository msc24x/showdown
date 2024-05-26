sudo /usr/local/go/bin/go build -buildvcs=false -o ./build/showdown ./cmd/showdown
sudo build/showdown -start -w -c http://localhost:7070 -p 7071 -config .config