// Copyright 2013-2014 Apcera Inc. All rights reserved.

package test

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats"
)

func TestServerRestartReSliceIssue(t *testing.T) {
	srvA, srvB, optsA, optsB := runServers(t)
	defer srvA.Shutdown()

	urlA := fmt.Sprintf("nats://%s:%d/", optsA.Host, optsA.Port)
	urlB := fmt.Sprintf("nats://%s:%d/", optsB.Host, optsB.Port)

	// msg to send..
	msg := []byte("Hello World")

	servers := []string{urlA, urlB}

	opts := nats.DefaultOptions
	opts.Timeout = (5 * time.Second)
	opts.ReconnectWait = (50 * time.Millisecond)
	opts.MaxReconnect = 1000

	reconnects := int32(0)
	reconnectsDone := make(chan bool)
	opts.ReconnectedCB = func(nc *nats.Conn) {
		atomic.AddInt32(&reconnects, 1)
		reconnectsDone <- true
	}

	// Create 20 random clients.
	// Half connected to A and half to B..
	numClients := 20
	for i := 0; i < numClients; i++ {
		opts.Url = servers[i%2]
		nc, err := opts.Connect()
		defer nc.Close()

		if err != nil {
			t.Fatalf("Failed to create connection: %v\n", err)
		}
		// Create 10 subscriptions each..
		for x := 0; x < 10; x++ {
			subject := fmt.Sprintf("foo.%d", (rand.Int()%50)+1)
			nc.Subscribe(subject, func(m *nats.Msg) {
				// Just eat it..
			})
		}
		// Pick one subject to send to..
		subject := fmt.Sprintf("foo.%d", (rand.Int()%50)+1)
		go func() {
			time.Sleep(10 * time.Millisecond)
			for i := 1; 1 <= 100; i++ {
				if err := nc.Publish(subject, msg); err != nil {
					return
				}
				if i%10 == 0 {
					time.Sleep(time.Millisecond)
				}
			}
		}()
	}

	// Wait for a short bit..
	time.Sleep(20 * time.Millisecond)

	// Restart SrvB
	srvB.Shutdown()
	srvB = RunServer(optsB)
	defer srvB.Shutdown()

	select {
	case <-reconnectsDone:
		break
	case <-time.After(2 * time.Second):
		t.Fatalf("Expected %d reconnects, got %d\n", numClients/2, reconnects)
	}
}

// This will test queue subscriber semantics across a cluster in the presence
// of server restarts.
func TestServerRestartAndQueueSubs(t *testing.T) {
	srvA, srvB, optsA, optsB := runServers(t)

	urlA := fmt.Sprintf("nats://%s:%d/", optsA.Host, optsA.Port)
	urlB := fmt.Sprintf("nats://%s:%d/", optsB.Host, optsB.Port)

	// Client options
	opts := nats.DefaultOptions
	opts.Timeout = (5 * time.Second)
	opts.ReconnectWait = (50 * time.Millisecond)
	opts.MaxReconnect = 1000
	opts.NoRandomize = true

	// Allow us to block on a reconnect completion.
	reconnectsDone := make(chan bool)
	opts.ReconnectedCB = func(nc *nats.Conn) {
		reconnectsDone <- true
	}

	// Helper to wait on a reconnect.
	waitOnReconnect := func() {
		var rcs int64
		for {
			select {
			case <-reconnectsDone:
				atomic.AddInt64(&rcs, 1)
				if rcs >= 2 {
					return
				}
			case <-time.After(2 * time.Second):
				t.Fatalf("Expected a reconnect, timedout!\n")
			}
		}
	}

	// Create two clients..
	opts.Servers = []string{urlA}
	nc1, err := opts.Connect()
	if err != nil {
		t.Fatalf("Failed to create connection for nc1: %v\n", err)
	}

	opts.Servers = []string{urlB}
	nc2, err := opts.Connect()
	if err != nil {
		t.Fatalf("Failed to create connection for nc2: %v\n", err)
	}

	c1, _ := nats.NewEncodedConn(nc1, "json")
	defer c1.Close()
	c2, _ := nats.NewEncodedConn(nc2, "json")
	defer c2.Close()

	// Flusher helper function.
	flush := func() {
		// Wait for processing.
		c1.Flush()
		c2.Flush()
		// Wait for a short bit for cluster propogation.
		time.Sleep(50 * time.Millisecond)
	}

	// To hold queue results.
	results := make(map[int]int)
	var mu sync.Mutex

	// This corresponds to the subsriptions below.
	const ExpectedMsgCount = 3

	// Make sure we got what we needed, 1 msg only and all seqnos accounted for..
	checkResults := func(numSent int) {
		mu.Lock()
		defer mu.Unlock()

		for i := 0; i < numSent; i++ {
			if results[i] != ExpectedMsgCount {
				t.Fatalf("Received incorrect number of messages, [%d] for seq: %d\n", results[i], i)
			}
		}

		// Auto reset results map
		results = make(map[int]int)
	}

	subj := "foo.bar"
	qgroup := "workers"

	cb := func(seqno int) {
		mu.Lock()
		defer mu.Unlock()
		results[seqno] = results[seqno] + 1
	}

	// Create queue subscribers
	c1.QueueSubscribe(subj, qgroup, cb)
	c2.QueueSubscribe(subj, qgroup, cb)

	// Do a wildcard subscription.
	c1.Subscribe("foo.*", cb)
	c2.Subscribe("foo.*", cb)

	// Wait for processing.
	flush()

	sendAndCheckMsgs := func(numToSend int) {
		for i := 0; i < numToSend; i++ {
			if i%2 == 0 {
				c1.Publish(subj, i)
			} else {
				c2.Publish(subj, i)
			}
		}
		// Wait for processing.
		flush()
		// Check Results
		checkResults(numToSend)
	}

	////////////////////////////////////////////////////////////////////////////
	// Base Test
	////////////////////////////////////////////////////////////////////////////

	// Now send 10 messages, from each client..
	sendAndCheckMsgs(10)

	////////////////////////////////////////////////////////////////////////////
	// Now restart SrvA and srvB, re-run test
	////////////////////////////////////////////////////////////////////////////

	srvA.Shutdown()
	srvA = RunServer(optsA)
	defer srvA.Shutdown()

	srvB.Shutdown()
	srvB = RunServer(optsB)
	defer srvB.Shutdown()

	waitOnReconnect()

	time.Sleep(50 * time.Millisecond)

	// Now send another 10 messages, from each client..
	sendAndCheckMsgs(10)
}

