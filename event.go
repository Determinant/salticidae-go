package salticidae

// #include "salticidae/event.h"
// #include <signal.h>
import "C"

type EventContext = *C.struct_eventcontext_t

func NewEventContext() EventContext { return C.eventcontext_new() }
func (self EventContext) Free() { C.eventcontext_free(self) }
func (self EventContext) Dispatch() { C.eventcontext_dispatch(self) }
func (self EventContext) Stop() { C.eventcontext_stop(self) }

type ThreadCall = *C.struct_threadcall_t
type ThreadCallCallback = C.threadcall_callback_t

func NewThreadCall(ec EventContext) ThreadCall { return C.threadcall_new(ec) }

func (self ThreadCall) Free() { C.threadcall_free(self) }

func (self ThreadCall) AsyncCall(callback ThreadCallCallback, userdata rawptr_t) {
    C.threadcall_async_call(self, callback, userdata)
}

func (self ThreadCall) GetEC() EventContext { return C.threadcall_get_ec(self) }

type TimerEvent = *C.timerev_t
type TimerEventCallback = C.timerev_callback_t

func NewTimerEvent(ec EventContext, cb TimerEventCallback, userdata rawptr_t) TimerEvent {
    return C.timerev_new(ec, cb, userdata)
}

func (self TimerEvent) Free() { C.timerev_free(self) }
func (self TimerEvent) SetCallback(callback TimerEventCallback, userdata rawptr_t) {
    C.timerev_set_callback(self, callback, userdata)
}

func (self TimerEvent) Add(t_sec float64) { C.timerev_add(self, C.double(t_sec)) }
func (self TimerEvent) Clear() { C.timerev_clear(self) }

type SigEvent = *C.sigev_t
type SigEventCallback = C.sigev_callback_t
var SIGTERM = C.SIGTERM
var SIGINT = C.SIGINT

func NewSigEvent(ec EventContext, cb SigEventCallback, userdata rawptr_t) SigEvent {
    return C.sigev_new(ec, cb, userdata)
}

func (self SigEvent) Add(sig int) { C.sigev_add(self, C.int(sig)) }
func (self SigEvent) Free() { C.sigev_free(self) }
