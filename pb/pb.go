//go:generate protoc --go_out=. cmd.proto
//go:generate protoc --go_out=. config.proto
//go:generate protoc --go_out=. device_report.proto
//go:generate protoc --go_out=. highway.proto
//go:generate protoc --go_out=. message.proto
//go:generate protoc --go_out=. message_body.proto
//go:generate protoc --go_out=. message_head.proto
//go:generate protoc --go_out=. message_hummer.proto
//go:generate protoc --go_out=. message_service.proto
//go:generate protoc --go_out=. oidb.proto
//go:generate protoc --go_out=. online_push.proto

package pb
