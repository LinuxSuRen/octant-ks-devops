build:
	go build -o bin/octant-ks-devops cmd/main.go

copy: build
	mkdir -p ~/.config/octant/plugins && cp bin/octant-ks-devops ~/.config/octant/plugins

run: copy
	octant --disable-open-browser
