package pbsessions

//go:generate protoc -I $GOPATH/src/github.com/fernandezvara/backd/internal/pbsessions sessions.proto --go_out=plugins=grpc:.
