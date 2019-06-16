package v6

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"gitlab.com/wisdomvast/AiskaraokeServerGolang/models"
	"net/url"
	"time"
)

func (this *GlobalApi) AppConfig() {

	configScreenObj := models.GetconfigScreenObjBylogicSplashScreen()

	newFeatureVersion := 1
	if configScreenObj != nil {
		newFeatureVersion = configScreenObj.VersionDisplay
	}

	result := map[string]interface{}{
		KEY_JSON_APPCONFIG: map[string]interface{}{
			KEY_JSON_NEWFEATUREVERSION: newFeatureVersion,
			KEY_JSON_IMAGEPATH:         "",
			KEY_JSON_TERM_VERSION:      1,
		},
	}
	this.ResponseJSON(result, 200, "message")
	return

}

func (this *GlobalApi) NewFeatureImage() {

	configScreenObj := models.GetconfigScreenObjBylogicSplashScreen()
	resultImages := this.RenderImagesOnSplashScreenByConfigScreen(configScreenObj)

	var version int
	isShowOnce := true

	if configScreenObj != nil {
		version = configScreenObj.VersionDisplay
		isShowOnce = configScreenObj.IsShowOnce
	}

	result := map[string]interface{}{
		KEY_JSON_VERSION:    version,
		KEY_JSON_ISSHOWONCE: isShowOnce,
		KEY_JSON_IMAGES:     resultImages,
	}

	this.ResponseJSON(result, 200, "")
	return

}

func (this *GlobalApi) ApiPopup() {

	popupObj := models.GetConfigPopupByStartEndAndLastUpdateAndEnabled()

	popupIdStr := ""
	link := ""
	isShow := false
	if popupObj != nil {
		popupIdStr = Int64ToString(popupObj.Id)
		link = this.GetHostApiGolang() + "/v6/content/popup/view?" + PARAMS_POPUP_ID + "=" + popupIdStr
		isShow = true
	}
	result := map[string]interface{}{
		KEY_JSON_LINK:       link,
		KEY_JSON_ISSHOW:     isShow,
		KEY_JSON_ISSHOWONCE: popupObj.IsShowOnce,
	}
	this.ResponseJSON(result, 200, "message")
	return

}

