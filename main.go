package main

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type clickhouseConfig struct {
	Host     string `default:"localhost:9000"`
	Username string `default:"default"`
	Password string `default:""`
	Database string `default:"some_db_dev"`
}

// NewClickhouseClient returns a connection or error if failed to connect/ping
func NewClickhouseClient(ctx context.Context) (clickhouse.Conn, error) {
	cfg := clickhouseConfig{
		Host:     "localhost:9000",
		Username: "default",
		Password: "",
		Database: "some_db_dev",
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Host},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	return conn, nil
}

func main() {
	ctx := context.Background()
	clickhouseConn, err := NewClickhouseClient(ctx)
	if err != nil {
		panic(err)
	}
	if err := clickhouseConn.Exec(ctx, `
INSERT INTO raw_events (
	id,
	event_id,
	occurred_at,
	entity_id
) VALUES (
          	'8b5a1bc0-7bca-41f9-a4c8-1792f82493bf',
          	'7b5a1bc0-7bca-41f9-a4c8-1792f82493bf',
          	'2022-07-30 00:00:00',
          	'0c964b58-9249-414c-bdbc-b3e47a56f3aa'
)
`); err != nil {
		panic(err)
	}

	rawEventsRow := clickhouseConn.QueryRow(ctx, "SELECT COUNT(*) FROM raw_events")
	var count uint64
	if err = rawEventsRow.Scan(&count); err != nil {
		panic(err)
	}
	fmt.Println("count: ", count)
	fmt.Println("The expected value of count should be 1, when it is run the first time, but produces 0.")

	aggregatedEventsRow := clickhouseConn.QueryRow(ctx, "SELECT countMerge(event_count) FROM events_count")
	var aggregatedCount uint64
	if err = aggregatedEventsRow.Scan(&aggregatedCount); err != nil {
		panic(err)
	}
	fmt.Println("aggregatedCount: ", aggregatedCount)
}
