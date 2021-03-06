// Copyright 2017 Gitai<i@gitai.me> All rights reserved.
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify,
// merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall
// be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR
// ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var httpTpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Proc</title>
    <link href="https://cdn.bootcss.com/normalize/8.0.1/normalize.min.css" rel="stylesheet">
    <link href="https://cdn.bootcss.com/codemirror/5.42.2/codemirror.min.css" rel="stylesheet">
    <link href="https://cdn.bootcss.com/codemirror/5.42.2/theme/darcula.min.css" rel="stylesheet">
    <script src="https://cdn.bootcss.com/zepto/1.0rc1/zepto.min.js"></script>
    <style type="text/css">
        body {
            background-color: #2a2734;
            color: #dddddd;
        }
        .nav {
            width: 250px;
            padding: 2em 1em;
            box-sizing: border-box;
            text-align: center;
        }
        .nav > .logo {
            width: 250px;
            margin: -1em;
        }
        .nav code {
            background-color: black;
            border-radius: 5px;
            padding: 5px;
        }
        .editor {
            top: 0;
            position: absolute;
            float: left;
            padding-left: 250px;
            width: 100%;
            box-sizing: border-box;
            height: 100%;
        }
        .CodeMirror {
            height: 100%;
        }
        #forkme_banner {
            display: block;
            text-decoration: none;
            position: absolute;
            top: 0;
            right: 10px;
            z-index: 10;
            padding: 10px 50px 10px 10px;
            color: #fff;
            background: url(https://pages-themes.github.io/slate/assets/images/blacktocat.png) #0090ff no-repeat 95% 50%;
            font-weight: 700;
            box-shadow: 0 0 10px rgba(0,0,0,0.5);
            border-bottom-left-radius: 2px;
            border-bottom-right-radius: 2px;
        }
    </style>
