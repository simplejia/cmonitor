[中文 README](#中文)


# [cmonitor](http://github.com/simplejia/cmonitor)
## Functions
* Used to monitor and manage processes

## Installation
> go get -u github.com/simplejia/cmonitor

## Implementation
* Check each process by sending signal0 per 300ms, each monitored process is monitored in an independent goroutine.
* When cmonitor is started, it will listen on a tcp port which is used to receive management commands（start|stop|restart|status|...）

## Playing with cmonitor
* Configuration files: [conf.json](http://github.com/simplejia/cmonitor/tree/master/conf/conf.json) (json format, supporting annotations)，You can redefine the parameters in the configuration files via passing env and conf parameters, such as ：./cmonitor -env dev -conf='port=8080;clog.mode=1', multiple parameters are separated by `;`
```
{
    "dev": {
        "port": 29118, // listening port
        "rootpath": "/home/simplejia/tools/go/ws/src/github.com/simplejia",
        "environ": "ulimit -n 65536", // set env
        "svrs": {
            // demo
            "demo": "wsp/demo/demo" // key: name, value: relative path (with parameters)
        },
        "log": {
            "mode": 3, // 0: none, 1: localfile, 2: collector (bitmap)
            "level": 15 // 0: none, 1: debug, 2: warn 4: error 8: info (bitmap)
        }
    },
    "test": {
        "port": 29118, // 
        "rootpath": "/home/simplejia/tools/go/ws/src/github.com/simplejia",
        "environ": "ulimit -n 65536", // 
        "svrs": {
            // demo 
            "demo": "wsp/demo/demo"
        },
        "log": {
            "mode": 3, // 0: none, 1: localfile, 2: collector (bitmap)
            "level": 15 // 0: none, 1: debug, 2: warn 4: error 8: info (bitmap)
        }
    },
    "prod": {
        "port": 29118, // 
        "rootpath": "/home/simplejia/tools/go/ws/src/github.com/simplejia",
        "environ": "ulimit -n 65536", // 
        "svrs": {
            // demo 
            "demo": "wsp/demo/demo"
        },
        "log": {
            "mode": 2, // 0: none, 1: localfile, 2: collector (bitmap)
            "level": 14 // 0: none, 1: debug, 2: warn 4: error 8: info (bitmap)
        }
    }
}
```
* Run：cmonitor.sh [start|stop|restart|status|check]
* Management：cmonitor -[h|env|status|start|stop|restart] [dev|test|prod|all|["svrname"]]

## Notice
* The logs of cmonitor are reported via clog, they also can be recorded in local cmonitor.log. Notes: The logs in cmonitor.log will not be split, so try to keep the log less. And it is recommended to report the log by [clog](github.com/simplejia/clog).
* When cmonitor starts its monitor processes, the console logs (cmonitor.log) of the monitored processes will be output to corresponding process directory which will be saved up to 30 days.，the history logs is like cmonitor.{day}.log
* When cmonitor is starting, all monitored processes will be started according to its conf.json configuration. When the monitored processes have already been started and meet the configuration requirements, cmonitor will automatically add them to the monitor list.
* cmonitor will periodically check the process status. If there is abnormal process exit, it will repeatedly try to restart it and record the error log.
* When the monitored processes are multiple process model, cmonitor will only monitor and manage parent process, the child process should be able to check the exit status of parent process and and exit with its parent process.
* The monitored process is started in the ways of nohup, so you should not set your process daemon mode.
* cmonitor will check the process status per minute via ps command. If there is any abnormal status, for example, there are multiple started processes, it will record the error log.
* cmonitor will start its internal httpserver (bind inner ip), so it also supports remote management，such as input the url in your safari：http://xxx.xxx.xxx.xxx:29118/?command=status&service=all

## demo
```
$ cmonitor -env dev -status all

*****STATUS OK SERVICE LIST*****
demo PID:13539

*****STATUS FAIL SERVICE LIST*****

$ cmonitor -env dev -restart demo
SUCCESS
```

---
中文
===

# [cmonitor](http://github.com/simplejia/cmonitor)
## 功能
* 用于进程监控，管理

## 安装
> go get -u github.com/simplejia/cmonitor

## 实现
* 被监控进程启动后，按每300ms执行一次状态检测（通过发signal0信号检测），每个被监控进程在一个独立的协程里被监测。
* cmonitor启动后会监听一个http端口用于接收管理命令（start|stop|restart|status|...）

## 使用方法
* 配置文件：[conf.json](http://github.com/simplejia/cmonitor/tree/master/conf/conf.json) (json格式，支持注释)，可以通过传入自定义的env及conf参数来重定义配置文件里的参数，如：./cmonitor -env dev -conf='port=8080;clog.mode=1'，多个参数用`;`分隔
```
{
    "dev": {
        "port": 29118, // 配置监听端口
        "rootpath": "/home/simplejia/tools/go/ws/src/github.com/simplejia",
        "environ": "ulimit -n 65536", // 配置环境变量
        "svrs": {
            // demo
            "demo": "wsp/demo/demo" // key: 名字 value: 将与rootpath拼接在一起运行
        },
        "log": {
            "mode": 3, // 0: none, 1: localfile, 2: collector (数字代表bit位)
            "level": 15 // 0: none, 1: debug, 2: warn 4: error 8: info (数字代表bit位)
        }
    },
    "test": {
        "port": 29118, // 配置监听端口
        "rootpath": "/home/simplejia/tools/go/ws/src/github.com/simplejia",
        "environ": "ulimit -n 65536", // 配置环境变量
        "svrs": {
            // demo 
            "demo": "wsp/demo/demo"
        },
        "log": {
            "mode": 3, // 0: none, 1: localfile, 2: collector (数字代表bit位)
            "level": 15 // 0: none, 1: debug, 2: warn 4: error 8: info (数字代表bit位)
        }
    },
    "prod": {
        "port": 29118, // 配置监听端口
        "rootpath": "/home/simplejia/tools/go/ws/src/github.com/simplejia",
        "environ": "ulimit -n 65536", // 配置环境变量
        "svrs": {
            // demo 
            "demo": "wsp/demo/demo"
        },
        "log": {
            "mode": 2, // 0: none, 1: localfile, 2: collector (数字代表bit位)
            "level": 14 // 0: none, 1: debug, 2: warn 4: error 8: info (数字代表bit位)
        }
    }
}
```
* 运行方法：cmonitor.sh [start|stop|restart|status|check]
* 进程管理：cmonitor -[h|env|status|start|stop|restart] [dev|test|prod|all|["svrname"]]

## 注意
* cmonitor的运行日志通过clog上报，也可记录在本地cmonitor.log日志文件里，注意：此cmonitor.log日志文件不会被切分，所以尽量保持较少的日志输出，建议通过clog方式上报日志
* cmonitor启动监控进程后，被监控进程控制台日志cmonitor.log会输出到相应进程目录，最多保存30天，历史日志以cmonitor.{day}.log方式备份
* 当cmonitor启动时，会根据conf.json配置启动所有被监控进程，当被监控进程已经启动过，并且符合配置要求时，cmonitor会自动将其加入监控列表
* cmonitor会定期检查进程运行状态，如果进程异常退出，cmonitor会反复重试拉起，并且记录日志
* 当被监控进程为多进程运行模式，cmonitor只监控管理父进程(子进程应实现检测父进程运行状态，并随父进程退出而退出）
* 被监控进程以nohup方式启动，所以你的程序就不要自己设定daemon运行了
* 每分钟通过ps方式检测一次进程状态，如果出现任何异常，比如有多份进程启动等，记日志
* 由于cmonitor会同时启动内部httpserver（绑内网ip），所以也支持远程管理，比如在浏览器里输入：http://xxx.xxx.xxx.xxx:29118/?command=status&service=all

## demo
```
$ cmonitor -env dev -status all

*****STATUS OK SERVICE LIST*****
demo PID:13539

*****STATUS FAIL SERVICE LIST*****

$ cmonitor -env dev -restart demo
SUCCESS
```
