package track

import (
	"fmt"
	"github.com/jjggzz/kj/errors"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter/http"
)

// 构建zipkin示踪器
// zipkinAddress: zipkin地址(host:port)
// serviceName:当前服务名
func BuildZipkinTracer(zipkinAddress string, serviceName string) (*zipkin.Tracer, error) {
	reporter := http.NewReporter(fmt.Sprintf("http://%s/api/v2/spans", zipkinAddress))
	zEP, err := zipkin.NewEndpoint(serviceName, "")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP))
}
