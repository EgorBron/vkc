package vkc

import (
	"fmt"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/api/params"
	"github.com/SevereCloud/vksdk/v3/events"
	"github.com/SevereCloud/vksdk/v3/object"
)

// Контекст команды. Передается в каждый обработчик.
type CommandContext[DEPS any] struct {
	VK         *api.VK
	Message    object.MessagesMessage
	Arguments  []string
	RawEvent   events.MessageNewObject
	Dependency DEPS
}

// Параметры отправки сообщения. Используется в методах Send и SendMessageRaw.
type SendTextParams struct {
	Reply bool
	Fmt   []any
}

// Сокращение для SendTextParams со включенным ответом на сообщение.
var WithReplyParams = &SendTextParams{
	Reply: true,
}

// Сокращение для SendTextParams с подстановкой.
func WithFmtParams(subst ...any) *SendTextParams {
	return &SendTextParams{
		Fmt: subst,
	}
}

// Сокращение для SendTextParams с подстановкой и включенным ответом на сообщение.
func WithFmtAndReplyParams(subst ...any) (ret *SendTextParams) {
	ret = WithFmtParams(subst...)
	ret.Reply = true
	return
}

// Базовый метод отправки сообщения.
// Возвращает ошибки в случаях: ...
func SendMessageRaw(vk *api.VK, msg *object.MessagesMessage, peerID int, text string, sendParams *SendTextParams) error {
	if sendParams == nil {
		sendParams = &SendTextParams{}
	}
	b := params.NewMessagesSendBuilder()
	if sendParams.Reply && msg != nil {
		b.ReplyTo(msg.ID)
	}
	finalText := text
	if len(sendParams.Fmt) > 0 {
		finalText = fmt.Sprintf(finalText, sendParams.Fmt...)
	}
	b.Message(finalText)
	b.RandomID(0)
	b.PeerID(peerID)
	if _, err := vk.MessagesSend(b.Params); err != nil {
		return fmt.Errorf("send error: %v", err)
	}
	return nil
}

// Отправка сообщения с параметрами.
func (ctx CommandContext[DEPS]) Send(text string, sendParams *SendTextParams) error {
	return SendMessageRaw(ctx.VK, &ctx.Message, ctx.Message.PeerID, text, sendParams)
}

// Отправка сообщения с форматированием.
func (ctx CommandContext[DEPS]) SendText(text string, fmts ...any) error {
	return SendMessageRaw(ctx.VK, &ctx.Message, ctx.Message.PeerID, text, WithFmtParams(fmts...))
}

// Отправка ответа на команду с форматированием.
func (ctx CommandContext[DEPS]) Reply(text string, fmts ...any) error {
	return SendMessageRaw(ctx.VK, &ctx.Message, ctx.Message.PeerID, text, WithFmtAndReplyParams(fmts...))
}
