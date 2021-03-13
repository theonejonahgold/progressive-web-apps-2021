package serve

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

func Serve() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Starting serve on port", port)
	return fasthttp.ListenAndServe(":"+port, fastHTTPHandler)
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	wd, _ := os.Getwd()
	url := string(ctx.Path())
	fp := filepath.Join(wd, "dist", url, "index.html")
	buf, err := ioutil.ReadFile(fp)
	if err != nil {
		serveStaticFiles(ctx)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/html")
	ctx.SetBodyString(string(buf))
}

func serveStaticFiles(ctx *fasthttp.RequestCtx) {
	wd, _ := os.Getwd()
	url := string(ctx.Path())
	fp := filepath.Join(wd, "dist", "static", url)
	ctx.SendFile(fp)
}
