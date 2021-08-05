package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisSessionConfig struct {
	KeyPrefix string
	Expire    time.Duration
}

type RedisSession struct {
	redis     *redis.Client
	keyPrefix string
	expire    time.Duration
}

func NewRedisSession(redisClient *redis.Client, config RedisSessionConfig) *RedisSession {
	return &RedisSession{
		redis:     redisClient,
		keyPrefix: config.KeyPrefix,
		expire:    config.Expire,
	}
}

func (r *RedisSession) GenerateToken(userId int64) string {
	s := fmt.Sprintf("%d%d%d", time.Now().UnixNano()/1000, userId, RandInt64(1000, 9999))
	return Md5(s)
}

func (r *RedisSession) GetFullKey(userId int64, token string) string {
	// adminapi:login:用户id:token
	return fmt.Sprintf("%s%d:%s", r.keyPrefix, userId, token)
}

// 清除用户其它已登录的redis
func (r *RedisSession) CleanOtherLogined(userId int64) (err error) {
	// adminapi:login:用户id:*
	prefix := fmt.Sprintf("%s%d:*", r.keyPrefix, userId)

	results, errs := r.redis.Keys(context.Background(), prefix).Result()
	if errs != nil {
		return errors.New(fmt.Sprintf("查询用户其他redis token列表失败, prefix: %v, errs:%v ", prefix, errs))
	}

	if len(results) == 0 {
		return nil
	}

	fmt.Printf("开始删除redis Keys[%v]\n", results)

	_, errs1 := r.redis.Del(context.Background(), results...).Result()
	if errs1 != nil {
		return errors.New(fmt.Sprintf("删除失败 list[%v], errs: %v  \n", results, errs1))
	}

	return nil
}

func (r *RedisSession) Login(userId int64) (string, error) {
	token := r.GenerateToken(userId)
	fullKey := r.GetFullKey(userId, token)

	_, setErr := r.redis.Set(context.Background(), fullKey, userId, r.expire).Result()
	if setErr != nil {
		return "", setErr
	}

	return token, nil
}

func (r *RedisSession) Refresh(userId int64, token string) error {
	fullKey := r.GetFullKey(userId, token)

	_, err := r.redis.Expire(context.Background(), fullKey, r.expire).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisSession) Logout(userId int64, token string) error {
	fullKey := r.GetFullKey(userId, token)

	_, err := r.redis.Del(context.Background(), fullKey).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisSession) IsLogined(userId int64, token string) error {
	fullKey := r.GetFullKey(userId, token)

	redisUserId, err := r.redis.Get(context.Background(), fullKey).Int64()
	if err != nil {
		return err
	}

	if redisUserId != userId {
		return errors.New("不是当前登录用户token")
	}

	return nil
}
