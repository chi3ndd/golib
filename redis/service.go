package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	rd "github.com/go-redis/redis/v8"

	"github.com/chi3ndd/golib/clock"
	"github.com/chi3ndd/golib/log"
)

type service struct {
	logger log.Logger
	con    *rd.Client
	mutex  *sync.Mutex
	config Config
}

func NewService(conf Config, tlsConf *tls.Config) Service {
	logger, _ := log.New(Module, log.DebugLevel, true, os.Stdout)
	con := rd.NewClient(&rd.Options{
		Addr:      conf.Address,
		Password:  conf.Password,
		DB:        conf.Db,
		TLSConfig: tlsConf,
	})
	r := &service{logger: logger, con: con, mutex: &sync.Mutex{}, config: conf}
	// Monitor
	r.monitor()
	// Success
	return r
}

func (r *service) monitor() {
	// Reconnect connection
	go func() {
		for {
			if r.con != nil {
				if r.Ping() != nil {
					r.logger.Info("connection closed")
					r.mutex.Lock()
					for {
						r.con = rd.NewClient(&rd.Options{
							Addr:      r.config.Address,
							Password:  r.config.Password,
							DB:        r.config.Db,
							TLSConfig: nil,
						})
						if r.Ping() == nil {
							r.logger.Info("reconnect success!")
							break
						}
						r.logger.Info("reconnecting ...")
						// Sleep
						clock.Sleep(clock.Second * 3)
					}
					r.mutex.Unlock()
				}
			}
			// Sleep
			clock.Sleep(clock.Second * 10)
		}
	}()
}

func (r *service) Ping() error {
	// Success
	return r.con.Ping(context.Background()).Err()
}

func (r *service) Delete(keys ...string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Success
	return r.con.Del(context.Background(), keys...).Err()
}

func (r *service) Expire(key string, ttl clock.Duration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Success
	return r.con.Expire(context.Background(), key, time.Duration(ttl)).Err()
}

func (r *service) ExpireAt(key string, tm time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Success
	return r.con.ExpireAt(context.Background(), key, tm).Err()
}

func (r *service) Set(key, value string, ttl clock.Duration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Success
	return r.con.Set(context.Background(), key, value, time.Duration(ttl)).Err()
}

func (r *service) SetObject(key string, value interface{}, ttl clock.Duration) error {
	bts, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// Success
	return r.Set(key, string(bts), ttl)
}

func (r *service) Get(key string) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Success
	result, err := r.con.Get(context.Background(), key).Result()
	if err == rd.Nil {
		return result, errors.New(NotFoundError)
	}
	return result, nil
}

func (r *service) GetObject(key string, pointer interface{}) error {
	result, err := r.Get(key)
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(result), pointer); err != nil {
		return err
	}
	// Success
	return nil
}
