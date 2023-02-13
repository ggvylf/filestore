## 代码来源
https://github.com/samtake/filestore-server/

## 关于应用启动

- 在加入rabbitMQ实现文件异步转移之前，启动方式：

    - 启动上传应用程序:
```bash
# cd $GOPATH/<你的工程目录>
> cd $GOPATH/filestore-server
> go run main.go
```

- 在加入rabbitMQ实现文件异步转移阶段，启动方式(分裂成了两个独立程序)：

    - 启动上传应用程序:
```bash
# cd $GOPATH/<你的工程目录>
> cd $GOPATH/filestore-server
> go run service/upload/main.go
```

    - 启动转移应用程序:
```bash
# cd $GOPATH/<你的工程目录>
> cd $GOPATH/filestore-server
> go run service/transfer/main.go
```

-  微服务架构下启动方式(非容器化部署):

    - 一键启动微服务(start-all.sh):
```bash
> cd $GOPATH/filestore-server
> ./service/start-all.sh 
 编译完成:  service/bin/dbproxy
 编译完成:  service/bin/upload
 编译完成:  service/bin/download
 编译完成:  service/bin/transfer
 编译完成:  service/bin/account
 编译完成:  service/bin/apigw
 已启动  dbproxy
 已启动  upload
 已启动  download
 已启动  transfer
 已启动  account
 已启动  apigw
微服务启动完毕.
```

    - 一键关闭微服务(stop-all.sh):
```bash
> cd $GOPATH/filestore-server
> ./service/stop-all.sh 
 已关闭:  apigw
 已关闭:  account
 已关闭:  transfer
 已关闭:  download
 已关闭:  upload
 已关闭:  dbproxy
执行完毕.
```

-  微服务架构下启动方式(容器化部署):
```bash
> cd $GOPATH/filestore-server
# 脚本方式启动容器
> ./deploy/start-all.sh
# 脚本方式关闭容器
> ./deploy/stop-all.sh
# docker-compose方式启动容器
> cd ./deploy/service_dc
> sudo docker-compose up -d
# k8s方式启动微服务
> cd ./deploy/service_k8s
> kubectl apply -f svc_account.yaml
> kubectl apply -f svc_apigw.yaml
> kubectl apply -f svc_dbproxy.yaml
> kubectl apply -f svc_download.yaml
> kubectl apply -f svc_transfer.yaml
> kubectl apply -f svc_upload.yaml
> cd ./deploy/traefik_k8s
> kubectl apply -f service-ingress.yaml
```
