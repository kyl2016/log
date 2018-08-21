# 日志监控系统 （Golang）

三部曲：1、认识问题；2、抽象问题；3、解决问题。

## 1.认识问题：同步读取、处理、写入
最小原型，梳理流程，容易看清问题的概貌。

[1.sync](https://github.com/kyl2016/log/tree/master/1.sync "")中的简单实现：
- ReadFromFile：只是将具体的字符串写入chan中。
- Parse：将字符串变为全大写。
- WriteToInfluxDB：打印出来。

流程：
```go
	go lp.ReadFromFile() 
	go lp.Parse()
	go lp.WriteToInfluxDB()
```

## 2.搭建框架、提出接口
此时，如果写具体代码，可能会遇到很多问题，代码量也可能比较大，那越写可能约乱，越复杂。
所以，此时要控制住复杂度，忽略细节，抽象问题，从整体考虑。
最好能拆分出模块，便于多人并行开发。

[2.with_interface](https://github.com/kyl2016/log/tree/master/2.with_interface "")中：
- ReadFromFile     -> Read 接口
- WriteToInfluxDB  -> Write 接口

流程：
```go
   go l.r.Read(l.rc)
   go l.Parse()
   go l.w.Write(l.wc)
```

## 3.实现文件读取、解析、写入
具体实现，接口细节可能要修改，但不影响整体架构。

[3.src](https://github.com/kyl2016/log/tree/master/3.src "")的流程：增加了并发数，流程没变
```go
    go lp.read.Read(lp.rc)

	for i:=0; i<2; i++{
		go lp.Parse()
	}

	for i:=0; i<4; i++{
		go lp.write.Write(lp.wc)
	}
```

## 扩展 三生万物

- 数据源：可以来自MQ、TCP……
- 解析器：一类日志样式对应一种解析器
- 输出： 写入到文件、远端服务器……

有点像工作流，不同的场景，配置不同的Processor序列。

学习慕课网麦克老师课程的代码整理，进行了结构调整。