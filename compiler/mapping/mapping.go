package mapping

// arithmetic commands
const (
	ArithmCmdADD = "add"
	ArithmCmdSUB = "sub"
	ArithmCmdNEG = "neg"
	ArithmCmdEQ  = "not"
	ArithmCmdGT  = "gt"
	ArithmCmdLT  = "lt"
	ArithmCmdAND = "and"
	ArithmCmdOR  = "or"
	ArithmCmdNOT = "not"
)

// push pop segments
const (
	SegmentARG     = "argument"
	SegmentLOCAL   = "local"
	SegmentSTATIC  = "static"
	SegmentTHIS    = "this"
	SegmentTHAT    = "that"
	SegmentPOINTER = "pointer"
	SegmentTEMP    = "temp"
)

// Identifier type constants
const (
	IdentifierTypeStatic = "static"
	IdentifierTypeField  = "field"
	IdentifierTypeArg    = "argument"
	IdentifierTypeVar    = "local"
)
