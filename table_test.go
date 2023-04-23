package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleCreateTable() {
	fmt.Println(CreateTable("films",
		Column("code", Char(5), Constraint("firstkey").PrimaryKey()),
		Column("title", VarChar(40), NotNull),
		Column("did", Integer, NotNull),
		Column("date_prod", Date),
		Column("kind", VarChar(10)),
		Column("len", Interval(Hour, Minute)),
		Constraint("production").Unique("date_prod"),
	))
	// Output:
	// CREATE TABLE films (
	// 	code CHAR(5) CONSTRAINT firstkey PRIMARY KEY,
	// 	title VARCHAR(40) NOT NULL,
	// 	did INTEGER NOT NULL,
	// 	date_prod DATE,
	// 	kind VARCHAR(10),
	// 	len INTERVAL HOUR TO MINUTE,
	// 	CONSTRAINT production UNIQUE (date_prod)
	// )
}

func ExampleCreateTempTable() {
	fmt.Println(CreateTempTable("temp_cities",
		Column("name", VarChar(80), PrimaryKey, NotNull),
	).OnCommitDeleteRows())
	// Output:
	// CREATE TEMPORARY TABLE temp_cities (
	// 	name VARCHAR(80) PRIMARY KEY NOT NULL
	// ) ON COMMIT DELETE ROWS
}

func ExamplePeriodForSystemTime() {
	fmt.Println(CreateTable("Department",
		Column("DeptID", Int, NotNull, PrimaryKey),
		Column("DeptName", VarChar(50), NotNull),
		Column("ValidFrom", DateTime2, Generated.Always().AsRowStart(), NotNull),
		Column("ValidTo", DateTime2, Generated.Always().AsRowEnd(), NotNull),
		PeriodForSystemTime("ValidFrom", "ValidTo"),
	).WithSystemVersioningOn())
	// Output:
	// CREATE TABLE Department (
	// 	DeptID INT NOT NULL PRIMARY KEY,
	// 	DeptName VARCHAR(50) NOT NULL,
	// 	ValidFrom DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
	// 	ValidTo DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
	// 	PERIOD FOR SYSTEM_TIME (ValidFrom, ValidTo)
	// ) WITH (SYSTEM_VERSIONING = ON)
}

func ExamplePeriodFor() {
	fmt.Println(CreateTable("t1",
		Column("name", VarChar(50)),
		Column("date_1", Date),
		Column("date_2", Date),
		PeriodFor("date_period")("date_1", "date_2"),
	))
	// Output:
	// CREATE TABLE t1 (
	// 	name VARCHAR(50),
	// 	date_1 DATE,
	// 	date_2 DATE,
	// 	PERIOD FOR date_period (date_1, date_2)
	// )
}

func ExampleCreateTable_of() {
	fmt.Println(CreateTable("employees").Of("employee_type",
		PrimaryKey("name"),
		Column("salary").WithOptions(Literal("1000").AsDefault()),
	))
	// Output:
	// CREATE TABLE employees OF employee_type (
	// 	PRIMARY KEY (name),
	// 	salary WITH OPTIONS DEFAULT 1000
	// )
}
