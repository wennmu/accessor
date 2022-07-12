# golang类型的getter/setter代码生成器

## 安装

```shell
  go get github.com/wennmu/accessor
```

## 使用

```go
  //go:generate accessor -types=User
```

```shell
  go generate ./...
```

## 注意

- -types 参数目前只支持了指定一个类型名称，现阶段无法指定多个
- 生成的代码目前输出到os.Stdout了，对于直接生成在文件里，还没有较好的方案。可以把生成结果复制到你的代码中
- 如有相关问题，请通过issue提交，感谢
