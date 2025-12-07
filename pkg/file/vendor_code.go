package file

type VendorCode struct {
	index int
}

func (vc VendorCode) String() string {
	return string([]rune{vc.getFirstLetter(), vc.getSecondLetter(), vc.getThirdLetter()})
}

func (vc VendorCode) getFirstLetter() rune {
	first := 'A' + (vc.index / (26 * 26))
	return rune(first)
}

func (vc VendorCode) getSecondLetter() rune {
	return rune('A' + ((vc.index % (26 * 26)) / 26))
}

func (vc VendorCode) getThirdLetter() rune {
	return rune('A' + (vc.index % 26))
}
