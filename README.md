

> eRpc is determined to be a simple implementation of RPC protocol for starters.
> e means easy but not esphe.

# GliRPC

<<<<<<< HEAD
# Implement

## 序列化
实现了 Codec 用于数据的序列化和反序列化，计划支持 gob 和 json 两种方式。
gob 更为 go 原生一点，但是不支持跨语言的反序列化。
使用 json 可以解决跨语言的问题，但是在初期实现反射过程中存在一些问题，所以暂时搁置作为一项 TODO item。

## 传输协议

使用 tcp

## 方法注册

如果一个方法想要注册成远程调用过程，需要满足以下几个要求。

- the method’s type is exported.
- the method is exported.
- the method has two arguments, both exported (or builtin) types.
- the method’s second argument is a pointer.
- the method has return type error.

也就是说，对应的方法应该有如下格式

```golang
func (t *T) MethodName(argType T1, replyType *T2) error
```
=======

>>>>>>> Update README.md

# Milestone

- 2022.4.17 eRPC started
- 2022.4.19 eRPC transfer into GliRPC with both implement with go, java
- 2022.4.19 eRPC has implement transport protocol
- 


