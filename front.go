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

import (
	"strings"
	"syscall/js"
)

const document = "document"
const getElementById = "getElementById"
const innerText = "innerText"
const onclick = "onclick"
const submit = "submit"
const value = "value"

func alert(message string) {
	jsAlert := js.Global().Get("alert")
	if jsAlert.Truthy() {
		jsAlert.Invoke(message)
	}
}

func loginRegisterAction(this js.Value, args []js.Value) any {
	jsDoc := js.Global().Get(document)
	if !jsDoc.Truthy() {
		return nil
	}
	loginForm := jsDoc.Call(getElementById, "loginForm")
	if !loginForm.Truthy() {
		return nil
	}
	loginRegisterField := jsDoc.Call(getElementById, "loginRegisterField")
	if loginRegisterField.Truthy() {
		loginRegisterField.Set(value, "true")
		loginForm.Call(submit)
	}
	return nil
}

func saveRoleAction(this js.Value, args []js.Value) any {
	jsDoc := js.Global().Get(document)
	if !jsDoc.Truthy() {
		return nil
	}
	editRoleForm := jsDoc.Call(getElementById, "editRoleForm")
	if !editRoleForm.Truthy() {
		return nil
	}
	editRoleNameField := jsDoc.Call(getElementById, "loginRegisterField")
	if !editRoleNameField.Truthy() {
		return nil
	}
	roleName := editRoleNameField.Get(value).String()
	if strings.EqualFold(roleName, "new") {
		errorMessageSpan := jsDoc.Call(getElementById, "errorBadRoleNameMessage")
		if errorMessageSpan.Truthy() {
			alert(errorMessageSpan.Get(innerText).String())
		}
	} else {
		editRoleForm.Call(submit)
	}
	return nil
}

func main() {
	jsDoc := js.Global().Get(document)
	if !jsDoc.Truthy() {
		return
	}

	loginRegisterButton := jsDoc.Call(getElementById, "loginRegisterButton")
	if loginRegisterButton.Truthy() {
		loginRegisterButton.Set(onclick, js.FuncOf(loginRegisterAction))
	}

	saveRoleButton := jsDoc.Call(getElementById, "saveRoleButton")
	if saveRoleButton.Truthy() {
		saveRoleButton.Set(onclick, js.FuncOf(saveRoleAction))
	}

	// keep the program active to allow function call from HTML/JavaScript
	<-make(chan struct{})
}
