# Logging Guide

## Overview

TwelveDataClient uses a structured, color-coded logging system with support for different log levels, timestamps, and context-aware messages.

## Log Levels

The logging system supports 5 severity levels:

| Level | Name | Usage | Color |
|-------|------|-------|-------|
| **0** | DEBUG | Detailed diagnostic information | Dim |
| **1** | INFO | General informational messages | Cyan |
| **2** | WARN | Warning messages | Yellow |
| **3** | ERROR | Error messages | Red |
| **4** | FATAL | Fatal errors that stop execution | Red (Bold) |

## Basic Usage

### Simple Logging

```go
import "twelve_data_client/internal/logger"

// Using global functions
logger.Debug("Debug message with %s", "details")
logger.Info("Application started")
logger.Warn("This is a warning")
logger.Error("An error occurred: %v", err)
logger.Fatal("Critical error: %v", err)  // Exits with code 1
```

### Creating a Custom Logger

```go
// Create logger with custom configuration
lg := logger.NewLogger(&logger.Config{
    Level:        logger.InfoLevel,
    UseTimestamp: true,
    UseColor:     true,
    Prefix:       "MyApp",
})

lg.Info("Using custom logger")
```

## Advanced Logging Functions

### LogError - Log errors with context

```go
err := someFunction()
logger.LogError("database operation", err, "table: users", "action: insert")
// Output: [15:04:05] [ERROR] MyApp database operation failed: connection timeout (table: users, action: insert)
```

### LogSuccess - Log successful operations

```go
logger.LogSuccess("API request", "endpoint: /stocks", "status: 200")
// Output: [15:04:05] [INFO] MyApp API request successful: endpoint: /stocks, status: 200
```

### LogRequest - Log HTTP requests

```go
logger.LogRequest("GET", "https://api.twelvedata.com/stocks")
// Output: [15:04:05] [DEBUG] MyApp Request: GET https://api.twelvedata.com/stocks
```

### LogResponse - Log HTTP responses

```go
logger.LogResponse(200, "Stock list retrieved")
// Output: [15:04:05] [DEBUG] MyApp Response: 200 Stock list retrieved
```

### LogDuration - Log operation duration

```go
start := time.Now()
// ... do something
logger.LogDuration("API call", time.Since(start))
// Output: [15:04:05] [DEBUG] MyApp API call completed in 245ms
```

## Configuration

### Set Log Level

```go
// Set minimum log level
logger.SetLevel(logger.DebugLevel)

// Now all debug messages will be shown
logger.Debug("This debug message will be logged")
```

### Set Logger Prefix

```go
// Set a prefix for all messages
logger.SetPrefix("TwelveDataClient")

// All messages will include this prefix
logger.Info("Starting application")
// Output: [15:04:05] [INFO] TwelveDataClient Starting application
```

## Real-World Examples

### API Integration Logging

```go
func GetAllStocks(exchange string, apiKey string) ([]model.Stock, error) {
    // Log the operation
    logger.LogRequest("GET", "/stocks")
    
    resp, err := http.Get(url)
    if err != nil {
        // Log error with context
        logger.LogError("fetching stocks", err, exchange)
        return nil, err
    }
    
    // Log response
    logger.LogResponse(resp.StatusCode, "stocks retrieved")
    
    // Process response...
    logger.Debug("parsed %d stocks", len(stocks))
    
    return stocks, nil
}
```

### WebSocket Connection Logging

```go
func GetTwelveDataWebSocket(apiKey string, subscription model.Subscription) (*websocket.Conn, error) {
    logger.Debug("connecting to WebSocket: %s", url)
    
    connection, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        logger.LogError("WebSocket connection", err)
        return nil, err
    }
    
    logger.LogSuccess("WebSocket connected")
    
    logger.Debug("sending subscription for %d symbols", len(symbols))
    if err := connection.WriteJSON(subscription); err != nil {
        logger.LogError("sending subscription", err)
        return nil, err
    }
    
    logger.LogSuccess("subscription sent")
    return connection, nil
}
```

### Error Handling

