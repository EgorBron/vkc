// Пакет vkc предоставляет простую базу для обработки команд в ботах социальной сети ВКонтакте (VK).
// Работает с ботами сообществ в группах и ЛС сообщества. Поддерживается только Long Poll API.
//
// Построен на базе VK SDK (github.com/SevereCloud/vksdk), перехватывая новые сообщения и обрабатывая их в соответствии с зарегистрированными шаблонами команд.
//
// Под командами в модуле понимаются сообщения, построенные по шаблону "префикс название_команды аргументы".

package vkc

import (
	"context"
	"fmt"
	"strings"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/events"
	"github.com/SevereCloud/vksdk/v3/longpoll-bot"
)

// Структура для управления командами. Содержит обработчики, префикс и зависимости.
// Структуру можно типизировать дженериком для передачи зависимостей в обработчики команд, или передать any в дженерик, если это не требуется.
// Для обработчиков команд опеределено поле Handlers, принимающее массив указателей на обработчики. Рекомендуется получать адрес обработчика прямо при передаче в массив, как здесь:
//
//		var SomeCommandHandler = CommandHandler[any]{ ... }
//	 // позднее в коде
//		handlers := []*CommandHandler[any]{
//			&SomeCommandHandler,
//		}
//
// Причин для соблюдения именно этого стиля нет. Возможно, что это плохое решение.
//
// Также в структуре есть поля для колбеков на события: OnMessage, OnEmptyPrefix, OnUnknownCommand, OnNoPermissions, OnCommandError. Их передача необязательна, однако, если указать эти обработчики, то они будут вызваны при соответствующих событиях.
type Commands[DEPS any] struct {
	Prefix       PrefixMatcher
	Dependencies DEPS
	Handlers     []*CommandHandler[DEPS]

	OnMessage        *func(vk *api.VK, obj events.MessageNewObject)
	OnEmptyPrefix    *HandlerFunc[DEPS]
	OnUnknownCommand *HandlerFunc[DEPS]
	OnNoPermissions  *HandlerFunc[DEPS]
	OnCommandError   *func(ctx CommandContext[DEPS], err error)
}

// Поиск команды по строке. Возвращает найденный обработчик и остаток строки (т.е. без названия команды).
//
// Примеры срабатывания функции:
//
//	FindCommand("help", handlers), есть обработчик с шаблоном Text("help") -> (обработчик для Text("help"), остаток - "")
//	FindCommand("help me", handlers), есть обработчик с шаблоном Text("help") -> (обработчик для Text("help"), остаток - "me")
//	FindCommand("user list admin", handlers), есть обработчик с шаблоном Text("user list") -> (обработчик для Text("user list"), остаток - "admin")
//	FindCommand("unknown command", handlers), нет обработчика, подходящего под начало строки -> (nil, "")
//
// и так далее
func FindCommand[DEPS any](rawCmd string, commands []*CommandHandler[DEPS]) (*CommandHandler[DEPS], string) {
	for _, handler := range commands {
		if handler == nil {
			continue
		}

		words := strings.Fields(rawCmd)

		for i := len(words); i > 0; i-- {
			candidate := strings.Join(words[:i], " ")

			if handler.Pattern(candidate) {
				return handler, strings.Join(words[i:], " ")
			}
		}
	}

	return nil, ""
}

// Подключение обработчика команд к LongPoll VK SDK.
func (commands Commands[any]) AttachToLongPoll(vk *api.VK, lp *longpoll.LongPoll) error {
	if lp == nil {
		return fmt.Errorf("LongPoll was nil")
	}
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		text := strings.TrimSpace(obj.Message.Text)
		if text == "" {
			return
		}

		if commands.OnMessage != nil {
			(*commands.OnMessage)(vk, obj)
		}

		if commands.Prefix == nil {
			return
		}

		matched, rawCmd := commands.Prefix(text)
		if !matched {
			return
		}

		cmdCtx := CommandContext[any]{
			VK:         vk,
			Message:    obj.Message,
			Arguments:  []string{},
			RawEvent:   obj,
			Dependency: commands.Dependencies,
		}

		if rawCmd == "" {
			if commands.OnEmptyPrefix != nil {
				(*commands.OnEmptyPrefix)(cmdCtx)
			}
			return
		}

		handler, remaining := FindCommand(rawCmd, commands.Handlers)
		if handler == nil {
			if commands.OnUnknownCommand != nil {
				(*commands.OnUnknownCommand)(cmdCtx)
			}
			return
		}

		cmdCtx.Arguments = SplitArgs(remaining)

		if !handler.IsAccessAvailable(cmdCtx) {
			if commands.OnNoPermissions != nil {
				(*commands.OnNoPermissions)(cmdCtx)
			}
			return
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					Stacktrace(r)
				}
			}()
			err := handler.Executor(cmdCtx)
			if err != nil {
				if commands.OnCommandError != nil {
					(*commands.OnCommandError)(cmdCtx, err)
				} else {
					Stacktrace(err)
				}
			}
		}()
	})

	return nil
}
