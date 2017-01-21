// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package message_test

import (
	"bytes"
	"fmt"

	"github.com/heynemann/deepstream.io-client-go/errors"
	"github.com/heynemann/deepstream.io-client-go/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Message Package", func() {
	Describe("[Unit]", func() {
		Describe("Message Model", func() {
			It("Should create a new Message", func() {
				rawMessage := "R\u001fP\u001fuser/Lisa\u001f1\u001flastname\u001fSOwen"
				message, err := message.NewMessage(rawMessage)
				Expect(err).NotTo(HaveOccurred())
				Expect(message.Topic).To(Equal("R"))
				Expect(message.Action).To(Equal("P"))
				Expect(message.RawData).To(HaveLen(4))
				Expect(message.RawData[0]).To(Equal("user/Lisa"))
				Expect(message.RawData[1]).To(Equal("1"))
				Expect(message.RawData[2]).To(Equal("lastname"))
				Expect(message.RawData[3]).To(Equal("SOwen"))
			})

			It("Should fail on empty message", func() {
				message, err := message.NewMessage("")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.ErrEmptyRawMessage))
				Expect(message).To(BeNil())
			})

			It("Should parse many messages", func() {
				rawMessage := "R\u001fP\u001fuser/Lisa\u001f1\u001flastname\u001fSOwen\u001eR\u001fP\u001fuser/Leslie\u001f1\u001flastname\u001fSWilson"
				messages, err := message.ParseMessages(rawMessage)

				Expect(err).NotTo(HaveOccurred())
				Expect(messages).To(HaveLen(2))

				Expect(messages[0].Topic).To(Equal("R"))
				Expect(messages[0].Action).To(Equal("P"))
				Expect(messages[0].RawData).To(HaveLen(4))

				Expect(messages[1].Topic).To(Equal("R"))
				Expect(messages[1].Action).To(Equal("P"))
				Expect(messages[1].RawData).To(HaveLen(4))
				Expect(messages[1].RawData[3]).To(Equal("SWilson"))
			})

			It("Should fail on empty message when parsing many", func() {
				message, err := message.ParseMessages("")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.ErrEmptyRawMessage))
				Expect(message).To(BeNil())
			})

			It("Should fail on empty message when parsing many and one is empty", func() {
				rawMessage := "R\u001fP\u001fuser/Lisa\u001f1\u001flastname\u001fSOwen\u001e"
				message, err := message.ParseMessages(rawMessage)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.ErrEmptyRawMessage))
				Expect(message).To(BeNil())
			})

			Measure("it should parse messages efficiently", func(b Benchmarker) {
				var buffer bytes.Buffer

				for i := 0; i < 1000; i++ {
					buffer.WriteString(fmt.Sprintf("R\u001fP\u001fuser/Lisa\u001f1\u001flastname\u001fSOwen%d\u001e", i))
				}
				buffer.WriteString(fmt.Sprintf("R\u001fP\u001fuser/Lisa\u001f1\u001flastname\u001fSOwenLast"))
				raw := buffer.String()

				runtime := b.Time("runtime", func() {
					_, err := message.ParseMessages(raw)
					Expect(err).NotTo(HaveOccurred())
				})

				Expect(runtime.Seconds()).Should(BeNumerically("<", 0.01), "Parsing messages shouldn't take too long.")
			}, 200)
		})
	})
})
