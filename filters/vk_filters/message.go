// Пакет vkfilters предоставляет фильтры, специфичные для работы с объектами VK API.
package vkfilters

import (
	"github.com/SevereCloud/vksdk/v3/object"
	"txts.su/vkc/filters"
)

const (
	VkApiField   filters.FieldDescriptor = "vkapi"
	MessageField filters.FieldDescriptor = "message"
)

// Создает фильтр, проверяющий текст сообщения на точное совпадение.
func MessageText(match string) filters.Filter {
	return filters.NewFilter(MessageField, func(value object.MessagesMessage) (bool, map[string]any) {
		return value.Text == match, nil
	})
}

// Создает фильтр, проверяющий ID беседы (peer ID) у сообщения.
func MessagePeerID(match int) filters.Filter {
	return filters.NewFilter(MessageField, func(value object.MessagesMessage) (bool, map[string]any) {
		return value.PeerID == match, nil
	})
}
