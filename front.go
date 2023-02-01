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

func loginRegisterAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	loginForm := doc.Call(getElementById, "loginForm")
	if !loginForm.Truthy() {
		return nil
	}

	loginRegisterField := doc.Call(getElementById, "loginRegisterField")
	if loginRegisterField.Truthy() {
		loginRegisterField.Set(value, true)
		loginForm.Call(submit)
	}
	return nil
}

func saveRoleAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	editRoleForm := doc.Call(getElementById, "editRoleForm")
	editRoleNameField := doc.Call(getElementById, "loginRegisterField")
	if !(editRoleForm.Truthy() && editRoleNameField.Truthy()) {
		return nil
	}

	roleName := editRoleNameField.Get(value).String()
	if strings.EqualFold(roleName, "new") {
		errorMessageSpan := doc.Call(getElementById, "errorBadRoleNameMessage")
		if errorMessageSpan.Truthy() {
			alert(errorMessageSpan.Get(textContent).String())
		}
	} else {
		editRoleForm.Call(submit)
	}
	return nil
}

func wikiLinkConstructor(this js.Value, args []js.Value) any {
	global := js.Global()
	doc := global.Get(document)

	wikiAttr := this.Call(getAttribute, "wiki")
	langAttr := this.Call(getAttribute, "lang")
	title := this.Call(getAttribute, "title").String() // always set

	wiki, lang := extractUrlData(global.Get("location").Get(href).String())

	if wikiAttr.Truthy() {
		wiki = wikiAttr.String()
		if wiki[len(wiki)-1] != '/' {
			wiki += "/"
		}
	}

	if langAttr.Truthy() {
		lang = langAttr.String()
	}

	var linkBuilder strings.Builder
	linkBuilder.WriteString(wiki)
	linkBuilder.WriteString(lang)
	linkBuilder.WriteString("/view/")
	linkBuilder.WriteString(title)

	shadow := attachShadow(this)
	linkElem := doc.Call(createElement, "a")
	linkElem.Set(href, linkBuilder.String())
	linkElem.Set(textContent, this.Get(textContent))
	shadow.Call(appendChild, linkElem)

	// set the computed value on the Javascript Object
	this.Set("wiki", wiki)
	this.Set("lang", lang)
	this.Set("title", title)
	return this
}

func extractUrlData(url string) (string, string) {
	start := 0
	end := 0
	count := 0
	index := len(url) - 2
	for ; ; index-- {
		if url[index] == '/' {
			count++
			if count == 2 {
				end = index
			} else if count == 3 {
				start = index + 1
				break
			}
		}
	}
	return url[:start], url[start:end]
}

func main() {
	global := js.Global()
	doc := global.Get(document)

	loginRegisterButton := doc.Call(getElementById, "loginRegisterButton")
	if loginRegisterButton.Truthy() {
		loginRegisterButton.Set(onclick, js.FuncOf(loginRegisterAction))
	}

	saveRoleButton := doc.Call(getElementById, "saveRoleButton")
	if saveRoleButton.Truthy() {
		saveRoleButton.Set(onclick, js.FuncOf(saveRoleAction))
	}

	htmlElement := global.Get("HTMLElement")

	wikiLink := jsClass{
		parent:      htmlElement,
		constructor: wikiLinkConstructor,
		content:     jsObject{},
	}

	customElem := global.Get("customElements")
	customElem.Call(define, "wiki-link", wikiLink.toJs())

	// keep the program active to allow function call from HTML/JavaScript
	<-make(chan struct{})
}
