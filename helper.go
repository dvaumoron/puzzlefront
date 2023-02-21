/*
 *
 * Copyright 2023 puzzlefront authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package main

import "syscall/js"

const document = "document"
const getElementById = "getElementById"
const getAttribute = "getAttribute"
const toggle = "toggle"
const createElement = "createElement"
const appendChild = "appendChild"
const classList = "classList"
const textContent = "textContent"
const location = "location"
const action = "action"
const href = "href"
const onclick = "onclick"
const onchange = "onchange"
const submit = "submit"
const value = "value"
const visible = "visible"

const Object = "Object"

const getPrototypeOf = "getPrototypeOf"
const assign = "assign"
const define = "define"

var openMode = jsObject{"mode": "open"}.toJs()

type jsObject map[string]any

func (o jsObject) toJs() js.Value {
	res := js.Global().Get(Object).New()
	for k, v := range o {
		var converted any
		switch casted := v.(type) {
		case jsObject:
			converted = casted.toJs()
		case jsFunc:
			converted = casted.toJs()
		case jsClass:
			converted = casted.toJs()
		default:
			converted = casted
		}
		res.Set(k, converted)
	}
	return res
}

type jsFunc func(js.Value, []js.Value) any

func (f jsFunc) toJs() js.Func {
	return js.FuncOf(f)
}

type jsClass struct {
	parent      js.Value
	constructor jsFunc // without super call
	content     jsObject
}

func (c jsClass) toJs() js.Func {
	object := js.Global().Get(Object)

	var res js.Func
	var resProto js.Value
	if c.parent.Truthy() {
		res = js.FuncOf(func(this js.Value, args []js.Value) any {
			c.parent.Call("apply", this, args) // emulate super call
			c.constructor(this, args)
			return this
		})

		parentProto := object.Call(getPrototypeOf, c.parent)

		properties := jsObject{"constructor": jsObject{
			"value":        res,
			"writable":     true,
			"configurable": true,
		}}
		resProto = object.Call("setPrototypeOf", object.Call("create", parentProto, properties.toJs()))
	} else {
		res = c.constructor.toJs()

		resProto = object.Call(getPrototypeOf, res)
	}

	if len(c.content) != 0 {
		object.Call(assign, resProto, c.content.toJs())
	}

	return res
}

func alert(message string) {
	jsAlert := js.Global().Get("alert")
	if jsAlert.Truthy() {
		jsAlert.Invoke(message)
	}
}

func attachShadow(this js.Value) js.Value {
	return this.Call("attachShadow", openMode)
}
