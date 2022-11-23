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
	"sync"
	"time"

	"github.com/arbitrarystone/dbpool/mongo"
	"github.com/arbitrarystone/dbpool/pool"
)

func main() {
	p, err := pool.NewPool(5, 10, pool.PoolGetModeStrict)
	if err != nil {
		panic("new pool error")
	}
	g := mongo.NewGenerator("mongodb://localhost:27017")
	p.RegisterClientGenerator(g)
	if err := p.InitPool(); err != nil {
		fmt.Errorf("%s", err.Error())
		panic("init pool error")
	}
	wg := sync.WaitGroup{}
	wg.Add(4)
	//获取连接
	go func() {
		for i := 0; i < 1000; i++ {
			client, err := p.Get(2 * time.Second)
			if err != nil {
				fmt.Printf("get mongo client error:%v\n", err)
				continue
			}
			v, ok := client.(*mongo.MongoClient)
			if !ok {
				panic("client type error")
			}
			_, err = v.GetCount("FilmDB", "area")
			if err != nil {
				fmt.Printf("get count error:%v", err)
				return
			}
			//fmt.Printf("count:%d\n", count)
			//time.Sleep(1000 * time.Millisecond)
			client.Close()
		}
		wg.Done()
	}()
	//获取连接
	go func() {
		for i := 0; i < 1000; i++ {
			client, err := p.Get(2 * time.Second)
			if err != nil {
				fmt.Printf("get mongo client error:%v\n", err)
				continue
			}
			v, ok := client.(*mongo.MongoClient)
			if !ok {
				panic("client type error")
			}
			_, err = v.GetCount("FilmDB", "area")
			if err != nil {
				fmt.Printf("get count error:%v", err)
				return
			}
			//time.Sleep(1000 * time.Millisecond)
			client.Close()
		}
		wg.Done()
	}()
	//打印连接数
	go func() {
		for i := 0; i < 1000; i++ {
			p.Dump()
			time.Sleep(2 * time.Second)
		}
		wg.Done()
	}()
	//释放连接
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(10 * time.Second)
			c, err := p.Get(2 * time.Second)
			if err != nil {
				fmt.Printf("get mongo client error:%v\n", err)
				continue
			}
			c.Release()
		}
		wg.Done()
	}()
	wg.Wait()
}
