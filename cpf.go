func isCPF(cpf string) bool {

	cpf = onlyDigits(cpf)

	if len(cpf) != 11 {
		return false
	}

	if allDigitsEqual(cpf) {
		return false
	}

	firstCheck := cpfCheckDigit(cpf[:9], 10)
	secondCheck := cpfCheckDigit(cpf[:10], 11)

	return cpf[9] == byte(firstCheck+'0') &&
		cpf[10] == byte(secondCheck+'0')
}

func onlyDigits(value string) string {

	var b strings.Builder

	for _, r := range value {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func allDigitsEqual(cpf string) bool {

	first := cpf[0]

	for i := 1; i < len(cpf); i++ {
		if cpf[i] != first {
			return false
		}
	}

	return true
}

func cpfCheckDigit(base string, weightStart int) int {

	sum := 0
	weight := weightStart

	for i := 0; i < len(base); i++ {
		sum += int(base[i]-'0') * weight
		weight--
	}

	rest := sum % 11

	if rest < 2 {
		return 0
	}

	return 11 - rest
}