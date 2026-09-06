package meowbail

import (
	"context"

	"go.mau.fi/whatsmeow/proto/waAICommon"
	"go.mau.fi/whatsmeow/proto/waAICommonDeprecated"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// SendAITableResponse mengirim pesan tabel AI bergaya Meta AI / rich response message
func (c *Client) SendAITableResponse(ctx context.Context, chat types.JID, title string, headers []string, rows [][]string) error {
	var tableRows []*waAICommonDeprecated.AIRichResponseTableMetadata_AIRichResponseTableRow

	// Header row
	if len(headers) > 0 {
		tableRows = append(tableRows, &waAICommonDeprecated.AIRichResponseTableMetadata_AIRichResponseTableRow{
			Items:     headers,
			IsHeading: proto.Bool(true),
		})
	}

	// Data rows
	for _, r := range rows {
		tableRows = append(tableRows, &waAICommonDeprecated.AIRichResponseTableMetadata_AIRichResponseTableRow{
			Items:     r,
			IsHeading: proto.Bool(false),
		})
	}

	tableSubMsgType := waAICommonDeprecated.AIRichResponseSubMessageType_AI_RICH_RESPONSE_TABLE
	submessages := []*waAICommonDeprecated.AIRichResponseSubMessage{
		{
			MessageType: &tableSubMsgType,
			TableMetadata: &waAICommonDeprecated.AIRichResponseTableMetadata{
				Title: proto.String(title),
				Rows:  tableRows,
			},
		},
	}

	richMsgType := waAICommonDeprecated.AIRichResponseMessageType_AI_RICH_RESPONSE_TYPE_STANDARD
	botMsg := &waE2E.Message{
		BotForwardedMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				RichResponseMessage: &waE2E.AIRichResponseMessage{
					MessageType: &richMsgType,
					Submessages: submessages,
					ContextInfo: &waE2E.ContextInfo{
						IsForwarded: proto.Bool(true),
						ForwardedAiBotMessageInfo: &waAICommon.ForwardedAIBotMessageInfo{
							BotJID: proto.String("867051314767696@bot"),
						},
					},
				},
			},
		},
	}

	_, err := c.Client.SendMessage(ctx, chat, botMsg)
	return err
}

// SendAICodeResponse mengirim pesan format code block dengan syntax highlight khas Meta AI
func (c *Client) SendAICodeResponse(ctx context.Context, chat types.JID, title, code, language string) error {
	if language == "" {
		language = "go"
	}

	codeSubMsgType := waAICommonDeprecated.AIRichResponseSubMessageType_AI_RICH_RESPONSE_CODE
	defaultHighlight := waAICommonDeprecated.AIRichResponseCodeMetadata_AI_RICH_RESPONSE_CODE_HIGHLIGHT_DEFAULT

	submessages := []*waAICommonDeprecated.AIRichResponseSubMessage{
		{
			MessageType: &codeSubMsgType,
			CodeMetadata: &waAICommonDeprecated.AIRichResponseCodeMetadata{
				CodeLanguage: proto.String(language),
				CodeBlocks: []*waAICommonDeprecated.AIRichResponseCodeMetadata_AIRichResponseCodeBlock{
					{
						HighlightType: &defaultHighlight,
						CodeContent:   proto.String(code),
					},
				},
			},
		},
	}

	richMsgType := waAICommonDeprecated.AIRichResponseMessageType_AI_RICH_RESPONSE_TYPE_STANDARD
	botMsg := &waE2E.Message{
		BotForwardedMessage: &waE2E.FutureProofMessage{
			Message: &waE2E.Message{
				RichResponseMessage: &waE2E.AIRichResponseMessage{
					MessageType: &richMsgType,
					Submessages: submessages,
					ContextInfo: &waE2E.ContextInfo{
						IsForwarded: proto.Bool(true),
						ForwardedAiBotMessageInfo: &waAICommon.ForwardedAIBotMessageInfo{
							BotJID: proto.String("867051314767696@bot"),
						},
					},
				},
			},
		},
	}

	_ = waAICommon.AIRichResponseUnifiedResponse{} // ensure waAICommon import satisfies compiler

	_, err := c.Client.SendMessage(ctx, chat, botMsg)
	return err
}
