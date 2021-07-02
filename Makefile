build:
	go build -o bin/octant-ks-devops cmd/main.go

copy: build
	cp bin/octant-ks-devops /Users/rick/.config/octant/plugins
