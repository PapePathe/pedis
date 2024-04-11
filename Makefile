start-primary:
	rm -rf /tmp/pedis*
	go run main.go --cluster  http://127.0.0.1:12379

start-secondary:
	go run main.go --id 2 --join --port 9122 --pedis 127.0.0.1:6389 --cluster http://127.0.0.1:12379,http://127.0.0.1:12380

start-tertiary:
	go run main.go --id 3 --port 9123 --pedis 127.0.0.1:6390 --cluster http://127.0.0.1:12379,http://127.0.0.1:12380,http://127.0.0.1:12381

test:
	go test -v ./... -race
