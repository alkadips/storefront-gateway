package claimrole

import (
	"fmt"
	"net/http"
	"netsepio-api/db"
	"netsepio-api/middleware/auth/jwt"
	"netsepio-api/models"
	"netsepio-api/util/pkg/cryptosign"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/claimrole")
	{
		g.Use(jwt.JWT)
		g.POST("", postClaimRole)
	}
}

func postClaimRole(c *gin.Context) {
	var req ClaimRoleRequest
	c.BindJSON(&req)

	//Message containing flowId
	role, err := getRoleByFlowId(req.FlowId)
	if err == gorm.ErrRecordNotFound {
		c.String(http.StatusNotFound, "flow id not found")
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	message := role.Eula + "m"
	fmt.Println("signed with", message)
	walletAddress, isCorrect, err := cryptosign.CheckSign(req.Signature, req.FlowId, message)

	if err == cryptosign.ErrFlowIdNotFound {
		c.String(http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	if !isCorrect {
		c.Status(http.StatusForbidden)
		return
	}

	// Update user role
	db.Db.Model(&models.User{}).Where("wallet_address = ?", walletAddress).
		Update("roles", gorm.Expr("array_cat(roles,?)", pq.Array([]int{role.RoleId})))

}

func getRoleByFlowId(flowId string) (models.Role, error) {
	var flowIdRecord models.FlowId
	err := db.Db.Model(&models.FlowId{}).Where("flow_id = ?", flowId).First(&flowIdRecord).Error
	if err != nil {
		return models.Role{}, err
	}

	var role models.Role
	err = db.Db.Model(&models.Role{}).First(&role, flowIdRecord.RelatedRoleId).Error
	if err != nil {
		return models.Role{}, err
	}
	return role, nil
}
