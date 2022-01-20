package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/tencentyun/tencentcloud-exporter/pkg/metric"
)

const (
	DcdbNamespace     = "QCE/TDMYSQL"
	DcdbInstanceidKey = "InstanceId"
)

func init() {
	registerHandler(DcdbNamespace, defaultHandlerEnabled, NewDcdbHandler)
}

type dcdbHandler struct {
	baseProductHandler
}

func (h *dcdbHandler) GetNamespace() string {
	return DcdbNamespace
}

func (h *dcdbHandler) IsMetricVaild(m *metric.TcmMetric) bool {
	// ignore node/shard metric, bug for cloud monitor if filter dim
	if len(m.Meta.SupportDimensions) != 1 {
		return false
	}
	return true
}

func NewDcdbHandler(c *TcProductCollector, logger log.Logger) (handler ProductHandler, err error) {
	handler = &dcdbHandler{
		baseProductHandler{
			monitorQueryKey: DcdbInstanceidKey,
			collector:       c,
			logger:          logger,
		},
	}
	return

}
