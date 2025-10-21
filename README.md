# Gopher Client

A Go client library for interacting with the Gopher AI data collection and search API. This client provides easy-to-use methods for performing various types of searches and data collection across multiple platforms including web, social media, and other data sources.

## Installation

```bash
go get github.com/gopher-lab/gopher-client
```

## Quick Start

The client can be configured using environment variables or by passing parameters directly. By default, a single pooled HTTP client is reused across requests for performance; you can customize pooling and timeouts with options if needed (see section below).

### Environment Variables

```bash
export GOPHER_CLIENT_TOKEN="your-api-token"
export GOPHER_CLIENT_TIMEOUT="120s"  # Optional: default is 60s
export GOPHER_CLIENT_URL="https://data.gopher-ai.com/api" # Optional: default is present
```

**Note:** The `GOPHER_CLIENT_TIMEOUT` environment variable automatically configures timeouts for all `AndWait` methods and HTTP requests.

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
    c := client.NewClient("https://data.gopher-ai.com/api", "your-api-token")
    
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
results, err := client.PerformWebSearchAndWait("https://example.com")
```

### 🐦 Twitter Search
```go
// Submit job
result, err := client.PerformTwitterSearch("golang programming")
fmt.Printf("Job ID: %s\n", result.UUID)

// Wait for results
results, err := client.PerformTwitterSearchAndWait("golang programming")
```

### 👽 Reddit Search
```go
// Submit jobs
result, err := client.PerformRedditSearchPosts("golang", 10)
result, err := client.PerformRedditSearchUsers("username", 5)
result, err := client.PerformRedditScrapeURL("https://reddit.com/r/golang", 10)

// Wait for results
results, err := client.PerformRedditSearchPostsAndWait("golang", 10)
results, err := client.PerformRedditSearchUsersAndWait("username", 5)
results, err := client.PerformRedditScrapeURLAndWait("https://reddit.com/r/golang", 10)
```

### 💼 LinkedIn Search
```go
import ptypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"

// Submit job
result, err := client.PerformLinkedInSearch("software engineer", ptypes.ScraperModeShort)

// Wait for results
results, err := client.PerformLinkedInSearchAndWait("software engineer", ptypes.ScraperModeShort)
```

### 🎵 TikTok Search
```go
// Submit jobs
result, err := client.PerformTikTokSearch("golang tutorial", 10)
result, err := client.PerformTikTokSearchByTrending("views", 20)
result, err := client.PerformTikTokTranscription("https://tiktok.com/@user/video/123")

// Wait for results
results, err := client.PerformTikTokSearchAndWait("golang tutorial", 10)
results, err := client.PerformTikTokSearchByTrendingAndWait("views", 20)
results, err := client.PerformTikTokTranscriptionAndWait("https://tiktok.com/@user/video/123")
```

### 🔍 Similarity Search
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

## 🔧 Advanced Usage with Flexible Arguments

For more control over job parameters, you can use the `Post*JobAndWait` methods with flexible argument types. These methods allow you to customize all available options for each platform.

### 🌐 Advanced Web Search
```go
import "github.com/masa-finance/tee-worker/api/args/web/page"

// Create custom web search arguments
args := page.NewArguments()
args.URL = "https://example.com"
args.MaxDepth = 3
args.FollowRedirects = true
args.UserAgent = "CustomBot/1.0"

// Submit job with custom arguments and wait
results, err := client.PostWebJobAndWait(args)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d documents\n", len(results))
```

### 🐦 Advanced Twitter Search
```go
import "github.com/masa-finance/tee-worker/api/args/twitter/search"

// Create custom Twitter search arguments
args := search.NewArguments()
args.Query = "golang programming"
args.MaxResults = 100
args.Language = "en"
args.ResultType = "recent"
args.IncludeEntities = true

// Submit job with custom arguments and wait
results, err := client.PostTwitterJobAndWait(args)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d tweets\n", len(results))
```

### 👽 Advanced Reddit Search
```go
import "github.com/masa-finance/tee-worker/api/args/reddit/search"

// Search for posts with custom arguments
args := search.NewSearchPostsArguments()
args.Queries = []string{"golang", "programming"}
args.MaxItems = 50
args.Sort = "hot"
args.TimeFilter = "week"

