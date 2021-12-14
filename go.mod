module github.com/giantswarm/micrologger

go 1.17

require (
	github.com/giantswarm/microerror v0.4.0
	github.com/go-kit/log v0.2.0
	github.com/go-stack/stack v1.8.0
	github.com/google/go-cmp v0.5.6
)

require github.com/go-logfmt/logfmt v0.5.1 // indirect

// We do not directly use the websocket package but within the dependency graph
// this package is necessary. We have to make sure it is at least at v1.4.2 due
// to some security fixes that CI would complain about otherwise.
replace github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
