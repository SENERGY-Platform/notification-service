/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, request *http.Request) {
		response := &ResponseWriterWithStatusCodeLog{Parent: res, Status: 200}
		now := time.Now()
		defer l(request, response, now)
		next.ServeHTTP(response, request)
	})
}

func l(request *http.Request, response *ResponseWriterWithStatusCodeLog, t time.Time) {
	method := request.Method
	path := request.URL
	status := response.Status
	log.Printf("[%v] %v %v %v\n", method, path, status, time.Since(t))
}

type ResponseWriterWithStatusCodeLog struct {
	Parent http.ResponseWriter
	Status int
}

func (this *ResponseWriterWithStatusCodeLog) Header() http.Header {
	return this.Parent.Header()
}

func (this *ResponseWriterWithStatusCodeLog) Write(payload []byte) (int, error) {
	return this.Parent.Write(payload)
}

func (this *ResponseWriterWithStatusCodeLog) WriteHeader(statusCode int) {
	this.Status = statusCode
	this.Parent.WriteHeader(statusCode)
}

func (this *ResponseWriterWithStatusCodeLog) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := this.Parent.(http.Hijacker)
	if !ok {
		panic("websocket: response does not implement http.Hijacker")
	}
	return h.Hijack()
}
