package functions

func Capitalize(s string) string {
	var result string
	IsNewWord := true
	for _, l := range s {
		alph := (l >= 'a' && l <= 'z') || (l >= 'A' && l <= 'Z') || (l >= '0' && l <= '9')
		if alph {
			if IsNewWord {
				if l >= 'a' && l <= 'z' {
					l = l + -32
				}
				IsNewWord = false
			} else {
				if l >= 'A' && l <= 'Z' {
					l = l + 32
				}
			}
		} else {
			IsNewWord = true
		}
		result += string(l)
	}
	return result
}