```go
// Log and return error
if apiKey == "" {
    logger.Error("API key is missing")
    return nil, fmt.Errorf("API key required")
}

// Log error with details
if err != nil {
    logger.LogError("parsing JSON", err, "response from stocks endpoint")
    return nil, err
}

// Log and exit on fatal error
if connection == nil {
    logger.Fatal("failed to establish WebSocket connection")
}
```

## Output Examples

### With Timestamps and Colors

```
[15:04:05] [INFO] TwelveDataClient API key: ...abc12
[15:04:05] [DEBUG] TwelveDataClient endpoint: https://api.twelvedata.com
[15:04:05] [DEBUG] TwelveDataClient Request: GET /stocks
[15:04:05] [DEBUG] TwelveDataClient Response: 200 Stock list retrieved
[15:04:05] [INFO] TwelveDataClient Successfully fetched 20 stocks
[15:04:05] [DEBUG] TwelveDataClient parsed 20 stocks
[15:04:05] [WARN] TwelveDataClient No data available for SYMBOL
[15:04:05] [ERROR] TwelveDataClient parsing failed: unexpected token at line 1
```

### Without Colors (for logs)

```
[15:04:05] [INFO] TwelveDataClient Starting application
[15:04:05] [DEBUG] TwelveDataClient Connecting to API
[15:04:05] [ERROR] TwelveDataClient Connection failed: EOF
```

## Best Practices

### ✅ Do's

1. **Use appropriate log levels**
   ```go
   logger.Debug("Detailed info for developers")
   logger.Info("Important application events")
   logger.Warn("Something unexpected")
   logger.Error("Error occurred but continue")
   logger.Fatal("Critical error, must exit")
   ```

2. **Include context in errors**
   ```go
   logger.LogError("fetching user data", err, "userID: 123", "endpoint: /users")
   ```

3. **Log at service boundaries**
   ```go
   // Log when entering a service
   logger.Debug("GetAllStocks called with exchange: %s", exchange)
   
   // Log at important checkpoints
   logger.LogResponse(resp.StatusCode, "API response")
   
   // Log when exiting
   logger.LogSuccess("GetAllStocks", fmt.Sprintf("returned %d stocks", len(stocks)))
   ```

4. **Use LogSuccess for positive outcomes**
   ```go
   logger.LogSuccess("WebSocket connection", "connected to:", url)
   ```

### ❌ Don'ts

1. **Don't use logger.Fatal for expected errors**
   ```go
   // BAD
   logger.Fatal("user not found: %v", err)
   
   // GOOD
   logger.Error("user not found: %v", err)
   return nil, err
   ```

2. **Don't log sensitive information**
   ```go
   // BAD
   logger.Debug("API Key: %s", apiKey)
   
   // GOOD
   logger.Debug("API Key: ...%s", apiKey[len(apiKey)-4:])
   ```

3. **Don't log the same error twice**
   ```go
   // BAD
   if err != nil {
       logger.Error(err)
       return err  // Error will be logged again by caller
   }
   
   // GOOD
   if err != nil {
       return fmt.Errorf("operation failed: %w", err)  // Let caller log once
   }
   ```

4. **Don't use string concatenation for log messages**
   ```go
   // BAD
   logger.Info("User " + name + " logged in")
   
   // GOOD
   logger.Info("User %s logged in", name)
   ```

## Performance Considerations

- Logging with timestamps and colors has minimal performance impact
- Debug logs are still processed even if not printed (use level filter)
- For high-frequency operations, consider filtering to higher log levels

## Disabling Logs

Set log level to prevent log output:

```go
// Disable all logging below a certain level
logger.SetLevel(logger.WarnLevel)  // Only show WARN, ERROR, FATAL

// Or create a logger that writes to /dev/null
lg := logger.NewLogger(&logger.Config{
    Output: io.Discard,  // Discard all output
})
```

## Migration from Old Logging

### Before

```go
fmt.Println("Starting application")
log.Printf("Error: %v", err)
fmt.Printf("Success: %d items", count)
```

### After

```go
logger.Info("Starting application")
logger.LogError("operation", err)
logger.LogSuccess("Processing", fmt.Sprintf("returned %d items", count))
```

---

For more examples, see `main.go` and `internal/services/*.go`
