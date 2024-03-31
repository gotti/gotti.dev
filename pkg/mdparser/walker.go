package mdparser

func WalkObject(object Object, f func(Object)) {
	f(object)
	if c, ok := object.(Container); ok {
		cs := c.GetChildren()
		for _, c := range cs {
			WalkObject(c, f)
		}
	}
}

func GetAllInlineBlocks(object Object) []InlineBlocks {
	var ret []InlineBlocks
	f := func(o Object) {
		if c, ok := o.(EndObject); ok {
			co := c.GetContents()
			ret = append(ret, co)
		}
	}
	WalkObject(object, f)
	return ret
}

func GetSpecifiedInlineBlocks(object Object, t InlineType) []InlineBlock {
	var ret []InlineBlock
	b := GetAllInlineBlocks(object)
	for _, bb := range b {
		for _, bbb := range bb.Children {
			if bbb.GetType() == t {
				ret = append(ret, bbb)
			}
		}
	}
	return ret
}

func ReplaceInlineBlocks(object Object, f func(InlineBlock) InlineBlock) {
	if c, ok := object.(Container); ok {
		cs := c.GetChildren()
		for _, c := range cs {
			ReplaceInlineBlocks(c, f)
		}
	}
	if c, ok := object.(EndObject); ok {
		co := c.GetContents()
		for i, b := range co.Children {
			n := f(b)
			co.Children[i] = n
		}
	}
}
