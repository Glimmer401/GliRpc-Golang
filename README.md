> eRpc is determined to be a simple implementation of RPC protocol for starters.
> e means easy but not esphe.

# GliRPC

学 6.824 的时候感觉对 RPC 内容了解的不到位，所以想自己造个玩具看看。

## Implement

简单谈谈实现过程中比较有意思的地方。

### 序列化
实现了 Codec 适配器 (大概叫这个吧，不是很懂设计模式) 用于数据的序列化和反序列化，支持 gob 和 json 两种方式。

Codec 是一个接口，基于 Codec 可以特化成某种特定的序列化方式的 Codec。

gob 更为 go 原生一点，但是不太支持跨语言。

之后有时间可能会考虑自己实现序列化的过程，而不是直接用 gob 或者 json。

### 传输协议

**Option**

基于 TCP\IP 作为传输层，在 client 向 server 请求的服务过程中，会先发送一个 Option。

这个用于告知 server ① 这是 eRPC ② 后面的内容通过什么方式序列化。

在发送了 Option 后， client 会发送 Header + Body。而且在一次 TCP 连接中可能会发送多次 Header + Body。

Option 统一用 json 编码。

**Header**

Header 用于告知 server 本次调用哪个方法。

同时，由于单次 TCP 连接可能存在多次调用过程，所以也需要带一个序列号 seq。


**Body**

Body 承载经过序列化的内容，对于 client 发送给 server 的 request，是参数的序列化内容。而 server 发给 client 的 response 则是返回结果的序列化内容。

### 方法注册

用户可以将一个类和其所有导出方法注册成一个服务 service.

如果一个方法想要注册成远程调用过程，需要满足以下几个要求。

- the method’s type is exported.
- the method is exported.
- the method has two arguments, both exported (or builtin) types.
- the method’s second argument is a pointer.
- the method has return type error.

也就是说，对应的方法应该有如下格式 (其中 argType 可以是指针，而 replyType 必须是指针)

```golang
func (t *T) MethodName(argType T1, replyType *T2) error
```

注册时传入一个类实例，然后通过反射的方式获取所有方法，以及每个方法的参数类型。
> 这里不需要记录返回类型，因为我们规定了返回类型需要时一个 type error。

> 而用户需要的返回类型是在参数中以指针的形式被修改返回

当用户调用过程时，以 serviceName.methodName 的格式用于指定特定方法。比如说：`Calc.Add`

## Milestone

- 2022.4.17 eRPC started
- 2022.4.19 eRPC has implemented transport protocol
- 2022.4.21 eRPC finished simple client and server
- 


