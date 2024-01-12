USE some_db_dev;

-- 0001_init_schema.up.sql
CREATE TABLE IF NOT EXISTS raw_events
(
    `id` String CODEC(ZSTD(1)),
    `event_id` String CODEC(ZSTD(1)),
    `occurred_at` DateTime('UTC'),
    `entity_id` String CODEC(ZSTD(1))
)
ENGINE = MergeTree -- using MergeTree as CI clickhouse server is single node and cannot run MergeTree with replication, but we need it in prod
PARTITION BY toYYYYMM(occurred_at)
ORDER BY (entity_id, occurred_at)
TTL occurred_at + INTERVAL 12 MONTH;

CREATE TABLE IF NOT EXISTS events_count
(
    `entity_id` String CODEC(ZSTD(1)),
    `event_count` AggregateFunction(COUNT) COMMENT 'count of events within a time slice, to query use countMerge or countMergeState',
    `occurred_at` DateTime('UTC') COMMENT 'occurred at rolled up to the earliest hour'
)
ENGINE = AggregatingMergeTree
PARTITION BY toYYYYMM(occurred_at) -- partitioning by YYYYMM as opposed YYYYMMDD to gives better performance as our seeks do at many times query over more than a day. This means we need to scan across more partitions if we partition by day. hence a YYYYMM partition for us makes more sense as we need to scan over less partitions (maybe at most 3).
ORDER BY (entity_id, occurred_at)
TTL occurred_at + INTERVAL 4 MONTH
SETTINGS index_granularity=64; -- use 64 for scanning lesser rows, useful for very low latency applications but increased insertion latency

CREATE MATERIALIZED VIEW IF NOT EXISTS events_count_mv
            TO events_count AS
SELECT
    entity_id,
    countState() AS event_count,
    toStartOfHour(toDateTime(occurred_at)) AS occurred_at
FROM raw_events
GROUP BY ALL;
