# dongtube-meowbail 🐱⚡

> **The Ultimate WhatsApp Library** — Unified Engine combining the raw speed, concurrency, and memory efficiency of **whatsmeow (Golang)** with the rich modern features (Interactive Buttons, SWGC Status Group, Custom Pairing Code) of **Baileys (Node.js)**.

Dibuat dari nol untuk menjadi pustaka WhatsApp nomor 1, independen, modern, dan multi-bahasa (**Golang & Node.js**).

---

## 🌟 Mengapa dongtube-meowbail?

| Fitur / Parameter | whatsmeow (Go) | Baileys (Node.js) | **dongtube-meowbail** |
|---|---|---|---|
| **Bahasa yang Didukung** | Golang saja | Node.js saja | **Golang & Node.js / TS** |
| **Konsumsi Memori (RAM)** | ~20 MB | ~120 MB+ | **~20 MB (Go) / Optimal (Node)** |
| **Custom Pairing Code (8 Chars)** | ❌ Terbatas auto | ✅ Ya (`alipclutch`) | ✅ **Ya (Kustom 8 Karakter)** |
| **Status Grup / Story Grup (SWGC)** | ❌ Manual proto | ✅ Ya (`groupStatusMessageV2`) | ✅ **Bawaan & Otomatis** |
| **Carousel & Slider Cards** | ❌ Manual proto | ✅ Ya | ✅ **SendCarousel Helper** |
| **Multi-Media Album** | ❌ Manual proto | ✅ Ya | ✅ **SendAlbum Helper** |
| **Menu Interactive / Dropdown** | ❌ Rumit | ✅ Ya (ViewOnce Document) | ✅ **Built-in Helper** |
| **Newsletter / Forward Saluran** | ⚠️ Parsial | ✅ Ya | ✅ **Full Forwarding Context** |

---

## 🚀 Penggunaan di Golang

### 1. Instalasi
```bash
go get github.com/hanzywaifu-coder/dongtube-meowbail
```

### 2. Custom Pairing Code
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
    container, _ := sqlstore.New(context.Background(), "sqlite3", "file:session.db")
    device, _ := container.GetFirstDevice(context.Background())
    waClient := whatsmeow.NewClient(device, nil)
    client := meowbail.NewClientFromWhatsmeow(waClient)

    // Request custom pairing code 8 karakter (contoh: DONG-TUBE)
    code, err := client.PairCustomPhone(context.Background(), "6283143961588", "DONGTUBE", true, whatsmeow.PairClientChrome, "Chrome (Linux)")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Pairing Code:", code)
}
```

### 3. Upload Status Grup (SWGC)
```go
err := client.SendGroupStatus(ctx, groupJID, meowbail.GroupStatusMedia{
    Image: imageBytes,
    Text:  "Status terbaru dari Dongtube Bot!",
})
```

### 4. Interactive Menu (Document Header + NativeFlow Buttons)
```go
sections := []meowbail.Section{
    {
        Title: "✨ MENU UTAMA",
        Rows: []meowbail.SectionRow{
            {Title: "📑 SEMUA FITUR", Description: "List perintah lengkap", ID: ".allmenu"},
            {Title: "👨💻 OWNER", Description: "Kontak pengembang", ID: ".owner"},
        },
    },
}

err := client.SendInteractiveMenu(ctx, chatJID, docBytes, thumbBytes, "Dongtube Bot 2026", sections, "Info Saluran", "https://whatsapp.com/channel/...", nil)
```

---

## 🟢 Penggunaan di Node.js / TypeScript

### 1. Struktur Import
```javascript
const { makeWASocket } = require('dongtube-meowbail/nodejs');

const sock = makeWASocket({
    phoneNumber: '6283143961588',
    newsletterJid: '120363xxx@newsletter',
    newsletterName: 'Dongtube Channel'
});

// 1. Custom Pairing Code
const code = await sock.requestPairingCode('6283143961588', 'DONGTUBE');
console.log('Pairing Code:', code);

// 2. Kirim Status Grup (SWGC)
await sock.sendGroupStatus('120363xxx@g.us', {
    image: imageBuffer,
    caption: 'Status grup berhasil diunggah!'
});

// 3. Menu Interactive
await sock.sendInteractiveMenu(m.chat, {
    body: 'Silakan pilih menu di bawah ini:',
    footer: 'Dongtube WhatsApp Bot',
    thumbnail: thumbBuffer,
    sections: [
        {
            title: 'KATEGORI FITUR',
            rows: [
                { title: '🎨 STICKERS', id: '.menu-stickers', description: 'Buat stiker' },
                { title: '📥 DOWNLOADER', id: '.menu-download', description: 'Download media' }
            ]
        }
    ],
    ctaText: 'Official Channel',
    ctaUrl: 'https://whatsapp.com/channel/0029Vb91qeW17Emm4TVqu53KJ'
});
```

---

## 📜 Lisensi & Kontribusi
- Dilisensikan di bawah **MIT License**.
- Dibuat secara mandiri menggabungkan standar protokol WhatsApp modern.
