package main

type BalancedQuinary string

func (b BalancedQuinary) String() string {
	return string(b)
}

func (b BalancedQuinary) Int() int {
	var result int
	for _, c := range b {
		result *= 5
		switch c {
		case '=':
			result -= 2
		case '-':
			result -= 1
		case '0':
			result += 0
		case '1':
			result += 1
		case '2':
			result += 2
		}
	}
	return result
}

func BalancedQuinaryFromInt(i int) BalancedQuinary {
	var result BalancedQuinary
	carry := 0
	for i != 0 || carry != 0 {
		m := (i % 5) + carry
		carry = 0

		switch m {
		case -2:
			result = "=" + result
		case -1:
			result = "-" + result
		case 0:
			result = "0" + result
		case 1:
			result = "1" + result
		case 2:
			result = "2" + result
		case 3:
			// This digit is = (-2), but we need to carry a 1
			result = "=" + result
			carry = 1
		case 4:
			// This digit is - (-1), but we need to carry a 1
			result = "-" + result
			carry = 1
		case 5:
			// This digit is 0 (0), but we need to carry a 1
			result = "0" + result
			carry = 1
		}
		i /= 5
	}
	return result
}
