package handler

import (
	"context"
)

func (s *Server) Background(bgCtx context.Context) {
	for {
		select {
		case <-bgCtx.Done():
			// 做一些清理工作...
			return
			// 其他的case...
		}
	}
}
