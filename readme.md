## 说明 

1. 基于 thrift-0.13.0
2. 基于 go 14
3. 需要安装tools工具， 进入tools/avalon 目录 执行go install；执行avalon -h 可以看到参数信息
    - i：输入thrift文件
    - o: 输出go文件
    - 此工具用于生成简单的server 文件，可以参考example/server or example/idgenerator
4. 生成工具后已经具有基础脚手架，需要配置一个base.yaml文件，可以参考base.yaml, 默认项目根目录，可以使用os.Setenv("base", "指定yaml地址")
5. go mod需要引入 github.com/941112341/avalon/sdk
6. 配置好hostPort, 检查zookeeper， enjoy yourself


## 项目规范
1. 保证工程idl/gen/client都放到common目录下，使用go mod 方式引入
2. 配置go-link


## 跨语言开发
> 底层是thrift，支持跨平台开发，但是中间层（服务发现）需要重新实现一遍