package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleStringType() {
	fmt.Println(Char)
	fmt.Println(Text)
	fmt.Println(VarChar(16))
	fmt.Println(NChar(16, Chars))
	fmt.Println(Char(10, CharSet("utf8mb4")).WithCollate("utf8mb4_unicode_ci"))
	// Output:
	// CHAR
	// TEXT
	// VARCHAR(16)
	// NCHAR(16 CHARACTERS)
	// CHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
}

func ExampleBinaryType() {
	fmt.Println(Binary)
	fmt.Println(VarBinary(3))
	// Output:
	// BINARY
	// VARBINARY(3)
}

func ExampleNumericType() {
	fmt.Println(Real)
	fmt.Println(Float(16))
	fmt.Println(Numeric(3, 1))
	// Output:
	// REAL
	// FLOAT(16)
	// NUMERIC(3, 1)
}

func ExampleIntType() {
	fmt.Println(SmallInt)
	// Output:
	// SMALLINT
}

func ExampleBoolType() {
	fmt.Println(Boolean)
	fmt.Println(Bit)
	// Output:
	// BOOLEAN
	// BIT
}

func ExampleDateType() {
	fmt.Println(Date)
	// Output:
	// DATE
}

func ExampleDateTimeType() {
	fmt.Println(Timestamp)
	fmt.Println(Timestamp.WithoutTimeZone())
	fmt.Println(Time(3).WithTimeZone())
	// Output:
	// TIMESTAMP
	// TIMESTAMP WITHOUT TIME ZONE
	// TIME(3) WITH TIME ZONE
}

func ExampleIntervalType() {
	fmt.Println(Day(3).To(Day))
	fmt.Println(Interval(Hour, Second))
	// Output:
	// INTERVAL DAY(3) TO DAY
	// INTERVAL HOUR TO SECOND
}
