package core

import (
	"strings"
	"sync"

	"go.mau.fi/whatsmeow/types"
)

// LIDResolver memelihara mapping dua arah antara JID Nomor Telepon (PN / s.whatsapp.net)
// dan JID Identitas Terhubung (LID / lid) untuk mencegah error dekripsi ("Bad MAC" / Session Mismatch).
type LIDResolver struct {
	mu     sync.RWMutex
	lidToPN map[string]string // "123456789@lid" -> "6283143961588@s.whatsapp.net"
	pnToLID map[string]string // "6283143961588@s.whatsapp.net" -> "123456789@lid"
}

func NewLIDResolver() *LIDResolver {
	return &LIDResolver{
		lidToPN: make(map[string]string),
		pnToLID: make(map[string]string),
	}
}

// RegisterMapping mendaftarkan pasangan JID LID dan Phone Number
func (lr *LIDResolver) RegisterMapping(lid, pn string) {
	if lid == "" || pn == "" {
		return
	}
	lr.mu.Lock()
	defer lr.mu.Unlock()

	lidNorm := strings.ToLower(strings.TrimSpace(lid))
	pnNorm := strings.ToLower(strings.TrimSpace(pn))

	lr.lidToPN[lidNorm] = pnNorm
	lr.pnToLID[pnNorm] = lidNorm
}

// ResolveToPN mengembalikan Phone Number JID jika input berupa LID JID
func (lr *LIDResolver) ResolveToPN(jid types.JID) types.JID {
	if jid.Server != "lid" {
		return jid
	}

	lr.mu.RLock()
	pnStr, exists := lr.lidToPN[jid.String()]
	lr.mu.RUnlock()

	if exists {
		parsed, err := types.ParseJID(pnStr)
		if err == nil {
			return parsed
		}
	}
	return jid
}

// ResolveToLID mengembalikan LID JID jika input berupa Phone Number JID
func (lr *LIDResolver) ResolveToLID(jid types.JID) types.JID {
	if jid.Server == "lid" {
		return jid
	}

	lr.mu.RLock()
	lidStr, exists := lr.pnToLID[jid.String()]
	lr.mu.RUnlock()

	if exists {
		parsed, err := types.ParseJID(lidStr)
		if err == nil {
			return parsed
		}
	}
	return jid
}
