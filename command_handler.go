package vkc

import (
	"log"
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
//		Pattern: Text("some"), // шаблон "только строка `some`"
//		Help: CommandHelp{ /* помощь по команде */ },
//		AccessCheck: &HandlerAccessCheck[DepsType]{ /* проверка доступа */ },
//		Executor: func(ctx CommandContext[DepsType]) error { /* логика команды */ },
//	}
type CommandHandler[DEPS any] struct {
	Pattern     CommandPattern
	Help        CommandHelp
	AccessCheck *HandlerAccessCheck[DEPS]
	Executor    HandlerFunc[DEPS]
}

// Метод для проверки доступности команды для пользователя.
func (handler *CommandHandler[any]) IsAccessAvailable(ctx CommandContext[any]) bool {
	log.Println(handler.AccessCheck == nil)
	if handler.AccessCheck != nil {
		log.Println(!handler.AccessCheck.Checker(handler, ctx))
	}
	return handler.AccessCheck == nil || handler.AccessCheck.Checker(handler, ctx)
}
