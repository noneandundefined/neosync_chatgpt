package analytics

import (
	"context"
	"database/sql"
	"encoding/json"
	"neomatica/neosync/infra/logger"
	"sync"
	"time"
)

const (
	collectorWorkers   = 4
	collectorQueueSize = 10000
	collectorBatchSize = 100
	collectorRetention = 400 * 24 * time.Hour
)

type Event struct {
	EventID    string
	OccurredAt time.Time
	UserUUID   string
	RoleCode   string
	SessionID  string
	EventName  string
	Category   string
	Source     string
	Path       string
	EntityType string
	EntityID   string
	DeviceID   *uint64
	Success    *bool
	DurationMs *int
	ErrorCode  string
	AppVersion string
	Properties map[string]interface{}
}

type Collector struct {
	db     *sql.DB
	queue  chan Event
	closed chan struct{}
	once   sync.Once
	wg     sync.WaitGroup
}

func NewCollector(db *sql.DB) *Collector {
	collector := &Collector{
		db:     db,
		queue:  make(chan Event, collectorQueueSize),
		closed: make(chan struct{}),
	}

	for i := 0; i < collectorWorkers; i++ {
		collector.wg.Add(1)
		go collector.worker()
	}
	collector.wg.Add(1)
	go collector.cleanupWorker()

	return collector
}

func (c *Collector) cleanupWorker() {
	defer c.wg.Done()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			_, err := c.db.ExecContext(ctx, `DELETE FROM product_analytics_events WHERE occurred_at < $1`, time.Now().UTC().Add(-collectorRetention))
			cancel()
			if err != nil {
				logger.Error("AnalyticsCollector: failed cleanup old events: %s", err.Error())
			}
		case <-c.closed:
			return
		}
	}
}

func (c *Collector) Enqueue(events []Event) int {
	accepted := 0
	for _, event := range events {
		select {
		case c.queue <- event:
			accepted++
		default:
			logger.Error("AnalyticsCollector: event queue is full")
			return accepted
		}
	}

	return accepted
}

func (c *Collector) Close() {
	c.once.Do(func() {
		close(c.closed)
		c.wg.Wait()
	})
}

func (c *Collector) worker() {
	defer c.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	batch := make([]Event, 0, collectorBatchSize)
	flush := func() {
		if len(batch) == 0 {
			return
		}

		if err := c.insertBatch(batch); err != nil {
			logger.Error("AnalyticsCollector: failed insert batch: %s", err.Error())
		}
		batch = batch[:0]
	}

	for {
		select {
		case event := <-c.queue:
			batch = append(batch, event)
			if len(batch) >= collectorBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-c.closed:
			for {
				select {
				case event := <-c.queue:
					batch = append(batch, event)
				default:
					flush()
					return
				}
			}
		}
	}
}

func (c *Collector) insertBatch(events []Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO product_analytics_events (
			event_id, occurred_at, user_uuid, role_code, session_id, event_name, category,
			source, path, entity_type, entity_id, device_id, success, duration_ms,
			error_code, app_version, properties
		)
		VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5, $6, $7, $8, NULLIF($9, ''),
			NULLIF($10, ''), NULLIF($11, ''), COALESCE($12::BIGINT, (SELECT id FROM devices WHERE imei = NULLIF($11, ''))), $13, $14, NULLIF($15, ''),
			NULLIF($16, ''), $17)
		ON CONFLICT (event_id) DO NOTHING
	`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, event := range events {
		properties, err := json.Marshal(event.Properties)
		if err != nil {
			properties = []byte("{}")
		}

		if _, err := stmt.ExecContext(ctx, event.EventID, event.OccurredAt, event.UserUUID,
			event.RoleCode, event.SessionID, event.EventName, event.Category, event.Source,
			event.Path, event.EntityType, event.EntityID, event.DeviceID, event.Success,
			event.DurationMs, event.ErrorCode, event.AppVersion, properties); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
