package meowbail

import (
	"context"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// NewsletterReaction mengirim reaksi emoji ke pesan dalam Saluran/Newsletter WhatsApp
func (c *Client) NewsletterReaction(ctx context.Context, channelJID types.JID, serverID types.MessageServerID, emoji string) error {
	return c.Client.NewsletterSendReaction(ctx, channelJID, serverID, emoji, "")
}

// LiveLocation mengirim lokasi real-time dengan update koordinat berkala
func (c *Client) SendLiveLocation(ctx context.Context, chat types.JID, lat, lng float64, accuracy uint32, speed float32, caption string, duration time.Duration) error {
	if duration == 0 {
		duration = 300 * time.Second
	}
	seqNumber := int64(1)
	msg := &waE2E.Message{
		LiveLocationMessage: &waE2E.LiveLocationMessage{
			DegreesLatitude:                proto.Float64(lat),
			DegreesLongitude:               proto.Float64(lng),
			AccuracyInMeters:               proto.Uint32(accuracy),
			SpeedInMps:                     proto.Float32(speed),
			DegreesClockwiseFromMagneticNorth: proto.Uint32(0),
			Caption:                        proto.String(caption),
			SequenceNumber:                 proto.Int64(seqNumber),
			TimeOffset:                     proto.Uint32(uint32(duration.Seconds())),
		},
	}
	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// ForwardMessage meneruskan pesan yang ada ke obrolan lain dengan forwarding context
func (c *Client) ForwardMessage(ctx context.Context, toChat types.JID, rawMsg *waE2E.Message, isMultiForward bool) error {
	if rawMsg == nil {
		return fmt.Errorf("pesan kosong")
	}

	clone := proto.Clone(rawMsg).(*waE2E.Message)
	score := uint32(1)
	if isMultiForward {
		score = 5
	}

	ctxInfo := &waE2E.ContextInfo{
		IsForwarded:     proto.Bool(true),
		ForwardingScore: proto.Uint32(score),
	}

	// Sisipkan context info ke jenis pesan yang sesuai
	switch {
	case clone.ExtendedTextMessage != nil:
		clone.ExtendedTextMessage.ContextInfo = ctxInfo
	case clone.ImageMessage != nil:
		clone.ImageMessage.ContextInfo = ctxInfo
	case clone.VideoMessage != nil:
		clone.VideoMessage.ContextInfo = ctxInfo
	case clone.DocumentMessage != nil:
		clone.DocumentMessage.ContextInfo = ctxInfo
	case clone.AudioMessage != nil:
		clone.AudioMessage.ContextInfo = ctxInfo
	case clone.StickerMessage != nil:
		clone.StickerMessage.ContextInfo = ctxInfo
	case clone.Conversation != nil:
		clone = &waE2E.Message{
			ExtendedTextMessage: &waE2E.ExtendedTextMessage{
				Text:        clone.Conversation,
				ContextInfo: ctxInfo,
			},
		}
	}

	_, err := c.Client.SendMessage(ctx, toChat, clone)
	return err
}

// SendPTTVoiceNote mengirim audio rekaman suara PTT (Push To Talk waveform)
func (c *Client) SendPTTVoiceNote(ctx context.Context, chat types.JID, oggOpusData []byte, waveform []byte) error {
	uploaded, err := c.UploadMedia(ctx, oggOpusData, whatsmeow.MediaAudio)
	if err != nil {
		return fmt.Errorf("upload audio: %w", err)
	}

	ptt := true
	msg := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String("audio/ogg; codecs=opus"),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uploaded.FileLength),
			PTT:           &ptt,
			Waveform:      waveform,
		},
	}

	_, err = c.Client.SendMessage(ctx, chat, msg)
	return err
}

// RequestPhoneNumber mengirim prompt tombol Native Flow request phone number WhatsApp
func (c *Client) RequestPhoneNumber(ctx context.Context, chat types.JID, bodyText string) error {
	btnName := "cta_phone_number"
	btnParams := `{"display_text":"Bagikan Nomor Saya"}`

	msg := &waE2E.Message{
		InteractiveMessage: &waE2E.InteractiveMessage{
			Body: &waE2E.InteractiveMessage_Body{
				Text: proto.String(bodyText),
			},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons: []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
						{
							Name:             proto.String(btnName),
							ButtonParamsJSON: proto.String(btnParams),
						},
					},
				},
			},
		},
	}

	bizNodes := buildBizAdditionalNodes()
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &bizNodes,
	})
	return err
}

// RequestLiveLocationPrompt meminta user mengirim live location secara native flow
func (c *Client) RequestLiveLocationPrompt(ctx context.Context, chat types.JID, bodyText string) error {
	btnName := "send_location"
	btnParams := `{}`

	msg := &waE2E.Message{
		InteractiveMessage: &waE2E.InteractiveMessage{
			Body: &waE2E.InteractiveMessage_Body{
				Text: proto.String(bodyText),
			},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons: []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
						{
							Name:             proto.String(btnName),
							ButtonParamsJSON: proto.String(btnParams),
						},
					},
				},
			},
		},
	}

	bizNodes := buildBizAdditionalNodes()
	_, err := c.Client.SendMessage(ctx, chat, msg, whatsmeow.SendRequestExtra{
		AdditionalNodes: &bizNodes,
	})
	return err
}

// SetDisappearingChat mengubah durasi disappearing timer obrolan (24 Jam, 7 Hari, 90 Hari, atau Mati)
func (c *Client) SetDisappearingChat(ctx context.Context, chat types.JID, timer time.Duration) error {
	msg := &waE2E.Message{
		ProtocolMessage: &waE2E.ProtocolMessage{
			Type:                      waE2E.ProtocolMessage_EPHEMERAL_SETTING.Enum(),
			EphemeralExpiration:       proto.Uint32(uint32(timer.Seconds())),
			EphemeralSettingTimestamp: proto.Int64(time.Now().Unix()),
		},
	}
	_, err := c.Client.SendMessage(ctx, chat, msg)
	return err
}

// RejectCall menolak panggilan masuk WhatsApp otomatis
func (c *Client) RejectCall(ctx context.Context, callFrom types.JID, callID string) error {
	return c.Client.RejectCall(ctx, callFrom, callID)
}

// GetContactProfilePicture mengambil link/buffer foto profil kontak atau grup
func (c *Client) GetContactProfilePicture(ctx context.Context, target types.JID, isCommunity bool) (*types.ProfilePictureInfo, error) {
	var extra *whatsmeow.GetProfilePictureParams
	if isCommunity {
		extra = &whatsmeow.GetProfilePictureParams{
			IsCommunity: true,
		}
	}
	return c.Client.GetProfilePictureInfo(ctx, target, extra)
}

// DownloadContactProfilePicture mengambil dan mengunduh byte gambar langsung dari avatar pengguna atau grup
func (c *Client) DownloadContactProfilePicture(ctx context.Context, target types.JID, isCommunity bool) ([]byte, error) {
	info, err := c.GetContactProfilePicture(ctx, target, isCommunity)
	if err != nil {
		return nil, err
	}
	if info == nil || info.URL == "" {
		return nil, fmt.Errorf("avatar tidak ditemukan")
	}
	return FetchURL(info.URL)
}
