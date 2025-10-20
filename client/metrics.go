package client

import (
	"fmt"

	"github.com/gopher-lab/gopher-client/log"
	"github.com/masa-finance/tee-worker/api/types"
)

func (c *Client) GetAllMetrics(refresh bool) ([]types.CollectionStats, error) {
	url := fmt.Sprintf("%s/v1/metrics?refresh=%t", c.BaseURL, refresh)

	var stats []types.CollectionStats
	err := c.doMetricsRequest(url, &stats)
	if err != nil {
		log.Error("Error while getting all metrics", "refresh", refresh, "error", err.Error())
		return nil, err
	}
	return stats, nil
}

func (c *Client) GetMetrics(source string, refresh bool) (*types.CollectionStats, error) {
	url := fmt.Sprintf("%s/v1/metrics/%s?refresh=%t", c.BaseURL, source, refresh)

	var stats types.CollectionStats
	err := c.doMetricsRequest(url, &stats)
	if err != nil {
		log.Error("Error while getting metrics", "source", source, "refresh", refresh, "error", err.Error())
		return nil, err
	}
	return &stats, nil
}
