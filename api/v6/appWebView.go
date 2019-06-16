package v6

import (
	"github.com/astaxie/beego"
	"gitlab.com/wisdomvast/AiskaraokeServerGolang/models"
)

func (this *GlobalApi) WebviewPopup() {

	popupId, _ := this.GetInt64(PARAMS_POPUP_ID)
	popupObj, err := models.GetConfigPopupById(popupId)

	if err != nil {
		beego.Debug(popupObj)
	}

	result := this.GenerateDeepLinkByPopupObj(popupObj)

	this.Data["hrefResult"] = result
	this.Data["path"] = "../../../"+popupObj.Image
	this.TplName = "appwebview/popupHome.html"
}
