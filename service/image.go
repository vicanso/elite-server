// Copyright 2020 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elton"
	"github.com/vicanso/tiny/pb"
	"google.golang.org/grpc"
)

var (
	grpcConn *grpc.ClientConn

	tinyConfig = config.GetTinyConfig()
)

type (
	// ImageOptimParams image optim params
	ImageOptimParams struct {
		Data       []byte
		Type       string
		SourceType string
		Quality    int
		Width      int
		Height     int
		Crop       int
	}
	// ImageSrv image service
	ImageSrv struct{}
)

func init() {
	done := make(chan int)
	go func() {
		opts := make([]grpc.DialOption, 0)
		opts = append(opts, grpc.WithInsecure())
		if util.IsProduction() {
			opts = append(opts, grpc.WithBlock())
		}
		target := fmt.Sprintf("%s:%d", tinyConfig.Host, tinyConfig.Port)
		conn, err := grpc.Dial(target, opts...)
		if err != nil {
			panic(err)
		}
		done <- 1
		grpcConn = conn
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		panic(errors.New("grpc dial timeout"))
	}
}

// optim image optim
func (*ImageSrv) optim(params *ImageOptimParams) (data []byte, err error) {
	client := pb.NewOptimClient(grpcConn)
	ctx, cancel := context.WithTimeout(context.Background(), tinyConfig.Timeout)
	defer cancel()
	in := &pb.OptimRequest{
		Data:    params.Data,
		Quality: uint32(params.Quality),
		Width:   uint32(params.Width),
		Height:  uint32(params.Height),
		Crop:    uint32(params.Crop),
	}
	switch params.Type {
	case "png":
		in.Output = pb.Type_PNG
	case "webp":
		in.Output = pb.Type_WEBP
	default:
		in.Output = pb.Type_JPEG
	}
	switch params.SourceType {
	case "png":
		in.Source = pb.Type_PNG
	case "webp":
		in.Source = pb.Type_WEBP
	default:
		in.Source = pb.Type_JPEG
	}
	reply, err := client.DoOptim(ctx, in)
	if err != nil {
		return
	}
	data = reply.Data
	return
}

// GetImageFromBucket get image from bucket
func (srv *ImageSrv) GetImageFromBucket(ctx context.Context, bucket, filename string, params ImageOptimParams) (data []byte, header http.Header, err error) {
	data, header, err = fileSrv.GetData(ctx, bucket, filename)
	if err != nil {
		return
	}
	contentType := header.Get("Content-Type")
	source := strings.Split(contentType, "/")[1]
	params.Data = data
	params.SourceType = source
	data, err = srv.optim(&params)
	if err != nil {
		return
	}
	header.Set(elton.HeaderContentType, "image/"+params.Type)

	return
}
