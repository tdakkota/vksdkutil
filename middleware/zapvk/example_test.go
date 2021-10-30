package zapvk_test

import (
	"fmt"

	"github.com/SevereCloud/vksdk/v2/api"
	sdkutil "github.com/tdakkota/vksdkutil/v3"
	"github.com/tdakkota/vksdkutil/v3/middleware/zapvk"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ExampleLog() {
	vk := sdkutil.BuildSDK("token").WithMiddleware(
		zapvk.Log(zap.L(), zapcore.DebugLevel, true),
	).Complete()

	users, err := vk.UsersGet(api.Params{})
	if err != nil {
		panic(err)
	}
	fmt.Println(users)
}
