package salticidae

// #include "salticidae/event.h"
// #include <signal.h>
import "C"
import "runtime"

//// begin EventContext def

// CEventContext is the C pointer type for an EventContext handle.
type CEventContext = *C.eventcontext_t
type eventContext struct {
	inner    CEventContext
	attached map[uintptr]interface{}
	autoFree bool
}

// EventContext is the handle for an event loop.
type EventContext = *eventContext

func eventContextSetFinalizer(res EventContext, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self EventContext) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (ec EventContext) Free() {
	C.eventcontext_free(ec.inner)
	if ec.autoFree {
		runtime.SetFinalizer(ec, nil)
	}
}

func (ec EventContext) attach(ptr RawPtr, x interface{}) { ec.attached[uintptr(ptr)] = x }
func (ec EventContext) detach(ptr RawPtr)                { delete(ec.attached, uintptr(ptr)) }

//// end EventContext def

//// begin ThreadCall def

// CThreadCall is the C pointer type for a ThreadCall handle.
type CThreadCall = *C.threadcall_t
type threadCall struct {
	inner    CThreadCall
	ec       EventContext
	autoFree bool
}

// ThreadCall is the handle for scheduling a function call executed by a
// particular EventContext event loop.
type ThreadCall = *threadCall

// ThreadCallCallback is the C function pointer type which takes
// threadcall_handle_t* and void* (passing in the custom user data allocated by
// C.malloc) as parameters.
type ThreadCallCallback = C.threadcall_callback_t

func threadCallSetFinalizer(res ThreadCall, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self ThreadCall) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (tc ThreadCall) Free() {
	tc.ec.detach(RawPtr(tc.inner))
	C.threadcall_free(tc.inner)
	if tc.autoFree {
		runtime.SetFinalizer(tc, nil)
	}
}

//// end ThreadCall def

//// begin TimerEvent def

// CTimerEvent is the C pointer type for TimerEvent handle.
type CTimerEvent = *C.timerev_t
type timerEvent struct {
	inner    CTimerEvent
	ec       EventContext
	autoFree bool
}

// TimerEvent is the handle for a timed event.
type TimerEvent = *timerEvent

// TimerEventCallback is the C function pointer type which takes timerev_t*
// (the C pointer to TimerEvent) and void* (the unsafe pointer to any userdata)
// as parameter.
type TimerEventCallback = C.timerev_callback_t

func timerEventSetFinalizer(res TimerEvent, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self TimerEvent) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (te TimerEvent) Free() {
	te.Clear()
	C.timerev_free(te.inner)
	if te.autoFree {
		runtime.SetFinalizer(te, nil)
	}
}

//// end TimerEvent def

//// begin SigEvent def

// CSigEvent is the C pointer type for a SigEvent.
type CSigEvent = *C.sigev_t
type sigEvent struct {
	inner    CSigEvent
	ec       EventContext
	autoFree bool
}

// SigEvent is the handle for a UNIX signal event.
type SigEvent = *sigEvent

// SigEventCallback is the callback function.
type SigEventCallback = C.sigev_callback_t

var (
	SIGTERM = C.SIGTERM
	SIGINT  = C.SIGINT
)

func sigEventSetFinalizer(res SigEvent, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self SigEvent) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (se SigEvent) Free() {
	se.Clear()
	C.sigev_free(se.inner)
	if se.autoFree {
		runtime.SetFinalizer(se, nil)
	}
}

//// end SigEvent def

//// begin MPSCQueue def

// CMPSCQueue is the C pointer type for a MPSCQueue object.
type CMPSCQueue = *C.mpscqueue_t
type mpscQueue struct {
	inner    CMPSCQueue
	ec       EventContext
	autoFree bool
}

// MPSCQueue is a Multi-Producer, Single-Consumer queue.
type MPSCQueue = *mpscQueue

// MPSCQueueCallback is the C function pointer type which takes mpscqueue_t*
// (the C pointer to MPSCQueue) and void* (the unsafe pointer to any userdata)
// as parameter, and returns either true (partial read from the queue, so it
// should be scheduled again), or false (the queue is drained, e.g. TryDequeue
// returns false).
type MPSCQueueCallback = C.mpscqueue_callback_t

func mpscQueueSetFinalizer(res MPSCQueue, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self MPSCQueue) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (q MPSCQueue) Free() {
	C.mpscqueue_free(q.inner)
	if q.autoFree {
		runtime.SetFinalizer(q, nil)
	}
}

//// end MPSCQueue def

//// begin EventContext methods

// NewEventContext creates an EventContext object.
func NewEventContext() EventContext {
	res := &eventContext{
		inner:    C.eventcontext_new(),
		attached: make(map[uintptr]interface{}),
	}
	eventContextSetFinalizer(res, true)
	return res
}

//// end MPSCQueue

//// end SigEvent def

// Dispatch starts the event loop. This is a blocking call that will hand over
// the control flow to the infinite loop which triggers callbacks upon new
// events.  The function will return when Stop() is called.
func (ec EventContext) Dispatch() { C.eventcontext_dispatch(ec.inner) }

