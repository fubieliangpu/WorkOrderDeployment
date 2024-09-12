package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PageRequest struct {
	// 分页的大小
	PageSize int `json:"page_size"`
	// 当前页码
	PageNumber int `json:"page_number"`
}

func NewPageRequest() *PageRequest {
	return &PageRequest{
		PageSize:   10,
		PageNumber: 1,
	}
}

func (req *PageRequest) Offset() int {
	return (req.PageNumber - 1) * req.PageSize
}

func NewPageRequestFromGinCtx(c *gin.Context) *PageRequest {
	p := NewPageRequest()
	pnStr := c.Query("page_number")
	psStr := c.Query("page_size")
	if pnStr != "" {
		pn, _ := strconv.ParseInt(pnStr, 10, 64)
		if pn != 0 {
			p.PageNumber = int(pn)
		}
	}

	if psStr != "" {
		ps, _ := strconv.ParseInt(psStr, 10, 64)
		if ps != 0 {
			p.PageSize = int(ps)
		}
	}

	return p
}
