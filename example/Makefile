#
# Make web page from photos and album config
#
PORT := 8200
WEB := /tmp/pweb
ASSETS := ${WEB}/pweb
GOROOT := $(shell go env GOROOT)

all: dir assets photos

run: all
	go run server.go --port=${PORT} --base=${WEB}

photos: dir
	(cd photos; pweb --assets=../../assets --base=${WEB} web)

assets: dir go # tinygo
	cp ../assets/css/*  ${ASSETS}
	cp ${GOROOT}/misc/wasm/wasm_exec.js ${ASSETS}

tinygo:
	(cd ../wasm; tinygo build -o ${ASSETS}/pweb.wasm -target wasm)

go:
	(cd ../wasm; GOOS=js GOARCH=wasm go build -o ${ASSETS}/pweb.wasm)

dir:
	mkdir -p ${ASSETS}
	mkdir -p ${WEB}

clean:
	rm -rf ${WEB} ${ASSETS}
