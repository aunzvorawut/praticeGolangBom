package v6

import "github.com/astaxie/beego"

const (
	SUCCESS_200_MESSAGE_TH  = "สำเร็จ"
	SUCCESS_200_MESSAGE_ENG = "success"
	ERROR_401_MESSAGE_TH    = "การใช้งานไม่ถูกต้อง"
	ERROR_401_MESSAGE_ENG   = "bad argument"
	ERROR_404_MESSAGE_TH    = "ข้อมูลผู้ใช้งานไม่ถูกต้อง"
	ERROR_404_MESSAGE_ENG   = "Invalid User"
	ERROR_410_MESSAGE_TH    = "Please make payment by the due date to ensure<br/>your smooth and continuous service"
	ERROR_410_MESSAGE_ENG   = "Please make payment by the due date to ensure<br/>your smooth and continuous service"
	ERROR_411_MESSAGE_TH    = "You have played 3 free song today, please subscribe package before continue playing."
	ERROR_411_MESSAGE_ENG   = "You have played 3 free song today, please subscribe package before continue playing. User"
	ERROR_412_MESSAGE_TH    = "nopack"
	ERROR_412_MESSAGE_ENG   = "nopack"
	ERROR_413_MESSAGE_TH    = "This song is no longer available."
	ERROR_413_MESSAGE_ENG   = "This song is no longer available."
	ERROR_500_MESSAGE_TH    = "ระบบภายในขัดข้อง"
	ERROR_500_MESSAGE_ENG   = "internal logic server error"
	ERROR_606_MESSAGE_TH    = "ตอนนี้คุณ Login ครบ 5 เครื่องแล้ว หากต้องการใช้งาน กรุณาติดต่อ 1175"
	ERROR_606_MESSAGE_ENG   = "Login limit only 5 device to apply, please contact 1175."
	ERROR_703_MESSAGE_TH    = "ไม่มีรหัสเครื่องในระบบ"
	ERROR_703_MESSAGE_ENG   = "no imei in keep device"
	ERROR_801_MESSAGE_TH    = "ไม่สามารถอัพเดทเครื่องได้"
	ERROR_801_MESSAGE_ENG   = "update device error"

	DEFAULT_MAX                  int = 10
	DEFAULT_OFFSET               int = 0
	DEFAULT_SONG_DETAIL_TIME_OUT     = 300000

	IMAGE_PATH_USER_PROFILE = "static/img/userProfile"

	KEY_JSON_1080P                    = "1080p"
	KEY_JSON_720P                     = "720p"
	KEY_JSON_480P                     = "480p"
	KEY_JSON_360P                     = "360p"
	KEY_JSON_240P                     = "240p"
	KEY_JSON_144P                     = "144p"
	KEY_JSON_auto                     = "auto"
	KEY_JSON_ACCESS_TOKEN      string = "accessToken"
	KEY_JSON_ACTION            string = "action"
	KEY_JSON_ALIAS_NAME               = "aliasName"
	KEY_JSON_APPCONFIG                = "appConfig"
	KEY_JSON_ARTIST                   = "artist"
	KEY_JSON_AVATAR                   = "avatar"
	KEY_JSON_BANNERS                  = "banners"
	KEY_JSON_CATEGORIES               = "categories"
	KEY_JSON_CATEGORY_SHELFS          = "categoryShelfs"
	KEY_JSON_CONTENT_SHELFS           = "contentShelfs"
	KEY_JSON_CONTENT_URL              = "contentUrl"
	KEY_JSON_COVER_URL                = "coverUrl"
	KEY_JSON_DURATION                 = "duration"
	KEY_JSON_ENABLED                  = "enabled"
	KEY_JSON_IS_NEW_USER              = "isNewUser"
	KEY_JSON_IMAGEPATH                = "imagePath"
	KEY_JSON_IMAGES                   = "images"
	KEY_JSON_IMG_URL                  = "imgUrl"
	KEY_JSON_ID                       = "id"
	KEY_JSON_IS_PURCHASED             = "isPurchased"
	KEY_JSON_IS_SUSPEND               = "isSuspend"
	KEY_JSON_IS_FREE                  = "isFree"
	KEY_JSON_ISSHOW                   = "isShow"
	KEY_JSON_ISSHOWONCE               = "isShowOnce"
	KEY_JSON_KARAOKE                  = "karaoke"
	KEY_JSON_LINK                     = "link"
	KEY_JSON_NEWFEATUREVERSION        = "newFeatureVersion"
	KEY_JSON_NUM_OF_LISTENER          = "numOfListener"
	KEY_JSON_NUM_OF_SONG              = "numOfSong"
	KEY_JSON_RECENTLY_SING            = "recentlySing"
	KEY_JSON_SONGS                    = "songs"
	KEY_JSON_STATUS_TAG               = "statusTag"
	KEY_JSON_SUBSCRIBE_PACKAGE        = "subscribedPackage"
	KEY_JSON_SUB_TITLE                = "subTitle"
	KEY_JSON_TERM_VERSION             = "termVersion"
	KEY_JSON_TITLE                    = "title"
	KEY_JSON_TYPE                     = "type"
	KEY_JSON_USER_PROFILE             = "userProfile"
	KEY_JSON_VERSION                  = "version"
	KEY_JSON_VIDEO                    = "video"
	KEY_JSON_VIEW_MORE                = "viewMore"
	KEY_JSON_VOCAL                    = "vocal"

	PARAMS_ACCESS_TOKEN        = "accessToken"
	PARAMS_ACTIONS             = "actions"
	PARAMS_BRAND               = "brand"
	PARAMS_CATEGORY_ID         = "categoryId"
	PARAMS_CLIENT_TYPE         = "clientType"
	PARAMS_CONFIRM_KICK        = "confirmKick"
	PARAMS_DEVICE_ID           = "deviceId"
	PARAMS_Device_TOKEN        = "deviceToken"
	PARAMS_DISPLAY_NAME        = "displayName"
	PARAMS_FFBID_BOX           = "ffbidBox"
	PARAMS_HEIGHT              = "height"
	PARAMS_ID_TYPE             = "idType"
	PARAMS_ID_VALUE            = "idValue"
	PARAMS_IS_CHROME_CAST      = "isChromeCast"
	PARAMS_IS_VOCAL            = "isVocal"
	PARAMS_LANGUAGE            = "language"
	PARAMS_MAX                 = "max"
	PARAMS_MODEL               = "model"
	PARAMS_MSISDN_MOBILE       = "msisdnMobile"
	PARAMS_NETWORK_TYPE        = "networkType"
	PARAMS_OFFSET              = "offset"
	PARAMS_OS_CLIENT           = "osClient"
	PARAMS_OS_VERSION          = "osVersion"
	PARAMS_POPUP_ID            = "popupId"
	PARAMS_PRIVATE_ID          = "privateId"
	PARAMS_PRIVATE_ID_PASSWORD = "privateIdPassword"
	PARAMS_PUBLIC_PLAYLIST_ID  = "publicPlaylistId"
	PARAMS_SECURE_TOKEN        = "secureToken"
	PARAMS_SERIAL              = "serial"
	PARAMS_SONG_ID             = "songId"
	PARAMS_UPLOAD_FILE         = "file"
	PARAMS_WIDTH               = "width"

	SECURE_AES128_KEY      = "tsavmodsiwkaraok"
	SECURE_SECURE_LINK_KEY = "aiskaraokeaiskaraoke"

	VALUE_DISPLAY_GENRE_TITLE_TH      = "เพลยลิสต์"
	VALUE_DISPLAY_GENRE_TITLE_ENG     = "Playlist"
	VALUE_DISPLAY_GENRE_SUB_TITLE_TH  = "เพลงที่คุณสนใจ"
	VALUE_DISPLAY_GENRE_SUB_TITLE_ENG = "Awesome song playlist for you"
	VALUE_DISPLAY_CAT_TITLE_TH        = "ประเภท"
	VALUE_DISPLAY_CAT_TITLE_ENG       = "Type of Song"
	VALUE_DISPLAY_CAT_SUB_TITLE_TH    = "เพลงตามประเภท"
	VALUE_DISPLAY_CAT_SUB_TITLE_ENG   = "Find your song by type"
	VALUE_MODEL_CAT_TYPE              = "songType"
	VALUE_MODEL_PLAYLIST_TYPE         = "Playlist"

	popupTypeSong      = "song"
	popupTypeCat       = "category"
	popupTypeListCat   = "listCategory"
	popupTypePurchased = "purchased"
	popupTypeWeb       = "web"

	QUERY_TRENDING = "Trending"

	CONTITION_WORD_TH  = "th"
	CONTITION_WORD_ENG = "en"
)

