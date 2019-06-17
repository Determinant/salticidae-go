package salticidae

// #include "salticidae/event.h"
// #include <signal.h>
import "C"
import "runtime"

type CEventContext = *C.eventcontext_t
type eventContext struct {
    inner CEventContext
    attached map[uintptr]interface{}
}
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
func (self EventContext) Dispatch() { C.eventcontext_dispatch(self.inner) }
func (self EventContext) Stop() { C.eventcontext_stop(self.inner) }

type CThreadCall = *C.threadcall_t
type threadCall struct { inner CThreadCall }
type ThreadCall = *threadCall

type ThreadCallCallback = C.threadcall_callback_t

func NewThreadCall(ec EventContext) ThreadCall {
    res := &threadCall{ inner: C.threadcall_new(ec.inner) }
    runtime.SetFinalizer(res, func(self ThreadCall) { self.free() })
    return res
}

func (self ThreadCall) free() { C.threadcall_free(self.inner) }

func (self ThreadCall) AsyncCall(callback ThreadCallCallback, userdata rawptr_t) {
    C.threadcall_async_call(self.inner, callback, userdata)
}

type CTimerEvent = *C.timerev_t
type timerEvent struct {
    inner CTimerEvent
    ec EventContext
}
type TimerEvent = *timerEvent

type TimerEventCallback = C.timerev_callback_t

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
func (self TimerEvent) SetCallback(callback TimerEventCallback, userdata rawptr_t) {
    C.timerev_set_callback(self.inner, callback, userdata)
}

func (self TimerEvent) Add(t_sec float64) { C.timerev_add(self.inner, C.double(t_sec)) }
func (self TimerEvent) Del() {
    self.ec.detach(rawptr_t(self.inner))
    C.timerev_del(self.inner)
}

func (self TimerEvent) Clear() {
    self.ec.detach(rawptr_t(self.inner))
    C.timerev_clear(self.inner)
}

type CSigEvent = *C.sigev_t
type sigEvent struct {
    inner CSigEvent
    ec EventContext
}
type SigEvent = *sigEvent

type SigEventCallback = C.sigev_callback_t

var (
    SIGTERM = C.SIGTERM
    SIGINT = C.SIGINT
)

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
func (self SigEvent) Add(sig int) { C.sigev_add(self.inner, C.int(sig)) }
func (self SigEvent) Del() {
    self.ec.detach(rawptr_t(self.inner))
    C.sigev_del(self.inner)
}

func (self SigEvent) Clear() {
    self.ec.detach(rawptr_t(self.inner))
    C.sigev_clear(self.inner)
}
