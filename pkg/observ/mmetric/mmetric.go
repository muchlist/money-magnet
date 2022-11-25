package mmetric

import (
	"github.com/lithammer/shortuuid/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
)

// Global meter
var meter = global.Meter("github.com/muchlist/moneymagnet")
var uniqueDeploymentCode = shortuuid.New()
var globalAtrs = []attribute.KeyValue{
	attribute.String("uid", uniqueDeploymentCode),
}
