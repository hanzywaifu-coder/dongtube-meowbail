# dongtube-meowbail 🐱

**Go WhatsApp Library** - Gabungan whatsmeow + Baileys dalam satu library Go yang powerful!

## Features

### Core Features (from whatsmeow)
- ✅ WhatsApp Web Multi-Device API
- ✅ End-to-End Encryption
- ✅ Signal Protocol
- ✅ Media Upload/Download
- ✅ Group Management
- ✅ QR Code / Pairing Code Login
- ✅ Auto Reconnect

### Baileys Features (Added)
- ✅ **Custom Pairing Code** - Custom 8-char pairing code (ala `alipclutch-baileys`)
- ✅ **Group Status (SWGC)** - Upload WhatsApp story langsung ke grup (`GroupStatusMessageV2`)
- ✅ **Interactive Menu** - ViewOnce + Document header + NativeFlow button dropdown (`SendInteractiveMenu`)
- ✅ **ButtonsMessage** - Quick Reply, CTA_URL, Phone Number, Copy Text
- ✅ **ListMessage** - Dropdown menus with sections
- ✅ **Newsletter Context** - Forwarded from channel
- ✅ **Interactive Messages** - NativeFlow buttons
- ✅ **Poll Creation** - Create polls
- ✅ **Reactions** - React to messages
- ✅ **Location Sharing** - Send locations
- ✅ **Contact Cards** - Send vCards
- ✅ **Stickers** - Send stickers

### Why dongtube-meowbail?

| Feature | whatsmeow | Baileys | dongtube-meowbail |
|---------|-----------|---------|-------------------|
| Language | Go | TypeScript | **Go** |
| Memory | ~20MB | ~100MB+ | **~20MB** |
| Stability | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | **⭐⭐⭐⭐⭐** |
| Buttons | ❌ | ✅ | **✅** |
| Newsletter | ❌ | ✅ | **✅** |
| API | Verbose | Easy | **Easy** |
| Type Safe | ✅ | ❌ | **✅** |

## Installation

```bash
go get github.com/hanzywaifu-coder/dongtube-meowbail
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    meowbail "github.com/hanzywaifu-coder/dongtube-meowbail"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/types/events"
)

func main() {
    // Initialize whatsmeow store
    container, _ := sqlstore.New(context.Background(), "sqlite3", "file:session.db")
    device, _ := container.GetFirstDevice(context.Background())

    // Create whatsmeow client
    waClient := whatsmeow.NewClient(device, nil)

    // Create dongtube-meowbail client
    client := meowbail.NewClientFromWhatsmeow(waClient)

    // Set newsletter context (optional)
    client.SetNewsletter("120363xxx@newsletter", "My Channel")
    client.SetBusinessOwner("628xxx@s.whatsapp.net")

    // Add event handler
    client.AddEventHandler(func(evt interface{}) {
        msg := meowbail.ParseMessageEvent(evt)
        if msg == nil || msg.IsFromMe {
            return
        }

        // Handle button responses
        if selectedID, isButton := meowbail.HandleButtonResponse(evt); isButton {
            switch selectedID {
            case "ping":
                client.SendText(context.Background(), msg.Chat, "Pong! 🏓")
            case "menu":
                // Send menu
            }
            return
        }

        // Handle text commands
        switch msg.Text {
        case ".menu":
            sendMenu(client, msg)
        case ".ping":
            client.SendText(context.Background(), msg.Chat, "Pong! 🏓")
        }
    })

    // Connect
    client.Connect(context.Background())
}

func sendMenu(client *meowbail.Client, msg *meowbail.MessageEvent) {
    // Send menu with buttons
    client.SendButtons(context.Background(), msg.Chat,
        "╭┈┈⬡「 Menu 」\n┃ 1. Ping\n┃ 2. Help\n╰┈┈┈┈┈┈┈┈⬡",
        []meowbail.Button{
            {Type: meowbail.ButtonQuickReply, ID: "ping", Text: "🏓 Ping"},
            {Type: meowbail.ButtonCTAURL, Text: "👤 Owner", URL: "https://wa.me/628xxx"},
        },
    )
}
```

## API Reference

### Client

```go
// Create client
client := meowbail.NewClientFromWhatsmeow(waClient)
client := meowbail.NewClientFromWhatsmeow(waClient, &meowbail.Config{
    NewsletterJID:   "120363xxx@newsletter",
    NewsletterName:  "My Channel",
    BusinessOwnerJID: "628xxx@s.whatsapp.net",
    AutoReconnect:   true,
})

// Set newsletter context
client.SetNewsletter("120363xxx@newsletter", "My Channel")

// Connect
client.Connect(ctx)
```

### Messages

```go
// Text
client.SendText(ctx, chat, "Hello!")

// Text with newsletter
client.SendTextWithNewsletter(ctx, chat, "Forwarded message")

// Buttons
client.SendButtons(ctx, chat, "Choose:", []meowbail.Button{
    {Type: meowbail.ButtonQuickReply, ID: "btn1", Text: "Button 1"},
    {Type: meowbail.ButtonCTAURL, Text: "Visit", URL: "https://example.com"},
})

// List/Dropdown
client.SendList(ctx, chat, "Title", "Description", "Select", []meowbail.Section{
    {Title: "Section 1", Rows: []meowbail.SectionRow{
        {Title: "Row 1", Description: "Desc", ID: "row1"},
    }},
})

// Media
client.SendImage(ctx, chat, imageData, "Caption")
client.SendVideo(ctx, chat, videoData, "Caption")
client.SendDocument(ctx, chat, docData, "file.pdf", "application/pdf", "Caption")
client.SendAudio(ctx, chat, audioData)
client.SendSticker(ctx, chat, stickerData)

// Other
client.SendLocation(ctx, chat, -6.2088, 106.8456, "Jakarta", "Indonesia")
client.SendContact(ctx, chat, "John", "628xxx")
client.SendReaction(ctx, chat, msgID, "👍")
client.SendPoll(ctx, chat, "Question", []string{"A", "B", "C"}, 1)
```

### Groups

```go
// Get info
info, _ := client.GetGroupInfo(ctx, groupJID)

// Admin operations
client.PromoteParticipant(ctx, groupJID, participantJID)
client.DemoteParticipant(ctx, groupJID, participantJID)
client.RemoveParticipant(ctx, groupJID, participantJID)
client.AddParticipant(ctx, groupJID, participantJID)

// Settings
client.SetGroupName(ctx, groupJID, "New Name")
client.SetGroupDescription(ctx, groupJID, "New Description")

// Check admin
isAdmin, _ := client.IsBotAdmin(ctx, groupJID)
```

### Utilities

```go
// Duration formatting
meowbail.FormatDuration(45 * time.Minute) // "45m 0s"

// Greeting
meowbail.GetGreeting() // "Selamat Pagi"

// Phone formatting
meowbail.FormatPhone("08123456789") // "628123456789"
```

## License

MIT License

## Credits

- [whatsmeow](https://github.com/tulir/whatsmeow) - Go WhatsApp library
- [Baileys](https://github.com/WhiskeySockets/Baileys) - WhatsApp Web API
- [dongtube-bot](https://github.com/hanzywaifu-coder/dongtube-wa) - WhatsApp Bot
