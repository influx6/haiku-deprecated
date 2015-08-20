# Scope

## Why (Use Cases)

OOP is hard for todays systems that need to be event based, but FRP makes it much easier to reason.




## Distributed Computation

The user does not write code bound to the physical topology.
The developer can then run a very concurrent program, and simulate it on their local machine, and the same code can be deployed to run across many computers.

### State
- We need all state queues to be durable & immutable across machines.
- We need all state queues to support Fan in, Fan Out
- We need all state queues to support back and forward pressure tolerance.
http
s://github.com/bitly/nsq ?

### Transport
Maybe libchan.
Maybe Grpc.



### Discovery
- We need them to be easy to be discovered. Reflextive
- We need them to be easy to bind to for processes that need to talk to us.

- grpc looks like the winner
