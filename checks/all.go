package checks

import "github.com/UrbanCompass/thriftlint"

var AllCheckers = thriftlint.Checks{
	CheckIndentation(),
	CheckNames(nil, nil),
	CheckOptional(),
	CheckDefaultValues(),
	CheckEnumSequence(),
	CheckMapKeys(),
	CheckTypeReferences(),
	CheckStructFieldOrder(),
}
