package main

import (
	"fmt"
	"github.com/donseba/expronaut"
	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/sse"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type App struct {
	HTMX *htmx.HTMX
}

var (
	sseManager sse.Manager
)

func main() {
	app := App{
		HTMX: htmx.New(),
	}

	sseManager = sse.NewManager(5)

	go func() {
		for {
			time.Sleep(1 * time.Second) // Send a message every second
			sseManager.Send(sse.NewMessage(fmt.Sprintf("%v", time.Now().Format(time.TimeOnly))).WithEvent("time"))
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(app.Home))
	mux.Handle("POST /calc", http.HandlerFunc(app.Calc))
	mux.Handle("GET /sse", http.HandlerFunc(app.SSE))

	err := http.ListenAndServe(":4321", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func (a *App) Calc(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	h := a.HTMX.NewHandler(w, r)

	in := r.PostFormValue("calc")
	if in == "" {
		h.TriggerError(fmt.Sprintf("error: %v", "missing input"))
		_, _ = h.Write([]byte{})
		return
	}

	ti := time.Now()
	// do some calculation
	out, err := expronaut.Evaluate(ctx, in)
	if err != nil {
		h.TriggerError(fmt.Sprintf("error: %v", err))
		_, _ = h.Write([]byte(fmt.Sprint(out)))
		return
	}

	calcTime := time.Since(ti).Microseconds()

	h.TriggerInfo(fmt.Sprintf("calculation took %d us", calcTime))

	_, _ = h.Write([]byte(fmt.Sprint(out)))
}

func (a *App) SSE(w http.ResponseWriter, r *http.Request) {
	cl := sse.NewClient(randStringRunes(10))

	sseManager.Handle(w, r, cl)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
