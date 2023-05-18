builder=go build -trimpath -gcflags "all=-N -l"
both: bin/botrunner bin/basic.so
bin/botrunner: */*/*.go 
	go build -a -trimpath -gcflags "all=-N -l" -o $@ ./cmd/botrunner
bin/basic.so: */*/*.go 
	go build -a -trimpath -gcflags "all=-N -l" --buildmode plugin -o $@ ./examples/basicmodule
run: both
	bin/botrunner -plugin bin/basic.so
clean:
	${RM} -r bin

