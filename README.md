# Gopher Client

A Go client library for interacting with the Gopher AI data collection and search API. This client provides easy-to-use methods for performing various types of searches and data collection across multiple platforms including web, social media, and other data sources.

## Installation

```bash
go get github.com/gopher-lab/gopher-client
```

## Quick Start

### Using Environment Variables (Recommended)
```go
package main

import (
    "fmt"
    "log"
    
    "github.com/gopher-lab/gopher-client/client"
    "github.com/masa-finance/tee-worker/api/types"
)

func main() {
    // Create client from environment variables
    c, err := client.NewClientFromConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // Perform a web search
    result, err := c.PerformWebSearch("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Job ID: %s\n", result.JobID)
}
```

### Using Explicit Configuration
```go
package main

import (
    "fmt"
    "log"
    
    "github.com/gopher-lab/gopher-client/client"
    "github.com/masa-finance/tee-worker/api/types"
)

func main() {
    // Create a new client with explicit configuration
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
export GOPHER_CLIENT_URL="https://data.gopher-ai.com"
export GOPHER_CLIENT_TOKEN="your-api-token"
```

### Programmatic Configuration

```go
import "github.com/gopher-lab/gopher-client/config"

// Method 1: Load config manually and create client
config, err := config.LoadConfig()
if err != nil {
    log.Fatal(err)
}

client := client.NewClient(config.BaseUrl, config.Token)

// Method 2: Create client directly from config (recommended)
client, err := client.NewClientFromConfig()
if err != nil {
    log.Fatal(err)
}

// Method 3: Create client from config with panic on error
client := client.MustNewClientFromConfig()
```

## Client Constructors

### Basic Constructor
```go
// Create client with explicit URL and token
client := client.NewClient("https://data.gopher-ai.com", "your-api-token")
```

### Configuration-Based Constructor (Recommended)
```go
// Create client from environment variables - handles config loading automatically
client, err := client.NewClientFromConfig()
if err != nil {
    log.Fatal(err)
}

// This method automatically:
// 1. Loads GOPHER_CLIENT_URL and GOPHER_CLIENT_TOKEN from environment
// 2. Creates the client with the loaded configuration
// 3. Returns an error if configuration is invalid or missing
```

### Panic-on-Error Constructor
```go
// Create client from config, panics on error - useful for initialization
client := client.MustNewClientFromConfig()

// This is equivalent to:
// client, err := client.NewClientFromConfig()
// if err != nil {
//     panic(err)
// }
```

## Client Methods

### Core Client Operations

#### Job Management
- `GetJobStatus(jobID string) (*types.IndexerJobResult, error)` - Get the status of a job
- `GetResult(jobID string, receiver any) error` - Get the result of a completed job

### AI Analysis

#### Analyze Data
```go
import "github.com/gopher-lab/gopher-client/types"

// Analyze data with full control
response, err := client.AnalyzeData(
    []string{"tweet1", "tweet2", "tweet3"}, // data to analyze
    "Analyze sentiment of these tweets",     // prompt
    "openai/gpt-4o-mini",                   // model (optional)
    false,                                  // app mode (optional)
    []types.ChatHistoryItem{                // chat history (optional)
        {Query: "previous query", Timestamp: "2024-01-01T00:00:00Z"},
    },
    "current search query",                 // current query (optional)
)
```

#### Simple Analysis
```go
// Analyze with default settings
response, err := client.AnalyzeDataSimple(
    []string{"tweet1", "tweet2"},
    "What is the sentiment?",
)
```

#### Get Available Models
```go
// Get list of available AI models
models, err := client.GetAvailableModels()
if err != nil {
    log.Fatal(err)
}

for _, model := range models {
    fmt.Printf("Available model: %s\n", model)
}
```

### Search Term Extraction

#### Extract Search Terms
```go
import "github.com/gopher-lab/gopher-client/types"

// Extract search terms with custom max terms
response, err := client.ExtractSearchTerms(
    "I want to find information about machine learning and AI", // user input
    5, // max terms (1-6, default 4)
)
```

#### Extract with Defaults
```go
// Extract search terms with default max terms (4)
response, err := client.ExtractSearchTermsWithDefaults(
    "Find articles about blockchain technology",
)
```

