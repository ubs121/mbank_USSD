package ussd

import (
	"log"
	"testing"
)

var (
	ctx *Context
)

func testSend(t *testing.T, input string) {
	log.Println("<<", input)

	if retmsg, err := _send(ctx, input); err != nil {
		t.Error(err)
		t.Fail()
	} else {
		log.Println(">>", retmsg)
	}
}

func TestHuulga(t *testing.T) {
	log.Println("************************************")
	ctx = NewContext("112233")

	testSend(t, "") // start
	testSend(t, "1234")
	testSend(t, "1")
	testSend(t, "4")

	if ctx.State != "." {
		t.Fail()
	}
}

func TestSelf(t *testing.T) {
	log.Println("************************************")
	ctx = NewContext("112244")

	testSend(t, "")
	testSend(t, "1234") // pin code
	testSend(t, "1")    // dans1
	testSend(t, "1")    // өөрийн данс хооронд
	testSend(t, "2")    // xuleen avax dans 2
	testSend(t, "100")  // mungun dun
	testSend(t, "1")    // confirm

	if ctx.State != "." {
		t.Fail()
	}
}

func TestGüilgee(t *testing.T) {
	log.Println("************************************")
	ctx = NewContext("112255")

	testSend(t, "")
	testSend(t, "1234") // pin code
	testSend(t, "1")    // dans1
	testSend(t, "2")    // guilgee

	if ctx.State != "." {
		t.Fail()
	}
}
