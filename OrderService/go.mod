module github.com/jefrryss/go-grpc-microservices/OrderService

go 1.25.1

require (
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/jefrryss/go-grpc-microservices/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.81.1
)

require (
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260608224507-4308a22a1bab // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/jefrryss/go-grpc-microservices/shared => ../shared
