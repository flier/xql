package xql

import (
	"fmt"
)

type SequenceGeneratorOption interface {
	fmt.Stringer

	Accepter

	sequenceGeneratorOption() SequenceGeneratorOption
}

var (
	_ SequenceGeneratorOption = SequenceGeneratorStartWithOption(0)
	_ SequenceGeneratorOption = SequenceGeneratorIncrementByOption(0)
	_ SequenceGeneratorOption = &SequenceGeneratorMaxValueOption{}
	_ SequenceGeneratorOption = &SequenceGeneratorMinValueOption{}
	_ SequenceGeneratorOption = SequenceGeneratorCycleOption(false)
)

type SequenceGeneratorStartWithOption int

func StartWith(value int) SequenceGeneratorStartWithOption {
	return SequenceGeneratorStartWithOption(value)
}

const kStartWith = Keyword("START WITH")

func (o SequenceGeneratorStartWithOption) sequenceGeneratorOption() SequenceGeneratorOption { return o }
func (o SequenceGeneratorStartWithOption) Accept(v Visitor) Visitor {
	return v.Visit(kStartWith, Int(int(o)))
}
func (o SequenceGeneratorStartWithOption) String() string { return XQL(o) }

type SequenceGeneratorIncrementByOption int

func IncrementBy(value int) SequenceGeneratorIncrementByOption {
	return SequenceGeneratorIncrementByOption(value)
}

const kIncrementBy = Keyword("INCREMENT BY")

func (o SequenceGeneratorIncrementByOption) sequenceGeneratorOption() SequenceGeneratorOption {
	return o
}
func (o SequenceGeneratorIncrementByOption) Accept(v Visitor) Visitor {
	return v.Visit(kIncrementBy, Int(int(o)))
}
func (o SequenceGeneratorIncrementByOption) String() string { return XQL(o) }

type SequenceGeneratorMaxValueOption struct{ Value *int }

func MaxValue(value int) SequenceGeneratorMaxValueOption {
	return SequenceGeneratorMaxValueOption{&value}
}

var NoMaxValue = &SequenceGeneratorMaxValueOption{nil}

const (
	kNoMaxValue = Keyword("NO MAXVALUE")
	kMaxValue   = Keyword("MAXVALUE")
)

func (o SequenceGeneratorMaxValueOption) sequenceGeneratorOption() SequenceGeneratorOption { return o }
func (o SequenceGeneratorMaxValueOption) Accept(v Visitor) Visitor {
	if o.Value == nil {
		return v.Visit(kNoMaxValue)
	}

	return v.Visit(kMaxValue, WS, Int(*o.Value))
}

func (o SequenceGeneratorMaxValueOption) String() string { return XQL(o) }

type SequenceGeneratorMinValueOption struct{ Value *int }

func MinValue(value int) SequenceGeneratorMinValueOption {
	return SequenceGeneratorMinValueOption{&value}
}

var NoMinValue = &SequenceGeneratorMinValueOption{nil}

const (
	kNoMinValue = Keyword("NO MINVALUE")
	kMinValue   = Keyword("MINVALUE")
)

func (o *SequenceGeneratorMinValueOption) sequenceGeneratorOption() SequenceGeneratorOption { return o }
func (o SequenceGeneratorMinValueOption) Accept(v Visitor) Visitor {
	if o.Value == nil {
		return v.Visit(kNoMinValue)
	}

	return v.Visit(kMinValue, WS, Int(*o.Value))
}

func (o SequenceGeneratorMinValueOption) String() string { return XQL(o) }

type SequenceGeneratorCycleOption bool

var (
	Cycle   = SequenceGeneratorCycleOption(true)
	NoCycle = SequenceGeneratorCycleOption(false)
)

const (
	kCycle   = Keyword("CYCLE")
	kNoCycle = Keyword("NO CYCLE")
)

func (o SequenceGeneratorCycleOption) sequenceGeneratorOption() SequenceGeneratorOption { return o }

func (o SequenceGeneratorCycleOption) Accept(v Visitor) Visitor {
	if o {
		return v.Visit(kCycle)
	}

	return v.Visit(kNoCycle)
}

func (o SequenceGeneratorCycleOption) String() string { return XQL(o) }
