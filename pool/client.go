/*
 * @Author: arbitrarystone arbitrarystone@163.com
 * @Date: 2022-11-21 22:21:45
 * @LastEditors: arbitrarystone arbitrarystone@163.com
 * @LastEditTime: 2022-11-22 18:45:04
 * @FilePath: /dbpool/pool/client.go
 * @Description: mongodb客户端
 *
 * Copyright (c) 2022 by arbitrarystone arbitrarystone@163.com, All Rights Reserved.
 */
package pool

type ClientGenerator interface {
	Generator() (Client, error)
}
type Client interface {
	Close()
	Release()
	SetPool(p *Pool)
}
