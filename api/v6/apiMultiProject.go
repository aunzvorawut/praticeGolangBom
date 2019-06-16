package v6

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var fileSeparator = string(filepath.Separator)

const (
	letterBytesnumber          = "1234567890"
	letterBytescharBig         = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterBytesCharAndNumBig   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterBytescharSmall       = "abcdefghijklmnopqrstuvwxyz"
	letterBytesCharAndNumSmall = "abcdefghijklmnopqrstuvwxyz1234567890"
	letterBytescharMix         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	letterBytesCharAndNumMix   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits
)

func Int64ToString(a int64) string {
	return strconv.FormatInt(a, 10)
}

func RandStringBytesMaskImprSrc(formatLetter string, n int) string {

	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(formatLetter) {
			b[i] = formatLetter[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func tryCatch() {

	fmt.Println("We started")
	Block{
		Try: func() {
			fmt.Println("I tried")
			ThrowTest("Oh,...sh...")
		},
		Catch: func(e Exception) {
			fmt.Printf("Caught %v\n", e)
		},
		Finally: func() {
		},
	}.Do()

}

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

func (tcf Block) Do() {
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

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

// StrPad returns the input string padded on the left, right or both sides using padType to the specified padding length padLength.
//
// Example:
// input := "Codes";
// StrPad(input, 10, " ", "RIGHT")        // produces "Codes     "
// StrPad(input, 10, "-=", "LEFT")        // produces "=-=-=Codes"
// StrPad(input, 10, "_", "BOTH")         // produces "__Codes___"
// StrPad(input, 6, "___", "RIGHT")       // produces "Codes_"
// StrPad(input, 3, "*", "RIGHT")         // produces "Codes"
func StrPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

func GetSafeFileName(input string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9_]+")
	if err != nil {
		beego.Error("")
	}

	safe := reg.ReplaceAllString(input, "-")
	safe = strings.ToLower(strings.Trim(safe, "-"))
	return safe
	//fmt.Println(safe)   // Output: a*-+fe5v9034,j*.ae6
}

func GetFileExtension(fileName string) string {
	var extension = filepath.Ext(fileName)
	return extension
}

func UploadProfileImageGlobal(id int64, file multipart.File, handler *multipart.FileHeader, err error, imagePathPa, filename string) (bool, string) {

	imgPath := imagePathPa                                // ทีี่เก็บรูป
	imgUrl := fileSeparator + imagePathPa + fileSeparator // ที่ลง db

	err2 := os.MkdirAll("./"+imgPath, 0775)
	if err2 != nil {
		beego.Error("UploadProfileImage: ", err2)
		return false, "Can't create directory"
	}

	defer file.Close()

	//====== start case upload หรือ upload sftp file ======

	imageLocation := "./" + imgPath + fileSeparator + filename

	f, err := os.OpenFile(imageLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		beego.Error("UploadProfileImage: ", err.Error())
		return false, "Can't Open File"
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		beego.Error("UploadProfileImage: ", err.Error())
		return false, "Can't write file"
	}

	if err != nil {
		beego.Error("UploadProfileImage: ", err.Error())
		return false, "Internal Server Error"
	}

	return true, imgUrl + filename

}

func TimeStampToTime(timeStamp string) time.Time {
	i, err := strconv.ParseInt(timeStamp, 10, 64)
	if err != nil {
		beego.Debug(err)
		return time.Now()
	}
	tm := time.Unix(i, 0)
	return tm
}

func SubstringWithIndex(sourceString, char string) (string, string) { //
	// sourceString = 12@34
	// char = @
	// result 12 , 34

	i := strings.Index(sourceString, char)
	if i > -1 {
		front := sourceString[:i]
		back := sourceString[i+1:]
		return front, back
	} else {
		return sourceString, sourceString
	}

}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//======== aes128 ECB ===========

func AesDecrypt(crypted, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("err is:", err)
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	fmt.Println("source is :", origData, string(origData))
	return origData
}

func AesEncrypt(src, key string) []byte {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return crypted
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}