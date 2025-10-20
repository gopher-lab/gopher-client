package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/masa-finance/tee-worker/api/types"
)

// GetAllMetrics retrieves metrics for all available collections.
//
// Returns:
//
//	A Map containing statistics for each source, or an error if the operation fails.
func (c *Client) GetAllMetrics(refresh bool) ([]types.CollectionStats, error) {
	url := fmt.Sprintf("%s/v1/metrics?refresh=%t", c.BaseURL, refresh)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status: %d; body: %s", resp.StatusCode, string(body))
	}

	var stats []types.CollectionStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return stats, nil
}

// GetMetrics retrieves metrics for a specific collection.
//
// Args:
//
//	collectionName: The name of the collection to retrieve metrics for.
//
// Returns:
//
//	A pointer to MetricsStats containing statistics for the provided source, or an error if the operation fails.
func (c *Client) GetMetrics(source string, refresh bool) (*types.CollectionStats, error) {
	url := fmt.Sprintf("%s/v1/metrics/%s?refresh=%t", c.BaseURL, source, refresh)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status: %d; body: %s", resp.StatusCode, string(body))
	}

	var stats types.CollectionStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &stats, nil
}
