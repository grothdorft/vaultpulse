// Package api provides an HTTP server exposing VaultPulse runtime status.
//
// Endpoints:
//
//	GET /healthz  — liveness check, returns {"status":"ok","time":"..."}
//	GET /leases   — lists all tracked leases with expiry and alert state
//	GET /metrics  — JSON snapshot of internal counters
//
// The server is created via New and started with Start. It supports
// graceful shutdown through Shutdown(ctx).
package api
