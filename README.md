# GO-STREAM

> 此工具迁移lal工程,原程序地址 https://github.com/q191201771/lal

## DEMO使用描述

- pullhttpflv [支持 https-flv拉流]


```
params:
- i flv拉流地址
- f 保存文件路径
- s 是否存储flv
- n 拉流goroutine数量
 go run app/demo/pullhttpflv/pullhttpflv.go -i [flv拉流地址]
```

##

- pushrtmp

```
 go run  -i [flv本地资源] -o [推流地址] -r [是否循环推flv资源 (true|false)]
```