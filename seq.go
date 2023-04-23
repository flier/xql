package xql

import (
	"fmt"
)

type SequenceGeneratorOption interface {
	fmt.Stringer

	sequenceGeneratorOption() SequenceGeneratorOption
}

var (
	_ SequenceGeneratorOption = SequenceGeneratorStartWithOption(0)
	_ SequenceGeneratorOption = SequenceGeneratorIncrementByOption(0)
	_ SequenceGeneratorOption = &SequenceGeneratorMaxValueOption{}
	_ SequenceGeneratorOption = &SequenceGeneratorMinValueOption{}
	_ SequenceGeneratorOption = SequenceGeneratorCycleOption(false)
)

type SequenceGeneratorStartWithOption int64

func StartWith(value int64) SequenceGeneratorStartWithOption {
	return SequenceGeneratorStartWithOption(value)
}

func (o SequenceGeneratorStartWithOption) sequenceGeneratorOption() SequenceGeneratorOption { return o }
func (o SequenceGeneratorStartWithOption) String() string                                   { return fmt.Sprintf("START WITH %d", o) }

type SequenceGeneratorIncrementByOption int64

func IncrementBy(value int64) SequenceGeneratorIncrementByOption {
	return SequenceGeneratorIncrementByOption(value)
}

func (o SequenceGeneratorIncrementByOption) sequenceGeneratorOption() SequenceGeneratorOption {
	return o
}
func (o SequenceGeneratorIncrementByOption) String() string { return fmt.Sprintf("INCREMENT BY %d", o) }

type SequenceGeneratorMaxValueOption struct{ Value *int64 }

func MaxValue(value int64) SequenceGeneratorMaxValueOption {
	return SequenceGeneratorMaxValueOption{&value}
}

var NoMaxValue = &SequenceGeneratorMaxValueOption{nil}

func (o SequenceGeneratorMaxValueOption) sequenceGeneratorOption() SequenceGeneratorOption { return o }
func (o SequenceGeneratorMaxValueOption) String() string {
	if o.Value == nil {
		return "NO MAXVALUE"
	}

	return fmt.Sprintf("MAXVALUE %d", *o.Value)
}

type SequenceGeneratorMinValueOption struct{ Value *int64 }

func MinValue(value int64) SequenceGeneratorMinValueOption {
	return SequenceGeneratorMinValueOption{&value}
}

var NoMinValue = &SequenceGeneratorMinValueOption{nil}

func (o *SequenceGeneratorMinValueOption) sequenceGeneratorOption() SequenceGeneratorOption { return o }
func (o SequenceGeneratorMinValueOption) String() string {
	if o.Value == nil {
		return "NO MINVALUE"
	}

	return fmt.Sprintf("MINVALUE %d", *o.Value)
}

type SequenceGeneratorCycleOption bool

var (
	Cycle   = SequenceGeneratorCycleOption(true)
	NoCycle = SequenceGeneratorCycleOption(false)
)

func (o SequenceGeneratorCycleOption) sequenceGeneratorOption() SequenceGeneratorOption { return o }

func (o SequenceGeneratorCycleOption) String() string {
	if o {
		return "CYCLE"
	}

	return "NO CYCLE"
}
