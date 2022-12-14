// Copyright © 2021 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sshutil

import (
	"bufio"
	"io"
	"strings"

	"github.com/ergoapi/log"

	"github.com/sirupsen/logrus"
)

// Cmd is in host exec cmd
func (ss *SSH) Cmd(host string, cmd string) []byte {
	ss.Log.Infof("[ssh][%s] %s", host, cmd)
	session, err := ss.Connect(host)
	defer func() {
		if r := recover(); r != nil {
			ss.Log.Errorf("[ssh][%s]Error create ssh session failed,%s", host, err)
		}
	}()
	if err != nil {
		panic(1)
	}
	defer session.Close()
	b, err := session.CombinedOutput(cmd)
	ss.Log.Debugf("[ssh][%s] command result is: %s", host, string(b))
	defer func() {
		if r := recover(); r != nil {
			ss.Log.Errorf("[ssh][%s]Error exec command failed: %s", host, err)
		}
	}()
	if err != nil {
		panic(1)
	}
	return b
}

func readPipe(slog log.Logger, host string, pipe io.Reader, isErr bool) {
	r := bufio.NewReader(pipe)
	for {
		line, _, err := r.ReadLine()
		if line == nil {
			return
		} else if err != nil {
			slog.Infof("[%s] %s", host, line)
			slog.Errorf("[ssh] [%s] %s", host, err)
			return
		} else {
			if isErr {
				logrus.Errorf("[%s] %s", host, line)
			} else {
				logrus.Infof("[%s] %s", host, line)
			}
		}
	}
}

func (ss *SSH) CmdAsync(host string, cmd string) error {
	ss.Log.Debugf("[%s] %s", host, cmd)
	session, err := ss.Connect(host)
	if err != nil {
		ss.Log.Errorf("[ssh][%s]Error create ssh session failed,%s", host, err)
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		ss.Log.Errorf("[ssh][%s]Unable to request StdoutPipe(): %s", host, err)
		return err
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		ss.Log.Errorf("[ssh][%s]Unable to request StderrPipe(): %s", host, err)
		return err
	}
	if err := session.Start(cmd); err != nil {
		ss.Log.Errorf("[ssh][%s]Unable to execute command: %s", host, err)
		return err
	}
	doneout := make(chan bool, 1)
	doneerr := make(chan bool, 1)
	go func() {
		readPipe(ss.Log, host, stderr, true)
		doneerr <- true
	}()
	go func() {
		readPipe(ss.Log, host, stdout, false)
		doneout <- true
	}()
	<-doneerr
	<-doneout
	return session.Wait()
}

// CmdToString is in host exec cmd and replace to spilt str
func (ss *SSH) CmdToString(host, cmd, spilt string) string {
	if data := ss.Cmd(host, cmd); data != nil {
		str := string(data)
		str = strings.ReplaceAll(str, "\r\n", spilt)
		return str
	}
	return ""
}
