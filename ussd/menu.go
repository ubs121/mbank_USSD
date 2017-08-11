package ussd

import (
	"errors"
	"regexp"
	"time"
)

const (
	Back          = "0"
	Settings      = "99"
	AccessDenied  = "Handah erhgui!"
	InvalidOption = "Buruu songolt!"
)

// handles user input
func _send(ctx *Context, input string) (string, error) {
	ctx.UpdateTime = time.Now()

	if input == Back {
		ctx.State = "login"
		// TODO: параметрүүд цэвэрлэх?
	} else {

		ctx.Prev = ctx.State

		if len(ctx.keyMap) > 0 && ctx.keyMap[input] != "" {
			input = ctx.keyMap[input]
			ctx.ClearKeyMap()
		} else {
			// normal text input
		}
	}

	msg, err := _process(ctx, input)

	return msg, err
}

// transit from current state to next state
func _process(ctx *Context, input string) (string, error) {

	switch ctx.State {
	case "":
		// TODO: msisdn бүртгэлтэй эсэхийг шалгах, байхгүй бол бүртгэлийн процесс явагдана
		return login(ctx, input)
	case "login":
		return accounts(ctx, input)
	case "accounts":
		if input == "99" {
			return settings(ctx, input)
		} else {
			return actions(ctx, input)
		}
	case "actions":
		switch input {
		case "opt_self":
			ctx.State = "self"
			ctx.Step = 0
			return _process(ctx, input)
		case "opt_trx":
			ctx.State = "trx"
			ctx.Step = 0
			return _process(ctx, input)
		case "opt_pay":
			ctx.State = "pay"
			ctx.Step = 0
			return _process(ctx, input)
		case "opt_stmt":
			return stmt(ctx, input)
		case "opt_topup":
			return topup(ctx, input)
		}
	case "self":
		return self(ctx, input)
	case "trx":
		return msg_with_back("Not implemented!")
	case "settings":
		switch input {
		case "1":
			ctx.State = "pin_change"
			ctx.Step = 0
			return _process(ctx, input)
		case "2":
			ctx.State = "add_acct"
			return choose_acct(ctx, input)
		case "3":
			ctx.State = "remove_acct"
			return choose_acct(ctx, input)
		case "4":
			return ecode(ctx, input)
		case "5":
			return off(ctx, input)
		}

	case "add_acct":
		return add_acct(ctx, input)
	case "remove_acct":
		return remove_acct(ctx, input)
	case "pin_change":
		return pin_change(ctx, input)
	case "off":
		return off(ctx, input)
	}

	return "Aldaa: " + ctx.Error + " \n0. Butsah", nil
}

func login(ctx *Context, input string) (string, error) {
	ctx.State = "login"
	return "COOL bank\nPIN kodoo oruulna uu?", nil
}

func accounts(ctx *Context, input string) (string, error) {
	if !ctx.IsAuthenticated {
		// pin код шалгах TODO: баазаас шалгах
		if input == "1234" {
			ctx.IsAuthenticated = true
			ctx.Params["cif"] = "11111" // TODO: set CIF
		} else {
			ctx.Reset()
			return AccessDenied, errors.New(AccessDenied)
		}
	}

	ctx.State = "accounts"

	// TODO: харилцагчийн цахим дансууд (3 хүртэлх) харуулах, 1 данс бол шууд сонгох

	return "Dansaa songono uu\n1. master card XXXXXX0023\n2. dans1 XXXXXXXX0636\n3. dans2 XXXXXX1142\n----\n99.Tohirgoo", nil
}

func actions(ctx *Context, input string) (string, error) {
	if !ctx.IsAuthenticated {
		return AccessDenied, errors.New(AccessDenied)
	}

	if ok, _ := regexp.MatchString("1|2|3", input); !ok {
		return msg_with_back("Buruu songolt")
	}

	ctx.State = "actions"
	ctx.Params["acct"] = input // сонгосон дансыг хувиргах

	// TODO: check CIF and `acct` are matches
	// TODO: acct дансны balance-g харуулах
	msg := "Uldegdel 1'000'0000.0 MNT"

	ctx.keyMap["1"] = "opt_self"
	ctx.keyMap["2"] = "opt_trx"
	ctx.keyMap["3"] = "opt_pay"
	ctx.keyMap["4"] = "opt_stmt"
	ctx.keyMap["5"] = "opt_topup"

	return msg + `

1. Ooriin dans ruu
2. Guilgee
3. Tolbor
4. Huulga
5. TopUp
-----
0. Butsah`, nil
}

func self(ctx *Context, input string) (string, error) {
	switch ctx.Step {
	case 0:
		ctx.Step++
		return "Huleen avah dans? 1. dans1 2. dans2", nil
	case 1:
		ctx.Step++
		return "Mongon dun?", nil
	case 2:
		ctx.Step++
		return "Batlah uu? \n1. tiim \n2. ugui", nil
	case 3:
		if input == "1" {
			ctx.State = "."
			return "Guilgee amjilttai hiigdlee", nil
		} else {
			ctx.State = "accounts"
			return "Guilgee tsutslagdlaa", nil
		}
	}

	return msg_with_back("Alhamiin aldaa")
}

func stmt(ctx *Context, input string) (string, error) {
	ctx.State = "." // сүүлийн дэлгэц
	// TODO: acct дансны хуулга харуулах
	return "Huulga: 123", nil
}

func topup(ctx *Context, input string) (string, error) {
	ctx.State = "topup"
	ctx.keyMap["1"] = "opt_100"
	ctx.keyMap["2"] = "opt_500"

	return "1. 100 negj \n2. 500 negj", nil
}

func settings(ctx *Context, input string) (string, error) {
	ctx.State = "settings"
	return "1. PIN solih\n2. Dans nemeh\n3. Dans hasah\n4 e-Code avah.\n5. Uilchilgee haah", nil
}

func pin_change(ctx *Context, input string) (string, error) {
	switch ctx.Step {
	case 0:
		ctx.Step++
		return "Huuchin PIN kod?", nil
	case 1:
		ctx.Step++
		return "Shine PIN kod?", nil
	case 2:
		ctx.Step++
		return "Batlah uu? 1. tiim 2. ugui", nil
	case 3:
		if input == "1" {
			ctx.State = "."
			return "PIN soligdloo", nil
		} else {
			return settings(ctx, input)
		}
	}

	return msg_with_back("Alhamiin aldaa")
}

func choose_acct(ctx *Context, input string) (string, error) {
	// TODO: нийт дансыг харуулах
	return `Dansaa songono uu?`, nil
}

func add_acct(ctx *Context, input string) (string, error) {
	// TODO: 3 хүртэлх данс нэмж болно

	return msg_with_back("Dans nemegdlee")
}

func remove_acct(ctx *Context, input string) (string, error) {
	// TODO: 1 данс үлдэх ёстой

	return msg_with_back("Dans xasagdlaa")
}

func ecode(ctx *Context, input string) (string, error) {
	ctx.State = "."

	return "Tanii e-code: 1234", nil
}

func off(ctx *Context, input string) (string, error) {
	ctx.State = "."

	return "Tanii erh haagdlaa", nil
}

func msg_with_back(msg string) (string, error) {
	return msg + "\n----\n0. Butsah ", nil
}
