package catdog_mongo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"
)

var resMap sync.Map

type mongoCfg struct {
	Uri      string // mongodb://localhost:27017/?connect=direct
	Username string
	Password string
}

type ClientAndOpts struct {
	// Current client implement a pool inside.
	*mongo.Client
	*options.ClientOptions
}

const (
	min_pool_size = 10
	max_pool_size = min_pool_size * 10
)

func Watch(prefix string, value []byte) (err error) {
	defer xerror.RespErr(&err)
	clientAndOpts, err := buildMongoClient(prefix, value)
	xerror.Panic(err)
	xlog.InfoF("Watcher.catdog_mongo_plugin: build catdog_mongo_plugin client succeeded: %s", string(value))
	resMap.Store(prefix, clientAndOpts)
	return
}

func PickupMongoClient(prefix string) (*ClientAndOpts, error) {
	val, ok := resMap.Load(prefix)
	if !ok {
		return nil, errors.New("Pick up catdog_mongo_plugin client error: not found prefix: " + prefix)
	}
	return val.(*ClientAndOpts), nil
}

/**
 * Access to refresh client outside.
 */
func RefreshClient(prefix string, newOpts *options.ClientOptions) error {
	val, ok := resMap.Load(prefix)
	if !ok {
		return errors.New("Pick up catdog_mongo_plugin client error: not found prefix: " + prefix)
	}
	oldClient := val.(*ClientAndOpts)

	newMongoOpts := options.MergeClientOptions(oldClient.ClientOptions, newOpts)
	newMongoClient, err := mongo.Connect(context.TODO(), newMongoOpts)
	if err != nil {
		return errors.New("Mongo 连接异常：" + err.Error())
	}

	newClient := newMongoClientAndOpts(newMongoClient, newMongoOpts)

	// Client not empty, then release previous connection
	if !reflect.DeepEqual(oldClient.Client, &mongo.Client{}) {
		defer func() {
			if err := oldClient.Disconnect(context.TODO()); err != nil {
				xlog.Error("Mongo old client disconnect error: " + err.Error())
			}
		}()
	}

	if err := newClient.Ping(context.TODO(), nil); err != nil {
		return errors.New("Mongo Ping 异常：" + err.Error())
	}

	resMap.Store(prefix, newClient)
	return nil
}

func newMongoClientAndOpts(client *mongo.Client, clientOptions *options.ClientOptions) *ClientAndOpts {
	newClient := &ClientAndOpts{
		Client:        client,
		ClientOptions: clientOptions,
	}
	newClient.ClientOptions.SetMinPoolSize(min_pool_size)
	newClient.ClientOptions.SetMaxPoolSize(max_pool_size)
	return newClient
}

// initCatDog or update default client
func buildMongoClient(prefix string, value []byte) (*ClientAndOpts, error) {
	if len(value) == 0 {
		return nil, errors.New("watch时读取json出错,prefix=" + prefix)
	}

	var mcfg mongoCfg
	err := json.Unmarshal(value, &mcfg)
	if err != nil {
		return nil, errors.New("json unmarshal, error=" + err.Error())
	}
	if mcfg.Uri == "" {
		return nil, errors.New("未找到 uri 配置， " + prefix + "无效")
	}

	newOpts := options.Client().ApplyURI(mcfg.Uri).SetAuth(options.Credential{Username: mcfg.Username, Password: mcfg.Password})
	newMongoClient, err := mongo.Connect(context.TODO(), newOpts)
	if err != nil {
		return nil, errors.New("Mongo 连接异常：" + err.Error())
	}

	return newMongoClientAndOpts(newMongoClient, newOpts), nil
}
