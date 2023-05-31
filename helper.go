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
const toggle = "toggle"
const classList = "classList"
const textContent = "textContent"
const location = "location"
const action = "action"
const href = "href"
const onclick = "onclick"
const onchange = "onchange"
const submit = "submit"
const value = "value"

func alertKey(messageSpanId string) {
	global := js.Global()
	jsAlert := global.Get("alert")
	errorMessageSpan := global.Get(document).Call(getElementById, messageSpanId)
	if jsAlert.Truthy() && errorMessageSpan.Truthy() {
		jsAlert.Invoke(errorMessageSpan.Get(textContent).String())
	}
}
