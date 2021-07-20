//go:generate protoc --go_out=. cmd.proto
//go:generate protoc --go_out=. config_push.proto
//go:generate protoc --go_out=. device_report.proto
//go:generate protoc --go_out=. highway.proto
//go:generate protoc --go_out=. message.proto
//go:generate protoc --go_out=. oidb.proto
//go:generate protoc --go_out=. online_push.proto

package pb
