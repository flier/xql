package xql

import (
	"fmt"
	"strings"
)

type DataType interface {
	fmt.Stringer

	ColumnDefOption

	dataType() DataType
}

var (
	_ DataType = &StringType{}
	_ DataType = CreateStringTypeFunc(nil)
	_ DataType = &BinaryType{}
	_ DataType = CreateBinaryTypeFunc(nil)
	_ DataType = &NumericType{}
	_ DataType = CreateExactNumericTypeFunc(nil)
	_ DataType = CreateApproximateNumericTypeFunc(nil)
	_ DataType = &IntType{}
	_ DataType = &BoolType{}
	_ DataType = &DateType{}
	_ DataType = &DateTimeType{}
	_ DataType = CreateDateTimeTypeFunc(nil)
	_ DataType = &IntervalType{}
)

//go:generate stringer -type=CharKind -linecomment

type CharKind int

const (
	KindCharacter                    CharKind = iota // CHARACTER
	KindChar                                         // CHAR
	KindCharacterVarying                             // CHARACTER VARYING
	KindCharVarying                                  // CHAR VARYING
	KindVarChar                                      // VARCHAR
	KindCharacterLargeObject                         // CHARACTER LARGE OBJECT
	KindCharLargeObject                              // CHAR LARGE OBJECT
	KindClob                                         // CLOB
	KindText                                         // TEXT
	KindNationalCharacter                            // NATIONAL CHARACTER
	KindNationalChar                                 // NATIONAL CHAR
	KindNChar                                        // NCHAR
	KindNationalCharacterVarying                     // NATIONAL CHARACTER VARYING
	KindNationalCharVarying                          // NATIONAL CHAR VARYING
	KindNCharVarying                                 // NCHAR VARYING
	KindNationalCharacterLargeObject                 // NATIONAL CHARACTER LARGE OBJECT
	KindNCharLargeObject                             // NCHAR LARGE OBJECT
	KindNClob                                        // NCLOB
)

type Length = uint

type CharLength = Length

//go:generate stringer -type=CharLengthUnit -linecomment

type CharLengthUnit int

const (
	UnitChars       CharLengthUnit = iota // CHARACTERS
	UnitOctets                            // OCTETS
	UnitCodeUnits32                       // CODEUNITS32
)

var (
	Chars       = charLengthUnit(UnitChars)
	Octets      = charLengthUnit(UnitOctets)
	CodeUnits32 = charLengthUnit(UnitCodeUnits32)
)

type CharSetName string

type StringType struct {
	Kind    CharKind
	Len     CharLength
	Unit    *CharLengthUnit
	CharSet CharSetName
	Collate *CollateClause
}

var (
	Character                    = stringType(KindCharacter)
	Char                         = stringType(KindChar)
	CharacterVarying             = stringType(KindCharacterVarying)
	CharVarying                  = stringType(KindCharVarying)
	VarChar                      = stringType(KindVarChar)
	CharacterLargeObject         = stringType(KindCharacterLargeObject)
	CharLargeObject              = stringType(KindCharLargeObject)
	Clob                         = stringType(KindClob)
	Text                         = &StringType{Kind: KindText}
	NationalCharacter            = stringType(KindNationalCharacter)
	NationalChar                 = stringType(KindNationalChar)
	NChar                        = stringType(KindNChar)
	NationalCharacterVarying     = stringType(KindNationalCharacterVarying)
	NationalCharVarying          = stringType(KindNationalCharVarying)
	NCharVarying                 = stringType(KindNCharVarying)
	NationalCharacterLargeObject = stringType(KindNationalCharacterLargeObject)
	NCharLargeObject             = stringType(KindNCharLargeObject)
	NClob                        = stringType(KindNClob)
)

type StringTypeOption interface {
	applyStringType(*StringType)
}

type applyStringTypeFunc func(*StringType)

func (f applyStringTypeFunc) applyStringType(t *StringType) { f(t) }

func CharSet(name string) StringTypeOption {
	return applyStringTypeFunc(func(t *StringType) {
		t.CharSet = CharSetName(name)
	})
}

