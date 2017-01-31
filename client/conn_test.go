// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package client_test

import (
	"time"

	"github.com/heynemann/deepstream.io-client-go/client"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Package", func() {
	//Describe("[Unit]", func() {
	//Describe("Client", func() {
	//var protocol *testing.MockProtocol

	//BeforeEach(func() {
	//protocol = testing.NewMockProtocol()
	//})

	//Describe("Connection", func() {
	//It("Should create a client", func() {
	//client, err := client.New("localhost:6020", protocol)
	//Expect(err).NotTo(HaveOccurred())
	//Expect(client).NotTo(BeNil())
	//Expect(client.ConnectionState).To(Equal(interfaces.ConnectionStateAwaitingAuthentication))
	//Expect(protocol.HasConnected).To(BeTrue())
	//})

	//It("Should be in error state when an error happens", func() {
	//expErr := fmt.Errorf("mock error")
	//protocol.Error = expErr

	//client, err := client.New("localhost:6020", protocol)
	//Expect(client).NotTo(BeNil())
	//Expect(err).To(MatchError(expErr))
	//Expect(client.ConnectionState).To(Equal(interfaces.ConnectionStateError))
	//Expect(protocol.HasConnected).To(BeFalse())
	//})

	//It("Should close a connection", func() {
	//client, err := client.New("localhost:6020", protocol)
	//Expect(err).NotTo(HaveOccurred())

	//err = client.Close()
	//Expect(err).NotTo(HaveOccurred())
	//Expect(client.ConnectionState).To(Equal(interfaces.ConnectionStateClosed))

	//Expect(protocol.IsClosed).To(BeTrue())
	//})

	//It("Should error when closing a connection", func() {
	//client, err := client.New("localhost:6020", protocol)
	//Expect(err).NotTo(HaveOccurred())

	//expErr := fmt.Errorf("mock error")
	//protocol.Error = expErr
	//err = client.Close()
	//Expect(err).To(MatchError(expErr))
	//Expect(client.ConnectionState).To(Equal(interfaces.ConnectionStateError))

	//Expect(protocol.IsClosed).To(BeFalse())
	//})

	//})

	//Describe("Authentication", func() {
	//It("Should authenticate", func() {
	//client, err := client.New("localhost:6020", protocol)
	//Expect(err).NotTo(HaveOccurred())
	//Expect(client).NotTo(BeNil())

	//authParams := map[string]interface{}{
	//"user":     "x",
	//"password": "y",
	//}
	//err = client.Login(authParams)
	//Expect(err).NotTo(HaveOccurred())

	//Expect(protocol.IsAuthenticated).To(BeTrue())
	//Expect(protocol.AuthParams).To(Equal(authParams))
	//})

	//It("Should return error when authentication fails", func() {
	//client, err := client.New("localhost:6020", protocol)
	//Expect(err).NotTo(HaveOccurred())
	//Expect(client).NotTo(BeNil())

	//expErr := fmt.Errorf("mock error")
	//protocol.Error = expErr

	//authParams := map[string]interface{}{
	//"user":     "x",
	//"password": "y",
	//}
	//err = client.Login(authParams)
	//Expect(err).To(HaveOccurred())
	//Expect(err).To(MatchError(expErr))

	//Expect(protocol.IsAuthenticated).To(BeFalse())
	//Expect(protocol.AuthParams).To(Equal(authParams))
	//})
	//})
	//})
	//})

	Describe("[Integration]", func() {
		Describe("Client", func() {
			Describe("Connection", func() {
				It("Should create a connected client", func() {
					cl, err := client.New("localhost:6020", "userA", "password")
					time.Sleep(10 * time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
					Expect(cl).NotTo(BeNil())
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateOpen))
				})

				It("Should create a connected client without logging in", func() {
					opts := client.DefaultOptions()
					opts.AutoLogin = false
					cl, err := client.New("localhost:6020", "userA", "password", opts)
					time.Sleep(10 * time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
					Expect(cl).NotTo(BeNil())
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateAwaitingConnection))
				})

				It("Should create a connected client then log in", func() {
					opts := client.DefaultOptions()
					opts.AutoLogin = false
					cl, err := client.New("localhost:6020", "userA", "password", opts)
					time.Sleep(10 * time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
					Expect(cl).NotTo(BeNil())
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateAwaitingConnection))

					err = cl.Login()
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateAuthenticating))
					time.Sleep(10 * time.Millisecond)
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateOpen))
				})

				It("Should create a disconnected client if invalid auth", func() {
					cl, err := client.New("localhost:6020", "qwe", "password!")
					time.Sleep(10 * time.Millisecond)
					Expect(err).NotTo(HaveOccurred())
					Expect(cl).NotTo(BeNil())
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateError))
				})
			})

			//Describe("Authentication", func() {
			//It("Should send authentication message", func() {
			//client, err := client.New("localhost:6020")
			//Expect(err).NotTo(HaveOccurred())

			//err = client.Login(map[string]interface{}{
			//"username": "userA",
			//"password": "password",
			//})
			//Expect(err).NotTo(HaveOccurred())

			//Expect(client.GetConnectionState()).To(Equal(interfaces.ConnectionStateOpen))
			//})
			//})
		})
	})
})
