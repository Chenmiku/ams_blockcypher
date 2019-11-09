package web

import (
	"net/http"
)

type IWebSuccess interface {
	StatusCode() int
}

type Created string

func (e Created) Success() string {
	return string(e)
}

func (e Created) StatusCode() int {
	return http.StatusCreated
}

type OK string

func (e OK) Success() string {
	return string(e)
}

func (e OK) StatusCode() int {
	return http.StatusOK
}