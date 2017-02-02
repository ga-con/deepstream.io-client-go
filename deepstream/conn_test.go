// deepstream.io-client-go
// https://github.com/heynemann/deepstream.io-client-go
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Bernardo Heynemann <heynemann@gmail.com>

package deepstream_test

import (
	"time"

	"github.com/heynemann/deepstream.io-client-go/deepstream"
	"github.com/heynemann/deepstream.io-client-go/interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deepstream Package", func() {
	Describe("[Unit]", func() {
		Describe("Client", func() {
		})
	})

	Describe("[Integration]", func() {
		Describe("Client", func() {
			Describe("Connection", func() {
				It("Should create a connected client", func() {
					opts := deepstream.DefaultOptions()
					opts.AutoLogin = true
					opts.Username = "userA"
					opts.Password = "password"

					cl, err := deepstream.New("localhost:6020", opts)
					Expect(err).NotTo(HaveOccurred())
					time.Sleep(10 * time.Millisecond)

					Expect(cl).NotTo(BeNil())
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateOpen))
				})

				It("Should create a connected client without logging in", func() {
					opts := deepstream.DefaultOptions()
					opts.AutoLogin = false
					opts.Username = "userA"
					opts.Password = "password"
					cl, err := deepstream.New("localhost:6020", opts)
					Expect(err).NotTo(HaveOccurred())

					time.Sleep(10 * time.Millisecond)
					Expect(cl).NotTo(BeNil())
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateAwaitingConnection))
				})

				It("Should create a disconnected client if invalid auth", func() {
					opts := deepstream.DefaultOptions()
					opts.Username = "invalid-user"
					opts.Password = "invalid-pass"
					cl, err := deepstream.New("localhost:6020", opts)
					Expect(err).NotTo(HaveOccurred())

					time.Sleep(10 * time.Millisecond)
					Expect(cl).NotTo(BeNil())
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateError))
				})
			})

			Describe("Authentication", func() {
				It("Should create a connected client then log in", func() {
					opts := deepstream.DefaultOptions()
					opts.AutoLogin = false
					opts.Username = "userA"
					opts.Password = "password"

					cl, err := deepstream.New("localhost:6020", opts)
					Expect(err).NotTo(HaveOccurred())

					time.Sleep(10 * time.Millisecond)
					Expect(cl).NotTo(BeNil())
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateAwaitingConnection))

					err = cl.Login()
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateAuthenticating))
					time.Sleep(10 * time.Millisecond)
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateOpen))
				})

				It("Should login once connection has been established if state is not valid", func() {
					opts := deepstream.DefaultOptions()
					opts.AutoLogin = false
					opts.Username = "userA"
					opts.Password = "password"

					cl, err := deepstream.New("localhost:6020", opts)
					Expect(err).NotTo(HaveOccurred())

					err = cl.Login()
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateInitializing))
					time.Sleep(10 * time.Millisecond)
					Expect(cl.GetConnectionState()).To(Equal(interfaces.ConnectionStateOpen))
				})
			})
		})
	})
})
