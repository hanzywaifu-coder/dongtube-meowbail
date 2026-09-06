# MeowBail

> **Next-Generation High-Performance WhatsApp Protocol Engine & Library**

MeowBail is an independent, production-grade WhatsApp protocol framework engineered for extreme reliability, minimal resource footprint, and enterprise-grade automation. Built to operate as a standalone, modern communication engine, MeowBail provides clean native interfaces in Go and Node.js/TypeScript while resolving low-level protocol inconsistencies, memory bloat, and delivery bottlenecks.

---

## Key Architectural Principles

- **True Standalone Identity:** MeowBail is designed from foundational protocol specifications with dedicated optimizations, isolated state sync, and robust session persistence.
- **Ultra-Low Memory Footprint:** Operates under continuous production load with an average resident set size of ~20 MB in Go and optimized buffers in Node.js.
- **Resilience by Design:** Built-in spiral loop prevention, automatic token refresh, transient exponential backoff, and transparent reconnect strategies.
- **Protocol Safety:** Native human-entropy pacing, proactive chatstate emulation, and compliance-aware payload formulation to minimize account health degradation.
- **Media Pipeline Integrity:** Automated container verification, WebP RIFF header alignment, strict single-pass ZIP serialization, and optimized MIME dispatching.

---

## Feature Comparison Matrix

| Capabilities | Standard Libraries | MeowBail Architecture |
|---|---|---|
| **Multi-Language Support** | Single platform locked | Native Go & Node.js / TypeScript interfaces |
| **Average Memory Footprint** | 120 MB - 300 MB+ | ~20 MB (Go) / Clean GC footprint (Node.js) |
| **Pairing Protocol** | Fixed pseudo-random tokens | Custom 8-character pairing codes (e.g., `DONGTUBE`) |
| **Group Stories (SWGC)** | Undocumented / Manual protobuf | Native `SendGroupStatus` with automatic type inference |
| **Identity Resolution** | Fragile LID/PN synchronization | Built-in self-healing `LIDResolver` |
| **Connection Stability** | Prone to disconnect storms | Built-in spiral loop detector & auto-reconnect |
| **Payload Delivery** | Frequent schema rejections | Validated NativeFlow v3, Carousels, & Album structures |
| **Call Management** | Limited call protocol coverage | Native WhatsApp Call Link generator (Audio & Video) |
| **Sticker Packs** | Corrupted frames in stream pipes | Strict RIFF repair, EXIF injection, & clean ZIP storage |
| **Privacy Configuration** | Incomplete IQ abstractions | Complete Privacy Manager (Read receipts, Online, Last seen) |

---

## Installation

### Go

```bash
go get github.com/hanzywaifu-coder/dongtube-meowbail
```

### Node.js / TypeScript

```bash
npm install github:hanzywaifu-coder/dongtube-meowbail
```

---

## Quick Start (Go)

### 1. Client Initialization & Custom Pairing

```go
package main

import (
	"context"
	"fmt"
	"log"

	meowbail "github.com/hanzywaifu-coder/dongtube-meowbail"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func main() {
	ctx := context.Background()

	container, err := sqlstore.New(ctx, "sqlite3", "file:session.db?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to initialize database container: %v", err)
	}

	device, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatalf("Failed to load device storage: %v", err)
	}

	baseClient := whatsmeow.NewClient(device, nil)
	client := meowbail.NewClientFromWhatsmeow(baseClient)

	// Custom 8-character pairing code request
	code, err := client.PairCustomPhone(
		ctx,
		"6283143961588",
		"DONGTUBE",
		true,
		whatsmeow.PairClientChrome,
		"Chrome (Linux)",
	)
	if err != nil {
		log.Fatalf("Pairing failed: %v", err)
	}

	fmt.Printf("Pairing code generated: %s\n", code)
}
```

### 2. Group Story Publication (SWGC)

Publish media directly to the ephemeral group status container:

```go
err := client.SendGroupStatus(ctx, groupJID, meowbail.GroupStatusMedia{
	Image: imageBytes,
	Text:  "Automated group announcement broadcast.",
})
if err != nil {
	log.Printf("Failed to publish group status: %v", err)
}
```

### 3. WhatsApp Call Link Generation

Generate authentic audio or video conference tokens:

```go
link, err := client.CreateCallLink(ctx, meowbail.CallMediaAudio, nil)
if err != nil {
	log.Printf("Call link generation failed: %v", err)
	return
}

fmt.Printf("Call URL: %s (Token: %s)\n", link.URL, link.Token)
```

### 4. Privacy Configuration

Control account telemetry and visibility parameters:

```go
// Disable blue read-receipt checkmarks
_, err := client.SetReadReceiptsPrivacy(ctx, false)

// Match online status strictly to Last Seen
_, err = client.SetOnlinePrivacy(ctx, true)

// Restrict incoming calls to known contacts
_, err = client.SetCallAddPrivacy(ctx, false)
```

---

## Quick Start (Node.js / TypeScript)

```typescript
import { makeWASocket } from 'dongtube-meowbail/nodejs'

async function start() {
	const sock = makeWASocket({
		phoneNumber: '6283143961588',
		customPairingCode: 'DONGTUBE',
		browser: ['MeowBail Engine', 'Chrome', '1.0.0'],
		antiBan: {
			enabled: true,
			humanEntropy: true,
		},
	})

	sock.ev.on('messages.upsert', async ({ messages }) => {
		for (const msg of messages) {
			if (!msg.message || msg.key.fromMe) continue

			const jid = msg.key.remoteJid!
			await sock.sendMessage(jid, {
				text: 'Message received and processed via MeowBail.',
			})
		}
	})
}

start().catch(console.error)
```

---

## Core Engine Modules

- `PrivacyManager` (`privacy_manager.go`): Unified access layer for all WhatsApp privacy namespaces.
- `CallLink` (`call_link.go`): Binary stanza query generator for official WhatsApp call links.
- `StickerPackPipeline` (`sticker_pack.go`, `sticker_pack_zip.go`): Strict compliance ZIP compiler and auto-repairing RIFF WebP muxer.
- `ChatActionsSync` (`chat_actions_sync.go`): App state synchronizer for stars, chat deletions, and pushname mutations.
- `GroupAnalytics` (`group_analytics.go`): High-throughput participant scoring and membership change tracking.
- `AntiBan` (`antiban_client.go`): Traffic shaping, jitter calculation, and behavioral entropy emulator.

---

## License

MeowBail is released under the MIT License. Developed and maintained by the MeowBail Engineering Team.