func charLengthUnit(unit CharLengthUnit) StringTypeOption {
	return applyStringTypeFunc(func(t *StringType) {
		t.Unit = &unit
	})
}

type CreateStringTypeFunc func(len CharLength, x ...StringTypeOption) *StringType

func (f CreateStringTypeFunc) dataType() DataType          { return f(0) }
func (f CreateStringTypeFunc) applyColumnDef(d *ColumnDef) { d.Type = f(0) }
func (f CreateStringTypeFunc) String() string {
	return f(0).String()
}
func (f CreateStringTypeFunc) With(x ...StringTypeOption) *StringType {
	return f(0).With(x...)
}
func (f CreateStringTypeFunc) WithUnit(unit CharLengthUnit) *StringType {
	return f(0).WithUnit(unit)
}
func (f CreateStringTypeFunc) WithCharSet(name CharSetName) *StringType {
	return f(0).WithCharSet(name)
}
func (f CreateStringTypeFunc) WithCollate(name string) *StringType {
	return f(0).WithCollate(name)
}

func stringType(kind CharKind) CreateStringTypeFunc {
	return func(len CharLength, x ...StringTypeOption) *StringType {
		t := &StringType{
			Kind: kind,
			Len:  len,
		}

		return t.With(x...)
	}
}

func (t *StringType) dataType() DataType          { return t }
func (t *StringType) applyColumnDef(d *ColumnDef) { d.Type = t }

func (t *StringType) With(x ...StringTypeOption) *StringType {
	for _, opt := range x {
		opt.applyStringType(t)
	}
	return t
}

func (t *StringType) WithUnit(unit CharLengthUnit) *StringType {
	t.Unit = &unit
	return t
}

func (t *StringType) WithCharSet(name CharSetName) *StringType {
	t.CharSet = name
	return t
}

func (t *StringType) WithCollate(name string) *StringType {
	t.Collate = &CollateClause{SchemaQName(name).LocalOrSchemaQName()}
	return t
}

func (t *StringType) String() string {
	var b strings.Builder

	b.WriteString(t.Kind.String())

	if t.Len > 0 {
		fmt.Fprintf(&b, "(%d", t.Len)
		if t.Unit != nil {
			fmt.Fprintf(&b, " %s", t.Unit.String())
		}
		b.WriteByte(')')
	}

	if len(t.CharSet) > 0 {
		fmt.Fprintf(&b, " CHARACTER SET %s", t.CharSet)
	}
	if t.Collate != nil {
		fmt.Fprintf(&b, " %s", t.Collate)
	}

	return b.String()
}

//go:generate stringer -type=BinaryKind -linecomment

type BinaryKind int

const (
	KindBinary         BinaryKind = iota // BINARY
	KindBinaryVarying                    // BINARY VARYING
	KindVarBinary                        // VARBINARY
	KindBigLargeObject                   // BINARY LARGE OBJECT
	KindBlob                             // BLOB
)

type BinaryType struct {
	Kind BinaryKind
	Len  Length
}

var (
	Binary         = binaryType(KindBinary)
	BinaryVarying  = binaryType(KindBinaryVarying)
	VarBinary      = binaryType(KindVarBinary)
	BigLargeObject = binaryType(KindBigLargeObject)
	Blob           = binaryType(KindBlob)
)

type CreateBinaryTypeFunc func(len Length) *BinaryType

func (f CreateBinaryTypeFunc) dataType() DataType          { return f(0) }
func (f CreateBinaryTypeFunc) applyColumnDef(d *ColumnDef) { d.Type = f(0) }
func (f CreateBinaryTypeFunc) String() string              { return f(0).String() }

func binaryType(kind BinaryKind) CreateBinaryTypeFunc {
	return func(len Length) *BinaryType {
		return &BinaryType{kind, len}
	}
}

func (t *BinaryType) dataType() DataType          { return t }
func (t *BinaryType) applyColumnDef(d *ColumnDef) { d.Type = t }

func (t *BinaryType) String() string {
	var b strings.Builder

	b.WriteString(t.Kind.String())

	if t.Len > 0 {
		fmt.Fprintf(&b, "(%d)", t.Len)
	}

	return b.String()
}

