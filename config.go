package main

type Configuration struct {
	Output       string
	Window       bool
	Outline      bool
	Shadow       bool
	CornerRadius int
	Padding      Sides
	Margin       Sides
	FontFamily   string
	FontSize     float64
	LineHeight   float64
}

func ConfigurationBase() Configuration {
	return Configuration{
		Output:       "out.svg",
		Window:       false,
		Outline:      false,
		Shadow:       false,
		CornerRadius: 0,
		Padding:      NewSides(20, 40, 20, 20),
		Margin:       NewSides(0),
		FontFamily:   "JetBrains Mono",
		FontSize:     14,
		LineHeight:   14 * 1.2,
	}
}

func ConfigurationDecoration() Configuration {
	return Configuration{
		Output:       "out.svg",
		Window:       true,
		Outline:      true,
		Shadow:       true,
		CornerRadius: 6,
		Padding:      NewSides(20, 40, 20, 20),
		Margin:       NewSides(40),
		FontFamily:   "JetBrains Mono",
		FontSize:     14,
		LineHeight:   14 * 1.2,
	}
}

type Sides struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func NewSides(sides ...int) Sides {
	switch len(sides) {
	case 1:
		return Sides{sides[0], sides[0], sides[0], sides[0]}
	case 2:
		return Sides{sides[0], sides[1], sides[0], sides[1]}
	case 4:
		return Sides{sides[0], sides[1], sides[2], sides[3]}
	default:
		return Sides{}
	}
}
