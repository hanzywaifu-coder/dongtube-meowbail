package meowbail

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// CachedMessage menyimpan pesan dalam memori untuk fungsi anti-delete, anti-edit, dan quote retrieval
type CachedMessage struct {
	ID        types.MessageID
	Chat      types.JID
	Sender    types.JID
	Timestamp time.Time
	Message   *waE2E.Message
	IsFromMe  bool
}

// MemoryMessageStore in-memory cache pesan mirip Baileys makeInMemoryStore
type MemoryMessageStore struct {
	messages map[types.MessageID]*CachedMessage
	order    []types.MessageID
	maxSize  int
	mu       sync.RWMutex
}

// NewMemoryMessageStore inisialisasi store pesan in-memory
func NewMemoryMessageStore(maxSize int) *MemoryMessageStore {
	if maxSize <= 0 {
		maxSize = 2500
	}
	return &MemoryMessageStore{
		messages: make(map[types.MessageID]*CachedMessage),
		order:    make([]types.MessageID, 0, maxSize),
		maxSize:  maxSize,
	}
}

// Put menyimpan pesan ke store
func (s *MemoryMessageStore) Put(msg *CachedMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.messages[msg.ID]; !exists {
		if len(s.order) >= s.maxSize {
			oldest := s.order[0]
			s.order = s.order[1:]
			delete(s.messages, oldest)
		}
		s.order = append(s.order, msg.ID)
	}
	s.messages[msg.ID] = msg
}

// Get mengambil pesan berdasarkan ID
func (s *MemoryMessageStore) Get(id types.MessageID) *CachedMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.messages[id]
}

// GetAlbumChildren mengambil semua pesan media anak yang berasosiasi dengan parent Album ID
func (s *MemoryMessageStore) GetAlbumChildren(parentID types.MessageID) []*CachedMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var children []*CachedMessage
	for _, m := range s.messages {
		if m.Message == nil || m.Message.MessageContextInfo == nil || m.Message.MessageContextInfo.MessageAssociation == nil {
			continue
		}
		assoc := m.Message.MessageContextInfo.MessageAssociation
		if assoc.ParentMessageKey != nil && assoc.ParentMessageKey.ID != nil {
			if *assoc.ParentMessageKey.ID == string(parentID) {
				children = append(children, m)
			}
		}
	}
	return children
}

// AttachMessageStore menghubungkan auto-caching event pesan ke whatsmeow client
func (c *Client) AttachMessageStore(store *MemoryMessageStore) {
	c.AddEventHandler(func(rawEvt interface{}) {
		switch evt := rawEvt.(type) {
		case *events.Message:
			if evt.Info.ID == "" || evt.Message == nil {
				return
			}

			// Cek apakah pesan membawa MessageAssociation ke parent album
			if evt.Message.MessageContextInfo != nil && evt.Message.MessageContextInfo.MessageAssociation != nil {
				assoc := evt.Message.MessageContextInfo.MessageAssociation
				if assoc.ParentMessageKey != nil && assoc.ParentMessageKey.ID != nil {
					log.Printf("[store] Stored child message %s associated with parent %s (hasImg=%v)",
						evt.Info.ID, *assoc.ParentMessageKey.ID, evt.Message.ImageMessage != nil)
				}
			}

			store.Put(&CachedMessage{
				ID:        evt.Info.ID,
				Chat:      evt.Info.Chat,
				Sender:    evt.Info.Sender,
				Timestamp: evt.Info.Timestamp,
				Message:   evt.Message,
				IsFromMe:  evt.Info.IsFromMe,
			})
		}
	})
}

// AntiDeleteWatcher fungsi utilitas mendeteksi pesan yang ditarik/dihapus oleh lawan bicara
func (c *Client) OnMessageRevoked(store *MemoryMessageStore, onRevoke func(deletedMsg *CachedMessage, revoker types.JID)) {
	c.AddEventHandler(func(rawEvt interface{}) {
		msgEvt, ok := rawEvt.(*events.Message)
		if !ok || msgEvt.Message == nil || msgEvt.Message.ProtocolMessage == nil {
			return
		}

		protoMsg := msgEvt.Message.ProtocolMessage
		if protoMsg.GetType() == waE2E.ProtocolMessage_REVOKE {
			key := protoMsg.GetKey()
			if key != nil && key.ID != nil {
				targetID := types.MessageID(*key.ID)
				if cached := store.Get(targetID); cached != nil {
					onRevoke(cached, msgEvt.Info.Sender)
				}
			}
		}
	})
}

// MarkChatRead mengirim tanda centang dua biru (Read Receipt / Seen) untuk chat
func (c *Client) MarkChatRead(ctx context.Context, chat types.JID, sender types.JID, msgIDs []types.MessageID) error {
	if len(msgIDs) == 0 {
		return nil
	}
	return c.Client.MarkRead(ctx, msgIDs, time.Now(), chat, sender)
}

// SetChatPresence mengirim status typing atau recording ke lawan bicara
func (c *Client) SetChatPresence(ctx context.Context, chat types.JID, state types.ChatPresence, media types.ChatPresenceMedia) error {
	return c.Client.SendChatPresence(ctx, chat, state, media)
}

// SendAlbumWithStore mengirim kumpulan foto atau video otomatis berurutan dengan jeda aman anti-spam
func (c *Client) SendAlbumBatch(ctx context.Context, chat types.JID, items [][]byte, isVideo bool, caption string) error {
	if len(items) == 0 {
		return fmt.Errorf("album kosong")
	}

	for i, itemData := range items {
		var itemCaption string
		if i == 0 {
			itemCaption = caption
		}

		if isVideo {
			if err := c.SendVideo(ctx, chat, itemData, itemCaption); err != nil {
				return err
			}
		} else {
			if err := c.SendImage(ctx, chat, itemData, itemCaption); err != nil {
				return err
			}
		}

		if i < len(items)-1 {
			time.Sleep(300 * time.Millisecond)
		}
	}
	return nil
}
