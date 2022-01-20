package instance

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentyun/tencentcloud-exporter/pkg/client"
	"github.com/tencentyun/tencentcloud-exporter/pkg/config"
)

func init() {
	registerRepository("QCE/LB", NewEIPTcInstanceRepository)
}

type EIPTcInstanceRepository struct {
	client *sdk.Client
	logger log.Logger
}

func (repo *EIPTcInstanceRepository) GetInstanceKey() string {
	return "eip"
}

func (repo *EIPTcInstanceRepository) Get(id string) (instance TcInstance, err error) {
	req := sdk.NewDescribeAddressesRequest()
	req.AddressIds = []*string{&id}
	resp, err := repo.client.DescribeAddresses(req)
	if err != nil {
		return
	}
	if len(resp.Response.AddressSet) != 1 {
		return nil, fmt.Errorf("Response instanceDetails size != 1, id=%s ", id)
	}
	meta := resp.Response.AddressSet[0]
	instance, err = NewEIPTcInstance(*meta.AddressIp, meta)
	if err != nil {
		return
	}
	return
}

func (repo *EIPTcInstanceRepository) ListByIds(id []string) (instances []TcInstance, err error) {
	return
}

func (repo *EIPTcInstanceRepository) ListByFilters(filters map[string]string) (instances []TcInstance, err error) {
	req := sdk.NewDescribeAddressesRequest()
	var offset int64 = 0
	var limit int64 = 100
	var total int64 = -1

	req.Offset = &offset
	req.Limit = &limit

getMoreInstances:
	resp, err := repo.client.DescribeAddresses(req)
	if err != nil {
		return
	}
	if total == -1 {
		total = *resp.Response.TotalCount
	}
	for _, meta := range resp.Response.AddressSet {
		ins, e := NewEIPTcInstance(*meta.AddressIp, meta)
		if e != nil {
			level.Error(repo.logger).Log("msg", "Create cdb instance fail", "id", *meta.InstanceId)
			continue
		}
		instances = append(instances, ins)
	}
	offset += limit
	if offset < total {
		req.Offset = &offset
		goto getMoreInstances
	}

	return
}

func NewEIPTcInstanceRepository(c *config.TencentConfig, logger log.Logger) (repo TcInstanceRepository, err error) {
	cli, err := client.NewVpvClient(c)
	if err != nil {
		return
	}
	repo = &EIPTcInstanceRepository{
		client: cli,
		logger: logger,
	}
	return
}