type PrecisionOption interface {
	NumericTypeOption
	DateTimeTypeOption
}

type precisionOption uint

func (o precisionOption) applyNumericType(t *NumericType) {
	t.Precision = uint(o)
}

func (o precisionOption) applyDateTimeType(t *DateTimeType) {
	t.Precision = uint(o)
}

func Precision(precision uint) PrecisionOption {
	return precisionOption(precision)
}

//go:generate stringer -type=NumericKind -linecomment

type NumericKind int

const (
	KindNumeric         NumericKind = iota // NUMERIC
	KindDecimal                            // DECIMAL
	KindDec                                // DEC
	KindFloat                              // FLOAT
	KindReal                               // REAL
	KindDoublePrecision                    // DOUBLE PRECISION
	KindDecFloat                           // DECFLOAT
	KindSmallSerial                        // SMALLSERIAL
	KindSerial                             // SERIAL
	KindBigSerial                          // BIGSERIAL
)

type NumericType struct {
	Kind      NumericKind
	Precision uint
	Scale     int
}

var (
	Numeric         = exactNumericType(KindNumeric)
	Decimal         = exactNumericType(KindDecimal)
	Dec             = exactNumericType(KindDec)
	Float           = approximateNumericType(KindFloat)
	DecFloat        = approximateNumericType(KindDecFloat)
	Real            = &NumericType{Kind: KindReal}
	DoublePrecision = &NumericType{Kind: KindDoublePrecision}
	SmallSerial     = &NumericType{Kind: KindSmallSerial}
	Serial          = &NumericType{Kind: KindSerial}
	BigSerial       = &NumericType{Kind: KindBigSerial}
)

type NumericTypeOption interface {
	applyNumericType(*NumericType)
}

type CreateExactNumericTypeFunc func(precision uint, scale int) *NumericType

func (f CreateExactNumericTypeFunc) dataType() DataType          { return f(0, 0) }
func (f CreateExactNumericTypeFunc) applyColumnDef(d *ColumnDef) { d.Type = f(0, 0) }
func (f CreateExactNumericTypeFunc) String() string              { return f(0, 0).String() }

func (f CreateExactNumericTypeFunc) With(x ...NumericTypeOption) *NumericType {
	return f(0, 0).With(x...)
}

func (f CreateExactNumericTypeFunc) WithPrecision(precision uint) *NumericType {
	return f(precision, 0)
}

func (f CreateExactNumericTypeFunc) WithScala(scala int) *NumericType {
	return f(0, scala)
}

func exactNumericType(kind NumericKind) CreateExactNumericTypeFunc {
	return func(precision uint, scale int) *NumericType {
		return &NumericType{kind, precision, scale}
	}
}

type CreateApproximateNumericTypeFunc func(precision uint) *NumericType

func (f CreateApproximateNumericTypeFunc) dataType() DataType          { return f(0) }
func (f CreateApproximateNumericTypeFunc) applyColumnDef(d *ColumnDef) { d.Type = f(0) }
func (f CreateApproximateNumericTypeFunc) String() string              { return f(0).String() }

func approximateNumericType(kind NumericKind) CreateApproximateNumericTypeFunc {
	return func(precision uint) *NumericType {
		return &NumericType{kind, precision, 0}
	}
}

func (t *NumericType) dataType() DataType          { return t }
func (t *NumericType) applyColumnDef(d *ColumnDef) { d.Type = t }

func (t *NumericType) With(x ...NumericTypeOption) *NumericType {
	for _, opt := range x {
		opt.applyNumericType(t)
	}
	return t
}

func (t *NumericType) WithPrecision(precision uint) *NumericType {
	t.Precision = precision
	return t
}

func (t *NumericType) WithScala(scala int) *NumericType {
	t.Scale = scala
	return t
}

func (t *NumericType) String() string {
	var b strings.Builder

	b.WriteString(t.Kind.String())

	if t.Precision > 0 {
		fmt.Fprintf(&b, "(%d", t.Precision)
		if t.Scale > 0 {
			fmt.Fprintf(&b, ", %d", t.Scale)
		}
		b.WriteByte(')')
	}

	return b.String()
}

