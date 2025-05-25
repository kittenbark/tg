package tg

type Button = InlineKeyboardButton

type ButtonT[T any] interface {
	Build(data T) *Button
}

var (
	_ ButtonT[any] = &buttonStates[any]{}
)

type buttonStates[T any] struct {
	pred   func(data T) int
	states []*Button
}

func (button buttonStates[T]) Build(data T) *Button {
	index := button.pred(data)
	if index < 0 || index >= len(button.states) {
		return nil
	}
	return button.states[index]
}

func ButtonStates[T any](pred func(data T) int, states ...*Button) ButtonT[T] {
	return &buttonStates[T]{
		pred:   pred,
		states: states,
	}
}

func Keyboard[T any](layout [][]ButtonT[T], data ...T) *InlineKeyboardMarkup {
	var defaultValue T
	result := [][]*Button{}
	for _, row := range layout {
		resultRow := []*Button{}
		for _, button := range row {
			if resultButton := button.Build(at(data, 0, defaultValue)); resultButton != nil {
				resultRow = append(resultRow, resultButton)
			}
		}
		result = append(result, resultRow)
	}
	return &InlineKeyboardMarkup{
		InlineKeyboard: result,
	}
}
