package main

const httpPrefix = "http://"

func ContainsHttpPrefix(message string, index int) bool {
	if index+len(httpPrefix) >= len(message) {
		return false
	}
	return message[index] == httpPrefix[0] &&
		message[index+1] == httpPrefix[1] &&
		message[index+2] == httpPrefix[2] &&
		message[index+3] == httpPrefix[3] &&
		message[index+4] == httpPrefix[4] &&
		message[index+5] == httpPrefix[5] &&
		message[index+6] == httpPrefix[6]
}

func ReplaceAllLinks(message string) string {
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

		if ContainsHttpPrefix(message, index) {
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
