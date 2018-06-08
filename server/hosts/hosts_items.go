package hosts

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/mds1455975151/cmdb/utils"
	"net/http"
)

func init() {

	handlerFuncList["/hosts"] = hostitem
}

//
// Request Hosts Item
//
type HostsItemReq struct {
	// The globalId
	//
	// Required: true
	GlobalId int64 `form:"globalId" json:"globalId" binding:"required"`
}

//
// 应答: 协议返回包
// swagger:response DoBuyItemRsp
// noinspection ALL
type HostsItemRsp struct {
	// in: body
	Body struct {
		// The response code
		//
		// Required: true
		Code int64 `json:"code"`
		// The response message
		//
		// Required: true
		Message string `json:"message"`
		// The payment sequence to identify the flow.
		//
		// Required: true
		Sequence string `json:"sequence"`
		// The response BuyItem
		//
		// Required: true
		BuyItemMidas struct {
			//
			// ret为0的时候，返回真正购买物品的url的参数，开发者需要把该参数
			// 传给sdk跳转到相关页面使用户完成真正的购买动作。
			//
			UrlParams string `json:"url_params"`
		} `json:"buyItemMidas"`
	}
}

//
// swagger:route POST /payment/buy_item payment payment_buy_item
//
// Return buy item for the given user:
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: DoBuyItemRsp
func hostitem(c *gin.Context) {

	c.JSON(http.StatusOK, resp.Body)
}