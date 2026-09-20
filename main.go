package main

func ReplaceAllLinks(message string) string {
	const httpPrefix = "http://"
	answer := make([]byte, 0, len(message))
	insideLink := false
	index := 0

	for index < len(message) {
		symbol := message[index]
		if insideLink {
			switch symbol {
			case ' ', '\t':
				insideLink = false
				answer = append(answer, symbol)
			default:
				answer = append(answer, '*')
			}
			index++
			continue
		}

		containsPrefix := true
		for i := range httpPrefix {
			if index+i >= len(message) || message[index+i] != httpPrefix[i] {
				containsPrefix = false
				break
			}
		}

		if containsPrefix {
			insideLink = true
			answer = append(answer, httpPrefix...)
			index += 6
			continue
		}

		answer = append(answer, symbol)
		index++
	}
	return string(answer)
}

func main() {
	myMessage := "http://yandex.com htt://wrong-link.com http://hereweare.com  here is my string http://google.com"
	result := ReplaceAllLinks(myMessage)
	println(result)
}
