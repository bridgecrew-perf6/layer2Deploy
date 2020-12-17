package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/ontio/layer2deploy/common"
	"github.com/ontio/layer2deploy/core"
	"github.com/ontio/ontology/common/log"
)

func RoutesApi(parent *gin.Engine) {
	apiRoute := parent.Group("/api")
	apiRoute.GET("/VerifyHash/:hash", VerifyHash)
	apiRoute.POST("/enableSendService/:enableSendService", EnableSendService)
}

func EnableSendService(c *gin.Context) {
	enable := c.Param("enableSendService")
	if enable == "true" {
		log.Infof("EnableSendService true")
		atomic.StoreUint32(&core.DefSendService.Enabled, 1)
	} else {
		log.Infof("EnableSendService false")
		atomic.StoreUint32(&core.DefSendService.Enabled, 0)
	}
}

func VerifyHash(c *gin.Context) {
	hash := c.Param("hash")

	if len(hash) != sha256.Size*2 {
		c.JSON(http.StatusOK, common.VerifyResponse{
			Code:    common.HASHLENERROR,
			Message: common.CodeMessageMap[common.HASHLENERROR],
		})
		return
	}

	log.Debugf("VerifyHash: Y.0 %s", hash)

	_, err := hex.DecodeString(hash)
	if err != nil {
		log.Debugf("VerifyHash: N.0 %s", err)
		c.JSON(http.StatusOK, common.VerifyResponse{
			Code:    common.HASHDATAERROR,
			Message: common.CodeMessageMap[common.HASHDATAERROR] + fmt.Sprintf("%v", err),
		})
		return
	}

	result, err := core.DefVerifyService.VerifyHashCore(hash)
	if err != nil {
		log.Debugf("VerifyHash: N.1 %s", err)
		c.JSON(http.StatusOK, common.VerifyResponse{
			Code:    common.SERVERERROR,
			Message: common.CodeMessageMap[common.SERVERERROR] + fmt.Sprintf(" %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, common.VerifyResponse{
		Code:    common.SUCCESS,
		Message: common.CodeMessageMap[common.SUCCESS],
		Result:  result,
	})
}
