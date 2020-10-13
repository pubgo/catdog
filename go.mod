module github.com/pubgo/catdog

go 1.14

require (
	github.com/apache/thrift v0.13.0
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gofiber/fiber v1.14.5
	github.com/gogo/protobuf v1.3.1
	github.com/gojektech/heimdall v5.0.2+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/go-version v1.2.1
	github.com/imdario/mergo v0.3.9
	github.com/jaegertracing/jaeger v1.19.2
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v3 v3.0.0-beta.0.20200828093547-6e30b5328036
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/common v0.10.0
	github.com/pubgo/dix v0.0.14
	github.com/pubgo/x v0.2.56 // indirect
	github.com/pubgo/xerror v0.2.11
	github.com/pubgo/xlog v0.0.7
	github.com/pubgo/xprocess v0.0.1
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0 // indirect
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/tebeka/strftime v0.1.5 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	go.mongodb.org/mongo-driver v1.3.2
	go.uber.org/zap v1.15.0
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	vitess.io/vitess v0.7.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
