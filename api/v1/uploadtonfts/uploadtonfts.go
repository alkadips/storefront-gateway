package uploadtonfts

import (
	"github.com/MyriadFlow/storefront_gateway/api/middleware/auth/paseto"
	"github.com/MyriadFlow/storefront_gateway/config/envconfig"
	"github.com/MyriadFlow/storefront_gateway/util/pkg/httphelper"
	"github.com/gin-gonic/gin"
	client "github.com/nftstorage/go-client"
	"context"
	"os"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/uploadtonfts")
	{
		g.Use(paseto.PASETO)
		g.POST("", uploadtonfts)
	}
}


func uploadtonfts(c *gin.Context) {

	token:=envconfig.EnvVars.NFT_API_KEY
	configuration := client.NewConfiguration()
	ctx := context.WithValue(context.Background(), client.ContextAccessToken, token)
	api_client := client.NewAPIClient(configuration)
	
	form, err := c.MultipartForm()
	if err != nil {
		httphelper.NewInternalServerError(c, "failed to parse multipart form", "failed to parse multipart form, error: %v", err.Error())
		return
	}

	responsePayload := make([]UploadToNftsPayload, 0)
	files:=form.File["file"]

	for _, file := range files {
		//tempporarily storing multipart file and then read as os file	
		filename := "./" + file.Filename
		err:=c.SaveUploadedFile(file, filename)
		if err != nil {
			httphelper.NewInternalServerError(c, "failed to SaveUpload file", "failed to open file, error: %v", err.Error())
			return
		}
		logrus.Info("\n File Saved at :",filename)
		bO, err := os.Open(filename)
		if err != nil {
			httphelper.NewInternalServerError(c, "failed to load file", "failed to open file, error: %v", err.Error())
			return
		}

		resp, r, err := api_client.NFTStorageAPI.Store(ctx).Body(bO).Execute()
		if err != nil {
			httphelper.NewInternalServerError(c, "failed to upload ", "Error when calling `NFTStorageAPI.Store``: error: %v", err.Error())
			httphelper.NewInternalServerError(c, "failed to upload ", "Full HTTP response, error: %v", r)
			return
		}
		bO.Close()
		err=os.Remove(filename)
		if err != nil {
			httphelper.NewInternalServerError(c, "failed to clear temporary file stored", "failed to clear temporary file stored, error: %v", err.Error())
			return
		}
		cid, _ := json.Marshal(resp.Value.Cid)

		responsePayload=append(responsePayload,UploadToNftsPayload{file.Filename,string(cid)})


	}

	httphelper.SuccessResponse(c, "file successfully uploaded to nft storage", responsePayload)
}

