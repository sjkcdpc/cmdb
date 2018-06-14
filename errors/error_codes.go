package errors

import "fmt"

const (
	Ok Code = iota
	Failed
	Faileda

	PartialContent

	TokenInvalidate           = 11001  // token无效
	UserNotExist              = 11002  // 用户不存在
	IllegalRequest            = 11003  // 非法请求
	GoToLogin                 = 11004  // 重新登录
	RepeatLogin               = 11005  // 顶号
	LoginAccountAlreadyActive = 11006  // 账户已经处于激活状态
	LoginRegisterByCoupon     = 11007  // 需要注册邀请码
	LoginRegisterClosed       = 11008  // 关闭账户注册
	LoginRegisterByMobile     = 11009  // 需要手机验证码
	ShopInfoNotFound          = 150001 // 商城信息获取失败
	ShopCallBackAgain         = 150002 // 轮询请求失败
	ShopDataException         = 150003 // 数据异常
	ShopCommitRmb             = 150004 // rmb购买下单失败
	DiamondBuyShop            = 150005 // 钻石购买道具失败
	QueryGeted                = 150006 // 已领取
	NotPay                    = 150007 // 查询未成功,米大师消费失败
	OrderMiss                 = 150008 // 订单不存在
)

type Code int

// errorString is a trivial implementation of error.
type Type struct {
	Code    Code
	Message string
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &Type{Failed, text}
}

// New returns an error that formats as the given text.
func NewCode(code Code) error {
	return NewCodeString(code, "")
}

// New returns an error that formats as the given text.
func NewCodeString(code Code, format string, a ...interface{}) error {
	return &Type{code, fmt.Sprintf(format, a...)}
}

// New returns an error that formats as the given text.
func NewCodeError(code Code, err error) error {
	return &Type{code, err.Error()}
}

func (e *Type) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
