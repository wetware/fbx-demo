//go:generate capnp compile -I$GOPATH/src/capnproto.org/go/capnp/std -ogo:llm llm.capnp
//go:generate capnp compile -I$GOPATH/src/capnproto.org/go/capnp/std -ogo:tiktok tiktok.capnp
package cap
