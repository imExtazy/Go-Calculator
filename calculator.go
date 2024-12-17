package main

import (
	"fmt"
	"math"
)

type Stack struct {
	elements []Nature
}

func (s *Stack) Push(value Nature) {
	s.elements = append(s.elements, value)
}

func (s *Stack) last() Nature {
	if len(s.elements) == 0 {
		return Nature{}
	}
	lastIndex := len(s.elements) - 1
	value := s.elements[lastIndex]
	return value
}

func (s *Stack) Pop() Nature {
	lastIndex := len(s.elements) - 1
	value := s.elements[lastIndex]
	s.elements = s.elements[:lastIndex]
	return value
}

func (s *Stack) IsEmpty() bool {
	return len(s.elements) == 0
}

type Nature struct {
	name     string  // digit, operator
	digit    float64 //если число, то будет хранится здесь
	symbol   string  //если не число, то будет хранится здесь
	priority int     //если оператор, то будет приоритет
}

func get_number(str string, pos *int) (float64, error) {
	ln := len(str)
	var numb float64 = 0
	var numb_after_dot float64 = 0
	have_dot := 0
	len_after_dot := 0.0
	for ; (*pos) < ln; (*pos)++ {
		if str[*pos] == '.' {
			if have_dot == 1 {
				err := fmt.Errorf("Неправильный ввод")
				return numb, err
			} else {
				have_dot = 1
			}
		} else if str[*pos] >= '0' && str[*pos] <= '9' {
			if have_dot == 1 {
				numb_after_dot = numb_after_dot*10 + float64(str[*pos]-'0')
				len_after_dot += 1
			} else {
				numb = numb*10 + float64(str[*pos]-'0')
			}
		} else {
			(*pos)--
			break
		}
	}
	if have_dot == 1 {
		if len_after_dot == 0 {
			return numb, fmt.Errorf("Неправильный вывод")
		} else {
			numb = numb + numb_after_dot/(math.Pow(10.0, len_after_dot))
		}
	}
	return numb, nil
}

func get_priority(symb byte) int {
	switch symb {
	case '(', ')':
		return 0
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 3
	}
}

func is_operator(symb byte) bool {
	switch symb {
	case '+', '-', '*', '/', '(', ')':
		return true
	default:
		return false
	}
}

func calculate_rpn(rpn_elements []Nature) (float64, error) {
	var s Stack
	for _, element := range rpn_elements {
		if element.name == "digit" {
			s.Push(element)
		} else {
			if element.priority == 3 {
				if s.IsEmpty() {
					return 0.0, fmt.Errorf("Неправильный ввод")
				} else if element.symbol == "~" {
					a := s.Pop()
					a.digit = -a.digit
					s.Push(a)
				}
			} else {
				if len(s.elements) < 2 {
					return 0.0, fmt.Errorf("Неправильный ввод")
				} else {
					b := s.Pop()
					a := s.Pop()
					switch element.symbol {
					case "+":
						a.digit = a.digit + b.digit
						s.Push(a)
					case "-":
						a.digit = a.digit - b.digit
						s.Push(a)
					case "*":
						a.digit = a.digit * b.digit
						s.Push(a)
					case "/":
						if b.digit == 0 {
							return 0.0, fmt.Errorf("Неправильный ввод")
						}
						a.digit = a.digit / b.digit
						s.Push(a)
					}
				}
			}
		}
	}
	if len(s.elements) != 1 {
		return 0.0, fmt.Errorf("Неправильный ввод")
	} else {
		return s.Pop().digit, nil
	}
}

func to_rpn(elements []Nature) ([]Nature, error) {
	var rpn_elements []Nature
	s := Stack{}
	for _, element := range elements {
		if element.name == "digit" {
			rpn_elements = append(rpn_elements, element)
		} else if element.symbol != "(" && element.symbol != ")" {
			if s.IsEmpty() {
				s.Push(element)
			} else if element.priority > s.last().priority {
				s.Push(element)
			} else if element.priority == s.last().priority {
				rpn_elements = append(rpn_elements, s.Pop())
				s.Push(element)
			} else {
				for {
					if s.IsEmpty() {
						s.Push(element)
						break
					} else if s.last().priority < element.priority {
						s.Push(element)
						break
					} else {
						rpn_elements = append(rpn_elements, s.Pop())
					}
				}
			}
		} else if element.symbol == "(" {
			s.Push(element)
		} else {
			for {
				if s.IsEmpty() {
					return []Nature{}, fmt.Errorf("Неправильный ввод")
				} else if s.last().symbol == "(" {
					s.Pop()
					break
				} else {
					rpn_elements = append(rpn_elements, s.Pop())
				}
			}
		}
	}
	for {
		if s.IsEmpty() {
			break
		} else if s.last().symbol == "(" {
			return []Nature{}, fmt.Errorf("Неправильный ввод")
		} else {
			rpn_elements = append(rpn_elements, s.Pop())
		}
	}
	return rpn_elements, nil
}

func Calc(expression string) (float64, error) {
	var elements []Nature
	i := 0
	for i < len(expression) {
		if expression[i] >= '0' && expression[i] <= '9' {
			number, err := get_number(expression, &i)
			if err != nil {
				return 0.0, err
			}
			elements = append(elements, Nature{name: "digit", digit: number})
		} else if is_operator(expression[i]) {
			if i != 0 {
				if is_operator(expression[i-1]) {
					if expression[i-1] == '(' {
						if expression[i] == '+' {
							expression = expression[:i] + "#" + expression[i+1:]
						} else if expression[i] == '-' {
							expression = expression[:i] + "~" + expression[i+1:]
						} else {
							return 0.0, fmt.Errorf("Неправильный ввод")
						}
					} else if expression[i-1] == ')' {
						elements = append(elements, Nature{name: "operator", symbol: string(expression[i]), priority: get_priority(expression[i])})
						i++
						continue
					} else if expression[i] == '(' {
						elements = append(elements, Nature{name: "operator", symbol: string(expression[i]), priority: get_priority(expression[i])})
						i++
						continue
					} else {
						return 0.0, fmt.Errorf("Неправильный ввод")
					}
				}
			} else {
				if expression[i] == '(' {
					elements = append(elements, Nature{name: "operator", symbol: string(expression[i]), priority: get_priority(expression[i])})
					i++
					continue
				} else if expression[i] == '+' {
					expression = "#" + expression[1:]
				} else if expression[i] == '-' {
					expression = "~" + expression[1:]
				} else {
					return 0.0, fmt.Errorf("Неправильный ввод")
				}
			}
			elements = append(elements, Nature{name: "operator", symbol: string(expression[i]), priority: get_priority(expression[i])})
		} else if expression[i] == ' ' || expression[i] == '\n' || expression[i] == '\t' {
			i++
			continue
		} else {
			return 0.0, fmt.Errorf("Неправильный ввод")
		}
		i++
	}
	//for _, element := range elements {
	//	println(element.name, " ", element.digit, " ", element.symbol, " ", element.priority)
	//}
	rpn_elements, err := to_rpn(elements)
	if err != nil {
		return 0.0, err
	}
	answer, err1 := calculate_rpn(rpn_elements)
	if err1 != nil {
		return 0.0, err1
	}
	return answer, nil
}

func main() {
	str := ""
	fmt.Scanln(&str)
	answer, err := Calc(str)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(answer)
	}
}
