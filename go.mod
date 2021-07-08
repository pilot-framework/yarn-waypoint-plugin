module github.com/pilot-framework/yarn-waypoint-plugin

go 1.14

require (
	github.com/golang/protobuf v1.5.2
	github.com/hashicorp/waypoint-plugin-sdk v0.0.0-20201021094150-1b1044b1478e
	github.com/mitchellh/go-glint v0.0.0-20201015034436-f80573c636de
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.27.1
)

// replace github.com/hashicorp/waypoint-plugin-sdk => ../../waypoint-plugin-sdk
