// Copyright (C) 2024 Matt Dunleavy. All rights reserved.
// Use of this source code is subject to the MIT license 
// that can be found in the LICENSE file.

package pluginmanager

import (
    "sync"
)

type Event interface {
    Name() string
}

type PluginLoadedEvent struct {
    PluginName string
}

func (e PluginLoadedEvent) Name() string {
    return "PluginLoaded"
}

type PluginUnloadedEvent struct {
    PluginName string
}

func (e PluginUnloadedEvent) Name() string {
    return "PluginUnloaded"
}

type PluginHotReloadedEvent struct {
    PluginName string
}

func (e PluginHotReloadedEvent) Name() string {
    return "PluginHotReloaded"
}

type EventHandler func(Event)

type EventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[string][]EventHandler),
    }
}

func (eb *EventBus) Subscribe(eventName string, handler EventHandler) {
    eb.mu.Lock()
    defer eb.mu.Unlock()
    eb.handlers[eventName] = append(eb.handlers[eventName], handler)
}

func (eb *EventBus) Publish(event Event) {
    eb.mu.RLock()
    defer eb.mu.RUnlock()
    for _, handler := range eb.handlers[event.Name()] {
        go handler(event)
    }
}