//go:generate stringer -type=IntKind -linecomment

type IntKind int

const (
	KindTinyInt   IntKind = iota // TINYINT
	KindSmallInt                 // SMALLINT
	KindMediumInt                // MEDIUMINT
	KindInt                      // INT
	KindInteger                  // INTEGER
	KindBigInt                   // BIGINT
)

type IntType struct {
	Kind IntKind
	Bits int
}

var (
	TinyInt   = &IntType{KindTinyInt, 8}
	SmallInt  = &IntType{KindSmallInt, 16}
	MediumInt = &IntType{KindMediumInt, 32}
	Int       = &IntType{KindInt, 32}
	Integer   = &IntType{KindInteger, 32}
	BigInt    = &IntType{KindBigInt, 64}
)

func (t *IntType) dataType() DataType          { return t }
func (t *IntType) applyColumnDef(d *ColumnDef) { d.Type = t }
func (t *IntType) String() string              { return t.Kind.String() }

//go:generate stringer -type=BoolKind -linecomment

type BoolKind int

const (
	KindBoolean BoolKind = iota // BOOLEAN
	KindBit                     // BIT
)

type BoolType struct {
	Kind BoolKind
}

var (
	Boolean = &BoolType{KindBoolean}
	Bit     = &BoolType{KindBit}
)

func (t *BoolType) dataType() DataType          { return t }
func (t *BoolType) applyColumnDef(d *ColumnDef) { d.Type = t }
func (t *BoolType) String() string              { return t.Kind.String() }

//go:generate stringer -type=DateTimeKind -linecomment

type DateTimeKind int

const (
	KindDate          DateTimeKind = iota // DATE
	KindSmallDateTime                     // SMALLDATETIME
	KindDateTime                          // DATETIME
	KindDateTime2                         // DATETIME2
	KindYear                              // YEAR
	KindTime                              // TIME
	KindTimestamp                         // TIMESTAMP
)

type DateType struct {
	Kind DateTimeKind
}

func (t *DateType) dataType() DataType          { return t }
func (t *DateType) applyColumnDef(d *ColumnDef) { d.Type = t }
func (t *DateType) String() string              { return t.Kind.String() }

type DateTimeType struct {
	DateType
	Precision uint
	TimeZone  *TimeZone
}

var (
	Date          = &DateType{Kind: KindDate}
	SmallDateTime = &DateType{Kind: KindSmallDateTime}
	DateTime      = &DateType{Kind: KindDateTime}
	DateTime2     = &DateType{Kind: KindDateTime2}
	DateYear      = &DateType{Kind: KindYear}
	Time          = dateTimeType(KindTime)
	Timestamp     = dateTimeType(KindTimestamp)
)

type DateTimeTypeOption interface {
	applyDateTimeType(*DateTimeType)
}

type dateTimeTypeOptionFunc func(*DateTimeType)

var (
	WithTimeZone = dateTimeTypeOptionFunc(func(t *DateTimeType) {
		tz := TimeZone(true)
		t.TimeZone = &tz
	})
	WithoutTimeZone = dateTimeTypeOptionFunc(func(t *DateTimeType) {
		tz := TimeZone(false)
		t.TimeZone = &tz
	})
)

func (f dateTimeTypeOptionFunc) applyDateTimeType(t *DateTimeType) { f(t) }

type CreateDateTimeTypeFunc func(precision uint, x ...DateTimeTypeOption) *DateTimeType

func (f CreateDateTimeTypeFunc) dataType() DataType                         { return f(0) }
func (f CreateDateTimeTypeFunc) applyColumnDef(d *ColumnDef)                { d.Type = f(0) }
func (f CreateDateTimeTypeFunc) String() string                             { return f(0).String() }
func (f CreateDateTimeTypeFunc) With(x ...DateTimeTypeOption) *DateTimeType { return f(0).With(x...) }
func (f CreateDateTimeTypeFunc) WithTimeZone() *DateTimeType                { return f(0).WithTimeZone() }
func (f CreateDateTimeTypeFunc) WithoutTimeZone() *DateTimeType             { return f(0).WithoutTimeZone() }

