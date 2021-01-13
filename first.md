### Hello World

使用Gin实现`Hello world`非常简单，创建一个router，然后使用其Run的方法：

```Go
import ( 
    "gopkg.in/gin-gonic/gin.v1" 
    "net/http" 
) 
func main(){ 
    router := gin.Default() 
    router.GET("/", func(c *gin.Context) { 
    c.String(http.StatusOK, "Hello World") 
    }) 
	router.Run(":8000") 
}
```

简单几行代码，就能实现一个web服务。使用gin的Default方法创建一个路由handler。然后通过HTTP方法绑定路由规则和路由函数。不同于net/http库的路由函数，gin进行了封装，把request和response都封装到gin.Context的上下文环境。最后是启动路由的Run方法监听端口。麻雀虽小，五脏俱全。当然，除了GET方法，gin也支持POST,PUT,DELETE,OPTION等常用的restful方法。

### restful路由

gin的路由来自httprouter库。因此httprouter具有的功能，gin也具有，不过gin不支持路由正则表达式：

```Go
func main(){ 
    router := gin.Default() 
    router.GET("/user/:name", func(c *gin.Context) { 
        name := c.Param("name") 
        c.String(http.StatusOK, "Hello %s", name) 
    }) 
}
```

冒号`:`加上一个参数名组成路由参数。可以使用c.Params的方法读取其值。当然这个值是字串string。诸如`/user/rsj217`，和`/user/hello`都可以匹配，而`/user/`和`/user/rsj217/`不会被匹配。

```Go
curl http://127.0.0.1:8000/user/rsj217 

Hello rsj217%                  

curl http://127.0.0.1:8000/user/rsj217/ 

404 page not found%            

 curl http://127.0.0.1:8000/user/ 
404 page not found%
```

除了`:`，gin还提供了`*`号处理参数，`*`号能匹配的规则就更多。

```Go
func main(){
    router := gin.Default() 
    router.GET("/user/:name/*action", func(c *gin.Context) { 
        name := c.Param("name") 
        action := c.Param("action") 
        message := name + " is " + action 
        c.String(http.StatusOK, message) 
    }) 
}
```

访问效果如下

```
curl http://127.0.0.1:8000/user/rsj217/
rsj217 is /%                       

curl http://127.0.0.1:8000/user/rsj217/中国
rsj217 is /中国%
```

### query string参数与body参数

web提供的服务通常是client和server的交互。其中客户端向服务器发送请求，除了路由参数，其他的参数无非两种，查询字符串query string和报文体body参数。所谓query string，即路由用，用`?`以后连接的`key1=value2&key2=value2`的形式的参数。当然这个key-value是经过urlencode编码。

**query string**

对于参数的处理，经常会出现参数不存在的情况，对于是否提供默认值，gin也考虑了，并且给出了一个优雅的方案：

```go
func main(){ 
    router := gin.Default() 
 
    router.GET("/welcome", func(c *gin.Context) { 
        firstname := c.DefaultQuery("firstname", "Guest") 
        lastname := c.Query("lastname") 
        c.String(http.StatusOK, "Hello %s %s", firstname, lastname) 
        }) 
 
    router.Run() 
}

```

使用c.DefaultQuery方法读取参数，其中当参数不存在的时候，提供一个默认值。使用Query方法读取正常参数，当参数不存在的时候，返回空字串：

```go
☁ ~ curl http://127.0.0.1:8000/welcome 
Hello Guest % 
 
☁ ~ curl http://127.0.0.1:8000/welcome\?firstname=中国 
Hello 中国 % 
 
☁ ~ curl http://127.0.0.1:8000/welcome\?firstname=中国&lastname=天朝 
Hello 中国 天朝% 
 
☁ ~ curl http://127.0.0.1:8000/welcome\?firstname\=&lastname=天朝 
Hello 天朝% 
 
☁ ~ curl http://127.0.0.1:8000/welcome\?firstname=%E4%B8%AD%E5%9B%BD 
Hello 中国 %

```

之所以使用中文，是为了说明urlencode。注意，当firstname为空字串的时候，并不会使用默认的Guest值，空值也是值，DefaultQuery只作用于key不存在的时候，提供默认值。

**body**

http的报文体传输数据就比query string稍微复杂一点，常见的格式就有四种。例如`application/json`，`application/x-www-form-urlencoded`, `application/xml`和`multipart/form-data`。后面一个主要用于图片上传。json格式的很好理解，urlencode其实也不难，无非就是把query string的内容，放到了body体里，同样也需要urlencode。默认情况下，c.PostFROM解析的是`x-www-form-urlencoded`或`from-data`的参数。

```go
func main(){ 
    router := gin.Default() 
    router.POST("/form_post", func(c *gin.Context) { 
        message := c.PostForm("message") 
        nick := c.DefaultPostForm("nick", "anonymous") 
        c.JSON(http.StatusOK, gin.H{ 
            "status": gin.H{ 
                "status_code": http.StatusOK, 
                "status": "ok", 
            }, 
            "message": message, 
            "nick": nick,         
        }) 
    }) 
}
```

与get处理query参数一样，post方法也提供了处理默认参数的情况。同理，如果参数不存在，将会得到空字串。

```go
☁ ~ curl -X POST http://127.0.0.1:8000/form_post -H "Content-Type:application/x-www-form-urlencoded" -d "message=hello&nick=rsj217" | python -m json.tool % Total % Received % Xferd 
 
Average Speed Time Time Time Current Dload Upload Total Spent Left Speed 
100 104 100 79 100 25 48555 15365 --:--:-- --:--:-- --:--:-- 79000 
{ 
    "message": "hello", 
    "nick": "rsj217", 
    "status": { 
        "status": "ok", 
        "status_code": 200 
    } 
}
```

前面我们使用c.String返回响应，顾名思义则返回string类型。content-type是plain或者text。调用c.JSON则返回json数据。其中gin.H封装了生成json的方式，是一个强大的工具。使用golang可以像动态语言一样写字面量的json，对于嵌套json的实现，嵌套gin.H即可。

发送数据给服务端，并不是post方法才行，put方法一样也可以。同时querystring和body也不是分开的，两个同时发送也可以：

```go
func main(){ 
    router := gin.Default() 
    router.PUT("/post", func(c *gin.Context) { 
        id := c.Query("id") 
        page := c.DefaultQuery("page", "0") 
        name := c.PostForm("name") 
        message := c.PostForm("message") 
        fmt.Printf("id: %s; page: %s; name: %s; message: %s \n", id, page, name, message)                   
        c.JSON(http.StatusOK, gin.H{ "status_code": http.StatusOK, }) 
    }) 
}

```

上面的例子，展示了同时使用查询字串和body参数发送数据给服务器。