module github.com/pubgo/catdog

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0

require (
	github.com/apache/thrift v0.13.0
	github.com/asim/nitro-plugins/client/grpc/v3 v3.3.1-0.20201031120104-4c96a26220fa
	github.com/asim/nitro-plugins/config/source/etcd/v3 v3.4.0
	github.com/asim/nitro-plugins/registry/mdns v0.0.0-20201101073154-04271fcbbf50
	github.com/asim/nitro-plugins/server/grpc/v3 v3.3.1-0.20201031120104-4c96a26220fa
	github.com/asim/nitro/v3 v3.3.0
	github.com/dave/jennifer v1.4.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gofiber/fiber v1.14.5
	github.com/gogo/protobuf v1.3.1
	github.com/gojektech/heimdall v5.0.2+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/gopherjs/gopherjs v0.0.0-20191106031601-ce3c9ade29de // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.1
	github.com/hashicorp/go-version v1.2.1
	github.com/imdario/mergo v0.3.9
	github.com/jaegertracing/jaeger v1.19.2
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pubgo/dix v0.1.0
	github.com/pubgo/xerror v0.2.12
	github.com/pubgo/xlog v0.0.9
	github.com/pubgo/xprocess v0.0.3
	github.com/pubgo/xprotogen v0.0.2
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	go.mongodb.org/mongo-driver v1.3.2
	go.uber.org/zap v1.16.0
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
	vitess.io/vitess v0.7.0
)
