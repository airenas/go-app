config: main.go
	go build -o $@

test-default: config
	./config
	
test-custom: config
	./config -c custom-config.yaml

test-env: config
	SAMPLE_NAME=OSGopher ./config

clean:
	rm -f config	