// Submit job with custom arguments and wait
results, err := client.PostRedditJobAndWait(args)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d Reddit posts\n", len(results))
```

### 💼 Advanced LinkedIn Search
```go
import (
    "github.com/masa-finance/tee-worker/api/args/linkedin/profile"
    ptypes "github.com/masa-finance/tee-worker/api/types/linkedin/profile"
)

// Create custom LinkedIn search arguments
args := profile.NewArguments()
args.Query = "software engineer"
args.ScraperMode = ptypes.ScraperModeShort
args.MaxResults = 25
args.Location = "San Francisco"
args.ExperienceLevel = "mid-senior"

// Submit job with custom arguments and wait
results, err := client.PostLinkedInJobAndWait(args)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d LinkedIn profiles\n", len(results))
```

### 🎵 Advanced TikTok Search
```go
import (
    "github.com/masa-finance/tee-worker/api/args/tiktok/query"
    "github.com/masa-finance/tee-worker/api/args/tiktok/trending"
    "github.com/masa-finance/tee-worker/api/args/tiktok/transcription"
)

// Custom TikTok search
searchArgs := query.NewArguments()
searchArgs.Search = []string{"golang tutorial", "go programming"}
searchArgs.MaxItems = 30
searchArgs.SortBy = "date"

results, err := client.PostTikTokSearchJobAndWait(searchArgs)
if err != nil {
    log.Fatal(err)
}

// Custom TikTok trending search
trendingArgs := trending.NewArguments()
trendingArgs.SortBy = "views"
trendingArgs.MaxItems = 50
trendingArgs.Region = "US"

trendingResults, err := client.PostTikTokTrendingJobAndWait(trendingArgs)
if err != nil {
    log.Fatal(err)
}

// Custom TikTok transcription
transcriptionArgs := transcription.NewArguments()
transcriptionArgs.VideoURL = "https://tiktok.com/@user/video/123"
transcriptionArgs.Language = "en"

transcriptionResults, err := client.PostTikTokTranscriptionJobAndWait(transcriptionArgs)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d search results, %d trending videos, %d transcriptions\n", 
    len(results), len(trendingResults), len(transcriptionResults))
```

### 🔄 Batch Processing Example
```go
// Process multiple sources in parallel
var wg sync.WaitGroup
results := make(map[string][]types.Document)

// Twitter search
wg.Add(1)
go func() {
    defer wg.Done()
    args := search.NewArguments()
    args.Query = "artificial intelligence"
    args.MaxResults = 50
    
    docs, err := client.PostTwitterJobAndWait(args)
    if err == nil {
        results["twitter"] = docs
    }
}()

// Reddit search
wg.Add(1)
go func() {
    defer wg.Done()
    args := search.NewSearchPostsArguments()
    args.Queries = []string{"AI", "machine learning"}
    args.MaxItems = 30
    
    docs, err := client.PostRedditJobAndWait(args)
    if err == nil {
        results["reddit"] = docs
    }
}()

// Web search
wg.Add(1)
go func() {
    defer wg.Done()
    args := page.NewArguments()
    args.URL = "https://example-ai-blog.com"
    args.MaxDepth = 2
    
    docs, err := client.PostWebJobAndWait(args)
    if err == nil {
        results["web"] = docs
    }
}()

wg.Wait()

// Process all results
totalDocs := 0
for source, docs := range results {
    fmt.Printf("%s: %d documents\n", source, len(docs))
    totalDocs += len(docs)
}
fmt.Printf("Total documents collected: %d\n", totalDocs)
```


## HTTP Client, Timeouts, and Connection Pooling

This client now reuses a single pooled `*http.Client` under the hood for all requests to improve performance and connection reuse. Existing constructors and usage remain unchanged.

- `NewClient(baseURL, token)` keeps the default 60s timeout.
- `NewClientFromConfig()` uses the timeout from your environment/config.
- For advanced control, use `NewClientWithOptions` and functional options.

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

### Environment timeouts

`GOPHER_CLIENT_TIMEOUT` affects both job polling (the `AndWait` helpers) and the HTTP client's request timeout when using `NewClientFromConfig()`.
