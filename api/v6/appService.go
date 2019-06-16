package v6

import (
	"github.com/astaxie/beego"
	"gitlab.com/wisdomvast/AiskaraokeServerGolang/models"
	"strconv"
	"time"
	//"time"
)

func (this *GlobalApi) KeepDevice() {

	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	imei := this.GetString(PARAMS_DEVICE_ID)
	deviceToken := this.GetString(PARAMS_Device_TOKEN)
	brand := this.GetString(PARAMS_BRAND)
	model := this.GetString(PARAMS_MODEL)
	osVersion := this.GetString(PARAMS_OS_VERSION)
	accessToken := this.GetString(PARAMS_ACCESS_TOKEN)
	serial := this.GetString(PARAMS_SERIAL)

	var osClient = this.CheckOsClient(clientType)

	if imei == "" {
		this.ResponseJSON(nil, 703, this.GetStringLanguage(ERROR_703_MESSAGE_TH, ERROR_703_MESSAGE_ENG))
		return
	}

	if clientType == "android" {
		clientType = "androidBox"
	}

	if clientType == "androidBox" {

		deviceUpdateObj := models.GetDeviceByImei(imei)
		if deviceUpdateObj != nil {

			var newOsClient string
			if deviceToken == "" && deviceUpdateObj.DeviceToken == "" {
				newOsClient = osClient + "false"
			} else {
				newOsClient = osClient
			}

			deviceUpdateObj.ClientType = clientType
			deviceUpdateObj.Serial = serial
			deviceUpdateObj.OsVersion = osVersion
			deviceUpdateObj.ModelName = model
			deviceUpdateObj.Brand = brand

			if deviceToken != "" {
				deviceUpdateObj.DeviceToken = deviceToken
			}

			deviceUpdateObj.MobileOs = newOsClient

			err := models.UpdateDeviceById(deviceUpdateObj)
			if err != nil {
				beego.Error("KeepDevice = ", this.GetStringLanguage(ERROR_801_MESSAGE_TH, ERROR_801_MESSAGE_ENG))
				this.ResponseJSON(nil, 801, this.GetStringLanguage(ERROR_801_MESSAGE_TH, ERROR_801_MESSAGE_ENG))
				return
			}
		}
	}

	userObj := this.GetUserByAccessToken(accessToken)
	checkDeviceObj := models.GetDeviceByImeiAndClientType(imei, clientType)

	if checkDeviceObj == nil {

		newOsClient := ""
		if deviceToken == "" {
			newOsClient = osClient + "false"
		} else {
			newOsClient = osClient
		}

		newDevice := models.Device{
			DeviceId:    imei,
			ClientType:  clientType,
			Serial:      serial,
			OsVersion:   osVersion,
			ModelName:   model,
			Brand:       brand,
			DeviceToken: deviceToken,
			MobileOs:    newOsClient,
		}
		idKeyDevice, err := models.AddDevice(&newDevice)
		if err != nil {
			beego.Error("err = ", err)
			beego.Error("Add Device on Keep device")
			this.ResponseJSON(nil, 801, this.GetStringLanguage(ERROR_801_MESSAGE_TH, ERROR_801_MESSAGE_ENG))
			return
		}

		deviceObj, _ := models.GetDeviceById(idKeyDevice)

		newStatsPlaySong := models.StatsPlaySong{
			Action:     "5",
			DeviceId:   imei,
			Song:       nil,
			SecUser:    userObj,
			DeviceObj:  deviceObj,
			ClientType: deviceObj.ClientType,
		}
		_, err = models.AddStatsPlaySong(&newStatsPlaySong)
		if err != nil {
			beego.Error("err = ", err)
			beego.Error("Add stats on Keep device")
			this.ResponseJSON(nil, 801, this.GetStringLanguage(ERROR_801_MESSAGE_TH, ERROR_801_MESSAGE_ENG))
			return
		}

	} else {

		newOsClient := ""
		if deviceToken == "" && checkDeviceObj.DeviceToken == "" {
			newOsClient = osClient + "false"
		} else {
			newOsClient = osClient
		}

		checkDeviceObj.Serial = serial
		checkDeviceObj.OsVersion = osVersion
		checkDeviceObj.ModelName = model
		checkDeviceObj.Brand = brand
		checkDeviceObj.MobileOs = newOsClient
		if deviceToken != "" {
			checkDeviceObj.DeviceToken = deviceToken
		}

		err := models.UpdateDeviceById(checkDeviceObj)
		if err != nil {
			beego.Error("err = ", err)
			this.ResponseJSON(nil, 801, this.GetStringLanguage(ERROR_801_MESSAGE_TH, ERROR_801_MESSAGE_ENG))
			return
		}

		newStatsPlaySong := models.StatsPlaySong{
			Action:     "5",
			DeviceId:   imei,
			Song:       nil,
			SecUser:    userObj,
			DeviceObj:  checkDeviceObj,
			ClientType: checkDeviceObj.ClientType,
		}
		_, err = models.AddStatsPlaySong(&newStatsPlaySong)
		if err != nil {
			beego.Error("err = ", err)
			beego.Error("Add stats on Keep device")
			this.ResponseJSON(nil, 801, this.GetStringLanguage(ERROR_801_MESSAGE_TH, ERROR_801_MESSAGE_ENG))
			return
		}
	}

	this.ResponseJSON(nil, 200, "")
	return

}

