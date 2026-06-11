package vkc

import "fmt"

// Добавляет обработчик в конец списка.
//
// Возвращает индекс добавленного обработчика.
//
// Если обработчик был nil, возвращает ошибку.
//
// `Предупреждение!` Метод не потокобезопасный. Пользуйтесь внешними примитивами синхронизации. Будет исправлено в версии v2.
func (commands *Commands[DEPS]) Register(handler *CommandHandler[DEPS]) (int, error) {
	commands.cmdMutex.Lock()
	defer commands.cmdMutex.Unlock()

	if handler == nil {
		return -1, fmt.Errorf("nil handler")
	}
	commands.Handlers = append(commands.Handlers, handler)
	return len(commands.Handlers) - 1, nil
}

// Удаляет первое вхождение обработчика.
//
// Если обработчик не найден, возвращает ошибку.
// Также возвращает ошибку, если передан nil вместо обработчика.
//
// `Предупреждение!` Метод не потокобезопасный. Пользуйтесь внешними примитивами синхронизации. Будет исправлено в версии v2.
func (commands *Commands[DEPS]) Unregister(handler *CommandHandler[DEPS]) error {
	commands.cmdMutex.Lock()
	defer commands.cmdMutex.Unlock()

	if handler == nil {
		return fmt.Errorf("nil handler")
	}
	for i, h := range commands.Handlers {
		if h == handler {
			commands.Handlers = append(commands.Handlers[:i], commands.Handlers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("handler not found")
}

// Удаляет обработчик по индексу.
//
// Если индекс выходит за пределы внутреннего среза, возвращает ошибку.
//
// `Предупреждение!` Метод не потокобезопасный. Пользуйтесь внешними примитивами синхронизации. Будет исправлено в версии v2.
func (commands *Commands[DEPS]) UnregisterAt(index int) error {
	commands.cmdMutex.Lock()
	defer commands.cmdMutex.Unlock()

	if index < 0 || index >= len(commands.Handlers) {
		return fmt.Errorf("index out of range")
	}
	commands.Handlers = append(commands.Handlers[:index], commands.Handlers[index+1:]...)
	return nil
}

// Заменяет обработчик по индексу.
//
// Если индекс выходит за пределы внутреннего среза, возвращает ошибку.
// Также возвращает ошибку, если передан nil вместо обработчика.
//
// `Предупреждение!` Метод не потокобезопасный. Пользуйтесь внешними примитивами синхронизации. Будет исправлено в версии v2.
func (commands *Commands[DEPS]) Replace(index int, handler *CommandHandler[DEPS]) error {
	commands.cmdMutex.Lock()
	defer commands.cmdMutex.Unlock()

	if handler == nil {
		return fmt.Errorf("nil handler")
	}
	if index < 0 || index >= len(commands.Handlers) {
		return fmt.Errorf("index out of range")
	}
	commands.Handlers[index] = handler
	return nil
}

// Ищет обработчик во внутреннем срезе и возвращает его индекс.
// Если обработчик не найден, возвращает -1.
//
// `Предупреждение!` Метод не потокобезопасный. Пользуйтесь внешними примитивами синхронизации. Будет исправлено в версии v2.
func (commands *Commands[DEPS]) FindHandlerIndex(handler *CommandHandler[DEPS]) int {
	commands.cmdMutex.RLock()
	defer commands.cmdMutex.RUnlock()

	for i, h := range commands.Handlers {
		if h == handler {
			return i
		}
	}
	return -1
}

// Возвращает количество зарегистрированных обработчиков.
//
// `Предупреждение!` Метод не потокобезопасный. Пользуйтесь внешними примитивами синхронизации. Будет исправлено в версии v2.
func (commands *Commands[DEPS]) HandlersCount() int {
	commands.cmdMutex.RLock()
	defer commands.cmdMutex.RUnlock()

	return len(commands.Handlers)
}

// Возвращает копию среза обработчиков.
//
// Изменение возвращённого среза не повлияет на внутреннее состояние Commands.
//
// `Предупреждение!` Метод не потокобезопасный. Пользуйтесь внешними примитивами синхронизации. Будет исправлено в версии v2.
func (commands *Commands[DEPS]) GetHandlers() []*CommandHandler[DEPS] {
	commands.cmdMutex.RLock()
	defer commands.cmdMutex.RUnlock()

	out := make([]*CommandHandler[DEPS], len(commands.Handlers))
	copy(out, commands.Handlers)
	return out
}

// Устанавливает список обработчиков в пустой срез.
//
// `Предупреждение!` Метод не потокобезопасный. Пользуйтесь внешними примитивами синхронизации. Будет исправлено в версии v2.
func (commands *Commands[DEPS]) ClearHandlers() {
	commands.cmdMutex.Lock()
	defer commands.cmdMutex.Unlock()

	commands.Handlers = make([]*CommandHandler[DEPS], 0)
}
