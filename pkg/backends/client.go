package backends

import (
	"errors"
	"fmt"
	"strings"

	"github.com/abtreece/confd/pkg/backends/consul"
	"github.com/abtreece/confd/pkg/backends/dynamodb"
	"github.com/abtreece/confd/pkg/backends/env"
	"github.com/abtreece/confd/pkg/backends/etcd"
	"github.com/abtreece/confd/pkg/backends/file"
	"github.com/abtreece/confd/pkg/backends/redis"
	"github.com/abtreece/confd/pkg/backends/s3"
	"github.com/abtreece/confd/pkg/backends/ssm"
	"github.com/abtreece/confd/pkg/backends/vault"
	"github.com/abtreece/confd/pkg/backends/zookeeper"
	"github.com/abtreece/confd/pkg/log"
	"github.com/abtreece/confd/pkg/util"
)

// The StoreClient interface is implemented by objects that can retrieve
// key/value pairs from a backend store.
type StoreClient interface {
	GetValues(keys []string) (map[string]string, error)
	WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error)
}

// New is used to create a storage client based on our configuration.
func New(config Config, revision *util.Revision) (StoreClient, error) {

	if config.Backend == "" {
		config.Backend = "etcd"
	}
	backendNodes := config.BackendNodes

	switch config.Backend {
	case "consul":
		log.Info("Backend source(s) set to " + strings.Join(backendNodes, ", "))
		return consul.New(config.BackendNodes, config.Scheme,
			config.ClientCert, config.ClientKey,
			config.ClientCaKeys,
			config.BasicAuth,
			config.Username,
			config.Password,
		)
	case "etcd":
		log.Info("Backend source(s) set to " + strings.Join(backendNodes, ", "))
		return etcd.NewEtcdClient(backendNodes, config.ClientCert, config.ClientKey, config.ClientCaKeys, config.ClientInsecure, config.BasicAuth, config.Username, config.Password)
	case "zookeeper":
		log.Info("Backend source(s) set to " + strings.Join(backendNodes, ", "))
		return zookeeper.NewZookeeperClient(backendNodes)
	case "redis":
		log.Info("Backend source(s) set to " + strings.Join(backendNodes, ", "))
		return redis.NewRedisClient(backendNodes, config.ClientKey, config.Separator)
	case "env":
		return env.NewEnvClient()
	case "file":
		log.Info("Backend source(s) set to " + strings.Join(config.YAMLFile, ", "))
		return file.NewFileClient(config.YAMLFile, config.Filter)
	case "vault":
		log.Info("Backend source(s) set to " + strings.Join(backendNodes, ", "))
		vaultConfig := map[string]string{
			"app-id":    config.AppID,
			"user-id":   config.UserID,
			"role-id":   config.RoleID,
			"secret-id": config.SecretID,
			"username":  config.Username,
			"password":  config.Password,
			"token":     config.AuthToken,
			"cert":      config.ClientCert,
			"key":       config.ClientKey,
			"caCert":    config.ClientCaKeys,
			"path":      config.Path,
		}
		return vault.New(backendNodes[0], config.AuthType, vaultConfig)
	case "dynamodb":
		table := config.Table
		log.Info("DynamoDB table set to " + table)
		return dynamodb.NewDynamoDBClient(table)
	case "ssm":
		return ssm.New()
	case "s3":
		bucket := config.Bucket
		key := config.Key
		log.Info(fmt.Sprintf("Backend bucket set to %s, key set to %s", bucket, key))
		return s3.NewS3Client(bucket, key, revision)
	}
	return nil, errors.New("Invalid backend")
}
