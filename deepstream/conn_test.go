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
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("deepstream Package", func() {
	Describe("[Unit]", func() {
		Describe("Client", func() {
		})
	})

	Describe("[Integration]", func() {
		var authOpts *deepstream.ClientOptions
		BeforeEach(func() {
			authOpts = deepstream.DefaultOptions()
			authOpts.AutoLogin = true
			authOpts.Username = "userA"
			authOpts.Password = "password"
		})

		Describe("Client", func() {
			Describe("Connection", func() {
				It("Should create a connected client", func() {
					cl, err := deepstream.New("localhost:6020", authOpts)
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

			Describe("Events", func() {
				var client *deepstream.Client
				BeforeEach(func() {
					var err error
					client, err = deepstream.New("localhost:6020", authOpts)
					Expect(err).NotTo(HaveOccurred())
					time.Sleep(10 * time.Millisecond)
				})

				Describe("Subscriptions", func() {
					It("Should subscribe to events", func() {
						onMessage := func(msg *deepstream.EventMessage) error { return nil }
						topic := uuid.NewV4().String()
						client.Event.Subscribe(topic, onMessage)
						Expect(client.Event.Subscriptions[topic].Event).To(Equal(topic))
						Expect(client.Event.Subscriptions[topic].Handlers).To(HaveLen(1))
						Expect(client.Event.Subscriptions[topic].Acked).To(BeFalse())

						time.Sleep(10 * time.Millisecond)

						Expect(client.Event.Subscriptions[topic].Acked).To(BeTrue())
					})

					It("Should subscribe twice only causes one trip to the server", func() {
						onMessage := func(msg *deepstream.EventMessage) error { return nil }
						topic := uuid.NewV4().String()
						client.Event.Subscribe(topic, onMessage)

						time.Sleep(10 * time.Millisecond)
						Expect(client.Event.Subscriptions[topic].Acked).To(BeTrue())

						client.Event.Subscriptions[topic].Acked = false
						client.Event.Subscribe(topic, onMessage)

						time.Sleep(10 * time.Millisecond)
						//Did not go to deepstream server again, thus did not receive another ack
						Expect(client.Event.Subscriptions[topic].Acked).To(BeFalse())
						Expect(client.Event.Subscriptions[topic].Handlers).To(HaveLen(2))
					})
				})

				Describe("Publishing", func() {
					It("Should publish events", func() {
						_client, err := deepstream.New("localhost:6020", authOpts)
						Expect(err).NotTo(HaveOccurred())
						time.Sleep(10 * time.Millisecond)

						var msg *deepstream.EventMessage
						onMessage := func(_msg *deepstream.EventMessage) error {
							msg = _msg
							return nil
						}
						topic := uuid.NewV4().String()
						_client.Event.Subscribe(topic, onMessage)
						time.Sleep(10 * time.Millisecond)

						m := map[string]interface{}{"qwe": 123}
						err = client.Event.Publish(topic, "yetAnotherValue", 10, true, false, m)
						Expect(err).NotTo(HaveOccurred())
						time.Sleep(10 * time.Millisecond)

						Expect(msg).NotTo(BeNil())
						Expect(msg.Data[0]).To(BeEquivalentTo("yetAnotherValue"))
						Expect(msg.Data[1]).To(BeEquivalentTo(10))
						Expect(msg.Data[2]).To(BeTrue())
						Expect(msg.Data[3]).To(BeFalse())

						obj := msg.Data[4].(map[string]interface{})
						Expect(obj["qwe"]).To(BeEquivalentTo(123))
					})
				})
			})
		})
	})
})
