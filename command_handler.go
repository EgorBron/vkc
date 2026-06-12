package vkc

import (
	"txts.su/vkc/filters"
	vkfilters "txts.su/vkc/filters/vk_filters"
)

// Функция обработчика команды. Получает контекст и возвращает ошибку или nil.
type HandlerFunc[DEPS any] func(ctx CommandContext[DEPS]) error

// Проверка доступа к команде. Используется для проверки прав пользователя.
//
// Пример использования:
//
//	var CheckAdmin = HandlerAccessCheck[any]{
//		Checker: func(handler *CommandHandler[any], ctx CommandContext[any]) bool {
//			return ctx.SenderID == 1
//		},
//	}
//	// позднее в обработчике:
//	&CommandHandler[any]{
//		Pattern:     Text("admincmd"),
//		AccessCheck: &CheckAdmin,
//		Executor:    func(ctx CommandContext[any]) error { /* ... */ },
//	}
//
// Если для пользователя, вызвавшего команду, не проходит проверка, то вызывается обработчик OnNoPermissions (если он задан в объекте команд).
type HandlerAccessCheck[DEPS any] struct {
	Checker func(handler *CommandHandler[DEPS], ctx CommandContext[DEPS]) bool
}

// Обработчик команды.
// Содержит необходимый минимум для большинства ботов: шаблон, по которому вызывается команда, объект помощи, проверку на доступ и функцию-исполнитель.
//
// Пример обработчика:
//
//	/* в качестве DepsType указывается структура или интерфейс
//	из дженерика в объекте команд;
//	впоследствии в контекст обработчика будет передано значение этого типа в поле Dependency */
//	var HandleSomeCommand = CommandHandler[DepsType]{
//		Filter: filter.PrefixF("some") // шаблон "только строка `some`"
//		Help: CommandHelp{ /* помощь по команде */ },
//		AccessCheck: &HandlerAccessCheck[DepsType]{ /* проверка доступа */ },
//		Executor: func(ctx CommandContext[DepsType]) error { /* логика команды */ },
//	}
type CommandHandler[DEPS any] struct {
	// Фильтр, совпадение с которым запустит команду
	Filter filters.Filter
	// Deprecated: будет удалено в v2. Используйте [CommandHandler.Filter].
	Pattern CommandPattern
	// Помощь по команде
	Help CommandHelp
	// Проверка на доступ к команде
	AccessCheck *HandlerAccessCheck[DEPS]
	// Исполнитель команды
	Executor HandlerFunc[DEPS]
}

// Метод для проверки совпадения контекста с фильтром команды.
func (handler *CommandHandler[T]) IsNotFiltered(ctx CommandContext[T], prefixRemainder string) (result bool, sideValues map[string]any) {
	return handler.Filter.Eval(map[filters.FieldDescriptor]any{
		vkfilters.MessageField: ctx.Message,
		vkfilters.VkApiField:   ctx.VK,
		filters.TextField:      prefixRemainder,
	})
}

// Метод для проверки доступности команды для пользователя.
func (handler *CommandHandler[any]) IsAccessAvailable(ctx CommandContext[any]) bool {
	return handler.AccessCheck == nil || handler.AccessCheck.Checker(handler, ctx)
}
