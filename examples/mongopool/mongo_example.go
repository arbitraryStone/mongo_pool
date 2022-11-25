/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-22 15:57:24
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 15:57:59
 * @FilePath: /dbpool/examples/mongopool/main.go
 * @Description: mongo连接池使用示例
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/arbitrarystone/dbpool/mongo"
	"github.com/arbitrarystone/dbpool/pool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AreaInfo struct {
	ID       string `bson:"_id,omitempty"` //id唯一键
	AreaName string `bson:"area_name"`
}

func main() {
	//新建连接池
	p, err := pool.NewPool("mongo", 5, 10, pool.PoolGetModeStrict)
	if err != nil {
		panic("new pool error")
	}
	defer p.Close()
	//创建mongo客户端生成器
	g := mongo.NewGenerator("mongodb://localhost:27017")
	//注册mongo客户端生成器
	p.RegisterClientGenerator(g)
	//初始化连接池
	if err := p.InitPool(); err != nil {
		panic("init pool error")
	}
	//获取客户端连接
	client, err := p.Get(2 * time.Second)
	if err != nil {
		panic("get client error")
	}
	defer client.Close()
	v, ok := client.(*mongo.MongoClient)
	if !ok {
		panic("client type error")
	}
	//获取文档数量
	count, err := v.GetCount("FilmDB", "area")
	if err != nil {
		fmt.Printf("get count error:%v", err)
		return
	}
	fmt.Printf("count:%d\n", count)

	//查询所有数据
	findOption := options.Find().SetProjection(bson.M{"_id": 1, "area_name": 1})
	cur, err := v.Find("FilmDB", "area", bson.D{}, findOption)
	if err != nil {
		fmt.Printf("find error:%v", err)
		return
	}
	var list []*AreaInfo
	for cur.Next(context.Background()) {
		var info AreaInfo
		if err = cur.Decode(&info); err != nil {
			fmt.Printf("decode error:%v", err)
			continue
		}
		list = append(list, &info)
	}
	fmt.Printf("find result:%+v", list)

	//查询一条数据
	filter := bson.M{"area_name": "xxx"}
	findOneOption := options.FindOne().SetProjection(bson.M{"area_name": 1})
	res2 := v.FindOne("FilmDB", "area", filter, findOneOption)
	if err != nil {
		fmt.Printf("find one error:%v", err)
		return
	}
	var findOneResult AreaInfo
	res2.Decode(&findOneResult)
	fmt.Printf("find one result:%+v", findOneResult)

	//插入一条数据
	insertOneData := &AreaInfo{
		AreaName: "xxx",
	}
	insertResult, err := v.InsertOne("FilmDB", "area", insertOneData)
	if err != nil {
		fmt.Printf("insert one error:%v", err)
		return
	}
	fmt.Printf("insert one result:%+v", insertResult)

}
