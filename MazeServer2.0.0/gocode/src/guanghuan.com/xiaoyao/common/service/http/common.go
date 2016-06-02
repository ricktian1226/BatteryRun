package xyhttpservice

type HttpOp int

const (
	HttpAny     HttpOp = iota
	HttpGet     HttpOp = iota
	HttpPost    HttpOp = iota
	HttpPut     HttpOp = iota
	HttpDelete  HttpOp = iota
	HttpPatch   HttpOp = iota
	HttpOptions HttpOp = iota
	HttpHead    HttpOp = iota
)
