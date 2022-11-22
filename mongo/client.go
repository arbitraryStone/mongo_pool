/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-22 10:01:39
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 15:49:27
 * @FilePath: /dbpool/mongo/client.go
 * @Description:
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package mongo

import (
	"context"

	"github.com/arbitrarystone/dbpool/pool"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoClient struct {
	client *mongo.Client
	pool   *pool.Pool
}

/**
 * @description: 关闭客户端连接 若属于连接池，则放回连接池中复用，不属于连接池，释放掉资源
 * @return {*}
 */
func (c *MongoClient) Close() {
	//不属于连接池，直接释放掉资源
	if c.pool == nil {
		c.Release()
		return
	}
	//连接池已关闭，直接释放
	if c.pool.IsClose() {
		c.Release()
		return
	}
	if c.client == nil {
		return
	}
	//属于连接池，放回所属连接池
	c.pool.Push(c)
}

func (c *MongoClient) Release() {
	c.client.Disconnect(context.TODO())
	if c.pool != nil {
		c.pool.DecSize()
	}
	c.pool = nil
	c.client = nil
}

func (c *MongoClient) SetPool(p *pool.Pool) {
	c.pool = p
}
