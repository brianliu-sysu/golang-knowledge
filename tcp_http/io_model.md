# I/O model
**what is an I/O models?**

An I/O model mainly answers two questions:
1. How does thread/coroutine wait while data is arriving?
(e.g.,when a socket receives data)
2. How is data copied from kernal space to user space?
    * who blocks?
    * who notifies?

## Blocking I/O
**Behavior**

```read()``` blocks untils data is ready + data is copied to user space, then returns

**Pros**
* simple code

**Cons**
* One connection/request may occupy a thread
* Poor scalability

**Typical use**
* Traditional one-thread-per-connection or one-process-per-connection servers

## Non-Blocking I/O
**Behavior**

After setting the file description(```fd```) to non-blocking, ```read()``` returns immediately with ```EAGAIN``` or ```EWOULDBLOCK``` if no data is available

**What you need to do**
* retry continuously(polling), or combine with event notification(which become I/O multiplexing)

**Pros**
* Thread are not blocked by single ```read```

**Cons**
* Pure polling wastes CPU

## I/O multiplexing(select/poll/epoll/kqueus)
**Behavior**

use ```select```/```poll```/```epoll```/```kqueue``` to wait for which file descriptions are readable/writable, then perform ```read```/```write``` on the ready ones

**Pros**
* One thread can handle many connections

**Cons**
* More complex programming
* Still requires copying data from kernal space to user space

**Typical use**
* Ngix/Redis(event loop)
* Go netpoll

## Singal-driven I/O
**Behavior**

Register a singal(e.g. SIGIO). when data is ready, the kernal send a singal to notify the process,  then you can ```read```

**Characteristic**
* Notification is asynchronous
* Data copy still need explicit ```read```

**In practice**
* Rarely used due to complexity of singal handling

## Asynchronous I/O (AIO)
**Behavior**

* submit a asynchronous read and return immediately
* where data is read and copy is complete, the kernal notifies you(via callback/event/completion queue)

**Key points**
* unlike multiplexing, the data copy to user space is alse done by the kernal
* while notified, the data is typically ready to use

**Typical Use**
* Windows IOCP
* Linux io_uring(modern high-performance async I/O)

**Summary**
* Blocking: read() waits until everything is done
* Non-blocking: read() returns immediately if no data
* Multiplexing: wait for “ready events”, then read() and copy
* Signal-driven: kernel signals “ready”, then read()
* Async (AIO): submit and return; notified when fully done (including copy)
