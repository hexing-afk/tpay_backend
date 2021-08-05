package utils

import (
	"context"
	"github.com/mojocn/base64Captcha"
	"time"

	"github.com/go-redis/redis/v8"
)

type CaptchaConfig struct {
	KeyPrefix string
	Expire    time.Duration
}

type Captcha struct {
	redis     *redis.Client
	keyPrefix string
	expire    time.Duration
}

func NewCaptcha(redisClient *redis.Client, config CaptchaConfig) *Captcha {
	return &Captcha{
		redis:     redisClient,
		keyPrefix: config.KeyPrefix,
		expire:    config.Expire,
	}
}

// 创建图片验证码并将值存入redis中
func (r *Captcha) GenerateCaptcha() (id string, b64s string, err error) {

	driver := &base64Captcha.DriverDigit{
		Height:   80,
		Width:    240,
		Length:   4,
		MaxSkew:  0.7,
		DotCount: 80,
	}

	store := base64Captcha.DefaultMemStore
	c := base64Captcha.NewCaptcha(driver, store)

	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}

	// 存入到redis中
	_, setErr := r.redis.Set(context.Background(), r.keyPrefix+id, answer, r.expire).Result()
	if setErr != nil {
		return "", "", setErr
	}

	b64s = item.EncodeB64string()

	return id, b64s, err
}

// 校验验证码,成功后会删除redis中的验证码
func (r *Captcha) VerifyCaptcha(id, answer string) (bool, error) {
	key := r.keyPrefix + id

	val, err := r.redis.Get(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	if val == answer {
		r.redis.Del(context.Background(), key)
		return true, nil
	}

	return false, nil
}
