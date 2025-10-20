# Gopher Client

A Go client library for interacting with the Gopher AI data collection and search API. This client provides easy-to-use methods for performing various types of searches and data collection across multiple platforms including web, social media, and other data sources.

## Installation

```bash
go get github.com/gopher-lab/gopher-client
```

## Quick Start

The client can be configured using environment variables or by passing parameters directly:

### Environment Variables

```bash
export GOPHER_CLIENT_TOKEN="your-api-token"
export GOPHER_CLIENT_TIMEOUT="120s"  # Optional: default is 60s
export GOPHER_CLIENT_URL="https://data.gopher-ai.com" # Optional: default is present
```

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
    
    fmt.Printf("Job ID: %s\n", result.UUID)
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
    
    fmt.Printf("Job ID: %s\n", result.UUID)
}
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

## Client Methods

### 🌐 Web Search
```go
// Submit job
result, err := client.PerformWebSearch("https://example.com")
fmt.Printf("Job ID: %s\n", result.UUID)

// Wait for results
results, err := client.PerformWebSearchAndWait("https://example.com", 2*time.Minute)
```

### 🐦 Twitter Search
```go
// Submit job
result, err := client.PerformTwitterSearch("golang programming")
fmt.Printf("Job ID: %s\n", result.UUID)

// Wait for results
results, err := client.PerformTwitterSearchAndWait("golang programming", 2*time.Minute)
```

### 👽 Reddit Search
```go
// Submit jobs
result, err := client.PerformRedditSearchPosts("golang", 10)
result, err := client.PerformRedditSearchUsers("username", 5)
result, err := client.PerformRedditScrapeURL("https://reddit.com/r/golang", 10)

// Wait for results
results, err := client.PerformRedditSearchPostsAndWait("golang", 10, 2*time.Minute)
results, err := client.PerformRedditSearchUsersAndWait("username", 5, 2*time.Minute)
results, err := client.PerformRedditScrapeURLAndWait("https://reddit.com/r/golang", 10, 2*time.Minute)
```

### 💼 LinkedIn Search
```go
import ptypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"

// Submit job
result, err := client.PerformLinkedInSearch("software engineer", ptypes.ScraperModeShort)

// Wait for results
results, err := client.PerformLinkedInSearchAndWait("software engineer", ptypes.ScraperModeShort, 2*time.Minute)
```

### 🎵 TikTok Search
```go
// Submit jobs
result, err := client.PerformTikTokSearch("golang tutorial", 10)
result, err := client.PerformTikTokSearchByTrending("views", 20)
result, err := client.PerformTikTokTranscription("https://tiktok.com/@user/video/123")

// Wait for results
results, err := client.PerformTikTokSearchAndWait("golang tutorial", 10, 2*time.Minute)
results, err := client.PerformTikTokSearchByTrendingAndWait("views", 20, 2*time.Minute)
results, err := client.PerformTikTokTranscriptionAndWait("https://tiktok.com/@user/video/123", 3*time.Minute)
```

### 🔍 Advanced Search
```go
import "github.com/masa-finance/tee-worker/api/types"

sources := []types.Source{types.WebSource, types.TwitterSource, types.RedditSource}

// Similarity search (immediate results)
results, err := client.PerformSimilaritySearch(
    "artificial intelligence",
    sources,
    []string{"AI", "machine learning"},
    "and",
    10,
)

// Hybrid search (immediate results)
results, err := client.PerformHybridSearch(
    "machine learning",
    sources,
    "artificial intelligence",
    0.7, 0.3,
    []string{"AI", "ML"},
    "or",
    15,
)
```

### 🤖 AI Analysis
```go
// Analyze data
response, err := client.AnalyzeDataSimple(
    []string{"tweet1", "tweet2"},
    "What is the sentiment?",
)

// Get available models
models, err := client.GetAvailableModels()
```

### 🔧 Search Tools
```go
// Extract search terms
response, err := client.ExtractSearchTermsWithDefaults(
    "Find articles about blockchain technology",
)

// Contextualize query
response, err := client.ContextualizeQueryWithDefaults(
    "Tell me more about that",
    chatHistory,
)
```

### 📊 Metrics
```go
// Get all metrics
stats, err := client.GetAllMetrics(false)

// Get specific collection metrics
stats, err := client.GetMetrics("web", true)
```