### Query Contextualization

#### Contextualize Query
```go
import "github.com/gopher-lab/gopher-client/types"

// Contextualize query with chat history
response, err := client.ContextualizeQuery(
    "Tell me more about that",              // current query
    []types.ChatHistoryItem{                // chat history
        {Query: "What is machine learning?", Timestamp: "2024-01-01T00:00:00Z"},
        {Query: "How does it work?", Timestamp: "2024-01-01T00:01:00Z"},
    },
    5, // max history items (1-10, default 5)
)
```

#### Contextualize with Defaults
```go
// Contextualize with default max history items (5)
response, err := client.ContextualizeQueryWithDefaults(
    "What are the applications?",
    chatHistory,
)
```

### Web Search

#### Basic Web Search
```go
// Simple web page scraping
result, err := client.PerformWebSearch("https://example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Job ID: %s\n", result.JobID)
```

#### Advanced Web Search
```go
import "github.com/masa-finance/tee-worker/api/args/web/page"

// Custom web job with specific arguments
args := page.NewArguments()
args.URL = "https://example.com"
// Set additional arguments as needed

result, err := client.PostWebJob(args)
if err != nil {
    log.Fatal(err)
}
```

### Social Media Search

#### Twitter
```go
// Search Twitter for posts
result, err := client.PerformTwitterSearch("golang programming")
if err != nil {
    log.Fatal(err)
}

// Advanced Twitter search with custom arguments
import "github.com/masa-finance/tee-worker/api/args/twitter/search"

args := search.NewArguments()
args.Query = "golang programming"
// Set additional arguments as needed

result, err := client.PostTwitterJob(args)
```

#### Reddit
```go
// Search Reddit posts
result, err := client.PerformRedditSearchPosts("golang", 10)
if err != nil {
    log.Fatal(err)
}

// Search Reddit users
result, err := client.PerformRedditSearchUsers("username", 5)
if err != nil {
    log.Fatal(err)
}

// Search Reddit communities
result, err := client.PerformRedditSearchCommunities("programming", 20)
if err != nil {
    log.Fatal(err)
}

// Scrape specific Reddit URLs
result, err := client.PerformRedditScrapeURL("https://reddit.com/r/golang", 10)
if err != nil {
    log.Fatal(err)
}

// Advanced Reddit search with custom arguments
import "github.com/masa-finance/tee-worker/api/args/reddit/search"

args := search.NewSearchPostsArguments()
args.Queries = []string{"golang", "programming"}
args.MaxItems = 10
// Set additional arguments as needed

result, err := client.PostRedditJob(args)
```

#### LinkedIn
```go
import ptypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"

// Search LinkedIn profiles
result, err := client.PerformLinkedInSearch("software engineer", ptypes.ScraperMode)
if err != nil {
    log.Fatal(err)
}

// Advanced LinkedIn search with custom arguments
import "github.com/masa-finance/tee-worker/api/args/linkedin/profile"

args := profile.NewArguments()
args.ScraperMode = ptypes.ScraperMode
args.Query = "software engineer"
// Set additional arguments as needed

result, err := client.PostLinkedInJob(args)
```

#### TikTok
```go
// Search TikTok videos
result, err := client.PerformTikTokSearch("golang tutorial", 10)
if err != nil {
    log.Fatal(err)
}

// Get trending TikTok videos
result, err := client.PerformTikTokSearchByTrending("views", 20)
if err != nil {
    log.Fatal(err)
}

// Transcribe TikTok video
result, err := client.PerformTikTokTranscription("https://tiktok.com/@user/video/123")
if err != nil {
    log.Fatal(err)
}

// Advanced TikTok search with custom arguments
import "github.com/masa-finance/tee-worker/api/args/tiktok/query"

args := query.NewArguments()
args.Search = []string{"golang tutorial"}
args.MaxItems = 10
// Set additional arguments as needed

result, err := client.PostTikTokJob(map[string]any{
    "search": args.Search,
    "maxItems": args.MaxItems,
})
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
    
    "github.com/gopher-lab/gopher-client/client"
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
    
    "github.com/gopher-lab/gopher-client/client"
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
