package v6

import (
	"github.com/iballbar/beegoAPI"
	"gitlab.com/wisdomvast/AiskaraokeServerGolang/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type GlobalApi struct {
	beegoAPI.API
}

func (this *GlobalApi) GetUserByAccessToken(accessToken string) *models.SecUser {

	userTokenImeiObj := models.GetUserTokenImeiByAceessTokenAndExpired(accessToken)

	if accessToken == "" || userTokenImeiObj == nil {
		return nil
	}

	return userTokenImeiObj.SecUser
}

func (this *GlobalApi) CheckClientType(clientType string) string {
	var clientDevice = "box"
	if clientType != "android" && clientType != "androidBox" && clientType != "" {
		clientDevice = "mobile"
	}
	return clientDevice
}

func (this *GlobalApi) CheckOsClient(clientType string) string {
	result := ""
	if (strings.Contains(clientType, "android") && clientType != "android" &&
		clientType != "androidBox" && clientType != "") ||
		clientType == "androidMobile" || clientType == "android" || clientType == "androidBox" {
		result = "android"
	} else {
		result = "ios"
	}
	return result
}

func (this *GlobalApi) GetCountAllSongInCat(catObj *models.Category, clientTypeFunc string) int64 {

	var result int64 = 0
	allSong := models.GetAllSongObjByCatIdAndClientType(catObj.Id, clientTypeFunc)

	for _, v := range allSong {
		result += v.Song.CountView
	}

	return result

}

func (this *GlobalApi) GetCountAllSongInGenre(genreObj *models.Genre, clientTypeFunc string) int64 {

	var result int64 = 0
	allSong := models.GetAllSongObjByGenreIdAndClientType(genreObj.Id, clientTypeFunc)

	for _, v := range allSong {
		result += v.Song.CountView
	}

	return result

}

func (this *GlobalApi) GenerateDeepLinkByPopupObj(popupObj *models.ConfigPopup) string {

	clickTo := popupObj.GoTo
	content := popupObj.ContentRef

	result := "data?nextTo=" + clickTo + "&contentUrl=" + content

	return result

}

//=============================== API RENDER ================================

func (this *GlobalApi) RenderImagesOnSplashScreenByConfigScreen(configScreenObj *models.ConfigScreen) []string {

	width, _ := this.GetInt64(PARAMS_WIDTH)
	height, _ := this.GetInt64(PARAMS_HEIGHT)

	result := []string{}
	if configScreenObj == nil || height == 0 || width == 0 {
		return result
	}

	resolutionImage := float64(height) / float64(width)

	allScreenImageObj := models.GetAllConfigScreenImageByconfigScreenObj(-1, 0, configScreenObj)
	for _, v := range allScreenImageObj {

		if resolutionImage <= 1.7 {
			result = append(result, this.GetHostStaticGrails()+"/"+v.ImagePathFat)
		} else {
			result = append(result, this.GetHostStaticGrails()+"/"+v.ImagePathThin)
		}
	}
	return result
}

func (this *GlobalApi) HashBcryptPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}
