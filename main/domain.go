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
	"os"
	"io"
	"strings"
	"log"
	"path/filepath"
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
	data []byte
	tree *Node
}

type States struct {
	serveDNS *ServeDNS
}

// private
var instance *States

// public
func GetInstance() *States {
	if instance == nil {
		instance = &States{}     // not thread safe
	}
	return instance
}

func (s *ServeDNS) clear()  {
	s.tree = &Node{}
}


func (s *ServeDNS) load(){
	f, err := os.Open("hosts.conf")
	path, _ := filepath.Abs(f.Name())
	log.Println("Load Config from " + path)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	s._load(br)
}

func (s *ServeDNS) read() {

}

func (s *ServeDNS) _load(br *bufio.Reader){
	s.tree = &Node{}

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
					s.tree.addRoute(host, strings.TrimSpace(pieces[0]))
				}
			}
		}
	}
}

func (s *ServeDNS) find(domain string) string {
	handler, _, _ := s.tree.getValue(domain, nil, true)

	if handler != "" {
		return handler
	}

	return domain
}
