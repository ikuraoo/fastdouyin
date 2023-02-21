package configure

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var ctx = context.Background()

var rdb *redis.Client

// redis 初始化
func InitRedis() {
	ip := viper.GetString("Redis.ip")
	port := viper.GetString("Redis.port")
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", ip, port),
		//Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}

// ProxyCache 缓存层
type ProxyCache struct {
}

// 代理层
var (
	proxyCacheOperation ProxyCache
)

func NewProxyIndexMap() *ProxyCache {
	return &proxyCacheOperation
}

// GetVideoFavor 获取点赞状态 ret： true点赞 false未点赞
func (p *ProxyCache) GetVideoFavor(uid, vid int64) bool {
	key := fmt.Sprintf("%dfavor:", uid)
	return rdb.SIsMember(ctx, key, vid).Val()
}

func (p *ProxyCache) SetFunc(key, value string) string {
	InitRedis()
	rdb.SAdd(ctx, key, value)
	v := rdb.Get(ctx, key).Val()
	fmt.Println(v)
	return v
}

// SetVideoFavor isFavor: true点赞 false取消点赞
func (p *ProxyCache) SetVideoFavor(uid, vid int64, isFavor bool) {
	key := fmt.Sprintf("%dfavor:", uid)
	if isFavor {
		rdb.SAdd(ctx, key, vid)
		return
	}
	rdb.SRem(ctx, key, vid)
}

// GetAFollowB 判断A是否关注了B
func (p *ProxyCache) GetAFollowB(a, b int64) bool {
	key := fmt.Sprintf("%dfollow:", a)
	return rdb.SIsMember(ctx, key, b).Val()
}

// SetAFollowB isFollowed：true已关注 false未关注
func (p *ProxyCache) SetAFollowB(a, b int64, isFollowed bool) {
	key := fmt.Sprintf("%dfollow:", a)
	if isFollowed {
		rdb.SAdd(ctx, key, b)
		return
	}
	rdb.SRem(ctx, key, b)
}
