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
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ParseHosts takes in hosts file content and returns a map of parsed results.
func ParseHosts(hostsFileContent []byte, err error) (map[string][]string, error) {
	if err != nil {
		return nil, err
	}
	hostsMap := map[string][]string{}
	for _, line := range strings.Split(strings.Trim(string(hostsFileContent), " \t\r\n"), "\n") {
		line = strings.Replace(strings.Trim(line, " \t"), "\t", " ", -1)
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}
		pieces := strings.SplitN(line, " ", 2)
		if len(pieces) > 1 && len(pieces[0]) > 0 {
			if names := strings.Fields(pieces[1]); len(names) > 0 {
				if _, ok := hostsMap[pieces[0]]; ok {
					hostsMap[pieces[0]] = append(hostsMap[pieces[0]], names...)
				} else {
					hostsMap[pieces[0]] = names
				}
			}
		}
	}
	return hostsMap, nil
}


type ServeDNS struct {
	tree *Node
}

func (s *ServeDNS) onChange(editor *Editor)  {
	if editor.data != nil {
		s._load(bufio.NewReader(bytes.NewReader(editor.data)))
		s.saveFile(editor.data)
	}
}

func loadConfig() *os.File {
	f, err := os.Open("hosts.conf")
	path, _ := filepath.Abs(f.Name())
	log.Println("Load Config from " + path)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return nil
	}
	return f
}

func (s *ServeDNS) loadFile(){
	defer func() {
		if err := recover();err != nil {
			fmt.Println(err)
		}
	}()
	f := loadConfig()
	br := bufio.NewReader(f)
	s._load(br)
	defer f.Close()
}

func (s *ServeDNS) saveFile(b []byte){
	defer func() {
		if err := recover();err != nil {
			fmt.Println(err)
		}
	}()
	f, err := os.OpenFile("hosts.conf", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	path, _ := filepath.Abs(f.Name())
	log.Println("Save Config to " + path)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()

	_, _ = f.Write(b)
}

func (s *ServeDNS) _load(br *bufio.Reader){
	log.Println("Reload rules... ")
	defer func() {
		if err := recover();err != nil {
			fmt.Println(err)
		}
	}()
	tree := &Node{}
	count := 0
	for {
		b, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := strings.TrimSpace(strings.Replace(strings.Trim(string(b), " \t"), "\t", " ", -1))
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}
		pieces := strings.SplitN(line, " ", 2)
		if len(pieces) > 1 && len(pieces[0]) > 0 {
			if names := strings.Fields(pieces[1]); len(names) > 0 {
				hosts := strings.Split(pieces[1], ",")
				for _, host := range hosts {
					host = strings.TrimSpace(host)
					tree.addRoute(host, strings.TrimSpace(pieces[0]))
					count++
				}
			}
		}
	}
	log.Println("Reload SUCCESS(", count, ")... ")
	s.tree = tree
}

func (s *ServeDNS) find(domain string) string {
	if s.tree == nil {
		return domain
	}
	handler, _, _ := s.tree.getValue(domain, nil, true)

	if handler != "" {
		return handler
	}

	return domain
}
