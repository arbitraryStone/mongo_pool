/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-22 17:14:52
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 22:54:52
 * @FilePath: /dbpool/mongo/options.go
 * @Description: mongo配置中心
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClientOptions struct {
	url string
}

func (opt *MongoClientOptions) GetConnOpts() *options.ClientOptions {
	return options.Client().ApplyURI(opt.url)
}

func MongoClientOptionsGenerator(url string) (*MongoClientOptions, error) {
	return &MongoClientOptions{url: url}, nil
}
