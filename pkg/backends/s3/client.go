package s3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abtreece/confd/pkg/log"
	"github.com/abtreece/confd/pkg/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

type Client struct {
	client       *s3.S3
	bucket       string
	key          string
	vars         map[string]string
	currRevision string
	revision     *util.Revision
}

// NewS3Client creates a new S3 backend client.
// Client credentials and region are configured through the usual environment
// variables, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_REGION.
func NewS3Client(bucket string, key string, revision *util.Revision) (*Client, error) {
	if bucket == "" {
		return nil, fmt.Errorf("Bucket must be defined")
	}
	if key == "" {
		return nil, fmt.Errorf("Key must be defined")
	}

	// Create session and validate credentials.
	session, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Could not configure session: %s", err)
	}

	// Create S3 client.
	client := s3.New(session)

	// Check if the bucket exists and is accessible.
	_, err = client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("Could not access bucket %s: %s", bucket, err)
	}

	// Check if the key exists and is accessible.
	_, err = client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("Could not access key %s: %s", key, err)
	}

	return &Client{client: client, bucket: bucket, key: key, revision: revision}, nil
}

// GetValues retrieves the values for the given keys from an S3 file.
// GetValues only fetches the file from S3 at most once per revision.
func (c *Client) GetValues(keys []string) (map[string]string, error) {
	rev := c.revision.Current()
	if c.currRevision != rev {
		log.Debug("New revision %s, fetching file from S3", rev)
		vars := make(map[string]string)
		result, err := c.client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(c.key),
		})
		if err != nil {
			return nil, err
		}
		defer result.Body.Close()
		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return nil, err
		}

		// Unmarshal the file based on extension.
		switch filepath.Ext(c.key) {
		case ".json":
			fileMap := make(map[string]interface{})
			err = json.Unmarshal(body, &fileMap)
			if err != nil {
				return nil, err
			}
			err = nodeWalk(fileMap, "/", vars)
		case "", ".yml", ".yaml":
			fileMap := make(map[interface{}]interface{})
			err = yaml.Unmarshal(body, &fileMap)
			if err != nil {
				return nil, err
			}
			err = nodeWalk(fileMap, "/", vars)
		default:
			err = fmt.Errorf("Invalid file extension, only json or yaml files allowed")
		}
		if err != nil {
			return nil, err
		}

		// Update cache.
		c.vars = vars
		c.currRevision = rev
	}

	// Filter out vars.
	filteredVars := make(map[string]string)
	for k := range c.vars {
		for _, key := range keys {
			if strings.HasPrefix(k, key) {
				filteredVars[k] = c.vars[k]
			}
		}
	}
	return filteredVars, nil
}

// nodeWalk recursively descends nodes, updating vars.
func nodeWalk(node interface{}, key string, vars map[string]string) error {
	switch node := node.(type) {
	case []interface{}:
		for i, j := range node {
			key := path.Join(key, strconv.Itoa(i))
			nodeWalk(j, key, vars)
		}
	case map[interface{}]interface{}:
		for k, v := range node {
			key := path.Join(key, k.(string))
			nodeWalk(v, key, vars)
		}
	case map[string]interface{}:
		for k, v := range node {
			key := path.Join(key, k)
			nodeWalk(v, key, vars)
		}
	case string:
		vars[key] = node
	case int:
		vars[key] = strconv.Itoa(node)
	case bool:
		vars[key] = strconv.FormatBool(node)
	case float64:
		vars[key] = strconv.FormatFloat(node, 'f', -1, 64)
	}
	return nil
}

// WatchPrefix is not implemented.
func (c *Client) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	<-stopChan
	return 0, nil
}
