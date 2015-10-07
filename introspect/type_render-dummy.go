package introspect

type TypeConverter interface {
	//ToPrimitive是唯一一个传递dbus signature的方法，用来得到基本类型
	//到对应语言的类型
	ToPrimitive(signature string) string

	//以下3个方法的参数均是通过ToPrimitive转换后的具体类型
	ToArray(element string) string
	ToDict(firstElement string, secondElement string) string
	ToStruct(elements ...string) string
}

type DummyConvert struct {
}

func (DummyConvert) ToPrimitive(signatures string) string {
	return ""
}

func (DummyConvert) ToArray(element string) string {
	return ""
}
func (DummyConvert) ToDict(first, second string) string {
	return ""
}

func (DummyConvert) ToStruct(elements ...string) string {
	return ""
}
