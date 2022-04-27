# Connection Architecture
In order to handler connections from multiple clients at a time, we need a multithreaded way of handling connections.

The general idea is a producer-consumer relationship between one and multiple threads. A single thread accepts the connections and then puts it onto a queue that a specified number of threads can then pull from. This ofcourse needs to be done in a thread safe manner. 

This is done in the code by the two classes: ConnectionHandler and ClientQueue. 

The ConnectionHandler is a configurable class that holds a reference to a client queue. This handler should be started from the server, whether the server itself accepts the clients or another thread does is not important. 

The ClientQueue is a thread-safe data structure the holds the active connections. 