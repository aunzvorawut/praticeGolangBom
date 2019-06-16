package v6

import (
	"github.com/astaxie/beego"
)

type BlockTest struct {
	Try     func()
	Catch   func(ExceptionTest)
	Finally func()
}

type ExceptionTest interface{}

func ThrowTest(up ExceptionTest) {
	panic(up)
}

func (this *GlobalApi) Test() {

	//======== test aes encrypt ============

	text := "aunz"

	encResult := AesEncrypt(text, SECURE_AES128_KEY)
	decResult := AesDecrypt(encResult, []byte(SECURE_AES128_KEY))



	beego.Error("encResult = ",string(encResult))
	beego.Error("decResult = ",string(decResult))

	//====== test substring =========

	//stringExample := "1234567890abcdefg"
	//
	//beego.Debug(stringExample)
	//
	//beego.Debug(stringExample[len(stringExample)-13:])

	// ======= test migrate db =========

	//today := time.Now()
	//beego.Debug(today)
	//
	//songObj, err := models.GetSongById(1)
	//categoryObj, err := models.GetCategoryById(1)
	//
	//beego.Debug(songObj)
	//beego.Debug(err)
	//beego.Debug(categoryObj)
	//beego.Debug(err)
	//
	//newSongCat := models.SongCategoryPosition{
	//	Song:     songObj,
	//	Category: categoryObj,
	//}
	//
	//id, err := models.AddSongCategoryPosition(&newSongCat)
	//beego.Debug("id = ", id)
	//beego.Debug("err = ", err)
	//
	//beego.Debug("===")

	//=========test try catch===========

	//fmt.Println("We started")
	//BlockTest{
	//	Try: func() {
	//		fmt.Println("I tried")
	//
	//		userObj, _ := models.GetSecUserById(1935701750371)
	//		beego.Debug(userObj)
	//		beego.Debug(userObj.Enabled)
	//
	//		ThrowTest("Oh,...sh...")
	//	},
	//	Catch: func(e ExceptionTest) {
	//		fmt.Printf("Caught %v\n", e)
	//	},
	//	Finally: func() {
	//		fmt.Println("Finally...")
	//	},
	//}.Do()
	//fmt.Println("We went on")

	this.ResponseJSON(nil, 200, "ssdd")
	return
}

func (tcf BlockTest) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}
