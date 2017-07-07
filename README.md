# toy

要求go1.8及以上版本

---

toy是一个简单的性能测试工具，其通过加载使用者提供的机器人插件，可针对不同系统执行性能测试。


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

bench命令
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

report命令
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

## 机器人插件
