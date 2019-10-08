package xfw

type XTypeEnum int

type XApiString interface {
	String() string
}

type XError interface {
	Error() string
}

type XHash 	 map[string]interface{}
type XArray  []interface{}
type XRecord []XHash
type XObject struct{}

const (
	Xunknown XTypeEnum = iota
	Xint
	Xint8
	Xint16
	Xint32
	Xint64
	Xuint
	Xuint8
	Xuint16
	Xuint32
	Xuint64
	Xbool
	Xfloat32
	Xfloat64
	Xdatetime
	Xstring
	Xobject
	Xarray
	Xrecord
)
