//%:{{ `
module gitlab.com/proemergotech/dliver-project-skeleton

//%: ` | replace "dliver-project-skeleton" .ProjectName | trim }}

go 1.15

//%: {{ regexReplaceAll "(?m)^.*indirect.*$\n" `
require (
	github.com/HdrHistogram/hdrhistogram-go v0.9.0 // indirect
	github.com/codahale/hdrhistogram v0.9.0 // indirect
	github.com/go-playground/validator/v10 v10.0.1
	github.com/gomodule/redigo v1.8.2
	github.com/json-iterator/go v1.1.7
	github.com/labstack/echo/v4 v4.1.11
	github.com/olivere/elastic v6.2.25+incompatible
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.1.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/uber/jaeger-client-go v2.20.1+incompatible
	github.com/uber/jaeger-lib v2.3.0+incompatible // indirect
	gitlab.com/proemergotech/apimd-generator-go v1.0.0
	gitlab.com/proemergotech/bind v1.0.0
	gitlab.com/proemergotech/errors v1.0.0
	gitlab.com/proemergotech/geb-client-go/v2 v2.0.0
	gitlab.com/proemergotech/log-go/v3 v3.0.3
	gitlab.com/proemergotech/microtime-go/v2 v2.0.1
	gitlab.com/proemergotech/retry v1.0.1
	gitlab.com/proemergotech/trace-go/v2 v2.1.0
	gitlab.com/proemergotech/uuid-go v1.0.0
	gitlab.com/proemergotech/yafuds-client-go v1.2.1
	go.uber.org/zap v1.10.0
	gopkg.in/h2non/gentleman.v2 v2.0.4
)

//%: ` "" | trim }}
