# Security Monitoring and Alerting Guide

## Overview

This guide provides implementation patterns for security monitoring, event correlation, and incident response for the panel.go application.

---

## Current Implementation

The application already has audit logging middleware that captures:
- Login attempts (success/failure)
- Logout events
- Resource access (CRUD operations)
- Permission checks
- Settings changes
- Failed authorization attempts

**Location**: `pkg/middleware/audit.go`

---

## Security Event Monitoring

### 1. Real-Time Security Alerts

```go
package monitoring

import (
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/pkg/middleware"
)

type SecurityMonitor struct {
	alertThresholds map[string]int
	eventCounts     map[string]int
	alertChannel    chan SecurityAlert
}

type SecurityAlert struct {
	Severity    string    // "critical", "high", "medium", "low"
	EventType   string
	Description string
	Timestamp   time.Time
	Metadata    map[string]interface{}
}

func NewSecurityMonitor() *SecurityMonitor {
	return &SecurityMonitor{
		alertThresholds: map[string]int{
			"login_failure":      5,  // 5 failed logins
			"permission_denied":  10, // 10 permission denials
			"rate_limit_hit":     3,  // 3 rate limit violations
		},
		eventCounts:  make(map[string]int),
		alertChannel: make(chan SecurityAlert, 100),
	}
}

func (sm *SecurityMonitor) ProcessEvent(event middleware.AuditEvent) {
	// Check for suspicious patterns
	if event.EventType == "login_failure" {
		sm.eventCounts[event.Email]++
		if sm.eventCounts[event.Email] >= sm.alertThresholds["login_failure"] {
			sm.SendAlert(SecurityAlert{
				Severity:    "high",
				EventType:   "brute_force_attempt",
				Description: fmt.Sprintf("Multiple failed login attempts for %s", event.Email),
				Timestamp:   time.Now(),
				Metadata: map[string]interface{}{
					"email":    event.Email,
					"ip":       event.IP,
					"attempts": sm.eventCounts[event.Email],
				},
			})
		}
	}

	// Check for permission violations
	if event.EventType == "permission_check" && !event.Success {
		sm.SendAlert(SecurityAlert{
			Severity:    "medium",
			EventType:   "unauthorized_access_attempt",
			Description: fmt.Sprintf("Unauthorized access attempt to %s", event.Resource),
			Timestamp:   time.Now(),
			Metadata: map[string]interface{}{
				"user_id":  event.UserID,
				"resource": event.Resource,
				"action":   event.Action,
				"ip":       event.IP,
			},
		})
	}
}

func (sm *SecurityMonitor) SendAlert(alert SecurityAlert) {
	sm.alertChannel <- alert
	// Log to console (can be extended to send to Slack, email, PagerDuty, etc.)
	fmt.Printf("[SECURITY ALERT] [%s] %s: %s\n", alert.Severity, alert.EventType, alert.Description)
}
```

### 2. Integration with Audit Logging

```go
// In pkg/middleware/audit.go, add monitoring integration

type MonitoringAuditLogger struct {
	baseLogger middleware.AuditLogger
	monitor    *monitoring.SecurityMonitor
}

func NewMonitoringAuditLogger(baseLogger middleware.AuditLogger, monitor *monitoring.SecurityMonitor) *MonitoringAuditLogger {
	return &MonitoringAuditLogger{
		baseLogger: baseLogger,
		monitor:    monitor,
	}
}

func (l *MonitoringAuditLogger) Log(event middleware.AuditEvent) error {
	// Log to base logger
	if err := l.baseLogger.Log(event); err != nil {
		return err
	}

	// Process for security monitoring
	l.monitor.ProcessEvent(event)

	return nil
}
```

---

## SIEM Integration

### Option 1: Splunk Integration

```go
package siem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ferdiunal/panel.go/pkg/middleware"
)

type SplunkLogger struct {
	endpoint string
	token    string
	client   *http.Client
}

func NewSplunkLogger(endpoint, token string) *SplunkLogger {
	return &SplunkLogger{
		endpoint: endpoint,
		token:    token,
		client:   &http.Client{},
	}
}

func (l *SplunkLogger) Log(event middleware.AuditEvent) error {
	payload := map[string]interface{}{
		"event": event,
		"sourcetype": "panel_security",
		"index": "security",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", l.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Splunk "+l.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("splunk returned status %d", resp.StatusCode)
	}

	return nil
}
```

### Option 2: ELK Stack Integration

```go
package siem

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ferdiunal/panel.go/pkg/middleware"
)

type ElasticsearchLogger struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchLogger(addresses []string, index string) (*ElasticsearchLogger, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticsearchLogger{
		client: client,
		index:  index,
	}, nil
}

func (l *ElasticsearchLogger) Log(event middleware.AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = l.client.Index(
		l.index,
		bytes.NewReader(data),
		l.client.Index.WithContext(context.Background()),
	)

	return err
}
```

