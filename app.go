package main

import (
	"changeme/ws"
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

var w *ws.Wser

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	if w != nil {
		_ = w.Stop(a.ctx)
	}

}

func (a *App) ShowErrDialog(msg string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "错误",
		Message: msg,
		Buttons: []string{"OK"},
	})
}

func (a *App) Start(port int) error {
	w = ws.NewWserver(a.ctx, port)
	err := w.Start(a.ctx, port)
	if err != nil {
		a.ShowErrDialog(fmt.Sprintf("启动服务端出错了: %s", err))
	}
	return err
}

func (a *App) Stop() error {
	if w != nil {
		err := w.Stop(a.ctx)
		if err != nil {
			a.ShowErrDialog(fmt.Sprintf("停止服务端出错了: %s", err))
		}
		w = nil
		return err
	}
	return nil
}

func (a *App) Status() bool {
	if w != nil {
		return w.Status
	}
	return false
}

func (a *App) All() []string {
	if w != nil {
		return w.All()
	}
	return []string{}
}

func (a *App) SendMsg(targets []string, message string) (err error) {
	if w == nil {
		return errors.New("server not run")
	}
	if len(targets) == 0 {
		err = w.SendAll(message)
		if err != nil {
			a.ShowErrDialog(fmt.Sprintf("发送信息出错了: %s", err))
		}
		return
	}
	for _, target := range targets {
		err = w.SendOne(target, message)
		if err != nil {
			a.ShowErrDialog(fmt.Sprintf("发送信息出错了: %s", err))
		}
	}

	return
}
