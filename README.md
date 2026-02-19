# VKC (VK Commands Go)

Небольшой модуль для обработки команд в ботах VK.

Работает на базе [VK SDK](https://github.com/SevereCloud/vksdk/).

## Установка

```sh
go get txts.su/vkc@latest
```

## Примеры

```go
package main

import (
	"log"
	"strconv"

	"txts.su/vkc"
	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/longpoll-bot"
)

// Структура для хранения зависимостей, которые прокидываются в обработчики команд из основной функции
type MyDependencies struct {
	MyStore     map[string]any
	AllHandlers *[]*vkc.CommandHandler[MyDependencies]
}

// Простая проверка на права. Понятно, для чего может быть полезным.
var AdminPermChecker = vkc.HandlerAccessCheck[MyDependencies]{
	Checker: func(handler *vkc.CommandHandler[MyDependencies], ctx vkc.CommandContext[MyDependencies]) bool {
		return ctx.Message.FromID == 1
	},
}

// Обработчик, реагирующий на !hello и !hi
var HandleHello = vkc.CommandHandler[MyDependencies]{
	Pattern: vkc.ListOf([]string{"hello", "hi"}),
	Executor: func(ctx vkc.CommandContext[MyDependencies]) error {
		if ctx.Dependency.MyStore[strconv.Itoa(ctx.Message.FromID)] == "banned" {
			return ctx.SendText("Вы забанены и не можете использовать эту команду.")
		}
		return ctx.SendText("Hello, world!")
	},
}

// Обработчик, реагирующий на !ban. Из-за проверки работает только с аккаунтом Дурова.
var HandleBan = vkc.CommandHandler[MyDependencies]{
	Pattern:     vkc.Text("ban"),
	AccessCheck: &AdminPermChecker,
	Executor: func(ctx vkc.CommandContext[MyDependencies]) error {
		if len(ctx.Arguments) == 0 {
			return ctx.SendText("Использование: !ban <айди>")
		}
		/* "баним" пользователя... */
		ctx.Dependency.MyStore[ctx.Arguments[0]] = "banned"
		return ctx.SendText("Пользователь %d забанен!", ctx.Arguments[0])
	},
}

func main() {
	// Инициализация команд
	handlers := vkc.Commands[MyDependencies]{
		Prefix:   vkc.PrefixText("!"),
		Handlers: &[]*vkc.CommandHandler[MyDependencies]{&HandleHello, &HandleBan},
		// тут также доступно несколько обработчиков для событий, например, OnMessage, OnCommandError и т.д.
	}
	handlers.Dependencies = MyDependencies{
		MyStore:     make(map[string]any),
		AllHandlers: handlers.Handlers,
	}

	vk := api.NewVK( /* Ваш токен */ )

	lp, err := longpoll.NewLongPollCommunity(vk)
	if err != nil {
		log.Fatal(err)
	}

	// Настраиваем перехват сообщений
	err = handlers.AttachToLongPoll(vk, lp)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Start Long Poll\n")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
```
