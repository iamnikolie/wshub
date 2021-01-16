package wshub

import "reflect"

func (h *WSHub) HandleMessage(fnc func(msg *[]byte)) {
	exists := false
	for _, exfnx := range h.handlers {
		exists = exists || reflect.ValueOf(exfnx) == reflect.ValueOf(fnc)
	}
	if !exists {
		h.handlers = append(h.handlers, fnc)
	}
}
