# US3-PACK

### 工具简介

US3-PACK 一个将US3的多个文件打包成一个文件并上传到US3的工具。将 US3-PACK 部署在本地或者云主机中提供HTTP服务，接受打包请求，处理打包任务。

### 预先准备

* 在[云主机控制台](https://console.ucloud.cn/uhost/uhost)上创建一台linux操作系统的Uhost云主机。
* 在[US3控制台](https://console.ucloud.cn/ufile/token)上拿到具有对于目标桶上传下载，以及列取权限的令牌。

### 基本原理

工具的工作原理如下：

1. 工具本身会拉起一个HTTP服务器
2. 用户发送POST请求以提交打包任务，您会在请求的response中拿到最终将会生成的归档文件的url
3. 打包工具根据您提交的请求中的信息，将US3中的文件下载到本地，并进行打包
4. 将本地打包好的归档文件上传到US3中的指定位置
5. 您可以根据第二条中获取的归档文件url下载文件

### 操作步骤

我们这里假定您已经创建了[云主机Uhost](https://console.ucloud.cn/uhost/uhost)，并且在[US3控制台](https://console.ucloud.cn/ufile/token)上拿到了对应的令牌。

1. 下载打包工具 [工具包](https://github.com/ufilesdk-dev/ufile-pack/releases/download/v1.0.0/US3-PACK.zip)

2. 将工具包解压缩  `unzip US3-PACK.zip`

3. 修改工具中的[server_conf.json](https://github.com/ufilesdk-dev/ufile-pack/blob/main/server_conf.json)配置文件，配置文件如下：

   ````json
   {
     "log": {
       "LogDir": "logs",
       "LogPrefix": "zip_",
       "LogSuffix": ".log",
       "LogSize": 50,
       "LogLevel": "DEBUG"
     },
     "http": {
       "Ip": "0.0.0.0", 
       "Port": 80
     },    //服务监听的端口和ip
     "us3_config": {
       "public_key":"xxxxxxxxxxxxxx", //Token中的公钥
       "private_key":"xxxxxxxxxxxxxxxxx" //Token中的私钥
     }
   }
   ````

   

4. 执行`./US3-PACK`以启动服务，您也可以使用后台进程来执行该服务`nohup ./US3-PACK &`

5. 此时您可以发送POST请求到服务的根url (例如http://xxx.xxx.xxx.xxx)，请求参数有两种类型，分别对应指定某个前缀下的所有文件进行打包的任务，以及指定具体文件进行打包的任务。

   > 注意，如果您申请的Uhost云主机只有内网IP，那么请您在同一台云主机上，或者同一VPC内部发送打包的POST请求。

#### 指定前缀进行打包

   ```json
{
    "action": "GetUFileZipRequest",
    "prefix": "prefix",
    "bucket_name":"BucketName",
    "file_host":"internal-cn-sh2-01.ufileos.com"
}
   ```

   其中：

   * action字段指定request类型

   * bucket_name对应您所需打包文件所在存储桶的桶名

   * prefix为您所需打包文件所在的前缀（文件夹）路径

   * file_host为您访问桶所使用的endpoint，请您参考 [地域和域名](https://docs.ucloud.cn/ufile/introduction/region)

#### 指定文件列表进行打包

   ```json
{
    "action": "GetUFileZipByListRequest",
    "file_list": "prefix/key1,prefix/key2,prefix/key3",
    "bucket_name":"BucketName",
    "file_host":"internal-cn-sh2-01.ufileos.com"
}
   ```

   其中： 

   * file_list字段指定要打包的文件名，注意此处文件名包括文件的前缀，但不包括桶名
   * 其他字段同上

> 我们在文件包中提供了请求的json示例，您可以在压缩包中找到[event.json](https://github.com/ufilesdk-dev/ufile-pack/blob/main/event.json)文件，根据上文中的请求格式以及您想要提交的打包任务情况来更改json中的内容，并在同一台云主机上通过这一命令进行测试: ```curl -X POST -d@event.json localhost```

在发送完请求后，您会在返回中收到压缩包的地址。请求返回格式如下：

````javascript
{
    "Action": "GetUFileZipByListRequest",
    "prefix": "prefix",
    "RetCode": 0,
    "ErrMsg":"",
    "Key":"output/xxxx-xxxx-xxxxx-xxxx-xxxxxxx.zip"
}
````

其中key字段即为打包请求处理完毕后，工具上传到US3中的压缩包的对象名

#### 获取压缩包

您可以通过logs/文件夹下的日志文件来查看任务是否完成

最后，您可以使用http客户端工具下载这一文件，例如：

`wget http://bucket.internal-cn-sh2-01.ufileos.com/output/xxxx-xxxx-xxxxx-xxxx-xxxxxxx.zip`

### 性能测试

在使用1核1G内存的uhost主机，内网传输数据的情况下，打包55个20M文件(总大小1.15G)，的时间大概为30S。