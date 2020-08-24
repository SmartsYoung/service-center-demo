package main

import (
    "net/url"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

func main() {
    e := echo.New()
    e.Use(middleware.Recover())
    e.Use(middleware.Logger())

    // Setup proxy
    url1, err := url.Parse("http://localhost:8081")
    if err != nil {
        e.Logger.Fatal(err)
    }
    url2, err := url.Parse("http://localhost:8082")
    if err != nil {
        e.Logger.Fatal(err)
    }
    targets := []*middleware.ProxyTarget{
        {
            URL: url1,
        },
        {
            URL: url2,
        },
    }
    e.Group("/v4")
    e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

    // 如果改成 g := e.Group("/v4")  g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets))) 则变成了轮询http接口


    e.Logger.Fatal(e.Start(":1323"))
}
