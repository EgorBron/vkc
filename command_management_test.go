package vkc

import (
	"testing"
)

func TestRegister(t *testing.T) {
	commands := &Commands[any]{
		Prefix: PrefixText("!"),
	}

	handler1 := &CommandHandler[any]{
		Pattern:  Text("help"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	handler2 := &CommandHandler[any]{
		Pattern:  Text("info"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	// Регистрируем обработчики
	if _, err := commands.Register(handler1); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if i, err := commands.Register(handler2); i != 1 || err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Проверяем количество обработчиков
	if count := commands.HandlersCount(); count != 2 {
		t.Errorf("Expected 2 handlers, got %d", count)
	}

	// Регистрируем nil, должно вернуть ошибку
	if _, err := commands.Register(nil); err == nil {
		t.Error("Register(nil) should return an error")
	}
}

func TestUnregister(t *testing.T) {
	commands := &Commands[any]{
		Prefix: PrefixText("!"),
	}

	handler1 := &CommandHandler[any]{
		Pattern:  Text("help"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	handler2 := &CommandHandler[any]{
		Pattern:  Text("info"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	commands.Register(handler1)
	commands.Register(handler2)

	// Удаляем первый обработчик
	if err := commands.Unregister(handler1); err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}

	if count := commands.HandlersCount(); count != 1 {
		t.Errorf("Expected 1 handler after unregister, got %d", count)
	}

	// Удаляем несуществующий обработчик должно вернуть ошибку
	if err := commands.Unregister(handler1); err == nil {
		t.Error("Unregister of non-existent handler should return an error")
	}
}

func TestUnregisterAt(t *testing.T) {
	commands := &Commands[any]{
		Prefix: PrefixText("!"),
	}

	handler1 := &CommandHandler[any]{
		Pattern:  Text("help"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	handler2 := &CommandHandler[any]{
		Pattern:  Text("info"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	commands.Register(handler1)
	commands.Register(handler2)

	// Удаляем по индексу
	if err := commands.UnregisterAt(0); err != nil {
		t.Fatalf("UnregisterAt failed: %v", err)
	}

	if count := commands.HandlersCount(); count != 1 {
		t.Errorf("Expected 1 handler, got %d", count)
	}

	// Удаляем с некорректным индексом
	if err := commands.UnregisterAt(10); err == nil {
		t.Error("UnregisterAt with out-of-range index should return an error")
	}
}

func TestReplace(t *testing.T) {
	commands := &Commands[any]{
		Prefix: PrefixText("!"),
	}

	oldHandler := &CommandHandler[any]{
		Pattern:  Text("help"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	newHandler := &CommandHandler[any]{
		Pattern:  Text("help"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	commands.Register(oldHandler)

	// Заменяем обработчик
	if err := commands.Replace(0, newHandler); err != nil {
		t.Fatalf("Replace failed: %v", err)
	}

	// Заменяем с некорректным индексом
	if err := commands.Replace(10, newHandler); err == nil {
		t.Error("Replace with out-of-range index should return an error")
	}

	// Заменяем на nil
	if err := commands.Replace(0, nil); err == nil {
		t.Error("Replace with nil handler should return an error")
	}
}

func TestFindHandlerIndex(t *testing.T) {
	commands := &Commands[any]{
		Prefix: PrefixText("!"),
	}

	handler1 := &CommandHandler[any]{
		Pattern:  Text("help"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	handler2 := &CommandHandler[any]{
		Pattern:  Text("info"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	commands.Register(handler1)
	commands.Register(handler2)

	// Ищем существующий обработчик
	if idx := commands.FindHandlerIndex(handler1); idx != 0 {
		t.Errorf("Expected index 0, got %d", idx)
	}

	if idx := commands.FindHandlerIndex(handler2); idx != 1 {
		t.Errorf("Expected index 1, got %d", idx)
	}

	// Ищем несуществующий обработчик
	nonExistent := &CommandHandler[any]{
		Pattern:  Text("unknown"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	if idx := commands.FindHandlerIndex(nonExistent); idx != -1 {
		t.Errorf("Expected index -1 for non-existent handler, got %d", idx)
	}
}

func TestClearHandlers(t *testing.T) {
	commands := &Commands[any]{
		Prefix: PrefixText("!"),
	}

	handler1 := &CommandHandler[any]{
		Pattern:  Text("help"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	handler2 := &CommandHandler[any]{
		Pattern:  Text("info"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	commands.Register(handler1)
	commands.Register(handler2)

	if count := commands.HandlersCount(); count != 2 {
		t.Errorf("Expected 2 handlers, got %d", count)
	}

	// Очищаем обработчики
	commands.ClearHandlers()

	if count := commands.HandlersCount(); count != 0 {
		t.Errorf("Expected 0 handlers after clear, got %d", count)
	}
}

// Deprecated: TODO: удалить в v2.
func TestBackwardsCompatibility(t *testing.T) {
	// Проверяем что старый способ инициализации через Handlers все еще работает
	handler1 := &CommandHandler[any]{
		Pattern:  Text("help"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	handler2 := &CommandHandler[any]{
		Pattern:  Text("info"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	commands := &Commands[any]{
		Prefix:   PrefixText("!"),
		Handlers: []*CommandHandler[any]{handler1, handler2},
	}

	// Проверяем что обработчики были установлены через поле Handlers
	if count := commands.HandlersCount(); count != 2 {
		t.Errorf("Expected 2 handlers from Handlers field, got %d", count)
	}

	// Проверяем что мы можем удалить обработчик
	if err := commands.Unregister(handler1); err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}

	if count := commands.HandlersCount(); count != 1 {
		t.Errorf("Expected 1 handler after unregister, got %d", count)
	}

	// Проверяем что можем добавить новый
	handler3 := &CommandHandler[any]{
		Pattern:  Text("status"),
		Executor: func(ctx CommandContext[any]) error { return nil },
	}

	if _, err := commands.Register(handler3); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if count := commands.HandlersCount(); count != 2 {
		t.Errorf("Expected 2 handlers after register, got %d", count)
	}
}
