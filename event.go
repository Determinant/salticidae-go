package salticidae

// #include "salticidae/event.h"
// #include <signal.h>
import "C"
import "runtime"

// The C pointer type for an EventContext handle.
type CEventContext = *C.eventcontext_t
type eventContext struct {
    inner CEventContext
    attached map[uintptr]interface{}
}
// The handle for an event loop.
type EventContext = *eventContext

func NewEventContext() EventContext {
    res := &eventContext{
        inner: C.eventcontext_new(),
        attached: make(map[uintptr]interface{}),
    }
    runtime.SetFinalizer(res, func(self EventContext) { self.free() })
    return res
}

func (self EventContext) attach(ptr rawptr_t, x interface{}) { self.attached[uintptr(ptr)] = x }
func (self EventContext) detach(ptr rawptr_t) { delete(self.attached, uintptr(ptr)) }
func (self EventContext) free() { C.eventcontext_free(self.inner) }

// Start the event loop. This is a blocking call that will hand over the
// control flow to the infinite loop which triggers callbacks upon new events.
// The function will return when Stop() is called.
func (self EventContext) Dispatch() { C.eventcontext_dispatch(self.inner) }

// Stop the event loop. This function is typically called in a callback. Notice
// that all methods of an EventContext should be invoked by the same thread
// which logically owns the loop. To schedule code executed in the event loop
// of a particular thread, see ThreadCall.
func (self EventContext) Stop() { C.eventcontext_stop(self.inner) }

// The C pointer type for a ThreadCall handle.
type CThreadCall = *C.threadcall_t
type threadCall struct { inner CThreadCall }
// The handle for scheduling a function call executed by a particular
// EventContext event loop.
type ThreadCall = *threadCall

// The C function pointer type which takes threadcall_handle_t* and void*
// (passing in the custom user data allocated by C.malloc) as parameters.
type ThreadCallCallback = C.threadcall_callback_t

// Create a ThreadCall handle attached to the given EventContext. Each
// invokcation of AsyncCall() will schedule a function call executed in the
// given EventContext event loop.
func NewThreadCall(ec EventContext) ThreadCall {
    res := &threadCall{ inner: C.threadcall_new(ec.inner) }
    runtime.SetFinalizer(res, func(self ThreadCall) { self.free() })
    return res
}

func (self ThreadCall) free() { C.threadcall_free(self.inner) }

// Schedule a function to be executed in the target event loop.
func (self ThreadCall) AsyncCall(callback ThreadCallCallback, userdata rawptr_t) {
    C.threadcall_async_call(self.inner, callback, userdata)
}

// The C pointer type for TimerEvent handle.
type CTimerEvent = *C.timerev_t
type timerEvent struct {
    inner CTimerEvent
    ec EventContext
}

// The handle for a timed event.
type TimerEvent = *timerEvent

// The C function pointer type which takes timerev_t* (the C pointer to
// TimerEvent) and void* (the unsafe pointer to any userdata) as parameter.
type TimerEventCallback = C.timerev_callback_t

// Create a TimerEvent handle attached to the given EventContext, with a
// registered callback.
func NewTimerEvent(_ec EventContext, cb TimerEventCallback, userdata rawptr_t) TimerEvent {
    res := &timerEvent{
        inner: C.timerev_new(_ec.inner, cb, userdata),
        ec: _ec,
    }
    _ec.attach(rawptr_t(res.inner), res)
    runtime.SetFinalizer(res, func(self TimerEvent) { self.free() })
    return res
}

func (self TimerEvent) free() { C.timerev_free(self.inner) }

// Change the callback.
func (self TimerEvent) SetCallback(callback TimerEventCallback, userdata rawptr_t) {
    C.timerev_set_callback(self.inner, callback, userdata)
}

// Schedule the timer to go off after t_sec seconds.
func (self TimerEvent) Add(t_sec float64) { C.timerev_add(self.inner, C.double(t_sec)) }


// Unschedule the timer if it was scheduled. The timer could still be rescheduled
// by calling Add() afterwards.
func (self TimerEvent) Del() {
    self.ec.detach(rawptr_t(self.inner))
    C.timerev_del(self.inner)
}

