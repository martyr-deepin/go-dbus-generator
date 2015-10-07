package introspect

/*
分析dbus的signature结构，分割成独立的最小元素并利用TypeConverter这个接口
去生成对应语言的具体类型。
*/

import (
	"fmt"
)

type TypeId int

const (
	ArrayTypeId TypeId = iota
	StructTypeId
	DictTypeId
	PrimitiveTypeId
)

type TypeMeta struct {
	Signature string
	Type      TypeId
}

func NewTypeMeta(sig string) TypeMeta {
	if len(sig) == 1 {
		return TypeMeta{sig, PrimitiveTypeId}
	}

	s1, s2 := sig[0], sig[1]
	if s1 == '(' {
		return TypeMeta{sig, StructTypeId}
	} else if s1 == 'a' {
		if s2 == '{' {
			return TypeMeta{sig, DictTypeId}
		} else {
			return TypeMeta{sig, ArrayTypeId}
		}
	}
	panic("Isn't a valid signature :" + sig)
}

type StructTypeMeta struct {
	TypeMeta
	Elements []TypeMeta
}

func NewStructTypeMeta(sig string) StructTypeMeta {
	subs := Split(sig[1 : len(sig)-1])
	s := StructTypeMeta{NewTypeMeta(sig), nil}
	for _, sig := range subs {
		s.Elements = append(s.Elements, NewTypeMeta(sig))
	}
	return s
}

type ArrayTypeMeta struct {
	TypeMeta
	Element TypeMeta
}

func NewArrayTypeMeta(sig string) ArrayTypeMeta {
	t := ArrayTypeMeta{
		NewTypeMeta(sig),
		NewTypeMeta(sig[1:len(sig)]),
	}
	return t
}

type DictTypeMeta struct {
	TypeMeta
	First  TypeMeta
	Second TypeMeta
}

func NewDictTypeMeta(sig string) DictTypeMeta {
	t := DictTypeMeta{
		NewTypeMeta(sig),
		NewTypeMeta(sig[2:3]),
		NewTypeMeta(sig[3 : len(sig)-1]),
	}
	return t
}

func one(s string) (string, string) {
	err, rem := validSingle(s, 0)
	if err != nil {
		panic("invliad signature: " + s)
	}
	return s[0 : len(s)-len(rem)], rem
}

func validSingle(s string, depth int) (err error, rem string) {
	if s == "" {
		return fmt.Errorf("empty signature"), ""
	}
	if depth > 64 {
		return fmt.Errorf("%q container nesting too deep", s), ""
	}
	switch s[0] {
	case 'y', 'b', 'n', 'q', 'i', 'u', 'x', 't', 'd', 's', 'g', 'o', 'v', 'h':
		return nil, s[1:]
	case 'a':
		if len(s) > 1 && s[1] == '{' {
			_, s, rem := ParseParenthess(s[1:], '{', '}')

			//key
			if err, _ = validSingle(s[:1], depth+1); err != nil {
				return err, ""
			}

			//value
			err, nr := validSingle(s[1:], depth+1)
			if err != nil {
				return err, ""
			}
			if nr != "" {
				return fmt.Errorf("too many types in dict:%s", s), ""
			}
			return nil, rem
		}
		return validSingle(s[1:], depth+1)
	case '(':
		_, s, rem = ParseParenthess(s, '(', ')')
		for err == nil && s != "" {
			err, s = validSingle(s, depth+1)
		}
		if err != nil {
			rem = ""
		}
		return
	default:
		return fmt.Errorf("invliad sig:%s", s), ""
	}
}

func Split(signature string) []string {
	var r []string
	//fmt.Println("Parse>...:", signature)
	s, rem := one(signature)
	r = append(r, s)
	if rem == "" {
		return r
	}
	r = append(r, Split(rem)...)
	return r
}

func ParseParenthess(s string, left, right rune) (string, string, string) {
	l, r, n := -1, -1, 0
	for i, c := range s {
		if c == left {
			if n == 0 {
				l = i
			}

			n++
		} else if c == right {
			n--

			if n == 0 {
				r = i
				return s[0:l], s[l+1 : r], s[r+1:]
			}
		}

	}
	panic("Didn't match " + s)
}

type TypesTable struct {
	conveter               TypeConverter
	compoundTypeSignatures map[string]struct{}
}

func NewTypesTable(conveter TypeConverter) TypesTable {
	return TypesTable{
		conveter,
		make(map[string]struct{}),
	}
}

func (table TypesTable) recordCompoundType(sigs ...string) {
	for _, sig := range sigs {
		table.compoundTypeSignatures[sig] = struct{}{}
	}
}

func (table TypesTable) GetCompoundTypeSignatures() []string {
	var r []string
	for k, _ := range table.compoundTypeSignatures {
		r = append(r, k)
	}
	return r
}

func (table TypesTable) Get(signature string) string {
	t := NewTypeMeta(signature)
	switch t.Type {
	case PrimitiveTypeId:
		return table.conveter.ToPrimitive(t.Signature)
	case StructTypeId:
		tt := NewStructTypeMeta(t.Signature)

		var elements []string
		for _, s := range tt.Elements {
			elements = append(elements, table.Get(s.Signature))
			table.recordCompoundType(s.Signature)
		}
		return table.conveter.ToStruct(elements...)
	case ArrayTypeId:
		tt := NewArrayTypeMeta(t.Signature)

		table.recordCompoundType(tt.Element.Signature)

		return table.conveter.ToArray(table.Get(tt.Element.Signature))
	case DictTypeId:
		tt := NewDictTypeMeta(t.Signature)
		table.recordCompoundType(tt.First.Signature, tt.Second.Signature)

		return table.conveter.ToDict(table.Get(tt.First.Signature), table.Get(tt.Second.Signature))
	default:
		panic("Invalid Type signautre " + signature)
	}
}
