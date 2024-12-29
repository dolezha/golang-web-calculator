package services

import (
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	if err := validateExpression(expression); err != nil {
		return 0, fmt.Errorf("ошибка в выражении: %v", err)
	}
	return evaluate(expression)
}

func validateExpression(expression string) error {
	openBrackets := 0
	for i, char := range expression {
		switch char {
		case '(':
			openBrackets++
		case ')':
			openBrackets--
			if openBrackets < 0 {
				return fmt.Errorf("неверно расставлены скобки")
			}
		case '+', '-', '*', '/':
			if i == 0 || i == len(expression)-1 || !isDigitOrBracket(expression[i-1]) || !isDigitOrBracket(expression[i+1]) {
				return fmt.Errorf("неверная расстановка операторов")
			}
		}
	}
	if openBrackets != 0 {
		return fmt.Errorf("несоответствие количества открывающих и закрывающих скобок")
	}
	return nil
}

func isDigitOrBracket(c byte) bool {
	return (c >= '0' && c <= '9') || c == '.' || c == ')' || c == '('
}

func evaluate(expression string) (float64, error) {
	var values []float64
	var operators []byte

	applyOperation := func() error {
		if len(values) < 2 {
			return fmt.Errorf("неправильное использование оператора")
		}

		right, left := values[len(values)-1], values[len(values)-2]
		values = values[:len(values)-2]
		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		switch op {
		case '+':
			values = append(values, left+right)
		case '-':
			values = append(values, left-right)
		case '*':
			values = append(values, left*right)
		case '/':
			if right == 0 {
				return fmt.Errorf("деление на ноль")
			}
			values = append(values, left/right)
		}

		return nil
	}

	precedence := func(op byte) int {
		switch op {
		case '+', '-':
			return 1
		case '*', '/':
			return 2
		}
		return 0
	}

	var number strings.Builder
	for i := 0; i < len(expression); i++ {
		char := expression[i]

		if isDigitOrBracket(char) && char != '(' && char != ')' {
			number.WriteByte(char)

			if i == len(expression)-1 || !isDigitOrBracket(expression[i+1]) || expression[i+1] == '(' || expression[i+1] == ')' {
				num, err := strconv.ParseFloat(number.String(), 64)

				if err != nil {
					return 0, fmt.Errorf("некорректное число")
				}

				values = append(values, num)
				number.Reset()
			}

		} else if char == '(' {
			operators = append(operators, char)

		} else if char == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				if err := applyOperation(); err != nil {
					return 0, err
				}
			}
			operators = operators[:len(operators)-1]

		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(char) {
				if err := applyOperation(); err != nil {
					return 0, err
				}
			}
			operators = append(operators, char)
		}
	}

	for len(operators) > 0 {
		if err := applyOperation(); err != nil {
			return 0, err
		}
	}

	if len(values) != 1 {
		return 0, fmt.Errorf("ошибка в выражении")
	}

	return values[0], nil
}
