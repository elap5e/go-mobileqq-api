//go:generate protoc --go_out=. config.proto
//go:generate protoc --go_out=. device_report.proto
//go:generate protoc --go_out=. message.proto
//go:generate protoc --go_out=. online.proto

package pb
