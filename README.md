# US3-PACK

## 工具简介
US3-PACK 一个将US3的多个文件打包成一个文件并上传到US3的工具。将 US3-PACK 部署在本地或者云主机中提供HTTP服务，接受打包请求，处理打包任务。

## 详细流程
* 用户US3-PACK，指定存储空间及待压缩的文件
* US3-PACK从US3中获取指定文件，并生成一个随机名称的ZIP压缩包
* US3-PACK将压缩包上传至US3
* US3-PACK将ZIP包的下载地址返回给用户
* 用户使用返回的下载地址从US3中下载文件

## 部署方式
#### 直接下载部署
* 下载文件 ```wget -O US3-PACK http://us3-release.cn-bj.ufileos.com/US3-PACK/US3-PACK```
* 参考注释编辑 ```server_conf.json```
* 启动US3-PACK


#### 下载代码编译部署
* 下载代码 ```git clone http://github.com/us3/us3-pack```
* 进入us3-pack编译 ```cd us3-pack & make```
* 参考注释编辑 ```server_conf.json```
* 启动US3-PACK

## 使用方式
* 编辑event.json指定打包参数
* 发送打包请求 ```curl -d @./event.json 123.123.123.123```