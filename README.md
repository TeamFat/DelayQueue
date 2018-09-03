# DelayQueue

基于Redis实现的延迟队列

## 应用场景
* 订单超过30分钟未支付,自动关闭
* 订单完成后, 如果用户一直未评价, 5天后自动好评
* 会员到期前15天, 到期前3天分别发送短信提醒

## 实现原理
> 利用Redis的有序集合，member为JobID,score为任务执行的时间戳,每秒扫描一次集合，取出执行时间小于等于当前时间的任务. 

## 源码安装
* `go`语言版本1.10+
* `go get -d github.com/TeamFat/DelayQueue`
* `go build`

## 运行
`./DelayQueue`  
> HTTP Server监听`0.0.0.0:8080`, Redis连接地址`127.0.0.1:6379`, 数据库编号`2`, 数据库密码`123456`
具体配置见conf/config.yaml,按需修改

## HTTP接口
* 请求方法 `POST`   
* 请求Body及返回值均为`json`

### 返回值
```json
{
  "code": 0,
  "message": "OK",
  "data": null
}
```

|  参数名 |     类型    |     含义     |        备注       |
|:-------:|:-----------:|:------------:|:-----------------:|
|   code  |     int     |    状态码    | 0: 成功 非0: 失败 |
| message |    string   | 状态描述信息 |                   |
|   data  | object, null |   附加信息   |                   |

### 添加任务   
URL地址 `/queue/push`   
```json
{
  "topic": "order",
  "delay": 30,
  "body": "{\"uid\": 10829378,\"created\": 1498657365 }"
}
```
|  参数名 |     类型    |     含义     |        备注       |
|:-------:|:-----------:|:------------:|:-----------------:|
|   topic  | string     |    Job类型                   |                     |
|   delay  | int        |    Job需要延迟的时间, 单位：秒    | 必须大于0          |
|   body   | string     |    Job的内容，供消费者做具体的业务处理，如果是json格式需转义 |                   |

### 轮询队列获取任务
服务端会Hold住连接, 直到队列中有任务或超时后返回默认超时10秒
 
URL地址 `/queue/pop`    
```json
{
  "topic": "order"
}
```
|  参数名 |     类型    |     含义     |        备注       |
|:-------:|:-----------:|:------------:|:-----------------:|
|   topic  | string     |    Job类型                   |多个topic 逗号分隔,每次仅pop一条|
|   timeout| int        |   队列超时时间, 单位：秒 |     可选项                       |


队列中有任务返回值
```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "topic": "order",
    "id":"DelayJob_932043ab-96d3-4419-bba1-67c88b059f1b",
    "delay": 1535977210,
    "body": "{\"uid\": 10829378,\"created\": 1498657365 }"
  }
}
```
队列为空返回值   
```json
{
  "code": 0,
  "message": "OK",
  "data": null
}
```

>参考资源
* [有赞延迟队列设计](http://tech.youzan.com/queuing_delay)
* [基于Redis实现的延迟队列](https://segmentfault.com/a/1190000010021748)
* [基于 Redis 的延时队列设计与实现](http://rust.love/2018/02/04/study_weekly_01.html)