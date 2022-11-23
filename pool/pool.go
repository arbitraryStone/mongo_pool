/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-21 22:25:02
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 22:19:01
 * @FilePath: /dbpool/pool/pool.go
 * @Description: 池化中心
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package pool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	clients  chan Client
	init     int32 //初始化客户端数
	idle     int32 //池中现有客户端数
	size     int32 //连接池管理客户端数数量
	capacity int32 //连接池最大容纳客户端数
	lock     sync.RWMutex
	mode     int //连接池模式
	ClientGenerator
}

/**
 * @description: 默认连接池 初始化客户端数：5 最大容量：10 模式：池满后不会新增
 * @return {pool实例指针 错误信息}
 */
func DefaultPool() (*Pool, error) {
	return NewPool(5, 10, PoolGetModeStrict)
}

/**
 * @description: 创建连接池
 * @param {int32} init 连接池大小(初始化连接数量)
 * @param {int32} cap 连接池容量(连接池最大容纳连接数)
 * @param {int} mode 连接池在池满时的连接状态(池满后继续创建 or 池满后不继续创建连接)
 * @return {*Pool, error} 创建的连接实例, 错误信息
 */
func NewPool(init int32, cap int32, mode int) (*Pool, error) {
	if init < 0 || cap <= 0 {
		return nil, ErrPoolInit
	}
	if init > cap {
		init = cap
	}
	pool := &Pool{
		clients:  make(chan Client, cap),
		init:     init,
		idle:     0,
		size:     0,
		capacity: cap,
		mode:     mode,
	}
	return pool, nil
}

/**
 * @description: 注册客户端生成器
 * @param {ClientGenerator} generator
 * @return {*}
 */
func (pool *Pool) RegisterClientGenerator(generator ClientGenerator) {
	pool.ClientGenerator = generator
}

/**
 * @description: 初始化连接池
 * @return {error}
 */
func (pool *Pool) InitPool() error {
	for i := int32(0); i < pool.init; i++ {
		client, err := pool.createClient()
		if err != nil {
			return ErrPoolInit
		}
		pool.Push(client)
	}
	return nil
}

/**
 * @description: 创建客户端
 * @return {*}
 */
func (pool *Pool) createClient() (Client, error) {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if pool.size >= pool.capacity {
		//池满后可以新建连接, 新建连接不受连接池管理, 使用完会立即关闭, 不会再放入连接池重用
		if pool.mode == PoolGetModeLoose {
			client, err := pool.ClientGenerator.Generator()
			if err != nil {
				return nil, ErrCreateClient
			}
			return client, err
		}
		return nil, ErrPoolIsFull
	}
	//池未满时，创建客户端
	client, err := pool.ClientGenerator.Generator()
	if err != nil {
		return nil, ErrCreateClient
	}
	client.SetPool(pool)
	//连接池管理数量增加
	pool.incSize()
	return client, nil
}

/**
 * @description: 从连接池中获取客户端
 * @param {time.Duration} timeout 获取客户端连接超时时间
 * @return {Client,error} 返回获取到的客户端，错误信息
 */
func (pool *Pool) Get(timeout time.Duration) (Client, error) {
	select {
	case client := <-pool.clients:
		pool.decIdle()
		return client, nil
	case <-time.After(timeout):
		return pool.createClient()
	}
}

func (pool *Pool) Push(c Client) {
	pool.clients <- c
	pool.incIdle()
}

func (pool *Pool) Close() {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if pool.IsClose() {
		return
	}
	t := time.NewTicker(1 * time.Millisecond)
	//异步释放连接
	go func() {
	loop:
		for {
			select {
			case client := <-pool.clients:
				pool.decIdle()
				client.Release()
			case <-t.C:
				if len(pool.clients) <= 0 {
					close(pool.clients)
					break loop
				}
			}
		}
	}()
}

func (pool *Pool) IsClose() bool {
	return pool == nil || pool.clients == nil
}

func (pool *Pool) incSize() {
	atomic.AddInt32(&pool.size, 1)
}

func (pool *Pool) DecSize() {
	atomic.AddInt32(&pool.size, -1)
}

func (pool *Pool) incIdle() {
	atomic.AddInt32(&pool.idle, 1)
}

func (pool *Pool) decIdle() {
	atomic.AddInt32(&pool.idle, -1)
}

func (pool *Pool) Size() int32 {
	return pool.size
}

func (pool *Pool) Idle() int32 {
	return pool.idle
}
func (pool *Pool) Capacity() int32 {
	return pool.capacity
}

func (pool *Pool) Dump() {
	fmt.Printf("pool.size:%d pool.idle:%d pool.capacity:%d\n",
		pool.Size(), pool.Idle(), pool.Capacity())
}
