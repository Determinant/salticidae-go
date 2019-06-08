package salticidae

// #include "salticidae/event.h"
// #include <signal.h>
import "C"

type EventContext = *C.struct_eventcontext_t

func NewEventContext() EventContext { return C.eventcontext_new() }
func (self EventContext) Free() { C.eventcontext_free(self) }
func (self EventContext) Dispatch() { C.eventcontext_dispatch(self) }
func (self EventContext) Stop() { C.eventcontext_stop(self) }

type SigEvent = *C.sigev_t
type SigEventCallback = C.sigev_callback_t
var SIGTERM = C.SIGTERM
var SIGINT = C.SIGINT

func NewSigEvent(ec EventContext, cb SigEventCallback) SigEvent {
    return C.sigev_new(ec, cb)
}

func (self SigEvent) Add(sig int) { C.sigev_add(self, C.int(sig)) }
func (self SigEvent) Free() { C.sigev_free(self) }
