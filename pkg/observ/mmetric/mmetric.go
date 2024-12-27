package mmetric

import (
	"github.com/lithammer/shortuuid/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Global meter
var meter = otel.Meter("github.com/muchlist/moneymagnet")
var uniqueDeploymentCode = "money-magnet." + shortuuid.New()
var uniquePerNodeID = attribute.String("uid", uniqueDeploymentCode)
