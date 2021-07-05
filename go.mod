module github.com/giantswarm/micrologger

go 1.14

require (
	github.com/giantswarm/microerror v0.3.0
	github.com/go-kit/kit v0.11.0
	github.com/go-stack/stack v1.8.0
	github.com/google/go-cmp v0.5.5
)

// We do not directly use the websocket package but within the dependency graph
// this package is necessary. We have to make sure it is at least at v1.4.2 due
// to some security fixes that CI would complain about otherwise.
replace github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
