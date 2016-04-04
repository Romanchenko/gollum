// Copyright 2015-2016 trivago GmbH
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

package core

import (
	"github.com/trivago/tgo/ttesting"
	"testing"
	"time"
)

func getMockMessage(data string) *Message {
	return &Message{
		StreamID:     1,
		PrevStreamID: 2,
		Timestamp:    time.Now(),
		Sequence:     4,
		Data:         []byte(data),
	}
}

func TestMessageEnqueue(t *testing.T) {
	expect := ttesting.NewExpect(t)
	msgString := "Test for Enqueue()"
	msg := getMockMessage(msgString)
	buffer := NewMessageQueue(0)

	expect.Equal(MessageStateDiscard, buffer.Push(msg, -1))

	go func() {
		expect.Equal(MessageStateOk, buffer.Push(msg, 0))
	}()

	retMsg, _ := buffer.Pop()
	expect.Equal(msgString, retMsg.String())

	retStatus := buffer.Push(msg, 10*time.Millisecond)
	expect.Equal(MessageStateTimeout, retStatus)

	go func() {
		expect.Equal(MessageStateOk, buffer.Push(msg, 1*time.Second))
	}()

	retMsg, _ = buffer.Pop()
	expect.Equal(msgString, retMsg.String())
}

func TestMessageRoute(t *testing.T) {
	expect := ttesting.NewExpect(t)
	msgString := "Test for Route()"
	msg := getMockMessage(msgString)

	mockDistributer := func(msg *Message) {
		expect.Equal(msgString, msg.String())
	}

	mockStream := StreamBase{}
	mockStream.filters = []Filter{&mockFilter{}}
	mockStream.distribute = mockDistributer
	mockStream.formatters = []Formatter{&mockFormatter{}}
	mockStream.boundStreamID = 1
	mockStream.AddProducer(&mockProducer{})
	StreamRegistry.Register(&mockStream, 1)

	msg.Route(1)

}
