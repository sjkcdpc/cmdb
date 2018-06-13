package cmdb

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mds1455975151/cmdb/storage"
	"github.com/mds1455975151/cmdb/utils"
	"github.com/mds1455975151/cmdb/errors"
	"fmt"
)

func init() {
	getHandlers["/users"] = users
}

// Binding from JSON
type RequestQueryUsersData struct {
	Id      string `form:"id" json:"id" binding:"required"`
	Expiration bool   `form:"expiration" json:"expiration"`
}

type ResponseQueryUsersData struct {
	// in: body
	Body struct {
		Code              int64  `json:"code"`
		Message           string `json:"message"`
		GlobalId          int64  `json:"globalId"`
		ExpirationSeconds int64  `json:"expiration_seconds"`
	}
}

func users(c *gin.Context) {

	var request RequestQueryUsersData

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "query_users BindJSON failed: %v", err)
		return
	}

	if request.Id == "" {

		utils.QuickReply(c, errors.Failed, "query_users Id is empty.")
		return
	}
	fmt.Println(request.Id)
	//UsersRecord := storage.QueryUsers(request.Id)
	//if err != nil {
	//
	//	code := errors.Failed
	//	if ec, ok := err.(*errors.Type); ok {
	//		code = ec.Code
	//	}
	//
	//	utils.QuickReply(c, code, "query_token QueryAccessToken failed: %v", err)
	//	return
	//}
	//
	//if UsersRecord.Id == 0 {
	//
	//	utils.QuickReply(c, errors.Failed, "query_token globalId is invalid.")
	//	return
	//}

	info := storage.QueryHost(1)
	c.JSON(http.StatusOK, info)
}