---

## Security Dashboards

### Metrics to Track

1. **Authentication Metrics**
   - Failed login attempts per hour
   - Successful logins per hour
   - Account lockouts per day
   - Password reset requests per day

2. **Authorization Metrics**
   - Permission denials per hour
   - Unauthorized access attempts per day
   - Resource access patterns

3. **Rate Limiting Metrics**
   - Rate limit violations per hour
   - Top IPs hitting rate limits
   - API endpoint usage patterns

4. **Security Events**
   - SQL injection attempts
   - CSRF token failures
   - CORS violations
   - Suspicious activity patterns

### Prometheus Metrics Example

```go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	LoginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "panel_login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"}, // "success" or "failure"
	)

	PermissionChecks = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "panel_permission_checks_total",
			Help: "Total number of permission checks",
		},
		[]string{"resource", "action", "result"}, // "granted" or "denied"
	)

	RateLimitViolations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "panel_rate_limit_violations_total",
			Help: "Total number of rate limit violations",
		},
		[]string{"endpoint"},
	)

	AccountLockouts = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "panel_account_lockouts_total",
			Help: "Total number of account lockouts",
		},
	)
)
```

---

## Incident Response Playbooks

### 1. Brute Force Attack Detected

**Trigger**: 5+ failed login attempts from same IP or email within 5 minutes

**Response**:
1. Automatically lock the account (already implemented)
2. Block the IP address temporarily (15 minutes)
3. Alert security team
4. Review audit logs for the IP address
5. Check for other accounts targeted from same IP

### 2. Unauthorized Access Attempt

**Trigger**: Multiple permission denials for same user

**Response**:
1. Log the event with full context
2. Alert security team if threshold exceeded (10+ denials)
3. Review user's recent activity
4. Consider temporary account suspension if suspicious

### 3. SQL Injection Attempt

**Trigger**: Invalid column names in filter parameters

**Response**:
1. Block the request (already implemented)
2. Log the attempt with full request details
3. Alert security team immediately
4. Block the IP address
5. Review all requests from the IP in last 24 hours

---

## Implementation Checklist

- [x] Audit logging middleware implemented
- [x] Security event types defined
- [ ] Security monitor with alerting implemented
- [ ] SIEM integration configured (Splunk/ELK)
- [ ] Prometheus metrics exposed
- [ ] Security dashboards created
- [ ] Incident response playbooks documented
- [ ] Alert notification channels configured (Slack/email/PagerDuty)
- [ ] Security team trained on playbooks
- [ ] Regular security review meetings scheduled

---

## Quick Start

### 1. Enable Security Monitoring

```go
// In main.go or app initialization
monitor := monitoring.NewSecurityMonitor()

// Start alert processor
go func() {
	for alert := range monitor.alertChannel {
		// Send to Slack, email, PagerDuty, etc.
		handleSecurityAlert(alert)
	}
}()

// Use monitoring audit logger
auditLogger := middleware.NewMonitoringAuditLogger(
	&middleware.ConsoleAuditLogger{},
	monitor,
)

app.Use(middleware.AuditMiddleware(auditLogger))
```

### 2. Configure SIEM Integration

```go
// For Splunk
splunkLogger := siem.NewSplunkLogger(
	os.Getenv("SPLUNK_ENDPOINT"),
	os.Getenv("SPLUNK_TOKEN"),
)

// For Elasticsearch
elkLogger, _ := siem.NewElasticsearchLogger(
	[]string{os.Getenv("ELASTICSEARCH_URL")},
	"panel-security",
)

// Use SIEM logger
app.Use(middleware.AuditMiddleware(splunkLogger))
```

### 3. Expose Prometheus Metrics

```go
import (
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Add metrics endpoint
app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
```

---

## Best Practices

1. **Log Everything Security-Related**: Authentication, authorization, data access, configuration changes
2. **Set Appropriate Alert Thresholds**: Balance between noise and missing real threats
3. **Automate Response**: Implement automatic blocking for clear threats (brute force, SQL injection)
4. **Regular Review**: Review security logs and alerts weekly
5. **Test Incident Response**: Run tabletop exercises quarterly
6. **Keep Audit Logs**: Retain for at least 90 days (longer for compliance)
7. **Monitor the Monitors**: Ensure monitoring systems are working correctly

---

## References

- [OWASP Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)
- [NIST Incident Response Guide](https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-61r2.pdf)
- [Splunk Security Essentials](https://www.splunk.com/en_us/software/splunk-security-essentials.html)
