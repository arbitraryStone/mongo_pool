/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-22 15:57:24
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 15:57:59
 * @FilePath: /dbpool/examples/mongopool/main.go
 * @Description: mongo连接池功能测试
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package main

import (
	"fmt"
	"time"

	"github.com/arbitrarystone/dbpool"
	"github.com/arbitrarystone/dbpool/mongo"
	"github.com/arbitrarystone/dbpool/pool"
)

func MongoClientTest() {
	url := "mongodb://localhost:27017"
	g := mongo.NewGenerator(url)
	c, err := g.Generator()
	if err != nil {
		panic(err)
	}
	v := c.(*mongo.MongoClient)
	count, err := v.GetCount("FilmDB", "area")
	if err != nil {
		panic(err)
	}
	fmt.Printf("get count:%d", count)
}

func main() {
	p, err := dbpool.NewPool("mongo", "mongodb://localhost:27017", 1, 10, pool.PoolGetModeStrict)
	if err != nil {
		panic(err)
	}
	//获取连接
	go func() {
		for {
			client, err := dbpool.GetMongoClient(p, 1*time.Second)
			if err != nil {
				continue
			}
			_, err = client.GetCount("FilmDB", "area")
			if err != nil {
				fmt.Printf("get count error:%v", err)
				return
			}
			client.Close()
		}
	}()
	//获取连接
	go func() {
		for {
			client, err := dbpool.GetMongoClient(p, 1*time.Second)
			if err != nil {
				continue
			}
			_, err = client.GetCount("FilmDB", "area")
			if err != nil {
				fmt.Printf("get count error:%v", err)
				return
			}
			client.Close()
		}
	}()
	//打印连接数
	go func() {
		for {
			p.Dump()
			time.Sleep(2 * time.Second)
		}
	}()
	//释放连接
	go func() {
		for {
			time.Sleep(20 * time.Second)
			c, err := p.Get(2 * time.Second)
			if err != nil {
				fmt.Printf("get mongo client error:%v\n", err)
				continue
			}
			fmt.Printf("release")
			c.Release()
		}
	}()
	for {
	}
}
