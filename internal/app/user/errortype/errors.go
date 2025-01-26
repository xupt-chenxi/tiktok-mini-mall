package errortype

import "github.com/pkg/errors"

var (
	ErrInputEmpty      = errors.New("用户输入不能为空")
	ErrEmailFormat     = errors.New("用户输入邮箱格式有误")
	PasswordMismatch   = errors.New("用户两次输入的密码不一致")
	ErrUserNotFound    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
)
