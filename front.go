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

const cssHidden = "hidden"

func loginRegisterAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	loginRegisterButton := doc.Call(getElementById, "loginRegisterButton")
	if !loginRegisterButton.Truthy() {
		return nil
	}

	loginRegisterButton.Set(visible, false)

	confirmPasswordBlockClasses := doc.Call(getElementById, "confirmPasswordBlock").Get(classList)
	if !confirmPasswordBlockClasses.Truthy() {
		return nil
	}

	confirmPasswordBlockClasses.Call(toggle, cssHidden)

	loginRegisterButton2Classes := doc.Call(getElementById, "loginRegisterButton2").Get(classList)
	if !loginRegisterButton2Classes.Truthy() {
		return nil
	}

	loginRegisterButton2Classes.Call(toggle, cssHidden)

	return nil
}

func loginRegisterAction2(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	loginForm := doc.Call(getElementById, "loginForm")
	if !loginForm.Truthy() {
		return nil
	}

	passwordField := doc.Call(getElementById, "passwordField")
	if !passwordField.Truthy() {
		return nil
	}

	confirmPasswordField := doc.Call(getElementById, "confirmPasswordField")
	if !confirmPasswordField.Truthy() {
		return nil
	}

	loginRegisterField := doc.Call(getElementById, "loginRegisterField")
	if loginRegisterField.Truthy() {
		if passwordField.Get(value).String() == confirmPasswordField.Get(value).String() {
			loginRegisterField.Set(value, true)
			loginForm.Call(submit)
		} else {
			errorMessageSpan := doc.Call(getElementById, "errorWrongConfimPasswordMessage")
			if errorMessageSpan.Truthy() {
				alert(errorMessageSpan.Get(textContent).String())
			}
		}
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

func disablePublishPost(this js.Value, args []js.Value) any {
	publishPostButton := js.Global().Get(document).Call(getElementById, "publishPostButton")
	publishPostButton.Set(onclick, js.FuncOf(displayPublishErrorAction))
	return nil
}

func publishPostAction(this js.Value, args []js.Value) any {
	publishPostForm := js.Global().Get(document).Call(getElementById, "publishPostForm")
	if publishPostForm.Truthy() {
		target := publishPostForm.Call(getAttribute, action).String()
		publishPostForm.Set(action, convertBlogPreviewUrlToPublish(target))
		publishPostForm.Call(submit)
	}
	return nil
}

func convertBlogPreviewUrlToPublish(url string) string {
	return url[:strings.LastIndexByte(url, '/')+1] + "save"
}

func displayPublishErrorAction(this js.Value, args []js.Value) any {
	errorMessageSpan := js.Global().Get(document).Call(getElementById, "errorModifiedMarkdownMessage")
	if errorMessageSpan.Truthy() {
		alert(errorMessageSpan.Get(textContent).String())
	}
	return nil
}

func wikiLinkConstructor(this js.Value, args []js.Value) any {
	global := js.Global()
	doc := global.Get(document)

	wikiAttr := this.Call(getAttribute, "wiki")
	langAttr := this.Call(getAttribute, "lang")
	title := this.Call(getAttribute, "title").String() // always set

	wiki, lang := extractWikiDataFromUrl(global.Get(location).Get(href).String())

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

func extractWikiDataFromUrl(url string) (string, string) {
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

	loginRegisterButton2 := doc.Call(getElementById, "loginRegisterButton2")
	if loginRegisterButton2.Truthy() {
		loginRegisterButton2.Set(onclick, js.FuncOf(loginRegisterAction2))
	}

	saveRoleButton := doc.Call(getElementById, "saveRoleButton")
	if saveRoleButton.Truthy() {
		saveRoleButton.Set(onclick, js.FuncOf(saveRoleAction))
	}

	postTitleField := doc.Call(getElementById, "postTitleField")
	postMarkdownField := doc.Call(getElementById, "postMarkdownField")
	publishPostButton := doc.Call(getElementById, "publishPostButton")
	if postTitleField.Truthy() && postMarkdownField.Truthy() && publishPostButton.Truthy() {
		postTitleField.Set(onchange, js.FuncOf(disablePublishPost))
		postMarkdownField.Set(onchange, js.FuncOf(disablePublishPost))
		publishPostButton.Set(onclick, js.FuncOf(publishPostAction))
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