// Stop the event loop. This function is typically called in a callback. Notice
// that all methods of an EventContext should be invoked by the same thread
// which logically owns the loop. To schedule code executed in the event loop
// of a particular thread, see ThreadCall.
func (ec EventContext) Stop() { C.eventcontext_stop(ec.inner) }

//// end EventContext methods

//// begin ThreadCall methods

// NewThreadCall creates a ThreadCall handle attached to the given
// EventContext. Each invokcation of AsyncCall() will schedule a function call
// executed in the given EventContext event loop.
func NewThreadCall(ec EventContext) ThreadCall {
	res := &threadCall{
		inner: C.threadcall_new(ec.inner),
		ec:    ec,
	}
	ec.attach(RawPtr(res.inner), res)
	threadCallSetFinalizer(res, true)
	return res
}

// AsyncCall schedules a function to be executed in the target event loop.
func (tc ThreadCall) AsyncCall(callback ThreadCallCallback, userdata RawPtr) {
	C.threadcall_async_call(tc.inner, callback, userdata)
}

//// end ThreadCall methods

//// begin TimerEvent methods

// NewTimerEvent creates a TimerEvent handle attached to the given
// EventContext, with a registered callback.
func NewTimerEvent(ec EventContext, cb TimerEventCallback, userdata RawPtr) TimerEvent {
	res := &timerEvent{
		inner: C.timerev_new(ec.inner, cb, userdata),
		ec:    ec,
	}
	ec.attach(RawPtr(res.inner), res)
	timerEventSetFinalizer(res, true)
	return res
}

// SetCallback changes the callback.
func (te TimerEvent) SetCallback(callback TimerEventCallback, userdata RawPtr) {
	C.timerev_set_callback(te.inner, callback, userdata)
}

// Add schedules the timer to go off after t_sec seconds.
func (te TimerEvent) Add(sec float64) { C.timerev_add(te.inner, C.double(sec)) }

// Del unschedules the timer if it was scheduled. The timer could still be
// rescheduled by calling Add() afterwards.
func (te TimerEvent) Del() {
	te.ec.detach(RawPtr(te.inner))
	C.timerev_del(te.inner)
}

// Clear the timer. It will be unscheduled and deallocated and its methods
// should never be called again.
func (te TimerEvent) Clear() {
	te.ec.detach(RawPtr(te.inner))
	C.timerev_clear(te.inner)
}

//// end TimerEvent methods

//// begin SigEvent methods

// NewSigEvent creates a SigEvent handle attached to the given EventContext,
// with a registered callback.
func NewSigEvent(ec EventContext, cb SigEventCallback, userdata RawPtr) SigEvent {
	res := &sigEvent{
		inner: C.sigev_new(ec.inner, cb, userdata),
		ec:    ec,
	}
	ec.attach(RawPtr(res.inner), res)
	sigEventSetFinalizer(res, true)
	return res
}

// Add registers the handle to listen for UNIX signal sig.
func (se SigEvent) Add(sig int) { C.sigev_add(se.inner, C.int(sig)) }

// Del unregisters the handle. The handle may be re-registered in the future.
func (se SigEvent) Del() {
	se.ec.detach(RawPtr(se.inner))
	C.sigev_del(se.inner)
}

// Clear the handle. Any methods of the handle should no longer be used
// and the handle will be deallocated.
func (se SigEvent) Clear() {
	se.ec.detach(RawPtr(se.inner))
	C.sigev_clear(se.inner)
}

//// end SigEvent methods

//// begin MPSCQueue methods

// NewMPSCQueue creates a MPSCQueue object.
func NewMPSCQueue() MPSCQueue {
	res := &mpscQueue{inner: C.mpscqueue_new(), ec: nil}
	mpscQueueSetFinalizer(res, true)
	return res
}

// RegHandler registers the callback invoked when there are new elements in the
// MPSC queue.
func (q MPSCQueue) RegHandler(ec EventContext, callback MPSCQueueCallback, userdata RawPtr) {
	C.mpscqueue_reg_handler(q.inner, ec.inner, callback, userdata)
	q.ec = ec
	ec.attach(RawPtr(q.inner), q)
}

// UnregHandler unregisters the callback.
func (q MPSCQueue) UnregHandler() {
	q.ec.detach(RawPtr(q.inner))
	C.mpscqueue_unreg_handler(q.inner)
}

// Enqueue an element (thread-safe). It returns true if successful. If
// unbounded is true the queue is expanded when full (and this function will
// return true).
func (q MPSCQueue) Enqueue(elem RawPtr, unbounded bool) bool {
	return bool(C.mpscqueue_enqueue(q.inner, elem, C.bool(unbounded)))
}

// TryDequeue tries to dequeue an element from the queue (should only be called
// from one thread). It returns true if successful.
func (q MPSCQueue) TryDequeue(elem *RawPtr) bool {
	return bool(C.mpscqueue_try_dequeue(q.inner, elem))
}

// SetCapacity sets the initial capacity of the queue. This should only be
// called before the first dequeue/enqueue operation.
func (q MPSCQueue) SetCapacity(capacity int) {
	C.mpscqueue_set_capacity(q.inner, C.size_t(capacity))
}

//// end MPSCQueue methods