func dateTimeType(kind DateTimeKind) CreateDateTimeTypeFunc {
	return func(precision uint, x ...DateTimeTypeOption) *DateTimeType {
		t := &DateTimeType{DateType{kind}, precision, nil}

		return t.With(x...)
	}
}

func (t *DateTimeType) dataType() DataType { return t }
func (t *DateTimeType) With(x ...DateTimeTypeOption) *DateTimeType {
	for _, opt := range x {
		opt.applyDateTimeType(t)
	}
	return t
}

func (t *DateTimeType) WithTimeZone() *DateTimeType {
	tz := TimeZone(true)
	t.TimeZone = &tz
	return t
}

func (t *DateTimeType) WithoutTimeZone() *DateTimeType {
	tz := TimeZone(false)
	t.TimeZone = &tz
	return t
}

func (t *DateTimeType) String() string {
	var b strings.Builder

	b.WriteString(t.Kind.String())

	if t.Precision > 0 {
		fmt.Fprintf(&b, "(%d)", t.Precision)
	}

	if t.TimeZone != nil {
		fmt.Fprintf(&b, " %s", t.TimeZone)
	}

	return b.String()
}

type TimeZone bool

func (tz TimeZone) String() string {
	if tz {
		return "WITH TIME ZONE"
	}

	return "WITHOUT TIME ZONE"
}

//go:generate stringer -type=DateTimeFieldKind -linecomment

type DateTimeFieldKind int

const (
	FieldYear        DateTimeFieldKind = iota // YEAR
	FieldMonth                                // MONTH
	FieldWeek                                 // WEEK
	FieldDay                                  // DAY
	FieldHour                                 // HOUR
	FieldMinute                               // MINUTE
	FieldSecond                               // SECOND
	FieldMillisecond                          // MILLISECOND
	FieldMicrosecond                          // MICROSECOND
)

var (
	Year        = dateTimeField(FieldYear)
	Week        = dateTimeField(FieldWeek)
	Month       = dateTimeField(FieldMonth)
	Day         = dateTimeField(FieldDay)
	Hour        = dateTimeField(FieldHour)
	Minute      = dateTimeField(FieldMinute)
	Second      = dateTimeField(FieldSecond)
	Millisecond = dateTimeField(FieldMicrosecond)
	Microsecond = dateTimeField(FieldMillisecond)
)

type CreateDateTimeFieldFunc func(precision uint) *DateTimeField

func (f CreateDateTimeFieldFunc) dateTimeField() *DateTimeField { return f(0) }
func (f CreateDateTimeFieldFunc) To(end ToDateTimeField) *IntervalType {
	return &IntervalType{*f(0), *end.dateTimeField()}
}
func (f CreateDateTimeFieldFunc) String() string { return f(0).String() }

func dateTimeField(kind DateTimeFieldKind) CreateDateTimeFieldFunc {
	return func(precision uint) *DateTimeField {
		return &DateTimeField{kind, precision}
	}
}

type ToDateTimeField interface {
	dateTimeField() *DateTimeField
}

type DateTimeField struct {
	Kind      DateTimeFieldKind
	Precision uint
}

func (f *DateTimeField) To(end ToDateTimeField) *IntervalType {
	return &IntervalType{*f, *end.dateTimeField()}
}

func (f *DateTimeField) dateTimeField() *DateTimeField { return f }

func (f *DateTimeField) String() string {
	var b strings.Builder

	b.WriteString(f.Kind.String())
	if f.Precision > 0 {
		fmt.Fprintf(&b, "(%d)", f.Precision)
	}

	return b.String()
}

type IntervalType struct {
	Start DateTimeField
	End   DateTimeField
}

func Interval(start, end ToDateTimeField) *IntervalType {
	return &IntervalType{*start.dateTimeField(), *end.dateTimeField()}
}

func (t *IntervalType) dataType() DataType          { return t }
func (t *IntervalType) applyColumnDef(d *ColumnDef) { d.Type = t }
func (t *IntervalType) String() string {
	return fmt.Sprintf("INTERVAL %s TO %s", &t.Start, &t.End)
}
