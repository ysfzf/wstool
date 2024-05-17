package ws

import (
	"context"
	"fmt"
	wsx "gitee.com/ysfzf/zerox/ws"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
)

type Wser struct {
	Port   int
	Status bool
	srv    *http.Server
	hub    *wsx.Hub
	ctx    context.Context
}

func NewWserver(ctx context.Context, port int) *Wser {

	return &Wser{
		Port: port,
		ctx:  ctx,
	}
}

func (s *Wser) Start(ctx context.Context, port int) (err error) {

	if !s.Status {

		s.hub = wsx.NewHub(s.ctx, wsx.SetOnMessage(OnMessage), wsx.SetOnConn(OnConn), wsx.SetOnClose(OnClose))
		go s.hub.Run()

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var t string
			for _, item := range s.hub.All() {
				t += fmt.Sprintf("<p>%s</p>", item)
			}
			_, _ = w.Write([]byte(t))
		})
		mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			wsx.ServeWs(s.hub, w, r)
		})
		s.srv = &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: mux,
		}

		go func() {
			s.Status = true
			runtime.EventsEmit(ctx, "server_start")
			err = s.srv.ListenAndServe()
			if err != nil {
				fmt.Println(err)
				s.Status = false
				runtime.EventsEmit(ctx, "server_stop")
				return
			}
		}()

	}
	return
}

func (s *Wser) Stop(ctx context.Context) (err error) {
	if s.srv != nil {
		err = s.srv.Shutdown(ctx)
		s.hub.CloseAll()
		if err == nil {
			s.Status = false
			runtime.EventsEmit(ctx, "server_stop")
		}
	}

	return
}

func (s *Wser) All() []string {
	return s.hub.All()
}

func (s *Wser) SendOne(target, message string) error {
	return s.hub.Send(target, []byte(message))
}

func (s *Wser) SendAll(message string) error {
	for _, item := range s.hub.All() {
		err := s.SendOne(item, message)
		if err != nil {
			return fmt.Errorf("send message to %s error:%s", item, err)
		}
	}
	return nil
}

func OnMessage(ctx context.Context, c *wsx.Client, msg []byte) {
	fmt.Println(c.Name, ":", string(msg))
	runtime.EventsEmit(ctx, "ws_message", c.Name, string(msg))
}

func OnConn(ctx context.Context, c *wsx.Client) {
	fmt.Println(c.Name + " connectioned")
	wsx.Send(c, []byte(c.Name))
	runtime.EventsEmit(ctx, "ws_conn", []string{c.Name})
}

func OnClose(ctx context.Context, c *wsx.Client) {
	runtime.EventsEmit(ctx, "ws_close", c.Name)
}
