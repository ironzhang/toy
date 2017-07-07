# toy

要求go1.8及以上版本

---

toy是一个简单的性能测试工具，通过加载不同的机器人插件，可针对不同系统执行性能测试。


## Install

```
go get -u github.com/ironzhang/toy
```

## Usage

```
Usage: toy COMMAND [arg...]

A extensible benchmark tool

Commands:
    bench       do benchmark test
    report      make benchmark report

run 'toy COMMAND --help' for more information on a command
```

bench命令执行性能测试，并输出测试记录
```
Usage: toy bench [OPTIONS]

start robots to do benchmark test

  -ask
        ask execute scheduler
  -output string
        the record file
  -robot-num int
        robot num (default 1)
  -robot-path string
        robot path (default "./robots/test-robot")
  -verbose int
        verbose level
```

report命令根据测试记录生成测试报告
```
Usage: toy report [OPTIONS] FILE [FILE...]

make benchmark report with test records

  -format string
        report format, html/text (default "html")
  -output-dir string
        output dir (default "output")
  -sample-size int
        sample size (default 500)
```

## Quick start

```
go get -u github.com/ironzhang/toy
cd $GOPATH/src/github.com/ironzhang/toy/
go build
cd ./robots/test-robot/; ./build.sh; cd ../..

./toy bench -robot-path ./robots/test-robot -robot-num 100 -output test.tbr
./toy report -format text test.tbr
```

## 机器人插件

toy要求加载的机器人插件实现robot.Robot接口

```go
// Robot 机器人接口
type Robot interface {
        OK() bool
        Do(name string) error
}
```

并在插件中实现`func NewRobots(n int, file string) ([]robot.Robot, error)`函数。

如下是test-robot的实现

```go
package main

import (
        "time"

        "github.com/ironzhang/toy/framework/robot"
)

func NewRobots(n int, file string) ([]robot.Robot, error) {
        robots := make([]robot.Robot, 0, n)
        for i := 0; i < n; i++ {
                robots = append(robots, &Robot{})
        }
        return robots, nil
}

type Robot struct {
}

func (r *Robot) OK() bool {
        return true
}

func (r *Robot) Do(name string) error {
        //fmt.Println(name)
        switch name {
        case "Connect":
                time.Sleep(10 * time.Millisecond)
        case "Prepare":
                time.Sleep(20 * time.Millisecond)
        case "Publish":
                time.Sleep(100 * time.Microsecond)
        case "Disconnect":
                time.Sleep(10 * time.Microsecond)
        }
        return nil
}
```

使用如下命令编译插件

```
go build -buildmode=plugin -o robot.so
```

### 调度器配置

在每一个插件目录下必须有一个`schedulers.json`文件来描述机器人的调度配置。如下是test-robot的调度配置。

```json
[
        {
                "Name": "Connect",
                "N": 1,
                "C": 10,
                "QPS": 1000
        },
        {
                "Name": "Prepare",
                "N": 1,
                "C": 10,
                "QPS": 1000
        },
        {
                "Name": "Publish",
                "N": 100,
                "C": 10,
                "QPS": 5000
        },
        {
                "Name": "Disconnect",
                "N": 1,
                "C": 10,
                "QPS": 1000
        }
]
```

`Name`表示执行的动作，`N`表示每个机器人的执行次数，`C`表示并发的协程数，`QPS`表示对吞吐量的控制参数。toy将按照调度配置的描述调用机器人执行性能测试。

### 已实现的机器人插件

* [test-robot](https://github.com/ironzhang/toy/tree/master/robots/test-robot) 用于验证toy工具的机器人插件，没有实际意义
* [mqtt-robot](https://github.com/ironzhang/toy/tree/master/robots/mqtt-robot) 可用于测试mqtt服务器性能的机器人插件
