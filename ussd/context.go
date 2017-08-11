package ussd

import (
	"time"
)

// context
type Context struct {
	ID              string
	IsAuthenticated bool
	CreationTime    time.Time
	UpdateTime      time.Time

	/*
	   Context params:
	   	cust_name - харилцагчийн нэр
	    cif - харилцагчийн ID
	    acct  - сонгогдсон данс
	    acct[1-n] - данснууд
	    acct[1-n].name - дансны нэрнүүд
	    acct[1-n].type - дансны төрлүүд
	    pin - PIN код
	    to_bank - хүлээн авах банк
	    to_acct - хүлээн авах данс
	    to_amt - хүлээн авах мөнгөн дүн
	    to_name - хүлээн авагчийн нэр

	*/
	Params map[string]string
	Prev   string
	State  string
	Step   int
	Error  string

	// options
	keyMap map[string]string
}

func NewContext(id string) *Context {
	ctx := new(Context)
	ctx.ID = id
	ctx.CreationTime = time.Now()
	ctx.UpdateTime = time.Now()
	ctx.Params = make(map[string]string)

	ctx.Prev = ""
	ctx.State = ""
	ctx.keyMap = make(map[string]string)

	return ctx
}

func (ctx *Context) Reset() {
	ctx.State = ""
	ctx.Prev = ""
	ctx.IsAuthenticated = false

	for k := range ctx.Params {
		delete(ctx.Params, k)
	}

}

func (ctx *Context) ClearKeyMap() {
	for k := range ctx.keyMap {
		delete(ctx.keyMap, k)
	}
}
