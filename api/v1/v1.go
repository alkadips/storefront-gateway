package apiv1

import (
	authenticate "github.com/MyriadFlow/storefront-gateway/api/v1/authenticate"
	creatorrole "github.com/MyriadFlow/storefront-gateway/api/v1/creatorRole"
	delegateassetcreation "github.com/MyriadFlow/storefront-gateway/api/v1/delegateAssetCreation"

	"github.com/MyriadFlow/storefront-gateway/api/v1/healthcheck"
	"github.com/MyriadFlow/storefront-gateway/api/v1/highlights"
	"github.com/MyriadFlow/storefront-gateway/api/v1/likes"
	"github.com/MyriadFlow/storefront-gateway/api/v1/nftstorage"
	"github.com/MyriadFlow/storefront-gateway/api/v1/profile"

	"github.com/MyriadFlow/storefront-gateway/api/v1/wishlist"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes Use the given Routes
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1.0")
	{
		authenticate.ApplyRoutes(v1)
		profile.ApplyRoutes(v1)

		creatorrole.ApplyRoutes(v1)
		delegateassetcreation.ApplyRoutes(v1)
		nftstorage.ApplyRoutes(v1)

		healthcheck.ApplyRoutes(v1)
		highlights.ApplyRoutes(v1)
		likes.ApplyRoutes(v1)
		wishlist.ApplyRoutes(v1)
	}
}
