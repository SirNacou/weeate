package services

import (
	"context"
	"fmt"
	"time"

	common_domain "github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/SirNacou/weeate/backend/internal/features/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

type UploadPayload struct {
	FileName          string `json:"fileName" doc:"The name of the file to be uploaded"`
	UseUniqueFileName string `json:"useUniqueFileName,omitempty" enum:"true,false" doc:"Whether to use a unique file name"`
	Tags              string `json:"tags,omitempty" doc:"The tags associated with the file"`
	Folder            string `json:"folder,omitempty" doc:"The folder path where the file will be uploaded"`
	IsPrivateFile     string `json:"isPrivateFile,omitempty" enum:"true,false" doc:"Whether the file is private"`
}

type ImageKitAuthParametersResponse struct {
	Token string `json:"token" doc:"The generated JWT token for ImageKit authentication"`
}

type GetIKTokenQueryHandler struct {
	privateKey string
	publicKey  string
}

func NewGetIKTokenQueryHandler(privateKey string, publicKey string) *GetIKTokenQueryHandler {
	return &GetIKTokenQueryHandler{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

type ImageKitClaims struct {
	jwt.RegisteredClaims
	UploadPayload
}

func (h *GetIKTokenQueryHandler) Handle(ctx context.Context, _ *struct{}) (*ImageKitAuthParametersResponse, error) {
	user, ok := domain.UserFromContext(ctx)
	if !ok {
		return nil, domain.ErrUserNotFoundInContext
	}

	now := time.Now()
	exp := now.Add(5 * 60 * time.Second)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ImageKitClaims{
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		UploadPayload{
			FileName:      fmt.Sprintf("avatar-%s", user.AppMetadata.DisplayName),
			Tags:          "avatar",
			Folder:        fmt.Sprintf("/%s/%s/%s/avatars", "weeate", "users", user.ID),
			IsPrivateFile: "false",
		},
	})
	token.Header["kid"] = h.publicKey

	signedToken, err := token.SignedString([]byte(h.privateKey))
	if err != nil {
		return nil, common_domain.NewError(common_domain.EInternal, "Failed to sign ImageKit token")
	}

	return &ImageKitAuthParametersResponse{
		Token: signedToken,
	}, nil
}
