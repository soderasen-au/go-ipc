module github.com/soderasen-au/go-ipc

go 1.20

require github.com/soderasen-au/go-common v0.2.0

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/rs/zerolog v1.29.1 // indirect
	golang.org/x/sys v0.9.0 // indirect
)

replace (
	github.com/soderasen-au/go-common v0.2.0 => ../../soderasen-au/go-common
)