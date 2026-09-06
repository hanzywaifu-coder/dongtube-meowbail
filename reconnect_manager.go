package meowbail

import (
	"context"
	"math"
	"sync"
	"time"
)

// AutoReconnectManager mengelola auto-reconnect cerdas dengan Exponential Backoff + Jitter
type AutoReconnectManager struct {
	client       *Client
	mu           sync.Mutex
	running      bool
	stopChan     chan struct{}
	retryCount   int
	maxRetries   int
	initialDelay time.Duration
	maxDelay     time.Duration
}

// NewAutoReconnectManager membuat instance baru AutoReconnectManager
func NewAutoReconnectManager(client *Client) *AutoReconnectManager {
	return &AutoReconnectManager{
		client:       client,
		stopChan:     make(chan struct{}),
		maxRetries:   50,
		initialDelay: 2 * time.Second,
		maxDelay:     60 * time.Second,
	}
}

// Start mengaktifkan supervisor loop yang memonitor koneksi WhatsApp secara real-time
func (m *AutoReconnectManager) Start(ctx context.Context) {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.mu.Unlock()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-m.stopChan:
				return
			case <-ticker.C:
				if m.client == nil || m.client.Client == nil {
					continue
				}

				// Hanya reconnect jika perangkat sudah terdaftar/logged in tapi koneksi terputus
				if m.client.IsLoggedIn() && !m.client.IsConnected() {
					m.attemptReconnect(ctx)
				} else if m.client.IsConnected() {
					m.mu.Lock()
					m.retryCount = 0
					m.mu.Unlock()
				}
			}
		}
	}()
}

func (m *AutoReconnectManager) attemptReconnect(ctx context.Context) {
	m.mu.Lock()
	if m.retryCount >= m.maxRetries {
		m.mu.Unlock()
		return
	}
	m.retryCount++
	attempt := m.retryCount

	// Exponential backoff: initialDelay * 2^(attempt-1)
	backoff := float64(m.initialDelay) * math.Pow(1.8, float64(attempt-1))
	delay := time.Duration(backoff)
	if delay > m.maxDelay {
		delay = m.maxDelay
	}
	m.mu.Unlock()

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return
	case <-m.stopChan:
		return
	}

	if !m.client.IsConnected() {
		_ = m.client.Connect(ctx)
	}
}

// Stop mematikan loop auto-reconnect
func (m *AutoReconnectManager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.running {
		return
	}
	m.running = false
	close(m.stopChan)
}

// EnableAutoReconnect mendaftarkan pengawas auto-reconnect pada Client MeowBail
func (c *Client) EnableAutoReconnect(ctx context.Context) *AutoReconnectManager {
	mgr := NewAutoReconnectManager(c)
	mgr.Start(ctx)
	return mgr
}

// SafeConnect melakukan koneksi dan otomatis mengaktifkan smart auto-reconnect
func (c *Client) SafeConnect(ctx context.Context) error {
	err := c.Connect(ctx)
	if c.config != nil && c.config.AutoReconnect {
		c.EnableAutoReconnect(ctx)
	}
	return err
}