// This will test request semantics across a route
func TestRequestsAcrossRoutes(t *testing.T) {
	srvA, srvB, optsA, optsB := runServers(t)
	defer srvA.Shutdown()
	defer srvB.Shutdown()

	urlA := fmt.Sprintf("nats://%s:%d/", optsA.Host, optsA.Port)
	urlB := fmt.Sprintf("nats://%s:%d/", optsB.Host, optsB.Port)

	nc1, err := nats.Connect(urlA)
	if err != nil {
		t.Fatalf("Failed to create connection for nc1: %v\n", err)
	}
	defer nc1.Close()

	nc2, err := nats.Connect(urlB)
	if err != nil {
		t.Fatalf("Failed to create connection for nc2: %v\n", err)
	}
	defer nc2.Close()

	ec2, _ := nats.NewEncodedConn(nc2, nats.JSON_ENCODER)

	response := []byte("I will help you")

	// Connect responder to srvA
	nc1.Subscribe("foo-req", func(m *nats.Msg) {
		nc1.Publish(m.Reply, response)
	})
	// Make sure the route and the subscription are propogated.
	nc1.Flush()

	var resp string

	for i := 0; i < 100; i++ {
		if err := ec2.Request("foo-req", i, &resp, 100*time.Millisecond); err != nil {
			t.Fatalf("Received an error on Request test [%d]: %s", i, err)
		}
	}
}

// This will test request semantics across a route to queues
func TestRequestsAcrossRoutesToQueues(t *testing.T) {
	srvA, srvB, optsA, optsB := runServers(t)
	defer srvA.Shutdown()
	defer srvB.Shutdown()

	urlA := fmt.Sprintf("nats://%s:%d/", optsA.Host, optsA.Port)
	urlB := fmt.Sprintf("nats://%s:%d/", optsB.Host, optsB.Port)

	nc1, err := nats.Connect(urlA)
	if err != nil {
		t.Fatalf("Failed to create connection for nc1: %v\n", err)
	}
	defer nc1.Close()

	nc2, err := nats.Connect(urlB)
	if err != nil {
		t.Fatalf("Failed to create connection for nc2: %v\n", err)
	}
	defer nc2.Close()

	ec1, _ := nats.NewEncodedConn(nc1, nats.JSON_ENCODER)
	ec2, _ := nats.NewEncodedConn(nc2, nats.JSON_ENCODER)

	response := []byte("I will help you")

	// Connect one responder to srvA
	nc1.QueueSubscribe("foo-req", "booboo", func(m *nats.Msg) {
		nc1.Publish(m.Reply, response)
	})
	// Make sure the route and the subscription are propogated.
	nc1.Flush()

	// Connect the other responder to srvB
	nc2.QueueSubscribe("foo-req", "booboo", func(m *nats.Msg) {
		nc2.Publish(m.Reply, response)
	})

	var resp string

	for i := 0; i < 100; i++ {
		if err := ec2.Request("foo-req", i, &resp, 100*time.Millisecond); err != nil {
			t.Fatalf("Received an error on Request test [%d]: %s", i, err)
		}
	}

	for i := 0; i < 100; i++ {
		if err := ec1.Request("foo-req", i, &resp, 100*time.Millisecond); err != nil {
			t.Fatalf("Received an error on Request test [%d]: %s", i, err)
		}
	}
}