func (this *GlobalApi) GetHostApiGrails() string {

	result := ""
	if beego.BConfig.RunMode == "dev" {
		result = "http://localhost:8080"
	} else if beego.BConfig.RunMode == "test" {
		result = "http://aiskaraoke.meevuu.com:8081"
	} else if beego.BConfig.RunMode == "prd" {
		result = ""
	}

	return result

}

func (this *GlobalApi) GetHostApiGolang() string {

	result := ""
	if beego.BConfig.RunMode == "dev" {
		result = "http://localhost:8090"
	} else if beego.BConfig.RunMode == "test" {
		result = "http://staging.aiskaraokego.meevuu.com"
	} else if beego.BConfig.RunMode == "prd" {
		result = ""
	}

	return result

}

func (this *GlobalApi) GetHostStaticGrails() string {

	result := ""
	if beego.BConfig.RunMode == "dev" {
		result = "http://localhost:8090"
	} else if beego.BConfig.RunMode == "test" {
		result = "http://aiskaraoke.meevuu.com:8081"
	} else if beego.BConfig.RunMode == "prd" {
		result = ""
	}

	return result

}

func (this *GlobalApi) GetHostStaticGolang() string {

	result := ""
	if beego.BConfig.RunMode == "dev" {
		result = "http://localhost:8090"
	} else if beego.BConfig.RunMode == "test" {
		result = "http://staging.aiskaraokego.meevuu.com"
	} else if beego.BConfig.RunMode == "prd" {
		result = ""
	}

	return result

}

func (this *GlobalApi) GetStringLanguage(stringTh, stringEng string) string {

	var result string

	if this.GetString(PARAMS_LANGUAGE) == CONTITION_WORD_ENG {
		result = stringEng
	} else {
		result = stringTh
	}
	return result
}