</head>
<body>
<a id="forkme_banner" href="https://github.com/GitaiQAQ/Env.beta">View on GitHub</a>
<div class="nav">
    <svg class="logo" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" height="189" width="256" viewBox="0,0,1024,756"><image xlink:href="data:image/png;base64,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" x="386.804" y="-70.126" width="87.165" height="87.165" transform="matrix(1.89167 0 0 1.89167 -494.544 432.68)"/><g style="line-height:1;text-align:center" font-weight="400" font-size="72" font-family="Chewy" text-anchor="middle" fill="#fff"><path style="line-height:1;text-align:center" d="M431.1 436.405q-2.26 0-3.59-2.128-1.33-1.995-1.996-4.655-.665-2.66-.93-5.32-.134-2.794-.134-4.257 0-3.99.4-7.98l1.063-7.98q.665-3.991 1.064-7.981.532-3.99.532-7.98 0-1.597-.133-3.459t-.665-3.458q-.399-1.73-1.463-2.793-1.064-1.064-2.793-1.064-1.197 0-2.66 1.995-1.463 1.862-2.793 4.389-1.197 2.394-2.128 4.921-.798 2.394-.798 3.591v5.986q.133 2.926.133 5.985v7.715q-.133 2.394-.4 4.921-.266 2.527-.798 4.921-.399 2.394-1.197 4.656-.798 2.128-2.128 3.857-1.197 1.596-2.926 2.66t-4.123 1.064q-2.66 0-4.39-1.73-1.729-1.728-2.793-4.388-1.064-2.794-1.596-6.119-.399-3.458-.665-6.65-.133-3.325-.133-5.985.133-2.794.133-4.39v-6.251q0-5.054-.266-9.976-.133-4.921-.133-9.842v-6.784q.133-4.123 1.064-7.98.931-3.99 3.06-6.784 2.128-2.793 6.25-2.793 1.863 0 3.326 1.197t2.527 2.926q1.197 1.73 1.862 3.725.798 1.995 1.197 3.458 1.33-1.862 3.06-3.591 1.862-1.73 4.123-3.06t5.054-2.128q2.793-.798 6.385-.798 6.783 0 10.507 3.326 3.857 3.192 5.587 8.113 1.862 4.788 2.26 10.375.4 5.586.4 10.108 0 2.261-.4 6.784-.398 4.389-1.33 9.842-.797 5.32-2.26 11.04-1.464 5.586-3.459 10.242-1.995 4.655-4.788 7.581-2.66 2.926-6.118 2.926zM514.23 385.064q-3.458 9.178-7.58 19.153-4.124 9.976-8.912 18.621-2.927 5.32-6.518 8.912-3.591 3.458-8.246 3.458-3.725 0-7.316-3.724-3.458-3.857-6.517-9.045-2.926-5.32-5.187-10.64-2.129-5.454-3.193-8.646-1.064-3.591-2.394-7.714-1.33-4.124-2.394-8.38-1.064-4.256-1.73-8.38-.664-4.256-.664-7.98 0-5.187 1.862-8.512 1.862-3.326 7.315-3.326 2.129 0 3.858 1.33 1.729 1.198 3.192 3.459 1.463 2.26 2.793 5.32 1.33 2.926 2.527 6.517 1.197 3.459 2.261 7.183 1.064 3.724 2.128 7.182 1.064 3.459 2.129 6.784 1.197 3.192 3.059 6.384.93-2.394 2.128-5.453 1.197-3.06 2.527-6.518 1.463-3.458 3.06-7.049 1.595-3.724 3.457-7.182 3.592-7.183 7.582-12.237 4.123-5.054 8.113-5.054 3.725 0 5.32 2.26 1.597 2.262 1.597 5.853 0 3.591-1.197 8.114-1.197 4.522-3.06 9.31zM532.852 426.297q0 1.729-.798 3.192-.798 1.33-2.129 2.26-1.197.932-2.793 1.464-1.463.532-3.059.532-1.995 0-4.39-.798-2.26-.665-4.255-1.995-1.863-1.33-3.193-3.06-1.33-1.861-1.33-4.123 0-3.59 3.06-5.32 3.059-1.862 6.25-1.862 1.996 0 4.124.665 2.261.665 4.123 1.995 1.995 1.197 3.193 3.06 1.197 1.728 1.197 3.99zM597.094 354.073q0 9.178-4.522 16.76-4.522 7.448-12.104 12.502 3.99.798 6.917 2.793 2.926 1.862 4.788 4.655 1.995 2.794 2.926 6.385 1.064 3.458 1.064 7.315 0 7.582-3.59 13.301-3.459 5.586-9.045 9.31-5.587 3.592-12.503 5.454-6.783 1.729-13.567 1.729-3.724 0-7.448-.532-3.725-.532-7.05-1.862-1.197-.4-1.995-3.192-.798-2.926-1.33-7.05-.399-4.123-.665-8.911-.266-4.789-.399-9.178-.133-4.39-.133-7.714v-4.656-9.31q0-4.788.266-9.444v-3.724q.133-2.527.133-5.586.133-3.192.266-6.65.266-3.459.532-6.385.4-3.06.931-5.187.665-2.128 1.463-2.66 3.06-1.996 6.784-3.725 3.724-1.862 7.714-3.192t7.98-2.128q3.99-.798 7.582-.798 4.921 0 9.444 1.463 4.655 1.33 7.98 4.123 3.458 2.793 5.453 6.917 2.128 3.99 2.128 9.177zm-22.079 7.05q0-2.262-1.463-3.326-1.33-1.064-3.591-1.064-1.862 0-3.99.798-1.995.798-3.459 2.129-.399.399-.798 2.26-.266 1.863-.532 4.257-.133 2.261-.266 4.39V375.353q.133.666.266 1.198.133.532.4.532 2.26 0 4.655-1.597 2.394-1.729 4.256-4.123 1.995-2.527 3.192-5.32 1.33-2.793 1.33-4.921zm-.532 38.439q0-2.261-1.596-3.192t-3.591-.931q-1.862 0-3.99.93-1.996.799-3.326 1.996-.399.266-.665 1.33-.133 1.064-.266 2.394 0 1.197-.133 2.394v4.257l.266 1.596q.266.665.532.665 1.862 0 4.124-1.064 2.26-1.065 4.123-2.66 1.995-1.597 3.192-3.592 1.33-2.128 1.33-4.123zM658.943 409.006q0 4.522-2.926 8.911-2.793 4.256-7.182 7.715-4.257 3.325-9.178 5.453-4.921 1.995-9.044 1.995-7.05 0-12.636-3.458-5.454-3.591-9.178-9.045-3.724-5.453-5.72-11.97-1.861-6.65-1.861-12.902 0-6.65 1.862-14.232 1.862-7.582 5.586-13.966 3.857-6.384 9.577-10.64 5.852-4.257 13.833-4.257 4.522 0 8.645 1.596 4.123 1.463 7.183 4.256 3.059 2.794 4.92 6.65 1.863 3.858 1.863 8.646 0 6.252-2.66 11.705t-7.05 9.444q-4.389 3.99-10.108 6.517-5.72 2.394-11.705 2.66.532 3.857 3.06 5.587 2.66 1.729 6.25 1.729 2.66 0 5.188-1.33 2.527-1.33 4.921-2.927 2.394-1.729 4.655-3.059 2.395-1.33 4.789-1.33 2.66 0 4.788 1.73 2.128 1.728 2.128 4.522zm-21.547-32.587q0-1.73-.931-3.592t-2.926-1.862q-2.926 0-4.922 1.596-1.995 1.463-3.325 3.858-1.33 2.26-2.128 4.92-.665 2.661-.798 4.922 1.862-.399 4.39-1.064 2.66-.665 4.92-1.729 2.395-1.197 3.99-2.926 1.73-1.73 1.73-4.123zM695.387 437.47q-9.443 0-14.763-2.528-5.188-2.394-7.848-6.916-2.527-4.523-3.06-10.907-.531-6.384-.531-14.232 0-6.251.399-12.503.399-6.384 1.197-13.566-1.862 0-2.926.133h-1.73q-4.256 0-7.315-1.863-3.192-1.862-3.192-5.719 0-2.394 1.862-4.256 1.73-1.862 4.522-3.192 2.794-1.33 5.986-2.128 3.325-.799 6.251-1.065 0-3.724.532-8.246t1.862-8.38q1.463-3.99 3.858-6.65 2.394-2.66 6.118-2.66 2.128 0 3.724 1.33 1.73 1.197 2.793 3.192 1.064 1.862 1.596 4.256.532 2.395.532 4.656 0 2.66-.399 5.453t-1.463 6.916q2.261-.266 3.06-.266.798-.133 1.463-.133 2.26 0 4.655.4 2.527.398 4.522 1.33 2.128.93 3.458 2.527 1.33 1.463 1.33 3.59 0 2.395-1.995 4.257-1.862 1.862-4.788 3.192t-6.251 2.128q-3.326.666-6.252.932-.798 5.453-1.197 10.108-.266 4.522-.665 9.178-.266 2.26-.532 4.655-.133 2.261-.133 4.655 0 1.73.133 4.124.133 2.26.931 4.389.798 2.128 2.527 3.724 1.73 1.463 4.789 1.463 1.862 0 3.857-.266 1.995-.399 3.99-.399 3.458 0 4.921 2.394 1.596 2.394 1.596 6.119 0 2.793-1.729 4.92-1.596 1.996-4.123 3.193-2.527 1.33-5.72 1.995-3.058.665-5.852.665z"/><path style="line-height:1;text-align:center" d="M772 413.927q0 2.394-.266 6.251-.133 3.857-1.064 7.582-.931 3.59-2.926 6.251-1.995 2.66-5.587 2.66-4.92 0-7.581-3.192-2.66-3.192-2.527-7.98-4.256 3.724-9.71 5.852-5.32 2.128-10.906 2.128-3.858 0-7.316-1.862-3.458-1.862-6.118-4.788-2.66-2.927-4.257-6.518-1.596-3.724-1.596-7.315 0-6.385 3.858-11.173 3.99-4.788 9.71-8.113 5.719-3.459 12.103-5.454 6.517-1.995 11.572-2.66V384.4q.133-.532.133-1.197 0-4.123-2.262-6.65-2.26-2.528-6.384-2.528-3.325 0-5.72 1.197-2.393 1.065-4.522 2.528-2.128 1.33-4.256 2.527-2.128 1.064-4.788 1.064-2.926 0-5.054-1.862-1.995-1.862-1.995-4.789 0-4.389 2.527-7.98 2.66-3.591 6.517-5.985 3.99-2.528 8.513-3.858 4.655-1.463 8.512-1.463 3.591 0 7.182.931 3.592.798 6.784 2.66 3.192 1.73 5.72 4.39 2.66 2.527 4.256 5.985 2.394 4.921 3.857 10.508 1.463 5.453 2.26 11.305.799 5.72 1.065 11.572.266 5.72.266 11.173zm-21.281-4.256l-.4-9.045h-.531q-.266-.133-.532-.133-1.73 0-4.39.931-2.527.798-4.92 2.394-2.395 1.463-4.124 3.325-1.596 1.863-1.596 3.858 0 1.862 1.33 2.527t3.059.665q2.926 0 6.251-1.463 3.326-1.596 5.853-3.06z"/></g></svg>
    <p>Simple domain redirection tool for debug.</p>
    <p><code>Ctrl+S</code>: Save and Apply</p>
