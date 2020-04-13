module github.com/aibotsoft/pin

go 1.14

require (
	github.com/aibotsoft/gen v0.0.0-00010101000000-000000000000
	github.com/aibotsoft/gen/pinapi v0.0.0-20200413081135-7a5abad9c9a4
	github.com/aibotsoft/micro v0.0.0-20200411114812-ccef30d833e9
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.14.1
	google.golang.org/grpc v1.28.0
)

replace github.com/aibotsoft/micro => ../micro

replace github.com/aibotsoft/gen => ../gen
