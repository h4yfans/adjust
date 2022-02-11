package tool

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	client "github.com/h4yfans/adjust/http"
)

func TestNewTool(t *testing.T) {
	type args struct {
		parallel uint
		client   client.HTTPClient
		address  []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				address:  []string{"http://adjust.com", "google.com", "facebook.com", "yahoo.com", "yandex.com", "twitter.com", "reddit.com/r/funny", "reddit.com/r/notfunny", "baroquemusiclibrary.com"},
				parallel: 10,
				client: &client.ClientMock{
					DoFunc: func(request *http.Request) (*http.Response, error) {
						return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("test-response"))}, nil
					},
				}},
			wantErr: false,
		},
		{
			name: "invalid url",
			args: args{
				address:  []string{"invalid-url"},
				parallel: 10,
				client: &client.ClientMock{
					DoFunc: func(request *http.Request) (*http.Response, error) {
						return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("test-response"))}, nil
					},
				}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewTool(tt.args.parallel, tt.args.client).Run(tt.args.address...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