</div>
<div class="editor">
    <textarea id="code" name="code">{{Data}}</textarea>
    <script src="https://cdn.bootcss.com/codemirror/5.42.2/codemirror.min.js"></script>
    <script>
        let form = document.getElementsByTagName("form")[0];
        CodeMirror.fromTextArea(document.getElementById("code"), {
            lineNumbers: true,
            theme: "darcula",
            extraKeys: {
                "Ctrl-S": function(cm) {
                    console.log(cm.getValue());
                    $.post("/", cm.getValue());
                    return;
                }
            }
        });
    </script>
</div>
</body>
</html>
`

type Editor struct {
	data       []byte
	lastModify time.Time
	handler Handler
}

type Handler interface {
	onChange(editor *Editor)
}

func (e * Editor)handleClientRequest(conn net.Conn, reader io.Reader) {
	request, err := http.ReadRequest(bufio.NewReader(reader))
	if err != nil {
		fmt.Print(err)
	}
	writer := httptest.NewRecorder()
	if request.Method == "GET" {
		e.get(request, writer)
	} else {
		e.post(request, writer)
	}
	_ = writer.Result().Write(conn)
}

func (e *Editor) getData() string {
	defer func() {
		if err := recover();err != nil {
			fmt.Println(err)
		}
	}()
	if e.data == nil {
		f := loadConfig()
		e.data, _ = ioutil.ReadAll(f)
		defer f.Close()
	}
	return string(e.data)
}

func (e * Editor)get(request *http.Request, response http.ResponseWriter) {
	var funs = template.FuncMap{
		"Data": e.getData,
	}
	tmpl, err := template.New("body").Funcs(funs).Parse(httpTpl)
	if err != nil {
		print(err)
	}
	var b = &strings.Builder{}
	err = tmpl.Execute(b, e)
	if err != nil {
		print(err)
	}
	_, _ = response.Write([]byte(b.String()))
}

func (e * Editor)post(request *http.Request, response http.ResponseWriter)  {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	e.data = body
	e.lastModify = time.Now()
	if e.handler != nil {
		e.handler.onChange(e)
	}
}