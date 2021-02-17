package server

//go:generate protoc -I=. --go_opt=paths=source_relative --go_out=plugins=grpc:. proto/tronimoes.proto tiles/proto/tiles.proto