// Empty the timer. It will be unscheduled and deallocated and its methods
// should never be called again.
func (self TimerEvent) Clear() {
    self.ec.detach(rawptr_t(self.inner))
    C.timerev_clear(self.inner)
}

// The C pointer type for a SigEvent.
type CSigEvent = *C.sigev_t
type sigEvent struct {
    inner CSigEvent
    ec EventContext
}

// The handle for a UNIX signal event.
type SigEvent = *sigEvent

type SigEventCallback = C.sigev_callback_t

var (
    SIGTERM = C.SIGTERM
    SIGINT = C.SIGINT
)

// Create a SigEvent handle attached to the given EventContext, with a
// registered callback.
func NewSigEvent(_ec EventContext, cb SigEventCallback, userdata rawptr_t) SigEvent {
    res := &sigEvent{
        inner: C.sigev_new(_ec.inner, cb, userdata),
        ec: _ec,
    }
    _ec.attach(rawptr_t(res.inner), res)
    runtime.SetFinalizer(res, func(self SigEvent) { self.free() })
    return res
}

func (self SigEvent) free() { C.sigev_free(self.inner) }

// Register the handle to listen for UNIX signal sig.
func (self SigEvent) Add(sig int) { C.sigev_add(self.inner, C.int(sig)) }


// Unregister the handle. The handle may be re-registered in the future.
func (self SigEvent) Del() {
    self.ec.detach(rawptr_t(self.inner))
    C.sigev_del(self.inner)
}

// Unregister the handle. Any methods of the handle should no longer be used
// and the handle will be deallocated.
func (self SigEvent) Clear() {
    self.ec.detach(rawptr_t(self.inner))
    C.sigev_clear(self.inner)
}

// The C pointer type for a MPSCQueue object.
type CMPSCQueue = *C.mpscqueue_t
type mpscQueue struct {
    inner CMPSCQueue
    ec EventContext
}

// The object for a Multi-Producer, Single-Consumer queue.
type MPSCQueue = *mpscQueue

// The C function pointer type which takes mpscqueue_t* (the C pointer to
// MPSCQueue) and void* (the unsafe pointer to any userdata) as parameter, and
// returns either true (partial read from the queue, so it should be scheduled
// again), or false (the queue is drained, e.g. TryDequeue returns false).
type MPSCQueueCallback = C.mpscqueue_callback_t

// Create a MPSCQueue object.
func NewMPSCQueue() MPSCQueue {
    res := &mpscQueue{ inner: C.mpscqueue_new(), ec: nil }
    runtime.SetFinalizer(res, func(self MPSCQueue) { self.free() })
    return res
}

func (self MPSCQueue) free() { C.mpscqueue_free(self.inner) }

// Register the callback invoked when there are new elements in the MPSC queue.
func (self MPSCQueue) RegHandler(_ec EventContext, callback MPSCQueueCallback, userdata rawptr_t) {
    C.mpscqueue_reg_handler(self.inner, _ec.inner, callback, userdata)
    self.ec = _ec
    _ec.attach(rawptr_t(self.inner), self)
}

// Unregister the callback.
func (self MPSCQueue) UnregHandler() {
    self.ec.detach(rawptr_t(self.inner))
    C.mpscqueue_unreg_handler(self.inner)
}

// Enqueue an element (thread-safe). It returns true if successful. If
// unbounded is true the queue is expanded when full (and this function will
// return true).
func (self MPSCQueue) Enqueue(elem rawptr_t, unbounded bool) bool {
    return bool(C.mpscqueue_enqueue(self.inner, elem, C.bool(unbounded)))
}

// Try to dequeue an element from the queue (should only be called from one
// thread). It returns true if successful.
func (self MPSCQueue) TryDequeue(elem *rawptr_t) bool {
    return bool(C.mpscqueue_try_dequeue(self.inner, elem))
}

// Set the initial capacity of the queue. This should only be called before the
// first dequeue/enqueue operation.
func (self MPSCQueue) SetCapacity(capacity int) {
    C.mpscqueue_set_capacity(self.inner, C.size_t(capacity))
}
