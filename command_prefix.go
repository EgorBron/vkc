package vkc

import (
	"regexp"
	"strings"
)

// Функция для проверки префикса команды. Возвращает успешность совпадения и остаток (сама команда и ее аргументы).
type PrefixMatcher func(input string) (matched bool, remaining string)

// Провайдер для создания PrefixMatcher из значения. Используется объектом команд и обработчиком нового сообщения.
//
// Пример использования (см. соответствующие провайдеры):
//
//	// Префикс, совпадающий только с "!"
//	Prefix: PrefixText("!")
//	// Префикс, совпадающий с любой строкой из среза ["!", "эй бот"]
//	Prefix: PrefixListOf([]string{"!", "эй бот"})
//	// Префикс, совпадающий со *скомпилированным* регулярным выражением
//	Prefix: PrefixRegex(regexp.MustCompile(`(?i)^(?:эй\s)?бот( .+)?`))
//	// Префикс, совпадающий с регулярным выражением в виде строки
//	Prefix: PrefixRegex(`(?i)^(?:эй\s)?бот( .+)?`)
//	// Префикс находится с помощью функции
//	Prefix: PrefixFunc(func(input string) (bool, string) { ... })
type PrefixMatcherProvider[T any] func(matcher T) PrefixMatcher

// Провайдер для точного совпадение с одной строкой.
var PrefixText PrefixMatcherProvider[string] = func(matcher string) PrefixMatcher {
	return func(input string) (bool, string) {
		if strings.HasPrefix(input, matcher) {
			remaining := strings.TrimSpace(input[len(matcher):])
			return true, remaining
		}
		return false, ""
	}
}

// Провайдер для совпадения с любой строкой из среза.
var PrefixListOf PrefixMatcherProvider[[]string] = func(matcher []string) PrefixMatcher {
	return func(input string) (bool, string) {
		for _, prefix := range matcher {
			if strings.HasPrefix(input, prefix) {
				remaining := strings.TrimSpace(input[len(prefix):])
				return true, remaining
			}
		}
		return false, ""
	}
}

// Провайдер для поиска совпадений *скомпилированным* регулярным выражением.
//
// Важно: в группу с индексом 1 обязательно должен попасть остаток, т.е. все, что идет после префикса.
// Если нужно сгруппировать части регулярки, то следует использовать группы без захвата (т.е. (?:...)).
//
// Также рекомендуется включить в регулярное выражение границу начала строки и отключить чувствительность к регистру символов. Например, так:
//
//	PrefixRegex(regex.MustCompile(`(?i)^вашерегулярноевыражение`))
var PrefixRegex PrefixMatcherProvider[*regexp.Regexp] = func(matcher *regexp.Regexp) PrefixMatcher {
	return func(input string) (bool, string) {
		matches := matcher.FindStringSubmatch(input)
		if matches == nil {
			return false, ""
		}
		remaining := strings.TrimSpace(matches[1])
		return true, remaining
	}
}

// Провайдер для поиска совпадений *строковым* регулярным выражением. Само выражение компилируется внутри провайдера.
//
// Важно: в группу с индексом 1 обязательно должен попасть остаток, т.е. все, что идет после префикса.
// Если нужно сгруппировать части регулярки, то следует использовать группы без захвата (т.е. (?:...)).
//
// Также рекомендуется включить в регулярное выражение границу начала строки и отключить чувствительность к регистру символов. Например, так:
//
//	PrefixRegexStr(`(?i)^вашерегулярноевыражение`)
//
// Из-за MustCompile может вызывать панику, если регулярное выражение составлено некорректно. Это поведение нельзя переопределить.
var PrefixRegexStr PrefixMatcherProvider[string] = func(matcher string) PrefixMatcher {
	re := regexp.MustCompile(matcher)
	return func(input string) (bool, string) {
		matches := re.FindStringSubmatch(input)
		if matches == nil {
			return false, ""
		}
		remaining := strings.TrimSpace(matches[1])
		return true, remaining
	}
}

// Провайдер для поиска совпадений с помощью функции.
//
// Функция должна возвращать два значения: булев (найден ли префикс в строке) и строку (весь текст после префикса).
//
// Изначально задумывалось, что провайдер будет получать дополнительный контекст для кастомизации префикса, но пока что данный функционал отсутствует.
var PrefixFunc PrefixMatcherProvider[func(string) (bool, string)] = func(matcher func(string) (bool, string)) PrefixMatcher {
	return matcher
}
