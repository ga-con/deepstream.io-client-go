// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client_test

import (
	"fmt"

	"github.com/heynemann/deepstream.io-client-go/client"
	"github.com/heynemann/deepstream.io-client-go/testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Package", func() {
	Describe("Client", func() {
		var protocol *testing.MockProtocol

		BeforeEach(func() {
			protocol = testing.NewMockProtocol()
		})

		It("Should create a client", func() {
			client := client.New("localhost:6020", protocol)
			Expect(client).NotTo(BeNil())
		})

		Describe("Authentication", func() {
			It("Should authenticate", func() {
				client := client.New("localhost:6020", protocol)
				Expect(client).NotTo(BeNil())

				err := client.Login()
				Expect(err).NotTo(HaveOccurred())

				Expect(protocol.IsAuthenticated).To(BeTrue())
			})

			It("Should return error when authentication fails", func() {
				expErr := fmt.Errorf("mock error")
				protocol.Error = expErr
				client := client.New("localhost:6020", protocol)
				Expect(client).NotTo(BeNil())

				err := client.Login()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(expErr))

				Expect(protocol.IsAuthenticated).To(BeFalse())
			})

		})
	})
})