func (this *GlobalApi) SubscribeMobile() {

	today := time.Now()
	networkType := this.GetString(PARAMS_NETWORK_TYPE)
	username := this.GetString(PARAMS_PRIVATE_ID)
	userNewPassword := this.GetString(PARAMS_PRIVATE_ID_PASSWORD)
	confirmKick := this.GetString(PARAMS_CONFIRM_KICK)
	deviceId := this.GetString(PARAMS_DEVICE_ID)
	brand := this.GetString(PARAMS_BRAND)
	model := this.GetString(PARAMS_MODEL)
	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	serial := this.GetString(PARAMS_SERIAL)
	osVersion := this.GetString(PARAMS_OS_VERSION)
	osClient := this.GetString(PARAMS_OS_CLIENT)
	msisdnMobile := this.GetString(PARAMS_MSISDN_MOBILE)
	ffbidBox := this.GetString(PARAMS_FFBID_BOX)
	idValue := this.GetString(PARAMS_ID_VALUE)
	idType := this.GetString(PARAMS_ID_TYPE)

	if ffbidBox == "" {
		ffbidBox = this.GetString(PARAMS_ID_VALUE)
	}

	if idValue == "" {
		idValue = this.GetString(PARAMS_FFBID_BOX)
	}

	deviceObj := models.GetDeviceByImeiAndClientType(deviceId, clientType)
	accessToken := this.GenerateUserToken()

	if deviceObj == nil {
		newDevice := models.Device{
			DeviceId:   deviceId,
			ClientType: clientType,
			Serial:     serial,
			OsVersion:  osVersion,
			ModelName:  model,
			Brand:      brand,
			MobileOs:   osClient,
		}
		_, errAddDevice := models.AddDevice(&newDevice)
		if errAddDevice != nil {
			beego.Error("errAddDevice = ", errAddDevice)
			beego.Error("Add Device on SubscribeMobile")
			this.ResponseJSON(nil, 801, this.GetStringLanguage(ERROR_401_MESSAGE_TH, ERROR_401_MESSAGE_ENG))
			return
		}
	}

	Block{
		Try: func() {

			oldUserObj := models.GetSecUserByUsernameAndEnabled(username, false)
			oldUserEnabledObj := models.GetSecUserByUsernameAndEnabled(username, true)

			if oldUserObj != nil {
				oldUserObj.Enabled = true
				oldUserObj.TypeTid = idType
				oldUserObj.ValueVid = idValue
				oldUserObj.Msisdn = msisdnMobile
				oldUserObj.FfbidBox = ffbidBox
				oldUserObj.NetworkType = networkType
				err := models.UpdateSecUserById(oldUserObj)
				if err != nil {
					beego.Error(err)
				}

				err = models.DeleteAllUserTokenByUserObj(oldUserObj)
				if err != nil {
					beego.Error(err)
				}

				newUserTokenImei := models.UserTokenImei{
					SecUser:     oldUserObj,
					Imei:        deviceId,
					AccessToken: accessToken,
					DateLogin:   today,
					DeviceObj:   deviceObj,
					DateExpired: today.AddDate(100, 0, 0),
				}

				_, err = models.AddUserTokenImei(&newUserTokenImei)
				if err != nil {
					beego.Error("AddUserTokenImei = ", err)
				}

				newStatePlaySong := models.StatsPlaySong{
					Action:     "6",
					DeviceId:   deviceId,
					Song:       nil,
					SecUser:    oldUserObj,
					DeviceObj:  deviceObj,
					ClientType: deviceObj.ClientType,
				}

				_, err = models.AddStatsPlaySong(&newStatePlaySong)
				if err != nil {
					beego.Error("AddnewStatePlaySong = ", err)
				}

				this.ResponseJSON(map[string]interface{}{
					KEY_JSON_ACCESS_TOKEN: accessToken,
					KEY_JSON_IS_NEW_USER:  false,
					KEY_JSON_USER_PROFILE: this.RenderUserProfile(oldUserObj),
				}, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
				return

			} else if oldUserEnabledObj != nil {

				oldUserEnabledObj.TypeTid = idType
				oldUserEnabledObj.ValueVid = idValue
				oldUserEnabledObj.Msisdn = msisdnMobile
				oldUserEnabledObj.FfbidBox = ffbidBox
				oldUserEnabledObj.NetworkType = networkType

				err := models.UpdateSecUserById(oldUserEnabledObj)
				if err != nil {
					beego.Error("UpdateSecUserById = ", err)
				}

				userTokenImeiObj := models.GetUserTokenImeiByUserObjAndImeiAndDeviceObj(oldUserEnabledObj, deviceId, deviceObj)

				if userTokenImeiObj != nil {

					userTokenImeiObj.DateLogin = today
					userTokenImeiObj.DateExpired = today.AddDate(100, 0, 0)

					err = models.UpdateUserTokenImeiById(userTokenImeiObj)
					if err != nil {
						beego.Error("UpdateUserTokenImeiById = ", err)
					}

					accessToken = userTokenImeiObj.AccessToken

					newStatsPlaySong := models.StatsPlaySong{
						Action:     "6",
						DeviceId:   deviceId,
						Song:       nil,
						SecUser:    oldUserEnabledObj,
						DeviceObj:  deviceObj,
						ClientType: deviceObj.ClientType,
					}

					_, err = models.AddStatsPlaySong(&newStatsPlaySong)
					if err != nil {
						beego.Error("AddStatsPlaySong 1 = ", err)
					}

					this.ResponseJSON(map[string]interface{}{
						KEY_JSON_ACCESS_TOKEN: accessToken,
						KEY_JSON_IS_NEW_USER:  false,
						KEY_JSON_USER_PROFILE: this.RenderUserProfile(oldUserEnabledObj),
					}, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
					return

				} else {

					allUserLogin := models.GetAllUserTokenImeiByUserObj(-1, 0, oldUserEnabledObj)
					userLoginSize := len(allUserLogin)

					if userLoginSize <= 7 || oldUserEnabledObj.Username == "YQ4oTDi3hML6t40BAwnoxgaz658i2URV1441969725265@ais.co.th" {

						newUserTokenImei := models.UserTokenImei{
							SecUser:     oldUserEnabledObj,
							Imei:        deviceId,
							AccessToken: accessToken,
							DateLogin:   today,
							DeviceObj:   deviceObj,
							DateExpired: today.AddDate(100, 0, 0),
						}

						_, err = models.AddUserTokenImei(&newUserTokenImei)
						if err != nil {
							beego.Error("AddUserTokenImei = ", err)
						}

						newStatsPlaySong := models.StatsPlaySong{
							Action:     "6",
							DeviceId:   deviceId,
							Song:       nil,
							SecUser:    oldUserEnabledObj,
							DeviceObj:  deviceObj,
							ClientType: deviceObj.ClientType,
						}

						_, err = models.AddStatsPlaySong(&newStatsPlaySong)
						if err != nil {
							beego.Error("AddStatsPlaySong 2 = ", err)
						}

						this.ResponseJSON(map[string]interface{}{
							KEY_JSON_ACCESS_TOKEN: accessToken,
							KEY_JSON_IS_NEW_USER:  false,
							KEY_JSON_USER_PROFILE: this.RenderUserProfile(oldUserEnabledObj),
						}, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
						return

					} else if confirmKick == "Y" {

						userLoginList := models.GetUserTokenImeiByUserObjAndminDateLogin(oldUserEnabledObj)

						userLoginList.Imei = deviceId
						userLoginList.AccessToken = accessToken
						userLoginList.DateLogin = today
						userLoginList.DateExpired = today.AddDate(100, 0, 0)

						err := models.UpdateUserTokenImeiById(userLoginList)
						if err != nil {
							beego.Error(err, " := models.UpdateUserTokenImeiById(userLoginList)")
						}

						newStatsPlaySong := models.StatsPlaySong{
							Action:     "6",
							DeviceId:   deviceId,
							Song:       nil,
							SecUser:    oldUserEnabledObj,
							DeviceObj:  deviceObj,
							ClientType: deviceObj.ClientType,
						}

						_, err = models.AddStatsPlaySong(&newStatsPlaySong)
						if err != nil {
							beego.Error("AddStatsPlaySong 3 = ", err)
						}

					} else {
						this.ResponseJSON(map[string]interface{}{
							KEY_JSON_ACCESS_TOKEN: accessToken,
							KEY_JSON_IS_NEW_USER:  false,
							KEY_JSON_USER_PROFILE: this.RenderUserProfile(oldUserEnabledObj),
						}, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
						return
					}
				}

			} else {

				passwordHash, errHashPassword := this.HashBcryptPassword(userNewPassword, 10)
				if errHashPassword != nil {
					beego.Error("errHashPassword = ", errHashPassword)
				}

				newSecUser := models.SecUser{
					Username:          username,
					Password:          passwordHash,
					PasswordNoEncrypt: userNewPassword,
					RegisterDate:      today,
					Enabled:           true,
					RoleString:        "roleUser",
					TypeTid:           idType,
					ValueVid:          idValue,
					Msisdn:            msisdnMobile,
					FfbidBox:          ffbidBox,
					NetworkType:       networkType,
				}

				_, err := models.AddSecUser(&newSecUser)
				if err != nil {
					beego.Error("AddSecUser 414 = ", err)
				}

				newUserTokenImei := models.UserTokenImei{
					SecUser:     &newSecUser,
					Imei:        deviceId,
					AccessToken: accessToken,
					DateLogin:   today,
					DeviceObj:   deviceObj,
					DateExpired: today.AddDate(100, 0, 0),
				}

				_, err = models.AddUserTokenImei(&newUserTokenImei)
				if err != nil {
					beego.Error("AddUserTokenImei 2 = ", err)
				}

				newStatsPlaySong := models.StatsPlaySong{
					Action:     "6",
					DeviceId:   deviceId,
					Song:       nil,
					SecUser:    oldUserEnabledObj,
					DeviceObj:  deviceObj,
					ClientType: deviceObj.ClientType,
				}

				_, err = models.AddStatsPlaySong(&newStatsPlaySong)
				if err != nil {
					beego.Error("AddStatsPlaySong 3 = ", err)
				}

				this.ResponseJSON(map[string]interface{}{
					KEY_JSON_ACCESS_TOKEN: accessToken,
					KEY_JSON_IS_NEW_USER:  false,
					KEY_JSON_USER_PROFILE: this.RenderUserProfile(oldUserEnabledObj),
				}, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
				return

			}

		},
		Catch: func(e Exception) {
			beego.Error("Caught = ", e)
			this.ResponseJSON(nil, 500, this.GetStringLanguage(ERROR_500_MESSAGE_TH, ERROR_500_MESSAGE_ENG))
			return
		},
		Finally: func() {
		},
	}.Do()

}

func (this *GlobalApi) UserPlaySession() {

	accessToken := this.GetString(PARAMS_ACCESS_TOKEN)
	deviceId := this.GetString(PARAMS_DEVICE_ID)
	clientType := this.GetString(PARAMS_CLIENT_TYPE)

	today := time.Now()

	userTokenImeiObj := models.GetUserByAccessToken(accessToken)
	deviceObj := models.GetDeviceByImeiAndClientType(deviceId, clientType)

	if userTokenImeiObj == nil {
		this.ResponseJSON(nil, 404, this.GetStringLanguage(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
		return
	}

	userObj := userTokenImeiObj.SecUser

	if deviceObj == nil {
		this.ResponseJSON(nil, 703, this.GetStringLanguage(ERROR_703_MESSAGE_TH, ERROR_703_MESSAGE_ENG))
		return
	}

	userTokenImei := models.GetUserTokenImeiByUserObjAndDeviceObj(userObj, deviceObj)

	if userTokenImei == nil {
		this.ResponseJSON(nil, 404, this.GetStringLanguage(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
		return
	}

	sessionUserPlaySongObj := models.GetSessionUserPlaySongBySecUserAndExpiredGreaterthan(userObj, today)

	if sessionUserPlaySongObj != nil {

		if sessionUserPlaySongObj.Device == deviceObj {

			sessionUserPlaySongObj.ExpiredSession = today.Add(time.Second * time.Duration(120))

			err := models.UpdateSessionUserPlaySongById(sessionUserPlaySongObj)
			if err != nil {
				beego.Error(err, " = models.UpdateSessionUserPlaySongById(sessionUserPlaySongObj)")
			}

			this.ResponseJSON(nil, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
			return

		} else {

			if sessionUserPlaySongObj.ExpiredSession.Nanosecond() < today.Nanosecond() {

				sessionUserPlaySongObj.Device = deviceObj
				sessionUserPlaySongObj.ExpiredSession = today.Add(time.Second * time.Duration(120))

				err := models.UpdateSessionUserPlaySongById(sessionUserPlaySongObj)
				if err != nil {
					beego.Error(err, " = models.UpdateSessionUserPlaySongById(sessionUserPlaySongObj)")
				}

				this.ResponseJSON(nil, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
				return

			} else {

				sessionUserPlaySongObj.Device = deviceObj
				sessionUserPlaySongObj.ExpiredSession = today.Add(time.Second * time.Duration(120))

				err := models.UpdateSessionUserPlaySongById(sessionUserPlaySongObj)
				if err != nil {
					beego.Error(err, " = models.UpdateSessionUserPlaySongById(sessionUserPlaySongObj)")
				}

				this.ResponseJSON(nil, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
				return

			}

		}

	} else {

		newSessionUserPlaySong := models.SessionUserPlaySong{
			SecUser:        userObj,
			Device:         deviceObj,
			ExpiredSession: today.Add(time.Second * time.Duration(120)),
		}

		_, err := models.AddSessionUserPlaySong(&newSessionUserPlaySong)
		if err != nil {
			beego.Error(err, " on AddSessionUserPlaySong")
		}

		this.ResponseJSON(nil, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
		return

	}

}

func (this *GlobalApi) StatsPlaySong() {

	accessToken := this.GetString(PARAMS_ACCESS_TOKEN)
	imei := this.GetString(PARAMS_DEVICE_ID)
	songId, _ := this.GetInt64(PARAMS_SONG_ID, 0)
	action := this.GetString(PARAMS_ACTIONS)
	clientType := this.GetString(PARAMS_CLIENT_TYPE)
	isChromeCast := this.GetString(PARAMS_IS_CHROME_CAST)
	isVocal := this.GetString(PARAMS_IS_VOCAL)

	if clientType == "android" {
		clientType = "androidBox"
	}

	deviceObj := models.GetDeviceByImeiAndClientType(imei, clientType)

	userObj := models.GetRealUserByAccessToken(accessToken)
	if userObj == nil {
		userObj = models.GetSecUserByUsername("wisdomvast@gmail.com")
	}

	if action != "5" {
		if deviceObj == nil || songId == 0 || action == "" {
			if deviceObj == nil {
				this.ResponseJSON(nil, 404, this.GetStringLanguage(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
				return
			}

			this.ResponseJSON(nil, 703, this.GetStringLanguage(ERROR_703_MESSAGE_TH, ERROR_703_MESSAGE_ENG))
			return

		}
	} else {
		if deviceObj == nil || action == "" {
			this.ResponseJSON(nil, 703, this.GetStringLanguage(ERROR_703_MESSAGE_TH, ERROR_703_MESSAGE_ENG))
			return
		}
	}

	songObj, _ := models.GetSongById(songId)

	if action != "5" {
		if songObj == nil {
			this.ResponseJSON(nil, 404, this.GetStringLanguage(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
			return
		}
	}

	if songObj.IsFree == true && songObj != nil && action == "1" {

		aCount := deviceObj.CountSong

		deviceObj.CountSong = aCount + 1
		errUpdateDevice := models.UpdateDeviceById(deviceObj)
		if errUpdateDevice != nil {
			beego.Error("errUpdateDevice = ", errUpdateDevice)
		}

		uCount := userObj.CountSong + 1
		userObj.CountSong = uCount + 1

		beego.Debug("uCount = ",uCount)
		beego.Debug("userObj.CountSong = ",userObj.CountSong)

		errUpdateUser := models.UpdateSecUserById(userObj)
		if errUpdateUser != nil {
			beego.Error("errUpdateUser = ", errUpdateUser)
		}

	} else {
		recordClient := deviceObj.ClientType

		if isChromeCast == "1" {
			recordClient = "chromecast"
		}

		newStats := models.StatsPlaySong{
			Action:         action,
			SecUser:        userObj,
			DeviceId:       imei,
			DeviceObj:      deviceObj,
			Song:           songObj,
			ClientType:     recordClient,
			IsVocal:        isVocal,
			CurrentPackage: userObj.CurrentPackage,
		}

		_, eerAddStats := models.AddStatsPlaySong(&newStats)
		if eerAddStats != nil {
			beego.Error("eerAddStats = ", eerAddStats)
		}
	}

	if clientType == "android" {
		clientType = "androidBox"
	}

	//checkDevice := models.GetDeviceByImeiAndClientType(imei, clientType)

	if songObj != nil && userObj != nil {

		songObj.CountView = songObj.CountView + 1
		errUpdateSong := models.UpdateSongById(songObj)

		if errUpdateSong != nil {
			beego.Error("errUpdateSong = ", errUpdateSong)
		}

	}

	this.ResponseJSON(nil, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}

//================== module ======================

func (this *GlobalApi) GenerateUserToken() string {
	token := RandStringBytesMaskImprSrc(letterBytesCharAndNumSmall, 30)
	userTokenObj := models.GetUserByAccessToken(token)
	if userTokenObj != nil {
		return this.GenerateUserToken()
	}

	return token
}

func (this *GlobalApi) EditProfile() {

	accessToken := this.GetString(PARAMS_ACCESS_TOKEN)
	nickname := this.GetString(PARAMS_DISPLAY_NAME)

	userTokenObj := models.GetUserByAccessToken(accessToken)
	if userTokenObj == nil {
		beego.Debug(accessToken)
		beego.Debug(userTokenObj)
		this.ResponseJSON(nil, 404, this.GetStringLanguage(ERROR_404_MESSAGE_TH, ERROR_404_MESSAGE_ENG))
		return
	}

	userObj := userTokenObj.SecUser

	if nickname == "" {
		nickname = userObj.NickNameSocial
	}

	file, handler, err := this.GetFile(PARAMS_UPLOAD_FILE)

	checkUpload := true
	reasonUpload := userObj.ImageProfile

	if file != nil {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		imgName := Int64ToString(userObj.Id) + GetSafeFileName(Int64ToString(userObj.Id)+"_userImage_"+timestamp)
		fileName := imgName + GetFileExtension(handler.Filename)
		checkUpload, reasonUpload = UploadProfileImageGlobal(userObj.Id, file, handler, err, IMAGE_PATH_USER_PROFILE, fileName)
	}

	if checkUpload == false {
		beego.Error("upload shopimage false")
	}

	userObj.ImageProfile = reasonUpload
	if nickname != "" {
		userObj.NickNameSocial = nickname
	}

	err = models.UpdateSecUserById(userObj)
	if err != nil {
		beego.Error("err := models.UpdateSecUserById(userObj) = ", err)
	}

	this.ResponseJSON(nil, 200, this.GetStringLanguage(SUCCESS_200_MESSAGE_TH, SUCCESS_200_MESSAGE_ENG))
	return

}
