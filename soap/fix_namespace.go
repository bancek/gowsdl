package soap

import (
	"github.com/beevik/etree"
)

func FixNamespace(bytes []byte) ([]byte, error) {
	doc := etree.NewDocument()

	if err := doc.ReadFromBytes(bytes); err != nil {
		return nil, err
	}

	setTagNamespace := func(e *etree.Element, space string) {
		e.Space = space
	}

	setAttrsNamespace := func(e *etree.Element, space string) {
		for i := range e.Attr {
			if e.Attr[i].Key == "xmlns" || e.Attr[i].Space == "xmlns" {
				e.Attr[i].Key = space
				e.Attr[i].Space = "xmlns"
			}
		}
	}

	setNamespace := func(e *etree.Element, space string) {
		setTagNamespace(e, space)
		setAttrsNamespace(e, space)
	}

	root := doc.Child[0].(*etree.Element)

	setNamespace(root, "S")

	for _, child := range root.Child {
		if e, ok := child.(*etree.Element); ok {
			setNamespace(e, "S")

			if e.Tag == "Body" {
				for _, child := range e.Child {
					if e, ok := child.(*etree.Element); ok {
						if e.Tag == "Fault" {
							setTagNamespace(e, "S")
							setAttrsNamespace(e, "ns4")

							for _, child := range e.Child {
								if e, ok := child.(*etree.Element); ok && e.Tag == "detail" {
									for _, child := range e.Child {
										if e, ok := child.(*etree.Element); ok {
											setNamespace(e, "ns2")
										}
									}
								}
							}
						} else {
							setNamespace(e, "ns1")
						}
					}
				}
			}
		}
	}

	newBytes, err := doc.WriteToBytes()
	if err != nil {
		return nil, err
	}

	return newBytes, nil
}
