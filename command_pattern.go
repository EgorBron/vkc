package vkc

import (
	"regexp"
	"slices"
)

// Функция для проверки шаблона команды.
type CommandPattern func(input string) bool

// Провайдер для создания CommandPattern из значения. Используется в обработчиках.
//
// Пример использования (см. соответствующие провайдеры):
//
//	// Команда, совпадающая только с "help"
//	Pattern: Text("help")
//	// Команда, совпадающая с любой строкой из среза ["help", "h"]
//	Pattern: ListOf([]string{"help", "h"})
//	// Команда, совпадающая со *скомпилированным* регулярным выражением
//	Pattern: Regex(regexp.MustCompile(`^user list( .+)?$`))
//	// Команда, совпадающая с регулярным выражением в виде строки
//	Pattern: RegexCmd(`^user list( .+)?$`)
type CommandPatternProvider[T any] func(matcher T) CommandPattern

// Провайдер для точного совпадение с одной строкой.
var Text CommandPatternProvider[string] = func(matcher string) CommandPattern {
	return func(input string) bool {
		return matcher == input
	}
}

// Провайдер для совпадения с любой строкой из среза.
var ListOf CommandPatternProvider[[]string] = func(matcher []string) CommandPattern {
	return func(input string) bool {
		return slices.Contains(matcher, input)
	}
}

// Провайдер для поиска совпадений *скомпилированным* регулярным выражением.
//
// Важно: в группу с индексом 1 обязательно должен попасть остаток, т.е. все, что идет после названия команды.
// Если нужно сгруппировать части регулярки, то следует использовать группы без захвата (т.е. (?:...)).
var Regex CommandPatternProvider[*regexp.Regexp] = func(matcher *regexp.Regexp) CommandPattern {
	return func(input string) bool {
		return matcher.Match([]byte(input))
	}
}

// Провайдер для поиска совпадений *строковым* регулярным выражением. Само выражение компилируется внутри провайдера.
//
// Важно: в группу с индексом 1 обязательно должен попасть остаток, т.е. все, что идет после названия команды.
// Если нужно сгруппировать части регулярки, то следует использовать группы без захвата (т.е. (?:...)).
//
// Из-за MustCompile может вызывать панику, если регулярное выражение составлено некорректно. Это поведение нельзя переопределить.
var RegexStr CommandPatternProvider[string] = func(matcher string) CommandPattern {
	re := regexp.MustCompile(matcher)
	return func(input string) bool {
		return re.Match([]byte(input))
	}
}
