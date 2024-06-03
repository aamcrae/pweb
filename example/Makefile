#
# Make web page from photos and album config
#
PORT := 8080
WEB := /tmp/pweb
ASSETS := ${WEB}/pweb
GOROOT := $(shell go env GOROOT)

all: dir assets photos

run: all
	go run server.go --port=${PORT} --base=${WEB}

photos: dir
	(cd photos; pweb --assets=../../assets --base=${WEB} web)

assets: dir
	(cd ../wasm/album; GOOS=js GOARCH=wasm go build -o ${ASSETS}/album.wasm)
	(cd ../wasm/gallery; GOOS=js GOARCH=wasm go build -o ${ASSETS}/gallery.wasm)
	cp ../assets/css/*  ${ASSETS}
	cp ${GOROOT}/misc/wasm/wasm_exec.js ${ASSETS}

dir:
	mkdir -p ${ASSETS}
	mkdir -p ${WEB}

clean:
	rm -rf ${WEB} ${ASSETS}
