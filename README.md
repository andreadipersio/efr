efr
===

[![GoDoc](https://godoc.org/github.com/andreadipersio/efr?status.png)](https://godoc.org/github.com/andreadipersio/efr)

*efr* stands for **Event Forwarder**.

Given a stream of unordered events, forward them
in the correct order to all connected clients, in a social graph fashion.

The system is written in **golang** and his composed by 4 main modules.

### usage
Install using `go install` and run by invoking `efr`

```shell
cd $GOPATH
go install github.com/andreadipersio/efr
bin/efr
```

This will start the program with default values:
```shell
--resequencerType=stream
--maxProcs=1
--eventSourcePort=9090
--clientPort=9099
--sequenceIndex=0
```

### event
It provide the `Event` and `Subscriber` interface

**event.Event**
An interface that describe behaviour of events.
Event should implement methods to access their Sequence number,
sender and recipient id.
Factory for this method is a function that take a string and parse it 
returning a concrete Event value.

**event.Subscriber**
An interface that describe behaviour of subscribers.
Subscriber should implement method to handle different event types
and access their ID.
Factory for this method is a function that take a string, representing
subscriber ID and returning a concrete User value.

### listener
Listen for TCP connection from an Event Source.
Once a connection is made it start reading CRLF terminated strings from
the event source and forward them to the **event-dispatcher**.

Before dispatching, events go through a **resequencer** which reorder them
based on their sequence ID.

Two type of resequencer are supported:
(`resequencerType` parameter)

- 'batch' resequencer: each event is appended in a buffer, once buffer length reach **resequencerCapacity**
event are ordered and then sent through the event channel.
The bigger the capacity, the more memory is consumed and client can timeout waiting
for data.

- 'stream' resequencer: Everytime an event is received, it order the buffer and verify that every
element in the buffer has a progressive sequence number.
CPU intensive, client can timeout waiting for data if event source randomness is high.

### subscription
Listen for TCP Connection from client.
Each connection should send a CRLF terminated strings containing
the ID of the subscriber.
Once ID is received a **SubscriptionRequest** is created, containing

- ID: ID of subscriber
- Connection: Connection to the client

SubscriptionRequest are then sent through subscription channel.

### dispatcher
Listen to the following channels:

- subscription channel: A new subscription request has been received. Subscribe the user
to the directory 

- event channel: A new event has been received. Get Sender and Recipient from the directory (or create
them as disconnected subscriber if they do not exist) and invoke ther `HandleEvent` method.

- EventSourceClosed channel: When a value is received through this channel, unsubscribe all the clients
and close the connection

[Diagram](https://www.dropbox.com/s/qe08veyzsurn0m1/eft-diagram.png)
