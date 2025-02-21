#!/bin/bash

#停止服务
docker stop upload-rpc
docker stop user-rpc
docker stop etcd-server

#删除容器
docker rm upload-rpc
docker rm user-rpc
docker rm etcd-server

#删除镜像
docker rmi upload-rpc:0.0.1
docker rmi user-rpc:0.0.1

#删除none镜像
docker rmi $(docker images | grep "none" | awk '{print $3}')

#构建服务
docker build -t upload-rpc:0.0.1 -f upload/Dockerfile upload
docker build -t user-rpc:0.0.1 -f user/Dockerfile user


#启动服务
docker run -itd -p 2379:2379 -p 2380:2380 --name etcd-server --env ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd
docker run -itd --network host --name upload-rpc upload-rpc:0.0.1
docker run -itd --network host --name user-rpc user-rpc:0.0.1


