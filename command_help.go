package vkc

// Помощь по команде. Содержит информацию о названии, описании и примерах использования.
//
// Также включает флаг Hidden для скрытия команды из списка помощи. Полезно при написании команды автопомощи.
//
// Пример использования:
//
//	help := CommandHelp{
//		Title:   "hello",
//		Brief:   "Приветствует пользователя",
//		Usage:   "!hello [имя]",
//		Aliases: "hello, hi",
//		Hidden:  false,
//	}
type CommandHelp struct {
	Title   string
	Brief   string
	Usage   string
	Aliases string
	Hidden  bool
}
