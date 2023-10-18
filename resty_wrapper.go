package gocom

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"sync"
)

type RestyWrapper struct {
	client *resty.Client
}

var onceRestyWrapper sync.Once
var restyWrapper RestyWrapper

func DefaultRestyClient() RestyWrapper {
	onceRestyWrapper.Do(func() {
		client := resty.New()
		restyWrapper = RestyWrapper{
			client: client,
		}
	})
	return restyWrapper
}

func (receiver RestyWrapper) GetHttpClient() *http.Client {
	return receiver.client.GetClient()
}

func (receiver RestyWrapper) GetClient() *resty.Client {
	return receiver.client
}
