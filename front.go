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

	fh "github.com/dvaumoron/puzzlefront/helper"
)

const cssHidden = "hide"

func loginRegisterAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(fh.Document)
	loginRegisterButtonClasses := doc.Call(fh.GetElementById, "loginRegisterButton").Get(fh.ClassList)
	confirmPasswordBlockClasses := doc.Call(fh.GetElementById, "confirmPasswordBlock").Get(fh.ClassList)
	loginRegisterButton2Classes := doc.Call(fh.GetElementById, "loginRegisterButton2").Get(fh.ClassList)
	if loginRegisterButtonClasses.Truthy() && confirmPasswordBlockClasses.Truthy() && loginRegisterButton2Classes.Truthy() {
		loginRegisterButtonClasses.Call(fh.Toggle, cssHidden)
		confirmPasswordBlockClasses.Call(fh.Toggle, cssHidden)
		loginRegisterButton2Classes.Call(fh.Toggle, cssHidden)
	}
	return nil
}

func loginRegisterAction2(this js.Value, args []js.Value) any {
	doc := js.Global().Get(fh.Document)
	loginForm := doc.Call(fh.GetElementById, "loginForm")
	loginField := doc.Call(fh.GetElementById, "loginField")
	passwordField := doc.Call(fh.GetElementById, "passwordField")
	confirmPasswordField := doc.Call(fh.GetElementById, "confirmPasswordField")
	loginRegisterField := doc.Call(fh.GetElementById, "loginRegisterField")
	if !(loginForm.Truthy() && loginField.Truthy() && passwordField.Truthy() && confirmPasswordField.Truthy() && loginRegisterField.Truthy()) {
		return nil
	}

	if loginField.Get(fh.Value).String() == "" {
		fh.AlertKey("errorEmptyLoginMessage")
		return nil
	}

	if passwordField.Get(fh.Value).String() == "" {
		fh.AlertKey("errorEmptyPasswordMessage")
		return nil
	}

	if passwordField.Get(fh.Value).String() == confirmPasswordField.Get(fh.Value).String() {
		loginRegisterField.Set(fh.Value, true)
		loginForm.Call(fh.Submit)
	} else {
		fh.AlertKey("errorWrongConfimPasswordMessage")
	}
	return nil
}

func disablePublishPost(this js.Value, args []js.Value) any {
	publishPostButton := js.Global().Get(fh.Document).Call(fh.GetElementById, "publishPostButton")
	publishPostButton.Set(fh.Onclick, js.FuncOf(displayPublishErrorAction))
	return nil
}

func publishPostAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(fh.Document)
	publishPostForm := doc.Call(fh.GetElementById, "publishPostForm")
	postTitleField := doc.Call(fh.GetElementById, "postTitleField")
	postMarkdownField := doc.Call(fh.GetElementById, "postMarkdownField")
	if !(publishPostForm.Truthy() && postTitleField.Truthy() && postMarkdownField.Truthy()) {
		return nil
	}

	if postTitleField.Get(fh.Value).String() == "" {
		fh.AlertKey("errorEmptyPostTitleMessage")
		return nil
	}

	if postMarkdownField.Get(fh.Value).String() == "" {
		fh.AlertKey("errorEmptyPostContentMessage")
		return nil
	}

	target := publishPostForm.Get(fh.Action).String()
	publishPostForm.Set(fh.Action, convertBlogPreviewUrlToPublish(target))
	publishPostForm.Call(fh.Submit)
	return nil
}

func convertBlogPreviewUrlToPublish(url string) string {
	return url[:strings.LastIndexByte(url, '/')+1] + "save"
}

func displayPublishErrorAction(this js.Value, args []js.Value) any {
	fh.AlertKey("errorModifiedMarkdownMessage")
	return nil
}

func buildWikiLink(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return "/?error=ErrorTechnicalProblem"
	}

	wikiArg := args[0]
	langArg := args[1]
	title := args[2].String() // always set

	wiki, lang := extractWikiDataFromUrl(js.Global().Get(fh.Location).Get(fh.Href).String())

	if wikiArg.Truthy() {
		wiki = wikiArg.String()
		if wiki[len(wiki)-1] != '/' {
			wiki += "/"
		}
	}

	if langArg.Truthy() {
		lang = langArg.String()
	}

	var linkBuilder strings.Builder
	linkBuilder.WriteString(wiki)
	linkBuilder.WriteString(lang)
	linkBuilder.WriteString("/view/")
	linkBuilder.WriteString(title)
	return linkBuilder.String()
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
	fh.RegisterValidationRules()
	fh.RegisterDisplayMessageAction()

	global := js.Global()
	doc := global.Get(fh.Document)

	fh.TruthyOnclick(doc.Call(fh.GetElementById, "loginRegisterButton"), loginRegisterAction)
	fh.TruthyOnclick(doc.Call(fh.GetElementById, "loginRegisterButton2"), loginRegisterAction2)

	postTitleField := doc.Call(fh.GetElementById, "postTitleField")
	postMarkdownField := doc.Call(fh.GetElementById, "postMarkdownField")
	publishPostButton := doc.Call(fh.GetElementById, "publishPostButton")
	if postTitleField.Truthy() && postMarkdownField.Truthy() && publishPostButton.Truthy() {
		postTitleField.Set(fh.Onchange, js.FuncOf(disablePublishPost))
		postMarkdownField.Set(fh.Onchange, js.FuncOf(disablePublishPost))
		publishPostButton.Set(fh.Onclick, js.FuncOf(publishPostAction))
	}

	global.Set("buildWikiLink", js.FuncOf(buildWikiLink))

	fh.KeepRunning()
}
