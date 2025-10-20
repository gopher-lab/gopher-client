# Gopher Client

A Go client library for interacting with the Gopher AI data collection and search API. This client provides easy-to-use methods for performing various types of searches and data collection across multiple platforms including web, social media, and other data sources.

## Installation

```bash
go get github.com/masa-finance/gopher-client
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/masa-finance/gopher-client/client"
    "github.com/masa-finance/tee-worker/api/types"
)

func main() {
    // Create a new client
    c := client.NewClient("https://data.gopher-ai.com", "your-api-token")
    
    // Perform a web search
    result, err := c.PerformWebSearch("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Job ID: %s\n", result.JobID)
}
```

## Configuration

The client can be configured using environment variables or by passing parameters directly:

### Environment Variables

```bash
export BASE_URL="https://data.gopher-ai.com"
export API_TOKEN="your-api-token"
```

### Programmatic Configuration

```go
import "github.com/masa-finance/gopher-client/config"

// Load from environment
config, err := config.LoadConfig()
if err != nil {
    log.Fatal(err)
}

client := client.NewClient(config.BaseUrl, "your-api-token")
```

## Client Methods

### Core Client Operations

#### Job Management
- `GetJobStatus(jobID string) (*types.JobResult, error)` - Get the status of a job
- `GetResult(jobID string, receiver any) error` - Get the result of a completed job

### Web Search

#### Basic Web Search
```go
// Simple web page scraping
result, err := client.PerformWebSearch("https://example.com")
```

#### Advanced Web Search
```go
// Custom web job with specific arguments
args := page.NewArguments()
args.URL = "https://example.com"
// Set additional arguments as needed

result, err := client.PostWebJob(args)
```

### Social Media Search

#### Twitter
```go
// Search Twitter for posts
result, err := client.PerformTwitterSearch("golang programming")
```

#### Reddit
```go
// Search Reddit posts
result, err := client.PerformRedditSearchPosts("golang", 10)

// Search Reddit users
result, err := client.PerformRedditSearchUsers("username", 5)

// Search Reddit communities
result, err := client.PerformRedditSearchCommunities("programming", 20)

// Scrape specific Reddit URLs
result, err := client.PerformRedditScrapeURL("https://reddit.com/r/golang", 10)
```

#### LinkedIn
```go
import ptypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"

// Search LinkedIn profiles
result, err := client.PerformLinkedInSearch("software engineer", ptypes.ScraperMode)
```

#### TikTok
```go
// Search TikTok videos
result, err := client.PerformTikTokSearch("golang tutorial", 10)

// Get trending TikTok videos
result, err := client.PerformTikTokSearchByTrending("views", 20)

// Transcribe TikTok video
result, err := client.PerformTikTokTranscription("https://tiktok.com/@user/video/123")
```

### Advanced Search

#### Similarity Search
```go
import "github.com/masa-finance/tee-worker/api/types"

// Define sources to search
sources := []types.Source{
    types.WebSource,
    types.RedditSource,
    types.TwitterSource,
}

// Perform similarity search
var results []types.SearchResult
err := client.PerformSimilaritySearch(
    "artificial intelligence",  // query
    sources,                    // sources to search
    []string{"AI", "machine learning"}, // keywords
    "AND",                      // keyword operator
    10,                         // max results
    &results,                   // receiver for results
)
```

#### Hybrid Search
```go
// Perform hybrid search combining text and similarity queries
var results []types.SearchResult
err := client.PerformHybridSearch(
    "machine learning",         // text query
    sources,                    // sources to search
    "artificial intelligence",  // similarity text
    0.7,                        // query weight
    0.3,                        // text weight
    []string{"AI", "ML"},       // keywords
    "OR",                       // keyword operator
    15,                         // max results
    &results,                   // receiver for results
)
```

### Metrics

#### Get All Metrics
```go
// Get metrics for all collections
stats, err := client.GetAllMetrics(false) // false = don't refresh cache
```

#### Get Specific Collection Metrics
```go
// Get metrics for a specific source
stats, err := client.GetMetrics("web", true) // true = refresh cache
```

## Response Types

### Job Response
Most search methods return a `*types.ResultResponse` containing:
- `JobID` - Unique identifier for the job
- `Status` - Current job status
- Other job metadata

### Job Status
Use `GetJobStatus()` to check job progress:
- `PENDING` - Job is queued
- `RUNNING` - Job is currently executing
- `COMPLETED` - Job finished successfully
- `FAILED` - Job encountered an error

### Job Results
Use `GetResult()` to retrieve completed job data:
```go
var results []types.SearchResult
err := client.GetResult(jobID, &results)
```

## Error Handling

The client provides detailed error messages for common issues:

```go
result, err := client.PerformWebSearch("invalid-url")
if err != nil {
    // Handle different error types
    if strings.Contains(err.Error(), "Status code") {
        // HTTP error
    } else if strings.Contains(err.Error(), "failed to unmarshal") {
        // JSON parsing error
    } else {
        // Other errors
    }
}
```

## Authentication

The client supports Bearer token authentication:

```go
// With token
client := client.NewClient("https://data.gopher-ai.com", "your-api-token")

// Without token (if API allows anonymous access)
client := client.NewClient("https://data.gopher-ai.com", "")
```

## Examples

### Complete Web Search Example
```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/masa-finance/gopher-client/client"
    "github.com/masa-finance/tee-worker/api/types"
)

func main() {
    c := client.NewClient("https://data.gopher-ai.com", "your-token")
    
    // Start a web search job
    result, err := c.PerformWebSearch("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Job started: %s\n", result.JobID)
    
    // Poll for completion
    for {
        status, err := c.GetJobStatus(result.JobID)
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Printf("Status: %s\n", status.Status)
        
        if status.Status == "COMPLETED" {
            break
        } else if status.Status == "FAILED" {
            log.Fatal("Job failed")
        }
        
        time.Sleep(5 * time.Second)
    }
    
    // Get results
    var searchResults []types.SearchResult
    err = c.GetResult(result.JobID, &searchResults)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d results\n", len(searchResults))
}
```

### Async Search with Similarity
```go
package main

import (
    "fmt"
    "log"
    
    "github.com/masa-finance/gopher-client/client"
    "github.com/masa-finance/tee-worker/api/types"
)

func main() {
    c := client.NewClient("https://data.gopher-ai.com", "your-token")
    
    // Define search sources
    sources := []types.Source{
        types.WebSource,
        types.RedditSource,
        types.TwitterSource,
    }
    
    // Perform similarity search
    var results []types.SearchResult
    err := c.PerformSimilaritySearch(
        "golang best practices",
        sources,
        []string{"go", "golang", "programming"},
        "AND",
        20,
        &results,
    )
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d similar results\n", len(results))
    for i, result := range results {
        fmt.Printf("%d. %s\n", i+1, result.Title)
    }
}
```

## Dependencies

- Go 1.24.6+
- `github.com/masa-finance/tee-worker` - Core API types and job definitions

## License

This project is part of the Masa Finance ecosystem. Please refer to the project license for usage terms.

## Support

For issues and questions:
1. Check the API documentation
2. Review error messages for troubleshooting
3. Contact the Masa Finance team for support
