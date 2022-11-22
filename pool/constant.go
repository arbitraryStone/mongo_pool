/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-20 00:06:32
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 11:06:08
 * @FilePath: /dbpool/pool/constant.go
 * @Description: 常量配置中心
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package pool

import "errors"

const (
	// PoolGetModeStrict 在实际创建连接数达上限后，池子中没有连接时不会新建连接
	PoolGetModeStrict = iota
	// PoolGetModeLoose 在实际创建连接数达上限后，池子中没有连接时会新建连接
	PoolGetModeLoose
)

var (
	// ErrPoolInit 连接池初始化出错
	ErrPoolInit = errors.New("Pool init occurred error")
	// ErrGetTimeout 获取连接超时
	ErrGetTimeout = errors.New("Get client timeout from dbpool")
	// ErrConn 创建连接发生错误
	ErrDialConn = errors.New("Client dailing connection failed")
	// ErrPoolIsClosed 连接池已关闭
	ErrPoolIsClosed = errors.New("Pool is closed")
	// ErrPoolMode 连接池模式错误
	ErrPoolMode = errors.New("Pool mode is unknow")
	//创建客户端错误
	ErrCreateClient = errors.New("Pool create client error")
	//连接池已满
	ErrPoolIsFull = errors.New("Pool is full")
)
