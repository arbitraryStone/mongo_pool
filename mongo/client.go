/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-22 10:01:39
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 23:29:09
 * @FilePath: /dbpool/mongo/client.go
 * @Description:
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package mongo

import (
	"context"
	"sync"

	"github.com/arbitrarystone/dbpool/pool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	pool   *pool.Pool
	// dbs    map[string]*mongo.Database
	dbs sync.Map
	//collections	map[string]map[string]*mongo.Collection
	collections sync.Map
	opts        *MongoClientOptions
}

/**
 * @description: mongo客户端生成器
 * @param {string} url
 * @return {*}
 */
func MongoClientGenerator(url string) (*MongoClient, error) {
	client := new(MongoClient)
	//获取配置生成器
	opts, err := MongoClientOptionsGenerator(url)
	if err != nil {
		return nil, err
	}
	client.opts = opts

	connOpt := opts.GetConnOpts()
	//获取mongo连接
	c, err := mongo.Connect(context.TODO(), connOpt)
	if err != nil {
		return nil, err
	}
	client.client = c
	return client, nil
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

/**
 * @description: 释放客户端连接
 * @return {*}
 */
func (c *MongoClient) Release() {
	c.client.Disconnect(context.TODO())
	if c.pool != nil {
		c.pool.DecSize()
	}
	c.pool = nil
	c.client = nil
}

/**
 * @description: 设置客户端所属连接池
 * @param {*pool.Pool} p
 * @return {*}
 */
func (c *MongoClient) SetPool(p *pool.Pool) {
	c.pool = p
}

func (c *MongoClient) GetCollection(dbName string, collectionName string) *mongo.Collection {
	v, ok := c.collections.Load(dbName)
	if ok {
		dbs := v.(sync.Map)
		vv, ook := dbs.Load(collectionName)
		if ook {
			coll := vv.(*mongo.Collection)
			return coll
		}
		db := c.GetDB(dbName)
		coll := db.Collection(collectionName)
		dbs.Store(collectionName, coll)
		return coll
	}
	db := c.GetDB(dbName)
	coll := db.Collection(collectionName)
	var m sync.Map
	m.Store(collectionName, coll)
	c.collections.Store(dbName, m)
	return coll
}

func (c *MongoClient) GetDB(dbName string) *mongo.Database {
	v, ok := c.dbs.Load(dbName)
	if ok {
		return v.(*mongo.Database)
	}
	db := c.client.Database(dbName)
	c.dbs.Store(dbName, db)
	return db
}

func (c *MongoClient) FindOne(dbName string, collectionName string,
	filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {

	coll := c.GetCollection(dbName, collectionName)
	return coll.FindOne(context.TODO(), filter, opts...)
}

func (c *MongoClient) Find(dbName string, collectionName string,
	filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {

	coll := c.GetCollection(dbName, collectionName)
	return coll.Find(context.TODO(), filter, opts...)
}

func (c *MongoClient) UpdateOne(dbName string, collectionName string,
	filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	coll := c.GetCollection(dbName, collectionName)
	return coll.UpdateOne(context.TODO(), filter, update, opts...)
}

func (c *MongoClient) DeleteOne(dbName string, collectionName string,
	filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {

	coll := c.GetCollection(dbName, collectionName)
	return coll.DeleteOne(context.TODO(), filter, opts...)
}

func (c *MongoClient) InsertOne(dbName string, collectionName string,
	data interface{}) (*mongo.InsertOneResult, error) {

	coll := c.GetCollection(dbName, collectionName)
	return coll.InsertOne(context.TODO(), data)
}

func (c *MongoClient) InsertMany(dbName string, collectionName string,
	data []interface{}) (*mongo.InsertManyResult, error) {

	coll := c.GetCollection(dbName, collectionName)
	return coll.InsertMany(context.TODO(), data)
}

func (c *MongoClient) GetCount(dbName string, collectionName string) (int64, error) {
	coll := c.GetCollection(dbName, collectionName)
	return coll.CountDocuments(context.TODO(), bson.D{})
}

func (c *MongoClient) CreateIndex(dbName string, collectionName string) {}
