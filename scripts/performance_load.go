package main

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "os"
    "sync"
    "time"
    
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/joho/godotenv"
)

// Simulated call record
type CallRecord struct {
    SIMID       int
    ModemID     int
    Destination string
    Duration    int
    Disposition string
}

func main() {
    // Load environment
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Connect to database
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL not set")
    }
    
    ctx := context.Background()
    pool, err := pgxpool.Connect(ctx, dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer pool.Close()
    
    fmt.Println("Connected to database")
    
    // Test parameters
    numGoroutines := 10
    callsPerGoroutine := 100
    totalCalls := numGoroutines * callsPerGoroutine
    
    fmt.Printf("Starting performance test: %d goroutines, %d calls each = %d total calls\n", 
        numGoroutines, callsPerGoroutine, totalCalls)
    
    // Start time
    startTime := time.Now()
    
    // Run concurrent inserts
    var wg sync.WaitGroup
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            insertCalls(ctx, pool, workerID, callsPerGoroutine)
        }(i)
    }
    
    wg.Wait()
    duration := time.Since(startTime)
    
    // Calculate metrics
    callsPerSecond := float64(totalCalls) / duration.Seconds()
    
    fmt.Printf("\nPerformance Test Results:\n")
    fmt.Printf("Total calls inserted: %d\n", totalCalls)
    fmt.Printf("Total time: %v\n", duration)
    fmt.Printf("Calls per second: %.2f\n", callsPerSecond)
    
    // Test query performance
    testQueries(ctx, pool)
}

func insertCalls(ctx context.Context, pool *pgxpool.Pool, workerID, numCalls int) {
    dispositions := []string{"ANSWERED", "NO ANSWER", "BUSY", "FAILED"}
    
    for i := 0; i < numCalls; i++ {
        // Generate random call data
        simID := rand.Intn(200) + 1
        modemID := rand.Intn(20) + 1
        destination := fmt.Sprintf("+1%010d", rand.Intn(10000000000))
        duration := rand.Intn(300) // 0-300 seconds
        disposition := dispositions[rand.Intn(len(dispositions))]
        
        // Insert call record
        query := `
            INSERT INTO call_detail_records (
                sim_card_id, modem_id, asterisk_unique_id, call_direction,
                destination_number, call_start_time, call_end_time,
                duration_seconds, disposition, created_at
            ) VALUES (
                $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
            )`
        
        uniqueID := fmt.Sprintf("test-%d-%d-%d", workerID, i, time.Now().Unix())
        startTime := time.Now().Add(-time.Duration(rand.Intn(86400)) * time.Second) // Random time in last 24h
        endTime := startTime.Add(time.Duration(duration) * time.Second)
        
        _, err := pool.Exec(ctx, query,
            simID, modemID, uniqueID, "outbound",
            destination, startTime, endTime,
            duration, disposition, time.Now(),
        )
        
        if err != nil {
            log.Printf("Worker %d: Error inserting call %d: %v", workerID, i, err)
        }
        
        // Small delay to not overwhelm the database
        if i%10 == 0 {
            time.Sleep(10 * time.Millisecond)
        }
    }
    
    fmt.Printf("Worker %d completed %d calls\n", workerID, numCalls)
}

func testQueries(ctx context.Context, pool *pgxpool.Pool) {
    fmt.Println("\nTesting query performance...")
    
    // Test 1: Count all calls
    start := time.Now()
    var count int64
    err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM call_detail_records").Scan(&count)
    if err != nil {
        log.Printf("Error counting calls: %v", err)
    }
    fmt.Printf("Total calls in database: %d (took %v)\n", count, time.Since(start))
    
    // Test 2: Analytics query (last 24 hours)
    start = time.Now()
    query := `
        SELECT 
            COUNT(*) as total_calls,
            COUNT(CASE WHEN disposition = 'ANSWERED' THEN 1 END) as successful_calls,
            COALESCE(SUM(duration_seconds) / 60.0, 0) as total_minutes
        FROM call_detail_records
        WHERE call_start_time >= NOW() - INTERVAL '24 hours'
    `
    
    var totalCalls, successfulCalls int64
    var totalMinutes float64
    err = pool.QueryRow(ctx, query).Scan(&totalCalls, &successfulCalls, &totalMinutes)
    if err != nil {
        log.Printf("Error running analytics query: %v", err)
    } else {
        fmt.Printf("24h Analytics: %d calls, %d successful, %.2f minutes (took %v)\n", 
            totalCalls, successfulCalls, totalMinutes, time.Since(start))
    }
    
    // Test 3: Hourly breakdown
    start = time.Now()
    query = `
        SELECT 
            EXTRACT(HOUR FROM call_start_time) as hour,
            COUNT(*) as calls
        FROM call_detail_records
        WHERE call_start_time >= NOW() - INTERVAL '24 hours'
        GROUP BY hour
        ORDER BY hour
    `
    
    rows, err := pool.Query(ctx, query)
    if err != nil {
        log.Printf("Error running hourly query: %v", err)
    } else {
        defer rows.Close()
        hourCount := 0
        for rows.Next() {
            hourCount++
        }
        fmt.Printf("Hourly breakdown: %d hours with data (took %v)\n", hourCount, time.Since(start))
    }
    
    // Test 4: Top destinations
    start = time.Now()
    query = `
        SELECT 
            SUBSTRING(destination_number FROM 1 FOR 5) as prefix,
            COUNT(*) as calls
        FROM call_detail_records
        WHERE call_start_time >= NOW() - INTERVAL '24 hours'
        GROUP BY prefix
        ORDER BY calls DESC
        LIMIT 10
    `
    
    rows, err = pool.Query(ctx, query)
    if err != nil {
        log.Printf("Error running top destinations query: %v", err)
    } else {
        defer rows.Close()
        destCount := 0
        for rows.Next() {
            destCount++
        }
        fmt.Printf("Top destinations: %d prefixes (took %v)\n", destCount, time.Since(start))
    }
}