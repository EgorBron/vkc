// Пакет vkc предоставляет простую базу для обработки команд в ботах социальной сети ВКонтакте (VK).
// Работает с ботами сообществ в группах и ЛС сообщества. Поддерживается только Long Poll API.
//
// Построен на базе VK SDK (github.com/SevereCloud/vksdk).
//
// Под командами в модуле понимаются сообщения, построенные по шаблону "префикс название_команды аргументы".

package vkc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/events"
	"github.com/SevereCloud/vksdk/v3/longpoll-bot"
)

func logDeprecationWarning(feature string) {
	log.Printf("обработчик %s скоро перестанет вызываться. См.: https://github.com/EgorBron/vkc/issues/1", feature)
}

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
	Prefix PrefixMatcher
	// Структура для передачи зависимостей в обработчики команд. Если зависимости не требуются, можно указать any в дженерике.
	//
	// Deprecated: Начиная с v2 будет удалено. Рекомендуется перейти на [context.Context] (см. https://github.com/EgorBron/vkc/issues/2 для просмотра обсуждения).
	Dependencies DEPS
	Handlers     []*CommandHandler[DEPS]

	// Deprecated: Начиная с v2 будет удалено. Рекомендуется переход на вызов [ProcessCommands].
	OnMessage *func(vk *api.VK, obj events.MessageNewObject)
	// Deprecated: Начиная с v2 будет удалено. Рекомендуется переход на вызов [ProcessCommands].
	OnEmptyPrefix *HandlerFunc[DEPS]
	// Deprecated: Начиная с v2 будет удалено. Рекомендуется переход на вызов [ProcessCommands].
	OnUnknownCommand *HandlerFunc[DEPS]
	// Deprecated: Начиная с v2 будет удалено. Рекомендуется переход на вызов [ProcessCommands].
	OnNoPermissions *HandlerFunc[DEPS]
	// Deprecated: Начиная с v2 будет удалено. Рекомендуется переход на вызов [ProcessCommands].
	OnCommandError *func(ctx CommandContext[DEPS], err error)
}

// Поиск команды по строке. Возвращает найденный обработчик и остаток строки (т.е. без названия команды).
//
// Поскольку срезы в Go являются упорядоченными, то при поиске приоритет будут иметь обработчики, расположенные ближе к началу среза.
//
// Примеры срабатывания функции:
//
//		FindCommand("help", handlers), есть обработчик с шаблоном Text("help") -> (обработчик для Text("help"), остаток - "")
//		FindCommand("help me", handlers), есть обработчик с шаблоном Text("help") -> (обработчик для Text("help"), остаток - "me")
//		FindCommand("user list admin", handlers), есть обработчик с шаблоном Text("user list") -> (обработчик для Text("user list"), остаток - "admin")
//		FindCommand("unknown command", handlers), нет обработчика, подходящего под начало строки -> (nil, "")
//	 FindCommand("tag me", handlers), есть обработчики Text("tag") и Text("tag me") -> (обработчик для Text("tag"), остаток - "me")
//
// и так далее.
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

// Поиск и выполнение обработчиков команд в сообщении.
//
// Метод следует вызывать из обработчика события [github.com/SevereCloud/vksdk/v3/events.FuncList.MessageNew].
//
// Процесс обработки команды включает следующие шаги:
//
//  1. Проверка наличия текста в сообщении. Если текст отсутствует, возвращается ошибка [ErrEmptyMessage].
//  2. (устарело) Вызов колбека [Commands.OnMessage] в горутине, если он указан, даже если в сообщении нет команды.
//  3. Проверка наличия префикса в начале текста с помощью функции [Commands.Prefix]. Если префикс не найден, возвращается ошибка [ErrNoPrefix].
//  4. Если после удаления префикса не остается текста, вызывается колбек [Commands.OnEmptyPrefix] и возвращается ошибка [ErrEmptyPrefix].
//  5. Поиск команды среди зарегистрированных обработчиков с помощью функции [FindCommand]. Если команда не найдена, вызывается колбек [Commands.OnUnknownCommand] и возвращается ошибка [ErrCommandNotFound].
//  6. Проверка прав доступа к команде с помощью метода [Commands.IsAccessAvailable] обработчика команды. Если доступ запрещен, вызывается колбек [Commands.OnNoPermissions] и возвращается ошибка [ErrNoPermissions].
//  7. Выполнение обработчика команды. Если во время выполнения возникает паника, она перехватывается и логируется с помощью функции [Stacktrace]. Если сам обработчик возвращает ошибку, вызывается колбек [Commands.OnCommandError] с этой ошибкой, и она же возвращается из метода.
//
// Все колбеки на события имеют несколько особенностей:
//   - они вызываются только если были установлены при создании структуры;
//   - они выполняются в отдельных горутинах;
//   - все они устарели и будут удалены в v2. Рекомендуется вместо этого обрабатывать ошибки метода ProcessCommands напрямую.
func (commands Commands[any]) ProcessCommands(ctx context.Context, vk *api.VK, msg events.MessageNewObject) error {
	text := strings.TrimSpace(msg.Message.Text)
	if text == "" {
		return ErrEmptyMessage
	}

	if commands.OnMessage != nil {
		logDeprecationWarning("OnMessage")
		go (*commands.OnMessage)(vk, msg)
	}

	if commands.Prefix == nil {
		return ErrNoPrefix
	}

	matched, rawCmd := commands.Prefix(text)
	if !matched {
		return ErrNoPrefix
	}

	cmdCtx := CommandContext[any]{
		VK:         vk,
		Message:    msg.Message,
		Arguments:  []string{},
		RawEvent:   msg,
		Dependency: commands.Dependencies,
	}

	if rawCmd == "" {
		if commands.OnEmptyPrefix != nil {
			logDeprecationWarning("OnEmptyPrefix")
			go (*commands.OnEmptyPrefix)(cmdCtx)
		}
		return ErrEmptyPrefix
	}

	handler, remaining := FindCommand(rawCmd, commands.Handlers)
	if handler == nil {
		if commands.OnUnknownCommand != nil {
			logDeprecationWarning("OnUnknownCommand")
			go (*commands.OnUnknownCommand)(cmdCtx)
		}
		return ErrCommandNotFound
	}

	cmdCtx.Arguments = SplitArgs(remaining)

	if !handler.IsAccessAvailable(cmdCtx) {
		if commands.OnNoPermissions != nil {
			logDeprecationWarning("OnNoPermissions")
			go (*commands.OnNoPermissions)(cmdCtx)
		}
		return ErrNoPermissions
	}

	defer func() {
		if r := recover(); r != nil {
			Stacktrace(r)
		}
	}()

	err := handler.Executor(cmdCtx)
	if err != nil {
		if commands.OnCommandError != nil {
			logDeprecationWarning("OnCommandError")
			go (*commands.OnCommandError)(cmdCtx, err)
		}
	}

	return err
}

// Подключение обработчика команд к LongPoll VK SDK.
//
// Deprecated: Будет удалено в v2. Рекомендуется вызывать [Commands.ProcessCommands] напрямую из обработчика сообщений, вместо использования этого метода.
func (commands Commands[any]) AttachToLongPoll(vk *api.VK, lp *longpoll.LongPoll) error {
	if lp == nil {
		return fmt.Errorf("LongPoll was nil")
	}
	lp.MessageNew(func(ctx context.Context, msg events.MessageNewObject) {
		commands.ProcessCommands(ctx, vk, msg)
	})

	return nil
}
