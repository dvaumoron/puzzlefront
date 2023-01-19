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

func loginRegisterAction(this js.Value, args []js.Value) any {
	jsDoc := js.Global().Get("document")
	if jsDoc.Truthy() {
		loginForm := jsDoc.Call("getElementById", "loginForm")
		if loginForm.Truthy() {
			loginRegisterField := jsDoc.Call("getElementById", "loginRegisterField")
			if loginRegisterField.Truthy() {
				loginRegisterField.Set("value", "true")
				loginForm.Call("submit")
			}
		}
	}
	return nil
}

func main() {
	jsDoc := js.Global().Get("document")
	if jsDoc.Truthy() {
		loginRegisterButton := jsDoc.Call("getElementById", "loginRegisterButton")
		if loginRegisterButton.Truthy() {
			loginRegisterButton.Set("onclick", js.FuncOf(loginRegisterAction))
		}
	}

	// keep the program active to allow function call from HTML/JavaScript
	<-make(chan struct{})
}
