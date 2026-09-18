# Logging and configuration history

## Goals

- Keep production file logs small enough to inspect directly.
- Stop using full server log archives to investigate one IMEI.
- Preserve meaningful per-device events for support.
- Preserve old tracker configurations and show exactly which UIDs changed.
- Keep synchronization failures attached to the configuration version that produced them.

## File logging

Routine ADM RX/TX packet lines no longer go to the normal TCP `server.log`.

Detailed packet tracing is written to `packets.log` only when:

```env
LOG_PACKET_TRACE=true
```

It is enabled automatically in DEV. Keep it disabled in production and enable it temporarily only while reproducing a protocol problem.

If the asynchronous logger queue is full, the logger no longer prints every dropped message to stdout. It prints an aggregated dropped-message counter at most once every five seconds. TCP file logging uses one writer to preserve line order.

## Per-device event log

Recent support events remain in the Redis stream:

```
adm:device:<imei>:events
```

Defaults:

- maximum 1000 events per IMEI;
- TTL: 7 days;
- stable SYNC payloads are sampled at most once per 30 minutes;
- periodic telemetry command responses are not written to the support event stream;
- meaningful user command answers, connect/disconnect, configuration and synchronization events remain available.

The existing device event screen reads this stream directly by IMEI.

A CSV export is available at:

```
GET /api/v1/devices/{imei}/logs/export
```

This avoids downloading global log archives and then grepping them locally.

## Configuration version history

Migration:

```
20260918143000_enhance_configuration_history.sql
```

The `configuration_history` table now keeps up to 100 binary snapshots per device.

Each version stores:

- timestamp;
- configuration hash;
- complete binary configuration;
- origin (`baseline`, `neosync`, `tracker`, or `unknown`);
- latest synchronization status for that version;
- last non-empty synchronization error observed for that version.

The migration also stores the currently active configuration as a baseline when needed. Existing old history rows are preserved.

### API

List versions and computed UID changes:

```
GET /api/v1/devices/{imei}/configuration/histories
```

The server parses adjacent binary snapshots and returns only changed fields with:

- numeric UID;
- field name;
- configuration section;
- previous value;
- new value.

Download a historical binary configuration:

```
GET /api/v1/devices/{imei}/configuration/histories/{historyId}/export
```

## UI

The device event page now has two tabs:

1. Terminal events
   - fast per-IMEI Redis events;
   - CSV export.

2. Configuration history
   - up to 100 versions;
   - hash, source and sync state;
   - synchronization error;
   - expandable UID-level diff;
   - historical binary download.

Configuration history is fetched only when the tab is opened.

## Operational result

The normal production log volume should drop substantially because periodic protocol RX/TX and telemetry traffic no longer dominate disk logs. Full packet tracing is still available on demand without making it the default logging mode.
