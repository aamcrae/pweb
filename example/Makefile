#
# Make web page from photos and album config
#
PORT := 8200
WEB := /tmp/pweb
ASSETS := ${WEB}/pweb
GOROOT := $(shell go env GOROOT)
TINYGOROOT := $(shell tinygo env TINYGOROOT)

all: dir assets photos

run: all
	go run server.go --port=${PORT} --base=${WEB}

photos: dir
	(cd photos; pweb --assets=../../assets --base=${WEB} web)

assets: dir ${ASSETS}/pweb.wasm
	cp ../assets/css/*  ${ASSETS}

${ASSETS}/pweb.wasm: *.go
	make tinygo # or go

tinygo:
	(cd ../wasm; tinygo build -target wasm -o ${ASSETS}/pweb.wasm)
	cp ${TINYGOROOT}/targets/wasm_exec.js ${ASSETS} # wasm_exec_tiny.js

go:
	(cd ../wasm; GOOS=js GOARCH=wasm go build -o ${ASSETS}/pweb.wasm)
	cp ${GOROOT}/misc/wasm/wasm_exec.js ${ASSETS}

dir:
	mkdir -p ${ASSETS}
	mkdir -p ${WEB}

clean:
	rm -rf ${WEB} ${ASSETS}
