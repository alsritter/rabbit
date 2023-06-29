# Rabbit
这是一个用于快速搭建 Kratos 服务的模板

```
(\(\ 
( -.-) 
o_(")(")
```

## 如何使用？
1. 全局替换 `alsritter.icu/rabbit-template` 为你的项目名
2. 执行 `make all` 编译项目项目

如何 debug？ 创建 `.vscode/launch.json` 文件，内容如下：

```json
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Package",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server/.",
      "cwd": "${workspaceRoot}",
      "args": [
        "-conf",
        "${workspaceRoot}/configs/local_config.yaml"
      ]
    }
  ]
}
```

## 配置环境

### win系统下安装工具支持gun命令

[MinGW](http://www.mingw.org/wiki/getting_started)