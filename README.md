# Gopher Client

A Go client library for interacting with the Gopher AI data collection and search API. This client provides easy-to-use methods for performing various types of searches and data collection across multiple platforms including web, social media, and other data sources.

## Installation

```bash
go get github.com/gopher-lab/gopher-client
```

## HTTP Client, Timeouts, and Connection Pooling

This client reuses a single pooled `*http.Client` under the hood for all requests to improve performance and connection reuse. The default timeout is 60s.

- `NewClient(baseURL, token)` keeps the defaults.
- `NewClientFromConfig()` uses your environment/config.
- `NewClientWithOptions` uses functional options.

### Environment Variables

```bash
export GOPHER_CLIENT_TOKEN="your-api-token"
export GOPHER_CLIENT_TIMEOUT="120s"  # Optional: default is 60s
export GOPHER_CLIENT_URL="https://data.gopher-ai.com/api" # Optional: default is present
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
    
    // Scrape a web page (async)
    job, err := c.ScrapeWebAsync("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Job ID: %s\n", job.UUID)
}
```

### Configure with Functional Options

```go
package main

import (
    "time"
    "github.com/gopher-lab/gopher-client/client"
)

func main() {
    c, err := client.NewClientWithOptions(
        "https://data.gopher-ai.com/api",
        "your-api-token",
        client.Timeout(90*time.Second),
        client.MaxIdleConnsPerHost(50),
        client.MaxConnsPerHost(200),
    )
    if err != nil { panic(err) }

    _ = c // use the client
}
```

## Client Methods

### 🌐 Web Scraping
```go
// Submit job (async)
job, err := client.ScrapeWebAsync("https://example.com")
fmt.Printf("Job ID: %s\n", job.UUID)

// Get results directly (sync)
results, err := client.ScrapeWeb("https://example.com")
```

### 🐦 Twitter Search
```go
// Submit job (async)
job, err := client.SearchTwitterAsync("golang programming")
fmt.Printf("Job ID: %s\n", job.UUID)

// Get results directly (sync)
results, err := client.SearchTwitter("golang programming")
```

### 👽 Reddit Operations
```go
// Scrape Reddit URL (async)
job, err := client.ScrapeRedditURLAsync("https://reddit.com/r/golang", 10)
fmt.Printf("Job ID: %s\n", job.UUID)

// Scrape Reddit URL (sync)
results, err := client.ScrapeRedditURL("https://reddit.com/r/golang", 10)

// Search Reddit posts (sync)
results, err := client.SearchRedditPosts("golang", 10)

// Search Reddit users (sync)
results, err := client.SearchRedditUsers("username", 5)

// Search Reddit communities (sync)
results, err := client.SearchRedditCommunities("golang", 10)
```

### 💼 LinkedIn Search
```go
// Submit job (async)
job, err := client.SearchLinkedInAsync("software engineer")
fmt.Printf("Job ID: %s\n", job.UUID)

// Get results directly (sync)
results, err := client.SearchLinkedIn("software engineer")
```

### 🎵 TikTok Operations
```go
// Submit jobs (async)
job, err := client.SearchTikTokAsync("golang tutorial", 10)
fmt.Printf("Job ID: %s\n", job.UUID)

job, err := client.SearchTikTokTrendingAsync("views", 20)
fmt.Printf("Job ID: %s\n", job.UUID)

job, err := client.TranscribeTikTokAsync("https://tiktok.com/@user/video/123")
fmt.Printf("Job ID: %s\n", job.UUID)

// Get results directly (sync)
results, err := client.SearchTikTok("golang tutorial", 10)
results, err := client.SearchTikTokTrending("views", 20)
results, err := client.TranscribeTikTok("https://tiktok.com/@user/video/123")
```

