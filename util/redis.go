package util

import (
	"gopkg.in/redis.v5"
	"log"
	"sync"
	"time"
)

var (
	redisConn *redis.Client
	syncOnce  sync.Once

	host  string = REDIS_host
	pswd  string = REDIS_pswd
	dbNum int    = REDIS_dbNum
)

const (
	DEFAULT_EXPIRATION = 60

	LOGFMT_ERROR                  = "error: %v \n"
	LOGFMT_REDIS_DEL_OK           = "redis del ok, key: %v \n"
	LOGFMT_REDIS_DEL_ERROR        = "redis del error, key: %v  error: %s \n"
	LOGFMT_REDIS_SET_EXPIRE_ERROR = "redis set expire error, key: %s  error: %s \n"
)

func GetClient() *redis.Client {
	syncOnce.Do(
		func() {
			redisConn = redis.NewClient(&redis.Options{
				Addr:     host,
				Password: pswd,
				DB:       dbNum,
			})
			log.Printf("INFO(%s): Connect to redis status.", redisConn.Ping().String())
		},
	)
	return redisConn
}

// -------- base funcs ----------------------

// --- op: string

func GetRedisValue(key string) (string, error) {
	return GetClient().Get(key).Result()
}

func SetRedisValue(key, value string, expiration int) error {
	return GetClient().Set(key, value, time.Second*time.Duration(expiration)).Err()
}

func delRedisValue(keys ...string) error {
	return GetClient().Del(keys...).Err()
}

func DelRedisKeys(keys ...string) error {
	err := delRedisValue(keys...)
	if err != nil {
		log.Printf(LOGFMT_REDIS_DEL_ERROR, keys, err)
	} else {
		log.Printf(LOGFMT_REDIS_DEL_OK, keys)
	}
	return err
}

func incrReidsValue(key string) (int64, error) {
	return GetClient().Incr(key).Result()
}

func setRedisKeyExpiration(key string, expiration int) error {
	return GetClient().Expire(key, time.Second*time.Duration(expiration)).Err()
}

func boolKeyExists(key string) bool {
	return GetClient().Exists(key).Val()
}

// --- op: list

// add to the end
func AddRedisListValueRight(key, value string, expiration int) (int64, error) {
	client := GetClient()
	nowCount, err := client.RPush(key, value).Result()
	if err != nil {
		return 0, err
	}
	if err := client.Expire(key, time.Second*time.Duration(expiration)).Err(); err != nil {
		log.Printf(LOGFMT_REDIS_SET_EXPIRE_ERROR, key, err)
	}

	return nowCount, err
}

// add to the beginning
func AddRedisListValueLeft(key, value string, expiration int) error {
	client := GetClient()
	pipe := client.Pipeline()
	defer pipe.Close()

	pipe.LPush(key, value)
	pipe.Expire(key, time.Second*time.Duration(expiration))
	if _, err := pipe.Exec(); err != nil {
		log.Printf(LOGFMT_ERROR, err)
		return err
	}
	return nil
}

func GetRedisListRangeByStartCount(key string, start, count int) ([]string, error) {
	client := GetClient()
	stop := start + count
	return client.LRange(key, int64(start), int64(stop)).Result()
}

func GetRedisListRangeByStartEnd(key string, start, end int) ([]string, error) {
	client := GetClient()
	return client.LRange(key, int64(start), int64(end)).Result()
}

func GetRedisListLen(key string) (int64, error) {
	client := GetClient()
	return client.LLen(key).Result()
}

// --- op: set

func stringSlice2InterfaceSlice(input []string) []interface{} {
	result := make([]interface{}, len(input))
	for i, s := range input {
		result[i] = s
	}
	return result
}

func AddRedisSetValue(key string, values ...string) error {
	client := GetClient()
	pipe := client.Pipeline()
	defer pipe.Close()

	pipe.SAdd(key, stringSlice2InterfaceSlice(values)...)
	pipe.Expire(key, time.Second*time.Duration(DEFAULT_EXPIRATION))

	if _, err := pipe.Exec(); err != nil {
		log.Printf(LOGFMT_ERROR, err)
		return err
	}
	return nil
}


func AddRedisSetValuePerm(key string, values ...string) error {
	client := GetClient()
	pipe := client.Pipeline()
	defer pipe.Close()

	pipe.SAdd(key, stringSlice2InterfaceSlice(values)...)
	//pipe.Expire(key, time.Second*time.Duration(second))

	if _, err := pipe.Exec(); err != nil {
		log.Printf(LOGFMT_ERROR, err)
		return err
	}
	return nil
}


func ScanRedisSetValue(key string, scan string) ([]string, error) {
	client := GetClient()

	keys,_,err := client.SScan(key, 0, scan, 10).Result()

	return keys, err
}

func CheckRedisSetValue(key, value string) bool {
	client := GetClient()
	return client.SIsMember(key, value).Val()
}

func GetRedisSetList(key string) []string {
	client := GetClient()
	return client.SMembers(key).Val()
}

func GetRedisSetCount(key string) int64 {
	client := GetClient()
	return client.SCard(key).Val()
}

// --- op: hash

func GetRedisHashValue(key, field string) (string, error) {
	return GetClient().HGet(key, field).Result()
}


func GetRedisHashValueAll(key string) (map[string]string, error) {
	return GetClient().HGetAll(key).Result()
}

func SetRedisHashValue(key, field, value string) error {
	return GetClient().HSet(key, field, value).Err()
}


