simple: main.go
	go build -o $@

test-info: simple
	./simple
	
test-debug-env: simple
	LOGGER_LEVEL=Debug ./simple

clean:
	rm -f simple	