### 🔍 Search & Analysis
```go
import "github.com/masa-finance/tee-worker/api/types"

sources := []types.Source{types.WebSource, types.TwitterSource, types.RedditSource}

// Similarity search (immediate results)
results, err := client.SearchSimilarity(
    "artificial intelligence",
    sources,
    []string{"AI", "machine learning"},
    "and",
    10,
)

// Hybrid search (immediate results)
results, err := client.SearchHybrid(
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
// Analyze data with simple prompt (default settings)
response, err := client.AnalyzeData(
    []string{"tweet1", "tweet2"},
    "What is the sentiment?",
)

// Analyze data with custom arguments
response, err := client.AnalyzeDataWithArgs(
    []string{"tweet1", "tweet2"},
    "What is the sentiment?",
    "openai/gpt-4o",  // Custom model
    true,            // App request
    chatHistory,     // Chat history
    "current query", // Current query
)

// Get available models
models, err := client.GetAvailableModels()
```

### 🔧 Search Tools
```go
// Extract search terms with default maxTerms (4)
response, err := client.ExtractSearchTerms(
    "Find articles about blockchain technology",
    0, // Use 0 for default
)

// Extract search terms with custom maxTerms
response, err := client.ExtractSearchTerms(
    "Find articles about blockchain technology",
    6, // Custom maxTerms
)

// Contextualize query with default maxHistoryItems (5)
response, err := client.ContextualizeQuery(
    "Tell me more about that",
    chatHistory,
    0, // Use 0 for default
)

// Contextualize query with custom maxHistoryItems
response, err := client.ContextualizeQuery(
    "Tell me more about that",
    chatHistory,
    8, // Custom maxHistoryItems
)
```

### 🔧 Advanced Operations with Custom Arguments
```go
// Web scraping with custom arguments
results, err := client.ScrapeWebWithArgs(pageArgs)

// Twitter search with custom arguments  
results, err := client.SearchTwitterWithArgs(searchArgs)

// Reddit search with custom arguments
results, err := client.SearchRedditWithArgs(redditArgs)

// LinkedIn search with custom arguments
results, err := client.SearchLinkedInWithArgs(linkedinArgs)

// TikTok operations with custom arguments
results, err := client.SearchTikTokWithArgs(queryArgs)
results, err := client.SearchTikTokTrendingWithArgs(trendingArgs)
results, err := client.TranscribeTikTokWithArgs(transcriptionArgs)

// Example, advanced usage with LinkedIn
import linkedin "github.com/masa-finance/tee-worker/api/args/linkedin/profile"
import "github.com/masa-finance/tee-worker/api/types/linkedin/profile"
import "github.com/masa-finance/tee-worker/api/types/linkedin/experiences"

args := linkedin.NewArguments()
args.Query = "software engineer"
args.ScraperMode = profile.ScraperModeFull
args.YearsOfExperience = []experiences.Id{experiences.ThreeToFiveYears}
args.SeniorityLevels = []seniorities.Id{seniorities.Senior}
args.Functions = []functions.Id{functions.Engineering}
args.Industries = []industries.Id{industries.SoftwareDevelopment}

// Submit with custom args (async)
job, err := client.SearchLinkedInWithArgsAsync(args)
fmt.Printf("Job ID: %s\n", job.UUID)

// Get results with custom args (sync)
results, err := client.SearchLinkedInWithArgs(args)
```



### 📊 Metrics
```go
// Get all metrics
stats, err := client.GetAllMetrics(false)

// Get specific collection metrics
stats, err := client.GetMetrics("web", true)
```

## Advanced Usage

### Inject a Custom `http.Client`

If you need full control (custom proxies, tracing, etc.), inject your own `*http.Client`. When provided, pool options are ignored in favor of your client.

```go
package main

import (
    "net/http"
    "time"
    "github.com/gopher-lab/gopher-client/client"
)

func main() {
    hc := &http.Client{ Timeout: 30 * time.Second }
    c, err := client.NewClientWithOptions(
        "https://data.gopher-ai.com/api",
        "your-api-token",
        client.HttpClient(hc),
    )
    if err != nil { panic(err) }

    _ = c
}
```

### Local development (self-signed certs)

```go
c, err := client.NewClientWithOptions(
    baseURL, token,
    client.IgnoreTLSCert(), // skip TLS verification (development only)
)
```