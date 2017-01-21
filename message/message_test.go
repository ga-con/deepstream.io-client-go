// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package message_test

import (
	"github.com/heynemann/deepstream.io-client-go/errors"
	"github.com/heynemann/deepstream.io-client-go/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Message Package", func() {
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
	})
})