func (this *GlobalApi) HomePage() {

	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	clientTypeFunc := this.CheckClientType(clientType)

	//shelf 1

	listBanner := models.GetAllBannerHomePage(-1, 0, clientTypeFunc)
	resultBanner := make([]map[string]interface{}, 0)

	for _, v := range listBanner {
		renderBanner := this.RenderBannerOnHome(v)
		resultBanner = append(resultBanner, renderBanner)
	}

	// shelf 2

	resultFreeSong := make([]map[string]interface{}, 0)
	categoryShelf2Obj := models.GetCategoryFreeRandom()
	var categoryShelf2Id int64 = 0
	if categoryShelf2Obj != nil {
		listSongId := models.GetAllSongIdByCatId(10, 0, categoryShelf2Obj.Id, clientTypeFunc)
		for _, v := range listSongId {
			freeSongObj, _ := models.GetSongById(v)
			renderFreeSong := this.RenderSongOnHome(freeSongObj)
			resultFreeSong = append(resultFreeSong, renderFreeSong)
		}
		categoryShelf2Id = categoryShelf2Obj.Id
	}

	// shelf 3

	resultTrendingSong := make([]map[string]interface{}, 0)
	categoryShelf3Obj := models.GetCategoryTrendingRandom()
	var categoryShelf3Id int64 = 0
	if categoryShelf3Obj != nil {
		listSongId := models.GetAllSongIdByCatId(10, 0, categoryShelf3Obj.Id, clientTypeFunc)
		for _, v := range listSongId {
			freeSongObj, _ := models.GetSongById(v)
			renderSong := this.RenderSongOnHome(freeSongObj)
			resultTrendingSong = append(resultTrendingSong, renderSong)
		}
		categoryShelf3Id = categoryShelf3Obj.Id
	}

	// shelf 4
	resultCategories := make([]map[string]interface{}, 0)
	allCategory := models.GetAllCategoryByClientType(-1, 0, clientTypeFunc)

	for _, catObj := range allCategory {
		renderCategory := this.RenderCatOnHome(catObj)
		resultCategories = append(resultCategories, renderCategory)
	}

	// final result

	var resultShelf []map[string]interface{}

	if len(resultBanner) != 0 {
		resultShelf = append(resultShelf, map[string]interface{}{
			KEY_JSON_BANNERS: resultBanner,
			KEY_JSON_TYPE:    "banner",
		})
	}

	if len(resultFreeSong) != 0 {
		resultShelf = append(resultShelf, map[string]interface{}{
			KEY_JSON_SONGS:     resultFreeSong,
			KEY_JSON_TITLE:     "Free Karaoke",
			KEY_JSON_SUB_TITLE: "Free Karaoke",
			KEY_JSON_TYPE:      "songHorizontalScroll",
			KEY_JSON_VIEW_MORE: this.RenderViewMore(popupTypeCat, Int64ToString(categoryShelf2Id)),
		})
	}

	if len(resultTrendingSong) != 0 {
		resultShelf = append(resultShelf, map[string]interface{}{
			KEY_JSON_SONGS:     resultTrendingSong,
			KEY_JSON_TITLE:     "Trending Karoke",
			KEY_JSON_SUB_TITLE: "Trending Karoke",
			KEY_JSON_TYPE:      "songRanking",
			KEY_JSON_VIEW_MORE: this.RenderViewMore(popupTypeCat, Int64ToString(categoryShelf3Id)),
		})
	}

	if len(resultCategories) != 0 {
		resultShelf = append(resultShelf, map[string]interface{}{
			KEY_JSON_CATEGORIES: resultCategories,
			KEY_JSON_TITLE:      "Genre Hits",
			KEY_JSON_SUB_TITLE:  "Genre Hits",
			KEY_JSON_TYPE:       "categoryGrid",
			KEY_JSON_VIEW_MORE:  this.RenderViewMore(popupTypeListCat, "noContentRef"),
		})
	}

	result := map[string]interface{}{
		KEY_JSON_CONTENT_SHELFS: resultShelf,
	}
	this.ResponseJSON(result, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}
func (this *GlobalApi) ListCategoryPage() {

	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	clientTypeFunc := this.CheckClientType(clientType)

	// shelf Cat
	resultCategories := make([]map[string]interface{}, 0)
	allCategory := models.GetAllCategoryByClientType(-1, 0, clientTypeFunc)

	for _, catObj := range allCategory {
		renderCategory := this.RenderCatOnHome(catObj)
		resultCategories = append(resultCategories, renderCategory)
	}

	// shelf Genre
	resultGenres := make([]map[string]interface{}, 0)
	allGenre := models.GetAllGenreByClientType(-1, 0, clientTypeFunc)

	for _, genreObj := range allGenre {
		resultGenre := this.RenderGenreOnHome(genreObj)
		resultGenres = append(resultGenres, resultGenre)
	}

	result := map[string]interface{}{
		KEY_JSON_CATEGORY_SHELFS: []map[string]interface{}{
			map[string]interface{}{
				KEY_JSON_CATEGORIES: resultGenres,
				KEY_JSON_TITLE:      this.GetStringLanguage(VALUE_DISPLAY_GENRE_TITLE_TH, VALUE_DISPLAY_GENRE_TITLE_ENG),
				KEY_JSON_SUB_TITLE:  this.GetStringLanguage(VALUE_DISPLAY_GENRE_SUB_TITLE_TH, VALUE_DISPLAY_GENRE_SUB_TITLE_ENG),
			},
			map[string]interface{}{
				KEY_JSON_CATEGORIES: resultCategories,
				KEY_JSON_TITLE:      this.GetStringLanguage(VALUE_DISPLAY_CAT_TITLE_TH, VALUE_DISPLAY_CAT_TITLE_ENG),
				KEY_JSON_SUB_TITLE:  this.GetStringLanguage(VALUE_DISPLAY_CAT_SUB_TITLE_TH, VALUE_DISPLAY_CAT_SUB_TITLE_ENG),
			},
		},
	}

	this.ResponseJSON(result, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}

func (this *GlobalApi) CategoryDetailPage() {

	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	clientTypeFunc := this.CheckClientType(clientType)

	categoryId, _ := this.GetInt64(PARAMS_CATEGORY_ID)
	max, _ := this.GetInt(PARAMS_MAX, DEFAULT_MAX)
	offset, _ := this.GetInt(PARAMS_OFFSET, DEFAULT_OFFSET)

	categoryObj, _ := models.GetCategoryById(categoryId)

	if categoryObj == nil {
		this.ResponseJSON(nil, 401, this.GetStringLanguage(ERROR_401_MESSAGE_TH, ERROR_401_MESSAGE_ENG))
		return
	}

	listSongId := models.GetAllSongIdByCatId(max, offset, categoryId, clientTypeFunc)
	resultSongs := make([]map[string]interface{}, 0)
	for _, v := range listSongId {

		songObj, _ := models.GetSongById(v)
		resultSong := this.RenderSongOnHome(songObj)

		resultSongs = append(resultSongs, resultSong)

	}

	result := map[string]interface{}{
		KEY_JSON_COVER_URL:       this.GetHostStaticGrails() + "/" + categoryObj.Cover,
		KEY_JSON_ID:              categoryObj.Id,
		KEY_JSON_NUM_OF_LISTENER: Int64ToString(this.GetCountAllSongInCat(categoryObj, clientTypeFunc)),
		KEY_JSON_NUM_OF_SONG:     models.GetCountSongIdByCatId(categoryObj.Id, clientTypeFunc)[0],
		KEY_JSON_STATUS_TAG:      0,
		KEY_JSON_TITLE:           this.GetStringLanguage(categoryObj.CategoryName, categoryObj.CategoryNameEng),
		KEY_JSON_SONGS:           resultSongs,
	}

	this.ResponseJSON(result, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}

func (this *GlobalApi) PublicPlaylistDetail() {

	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	clientTypeFunc := this.CheckClientType(clientType)

	genreId, _ := this.GetInt64(PARAMS_PUBLIC_PLAYLIST_ID)
	max, _ := this.GetInt(PARAMS_MAX, DEFAULT_MAX)
	offset, _ := this.GetInt(PARAMS_OFFSET, DEFAULT_OFFSET)

	genreObj, _ := models.GetGenreById(genreId)

	if genreObj == nil {
		this.ResponseJSON(nil, 401, this.GetStringLanguage(ERROR_401_MESSAGE_TH, ERROR_401_MESSAGE_ENG))
		return
	}

	listSongId := models.GetAllSongIdByGenreId(max, offset, genreId, clientTypeFunc)
	resultSongs := make([]map[string]interface{}, 0)
	for _, v := range listSongId {

		songObj, _ := models.GetSongById(v)
		resultSong := this.RenderSongOnHome(songObj)

		resultSongs = append(resultSongs, resultSong)

	}

	result := map[string]interface{}{
		KEY_JSON_COVER_URL:       this.GetHostStaticGrails() + "/" + genreObj.CoverImage,
		KEY_JSON_ID:              genreObj.Id,
		KEY_JSON_NUM_OF_LISTENER: Int64ToString(this.GetCountAllSongInGenre(genreObj, clientTypeFunc)),
		KEY_JSON_NUM_OF_SONG:     models.GetCountSongIdByGenreId(genreObj.Id, clientTypeFunc)[0],
		KEY_JSON_STATUS_TAG:      0,
		KEY_JSON_TITLE:           this.GetStringLanguage(genreObj.NameTh, genreObj.NameEng),
		KEY_JSON_SONGS:           resultSongs,
	}

	this.ResponseJSON(result, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}

func (this *GlobalApi) UserProfile() {

	accessToken := this.GetString(PARAMS_ACCESS_TOKEN)
	userObj := this.GetUserByAccessToken(accessToken)
	if userObj == nil {
		this.ResponseJSON(nil, 404, this.GetString(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
		return
	}

	result := this.RenderUserProfile(userObj)
	this.ResponseJSON(result, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}

func (this *GlobalApi) UserRecently() {

	accessToken := this.GetString(PARAMS_ACCESS_TOKEN)
	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	clientTypeFunc := this.CheckClientType(clientType)
	userObj := this.GetUserByAccessToken(accessToken)
	if userObj == nil {
		this.ResponseJSON(nil, 404, this.GetString(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
		return
	}

	songIdsList := models.GetAllRecentlySongIdByUserObj(50, 0, userObj, clientTypeFunc)
	resultRecentlySongs := make([]map[string]interface{}, 0)

	for _, v := range songIdsList {
		songObj, _ := models.GetSongById(v)
		resultRecentSong := this.RenderSongOnHome(songObj)
		resultRecentlySongs = append(resultRecentlySongs, resultRecentSong)
	}

	this.ResponseJSON(map[string]interface{}{
		KEY_JSON_RECENTLY_SING: resultRecentlySongs,
	}, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}

func (this *GlobalApi) AllProfile() {

	accessToken := this.GetString(PARAMS_ACCESS_TOKEN)
	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	clientTypeFunc := this.CheckClientType(clientType)
	userObj := this.GetUserByAccessToken(accessToken)
	if userObj == nil {
		this.ResponseJSON(nil, 404, this.GetString(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
		return
	}

	resultProfile := this.RenderUserProfile(userObj)
	songIdsList := models.GetAllRecentlySongIdByUserObj(10, 0, userObj, clientTypeFunc)
	resultRecentlySongs := make([]map[string]interface{}, 0)

	for _, v := range songIdsList {
		songObj, _ := models.GetSongById(v)
		resultRecentSong := this.RenderSongOnHome(songObj)
		resultRecentlySongs = append(resultRecentlySongs, resultRecentSong)
	}

	result := map[string]interface{}{
		KEY_JSON_USER_PROFILE:  resultProfile,
		KEY_JSON_RECENTLY_SING: resultRecentlySongs,
	}

	this.ResponseJSON(result, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}

func (this *GlobalApi) SongDetail() {
	songId, _ := this.GetInt64(PARAMS_SONG_ID)
	songObj, _ := models.GetSongById(songId)
	today := time.Now()

	accessToken := this.GetString(PARAMS_ACCESS_TOKEN)
	secureToken := this.GetString(PARAMS_SECURE_TOKEN)

	Block{
		Try: func() {

			secureToken, err := url.QueryUnescape(secureToken)
			if err != nil {
				this.ResponseJSON(nil, 500, this.GetString(ERROR_500_MESSAGE_TH, ERROR_500_MESSAGE_ENG))
				return
			}

			secureTokenByte, _ := base64.StdEncoding.DecodeString(secureToken)

			secureTokenDecryptByte := AesDecrypt(secureTokenByte, []byte(SECURE_AES128_KEY))
			beego.Debug(string(secureTokenDecryptByte))
			secureTokenDecrypt := string(secureTokenDecryptByte)
			secureTokenTocompare := secureTokenDecrypt[0 : len(secureTokenDecrypt)-13]
			timeStamp13digit := secureTokenDecrypt[len(secureTokenDecrypt)-13:]

			timeStamp10digit := timeStamp13digit[0:10]

			expireTokentime := TimeStampToTime(timeStamp10digit)

			if (expireTokentime.Unix()*1000-today.Unix()*1000 >= DEFAULT_SONG_DETAIL_TIME_OUT || today.Unix()*1000-expireTokentime.Unix()*1000 >= DEFAULT_SONG_DETAIL_TIME_OUT) ||
				secureTokenTocompare != accessToken {

				beego.Debug("expireTokentime.Unix()*1000 = ", expireTokentime.Unix()*1000)
				beego.Debug("today.Unix()*1000 = ", today.Unix()*1000)
				beego.Debug("secureTokenTocompare = ", secureTokenTocompare)
				beego.Debug("accessToken = ", accessToken)

				this.ResponseJSON(nil, 404, this.GetString(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
				return
			}

		},
		Catch: func(e Exception) {
			beego.Error("e = ", e)
			this.ResponseJSON(nil, 404, this.GetString(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
			return
		},
		Finally: func() {

			isChromeCast := this.GetString(PARAMS_IS_CHROME_CAST)
			userTokenObj := models.GetUserByAccessToken(accessToken)

			var userObj *models.SecUser

			if userTokenObj == nil {
				this.ResponseJSON(nil, 404, this.GetString(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
				return
			}

			userObj = userTokenObj.SecUser

			clientType := this.GetString(PARAMS_CLIENT_TYPE)
			newToday := today
			deviceId := this.GetString(PARAMS_DEVICE_ID)

			newTemporarySongDetailLog := models.TemporarySongDetailLog{
				DeviceId:    deviceId,
				ClientType:  clientType,
				AccessToken: accessToken,
				SongId:      Int64ToString(songId),
			}

			_, err := models.AddTemporarySongDetailLog(&newTemporarySongDetailLog)
			if err != nil {
				beego.Error("err = ", err)
			}

			deviceObj := models.GetDeviceByImeiAndClientType(deviceId, clientType)
			if deviceObj == nil && isChromeCast != "1" {
				beego.Error("deviceObj == nil && isChromeCast != 1")
				this.ResponseJSON(nil, 500, this.GetString(ERROR_500_MESSAGE_TH, ERROR_500_MESSAGE_ENG))
				return
			}

			dayInt := int(today.Weekday())
			var dayStr string

			if dayInt == 1 {
				dayStr = "monday"
			} else if dayInt == 2 {
				dayStr = "tuesday"
			} else if dayInt == 3 {
				dayStr = "wednesday"
			} else if dayInt == 4 {
				dayStr = "thursday"
			} else if dayInt == 5 {
				dayStr = "friday"
			} else if dayInt == 6 {
				dayStr = "saturday"
			} else {
				dayStr = "sunday"
			}

			timeCheck := today.Format("15:04")
			isFree := models.GetDetailFreeByStartDateLessAndEndDateGreaterAndDayname(timeCheck, dayStr)

			if songObj == nil {
				beego.Error("songObj = nil")
				this.ResponseJSON(nil, 500, this.GetString(ERROR_500_MESSAGE_TH, ERROR_500_MESSAGE_ENG))
				return
			}

			if userObj == nil {
				userObj = models.GetSecUserByUsername("wiscontent@gmail.com")
			}
			isPurchased := false
			isSuspend := false
			havePack := false

			buffetObjListSuspend := models.GetAllBuffetUserByUserObjAndDateExpiredGreaterAndFailCase(userObj.Id, newToday)

			for _, itBuffetListSuspend := range buffetObjListSuspend {
				if models.GetBuffetSongPositionByBuffetAndSong(itBuffetListSuspend.Buffet.Id, songObj.Id) != nil {
					isPurchased = true
					isSuspend = true
				}
			}

			buffetObjList := models.GetAllBuffetUserByUserObjAndDateExpiredGreaterAndSuccessCase(userObj.Id, newToday)

			for _, itBuffetList := range buffetObjList {
				if models.GetBuffetSongPositionByBuffetAndSong(itBuffetList.Buffet.Id, songObj.Id) != nil {
					isPurchased = true
					isSuspend = false
					havePack = true
				}
			}

			categoryFreeList := models.GetAllCategoryByisFree(-1, 0, true)

			for _, itCatFree := range categoryFreeList {

				if models.GetSongCategoryPositionBySongIdAndCatId(songObj.Id, itCatFree.Id) != nil &&
					songObj.IsFree && userObj.CountSong < 3 {
					isPurchased = true
					break
				} else if havePack {
					isPurchased = true
				} else {
					isPurchased = false
				}
			}

			if isFree != nil {
				isPurchased = true
				isSuspend = false
			}

			if models.GetWhiteListDataByPrivateIdAndEnabled(userObj.Username, true) != nil {
				isPurchased = true
				isSuspend = false
			}

			result := this.RenderSongOnHome(songObj)
			result[KEY_JSON_DURATION] = songObj.Duration
			result[KEY_JSON_VIDEO] = this.RenderResultVideoBySongObj(songObj)
			result[KEY_JSON_IS_PURCHASED] = isPurchased
			result[KEY_JSON_IS_SUSPEND] = isSuspend
			result[KEY_JSON_IS_FREE] = isFree

			if songObj.Enabled == false || (songObj.StartDate.Unix() >= today.Unix() && songObj.ExpiredDate.Unix() <= today.Unix()) {
				this.ResponseJSON(nil, 413, this.GetString(ERROR_413_MESSAGE_TH, ERROR_413_MESSAGE_ENG))
				return

			} else if isSuspend == true {
				this.ResponseJSON(nil, 410, this.GetString(ERROR_410_MESSAGE_TH, ERROR_410_MESSAGE_ENG))
				return

			} else if !isPurchased && songObj.IsFree == false {
				this.ResponseJSON(nil, 411, this.GetString(ERROR_411_MESSAGE_TH, ERROR_411_MESSAGE_ENG))
				return

			} else if !isPurchased {
				this.ResponseJSON(nil, 412, this.GetString(ERROR_412_MESSAGE_TH, ERROR_412_MESSAGE_ENG))
				return

			} else if isPurchased {

				this.ResponseJSON(result, 200, this.GetString(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
				return

			}
		},
	}.Do()

}

//=================== Module ====================

func (this *GlobalApi) RenderBannerOnHome(bannerObj *models.Banner) map[string]interface{} {
	result := map[string]interface{}{
		KEY_JSON_ID:          bannerObj.Id,
		KEY_JSON_ACTION:      bannerObj.NextTo,
		KEY_JSON_CONTENT_URL: bannerObj.ContentUrl,
		KEY_JSON_IMG_URL:     this.GetHostStaticGrails() + "/" + bannerObj.BannerCoverImage,
	}
	return result
}

func (this *GlobalApi) RenderSongOnHome(songObj *models.Song) map[string]interface{} {

	result := map[string]interface{}{
		KEY_JSON_ARTIST:          this.GetStringLanguage(songObj.ArtistName, songObj.ArtistNameEng),
		KEY_JSON_TITLE:           this.GetStringLanguage(songObj.SongName, songObj.SongNameEng),
		KEY_JSON_COVER_URL:       this.GetHostStaticGrails() + "/" + songObj.CoverPath,
		KEY_JSON_NUM_OF_LISTENER: "0k",
		KEY_JSON_STATUS_TAG:      0,
		KEY_JSON_ID:              songObj.Id,
	}
	return result
}

func (this *GlobalApi) RenderCatOnHome(catObj *models.Category) map[string]interface{} {

	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	clientTypeFunc := this.CheckClientType(clientType)

	songIds := models.GetAllSongIdByCatId(-1, 0, catObj.Id, clientTypeFunc)

	result := map[string]interface{}{
		KEY_JSON_TITLE:           this.GetStringLanguage(catObj.CategoryName, catObj.CategoryNameEng),
		KEY_JSON_COVER_URL:       this.GetHostStaticGrails() + "/" + catObj.Cover,
		KEY_JSON_NUM_OF_LISTENER: "0k",
		KEY_JSON_NUM_OF_SONG:     len(songIds),
		KEY_JSON_STATUS_TAG:      0,
		KEY_JSON_ID:              catObj.Id,
		KEY_JSON_TYPE:            VALUE_MODEL_CAT_TYPE,
	}
	return result
}

func (this *GlobalApi) RenderGenreOnHome(genreObj *models.Genre) map[string]interface{} {

	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	clientTypeFunc := this.CheckClientType(clientType)

	songIds := models.GetAllSongIdByGenreId(-1, 0, genreObj.Id, clientTypeFunc)

	result := map[string]interface{}{
		KEY_JSON_TITLE:           this.GetStringLanguage(genreObj.NameTh, genreObj.NameEng),
		KEY_JSON_COVER_URL:       this.GetHostStaticGrails() + "/" + genreObj.CoverImage,
		KEY_JSON_NUM_OF_LISTENER: "0k",
		KEY_JSON_NUM_OF_SONG:     len(songIds),
		KEY_JSON_STATUS_TAG:      0,
		KEY_JSON_ID:              genreObj.Id,
		KEY_JSON_TYPE:            VALUE_MODEL_PLAYLIST_TYPE,
	}
	return result
}

func (this *GlobalApi) RenderViewMore(pageType, contentRef string) map[string]interface{} {

	enabled := true

	if contentRef == "0" || contentRef == "" {
		enabled = false
	}

	result := map[string]interface{}{
		KEY_JSON_ACTION:      pageType,
		KEY_JSON_CONTENT_URL: contentRef,
		KEY_JSON_ENABLED:     enabled,
	}
	return result
}

func (this *GlobalApi) RenderUserProfile(userObj *models.SecUser) map[string]interface{} {

	result := map[string]interface{}{
		KEY_JSON_AVATAR:            this.GetUserImage(userObj),
		KEY_JSON_ALIAS_NAME:        this.GetUserNickname(userObj),
		KEY_JSON_SUBSCRIBE_PACKAGE: this.GetUserPackageName(userObj),
	}

	return result

}

func (this *GlobalApi) GetUserNickname(userObj *models.SecUser) string {

	if userObj.IsFacebook {
		return userObj.NickNameSocialFacebook
	} else {
		if userObj.NickNameSocial != "" {
			return userObj.NickNameSocial
		} else {
			return "AIS_" + StrPad(Int64ToString(userObj.Id), 8, "0", "LEFT")
		}
	}

}

func (this *GlobalApi) GetUserImage(userObj *models.SecUser) string {

	if userObj.IsFacebook {
		return "https://graph.facebook.com/" + userObj.Facebookid + "/picture?type=large"
	} else if userObj == nil || userObj.ImageProfile == "" {
		return this.GetHostStaticGrails() + "/defaultValue/profile.jpeg"
	} else {
		return this.GetHostStaticGolang() + "/" + userObj.ImageProfile
	}
}

func (this *GlobalApi) GetUserPackageName(userObj *models.SecUser) string {
	return userObj.CurrentPackage
}

func (this *GlobalApi) RenderResultVideoBySongObj(songObj *models.Song) map[string]interface{} {

	result := map[string]interface{}{
		KEY_JSON_VOCAL: map[string]interface{}{
			KEY_JSON_1080P: this.RenderEncryptm3u8(songObj.P1080link),
			KEY_JSON_720P:  this.RenderEncryptm3u8(songObj.Link),
			KEY_JSON_480P:  this.RenderEncryptm3u8(songObj.P480link),
			KEY_JSON_360P:  this.RenderEncryptm3u8(songObj.SdLink),
			KEY_JSON_240P:  this.RenderEncryptm3u8(songObj.P240link),
			KEY_JSON_144P:  this.RenderEncryptm3u8(songObj.P144sdLink),
			KEY_JSON_auto:  this.RenderEncryptm3u8(songObj.AutoLink),
		},
		KEY_JSON_KARAOKE: map[string]interface{}{
			KEY_JSON_1080P: this.RenderEncryptm3u8(songObj.P1080link),
			KEY_JSON_720P:  this.RenderEncryptm3u8(songObj.NoVocalLink),
			KEY_JSON_480P:  this.RenderEncryptm3u8(songObj.P480noVocalLink),
			KEY_JSON_360P:  this.RenderEncryptm3u8(songObj.SdNoVocalLink),
			KEY_JSON_240P:  this.RenderEncryptm3u8(songObj.P240noVocalLink),
			KEY_JSON_144P:  this.RenderEncryptm3u8(songObj.P144sdNoVocalLink),
			KEY_JSON_auto:  this.RenderEncryptm3u8(songObj.AutoNoVocalLink),
		},
	}

	return result
}

func (this *GlobalApi) CreateSecureLink(contentLink string) string {

	today := time.Now()

	expiredSec := 60 * 60 * 3
	endTime := Int64ToString(today.Unix() + int64(expiredSec))

	uri := this.GetUriFromFullm3u8(contentLink)
	secureLink := endTime + uri + " " + SECURE_SECURE_LINK_KEY
	md5 := GetMD5Hash(secureLink)
	//beego.Debug("output createSecureLink = " + contentLink + "?m=" + md5 + "&e=" + endTime)

	return contentLink + "?m=" + md5 + "&e=" + endTime
}

func (this *GlobalApi) GetUriFromFullm3u8(fullLink string) string {

	_, step1 := SubstringWithIndex(fullLink, "/")
	_, step2 := SubstringWithIndex(step1, "/")
	_, step3 := SubstringWithIndex(step2, "/")

	return step3

}

func (this *GlobalApi) RenderEncryptm3u8(contentLink string) string {

	if contentLink == "" {
		return ""
	} else {
		secureLink := this.CreateSecureLink(contentLink)

		//beego.Error("secureLink = ",secureLink)

		result := base64.StdEncoding.EncodeToString(AesEncrypt(secureLink, SECURE_AES128_KEY))
		return result
	}

}
