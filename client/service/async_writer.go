// Author: xufei
// Date: 2019-09-10 14:54

package service

import (
	"fmt"
	"gim/internal/lg"
	"gim/internal/util/fileutils"
	"gim/model"
	"os"
	"time"

	"go.uber.org/zap"
)

type MsgLogger interface {
	Log(msg string)
	Close()
}

type asyncWriter struct {
	config   *model.ClientConfig
	writerCh chan string
	exitChan chan int
}

func (w *asyncWriter) Log(msg string) {
	w.writerCh <- msg
}

func (w *asyncWriter) Close() {
	close(w.exitChan)
}

func (w *asyncWriter) logPump() {
	for {
		select {
		case msg := <-w.writerCh:
			w.writeMsg(msg)
		case <-w.exitChan:
			lg.Logger().Info("关闭消息文件的写入!")
			return
		}
	}
}

func NewWriter(cfg *model.ClientConfig) *asyncWriter {
	writer := &asyncWriter{
		config:   cfg,
		writerCh: make(chan string, 16), // 缓存通道，最多处理 16 条消息写入
		exitChan: make(chan int),
	}

	go writer.logPump()

	return writer
}

func (w *asyncWriter) writeMsg(msg string) {
	now := time.Now()
	dir := w.config.MsgLogPath + w.config.UserName + "/"
	fileName := fmt.Sprintf("%d%d%d.log", now.Year(), now.Month(), now.Day())

	_, ok := fileutils.FindDir(dir)
	if !ok {
		if err := fileutils.CreateDir(dir); err != nil {
			lg.Logger().Error("消息日志文件夹创建失败", zap.Error(err))
			return
		}
	}

	file, err := os.OpenFile(dir+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		lg.Logger().Error("消息日志文件打开失败", zap.Error(err))
		return
	}
	defer file.Close()

	fmt.Fprintln(file, msg)
}
