package: build
	test -d dist || mkdir dist
	tar cvf dist/agents.tar --owner=0 --group=0 --transform "s|dist/|plugins/|" dist/Teamspeak3
	tar cvf dist/checks.tar --owner=0 --group=0 --transform "s|.py||" Teamspeak3.py
	tar czvf dist/teamspeak3-$$(git describe --tags).mkp --owner=0 --group=0 --transform "s|dist/||" \
	dist/info info.json dist/agents.tar dist/checks.tar

build:
	go mod download
	test -d dist || mkdir dist
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -ldflags="-extldflags=-static" \
	-o dist/Teamspeak3 Teamspeak3.go
	python -c 'import json,sys,pprint;print(pprint.pformat(json.load(sys.stdin)))' < info.json > dist/info

clean:
	rm -r dist