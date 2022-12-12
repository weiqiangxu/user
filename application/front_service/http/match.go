package http

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/weiqiangxu/common-config/logger"
	"net/http"
	"strings"
	"time"

	"github.com/weiqiangxu/user/application/front_service/dtos"

	"github.com/weiqiangxu/protocol/order"
	pbUser "github.com/weiqiangxu/protocol/user"

	"github.com/gin-gonic/gin"
	"github.com/weiqiangxu/user/domain/user"
)

type UserAppHttpOption func(service *UserAppHttpService)

func WithUserDomainService(t user.DomainInterface) UserAppHttpOption {
	return func(service *UserAppHttpService) {
		service.userDomainSrv = t
	}
}

func WithOrderRpcClient(t order.OrderClient) UserAppHttpOption {
	return func(service *UserAppHttpService) {
		service.orderRpcClient = t
	}
}

func WithUserRpcClient(t pbUser.LoginClient) UserAppHttpOption {
	return func(service *UserAppHttpService) {
		service.userRpcClient = t
	}
}

type UserAppHttpService struct {
	userDomainSrv  user.DomainInterface
	orderRpcClient order.OrderClient
	userRpcClient  pbUser.LoginClient
}

func NewUserAppHttpService(options ...UserAppHttpOption) *UserAppHttpService {
	srv := &UserAppHttpService{}
	for _, o := range options {
		o(srv)
	}
	return srv
}

// GetUserList get user list
func (m *UserAppHttpService) GetUserList(c *gin.Context) {
	// 如果没有下面这一段执行的太快了在list查看不到
	ch := make(chan bool)
	go func() {
		var stringSlice []string
		for i := 0; i < 20; i++ {
			// pprof 显示在这里占用2MB的内存开销
			repeat := strings.Repeat("hello,world", 50000)
			stringSlice = append(stringSlice, repeat)
			time.Sleep(time.Millisecond * 500)
		}
		logger.Info(len(stringSlice))
		ch <- true
	}()
	<-ch
	info, _ := m.userDomainSrv.GetUserInfo(10)
	dto := &dtos.UserDto{
		Name: info.Name,
	}
	c.JSON(http.StatusOK, dto)
}

// GetUserInfo get user info
func (m *UserAppHttpService) GetUserInfo(c *gin.Context) {
	response, err := m.userRpcClient.GetUserInfo(c.Request.Context(), &pbUser.GetUserInfoRequest{
		UniqueId: "1",
		NameMain: "2",
		NameSub:  "3",
	})
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	c.JSON(http.StatusOK, response.UserInfo)
}

func validateStruct() error {
	address := &dtos.AddressRequest{
		Street: "Docks",
		Planet: "person",
		Phone:  "none",
	}
	u := &dtos.UserInfoRequest{
		FirstName: "a",
		LastName:  "Smith",
		Age:       135,
		Email:     "Badger.Smith@gmail.com",
		Addresses: []*dtos.AddressRequest{address},
	}
	validate := validator.New()
	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(u)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Info(err)
		}
		for _, err := range err.(validator.ValidationErrors) {
			logger.Info(err.Namespace())
			logger.Info(err.Field())
			logger.Info(err.StructNamespace())
			logger.Info(err.StructField())
			logger.Info(err.Tag())
			logger.Info(err.ActualTag())
			logger.Info(err.Kind())
			logger.Info(err.Type())
			logger.Info(err.Value())
			logger.Info(err.Param())
		}
		// from here you can create your own error messages in whatever language you wish
	}
	// save user to database
	return err
}

func validateStruct2() error {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	//var trans ut.Translator
	validate := validator.New()
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	type User struct {
		Username string `validate:"required"`
	}
	var user User

	err := validate.Struct(user)
	if err != nil {

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			// can translate each error one at a time.
			fmt.Println(e.Translate(trans))
		}
	}
	//u := &dtos.UserRequest{
	//	FirstName: "a",
	//}
	//validate := validator.New()
	//return validate.Struct(user)
	return nil
}
