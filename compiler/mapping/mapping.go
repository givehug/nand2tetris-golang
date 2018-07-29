package mapping

// arithmetic commands
const (
	ArithmCmdADD = "add"
	ArithmCmdSUB = "sub"
	ArithmCmdNEG = "neg"
	ArithmCmdEQ  = "eq"
	ArithmCmdGT  = "gt"
	ArithmCmdLT  = "lt"
	ArithmCmdAND = "and"
	ArithmCmdOR  = "or"
	ArithmCmdNOT = "not"
)

var ArithmSymbols = map[string]string{
	"+": ArithmCmdADD,
	"-": ArithmCmdSUB,
	"*": "call Math.multiply 2",
	"/": "call Math.divide 2",
	"&": ArithmCmdAND,
	"|": ArithmCmdOR,
	"<": ArithmCmdLT,
	">": ArithmCmdGT,
	"=": ArithmCmdEQ,
}

// push pop segments
const (
	SegmentARG    = "argument"
	SegmentLOCAL  = "local"
	SegmentSTATIC = "static"
	SegmentTHIS   = "this"
	SegmentTHAT   = "that"
	SegmentPOINT  = "pointer"
	SegmentTEMP   = "temp"
	SegmentCONST  = "constant"
)

// identifier type constants
const (
	IdentifierTypeStatic = "static"
	IdentifierTypeField  = "field"
	IdentifierTypeArg    = "argument"
	IdentifierTypeVar    = "local"
)
