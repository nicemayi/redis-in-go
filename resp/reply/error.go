package reply

var unknownerrbytes = []byte("-Err unknown\r\n")

func (u UnknownErrReply) Error() string {
	return "Err unknown"
}

func (u UnknownErrReply) ToBytes() []byte {
	return unknownerrbytes
}

type UnknownErrReply struct{}

func (r *ArgNumErrReply) Error() string {
	return "-ERR wrong number of arguments for '" + r.Cmd + "' command\r\n"
}

func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + r.Cmd + "' command\r\n")
}

type ArgNumErrReply struct {
	Cmd string
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{cmd}
}

type SyntaxErrReply struct{}

var syntaxerrbytes = []byte("-Err syntax error\r\n")
var thesyntaxerrreply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return thesyntaxerrreply
}

func (r *SyntaxErrReply) ToBytes() []byte {
	return syntaxerrbytes
}

func (r *SyntaxErrReply) Error() string {
	return "-Err syntax error"
}

type WrongTypeErrReply struct{}

var wrongtypeerrbytes = []byte("-WRONGTYPE Operation against a key hold the wrong kind of value\r\n")

func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongtypeerrbytes
}

func (r *WrongTypeErrReply) Error() string {
	return "-WRONGTYPE Operation against a key hold the wrong kind of value"
}

type ProtocalErrReply struct {
	Msg string
}

func (r *ProtocalErrReply) ToBytes() []byte {
	return []byte("-ERR Protocal error: '" + r.Msg + "'\r\n")
}

func (r *ProtocalErrReply) Error() string {
	return "ERR Protocal error: '" + r.Msg
}
