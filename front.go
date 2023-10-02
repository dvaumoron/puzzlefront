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

const cssHidden = "hide"

func registerValidation() {
	doc := js.Global().Get(document)
	buttons := doc.Call(querySelectorAll, "[validate-form]")
	if !buttons.Truthy() {
		return
	}
	size := buttons.Length()
	if size == 0 {
		return
	}

	for i := 0; i < size; i++ {
		if button := buttons.Index(i); button.Truthy() {
			formId := button.Get("validate-form").String()
			button.Set(onclick, js.FuncOf(func(this js.Value, args []js.Value) any {
				formToValidate := doc.Call(getElementById, formId)
				if !formToValidate.Truthy() {
					return nil
				}

				fields := formToValidate.Call(querySelectorAll, "[required-message]")
				if !fields.Truthy() {
					return nil
				}
				size := fields.Length()
				for i := 0; i < size; i++ {
					field := fields.Index(i)
					if !field.Truthy() {
						return nil
					}

					fieldValue := field.Get(value).String()
					if fieldValue == "" {
						alertKey(field.Get("required-message").String())
						return nil
					}

					if forbiddenValues := field.Get("forbidden-values").String(); forbiddenValues != "" {
						for _, forbiddenValue := range strings.Split(forbiddenValues, ",") {
							if strings.EqualFold(fieldValue, forbiddenValue) {
								alertKey(field.Get("required-message").String())
								return nil
							}
						}
					}
				}

				fields = formToValidate.Call(querySelectorAll, "[confirm-field]")
				if !fields.Truthy() {
					return nil
				}
				size = fields.Length()
				for i := 0; i < size; i++ {
					field := fields.Index(i)
					if !field.Truthy() {
						return nil
					}

					field2 := doc.Call(getElementById, field.Get("confirm-field").String())
					if !field2.Truthy() {
						return nil
					}

					if field.Get(value).String() != field2.Get(value).String() {
						alertKey(field.Get("confirm-message").String())
						return nil
					}
				}

				formToValidate.Call(submit)
				return nil
			}))
		}
	}
}

func loginRegisterAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	loginRegisterButtonClasses := doc.Call(getElementById, "loginRegisterButton").Get(classList)
	confirmPasswordBlockClasses := doc.Call(getElementById, "confirmPasswordBlock").Get(classList)
	loginRegisterButton2Classes := doc.Call(getElementById, "loginRegisterButton2").Get(classList)
	if !(loginRegisterButtonClasses.Truthy() && confirmPasswordBlockClasses.Truthy() && loginRegisterButton2Classes.Truthy()) {
		return nil
	}

	loginRegisterButtonClasses.Call(toggle, cssHidden)
	confirmPasswordBlockClasses.Call(toggle, cssHidden)
	loginRegisterButton2Classes.Call(toggle, cssHidden)
	return nil
}

func loginRegisterAction2(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	loginForm := doc.Call(getElementById, "loginForm")
	loginField := doc.Call(getElementById, "loginField")
	passwordField := doc.Call(getElementById, "passwordField")
	confirmPasswordField := doc.Call(getElementById, "confirmPasswordField")
	loginRegisterField := doc.Call(getElementById, "loginRegisterField")
	if !(loginForm.Truthy() && loginField.Truthy() && passwordField.Truthy() && confirmPasswordField.Truthy() && loginRegisterField.Truthy()) {
		return nil
	}

	if loginField.Get(value).String() == "" {
		alertKey("errorEmptyLoginMessage")
		return nil
	}

	if passwordField.Get(value).String() == "" {
		alertKey("errorEmptyPasswordMessage")
		return nil
	}

	if passwordField.Get(value).String() == confirmPasswordField.Get(value).String() {
		loginRegisterField.Set(value, true)
		loginForm.Call(submit)
	} else {
		alertKey("errorWrongConfimPasswordMessage")
	}
	return nil
}

func disablePublishPost(this js.Value, args []js.Value) any {
	publishPostButton := js.Global().Get(document).Call(getElementById, "publishPostButton")
	publishPostButton.Set(onclick, js.FuncOf(displayPublishErrorAction))
	return nil
}

func publishPostAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	publishPostForm := doc.Call(getElementById, "publishPostForm")
	postTitleField := doc.Call(getElementById, "postTitleField")
	postMarkdownField := doc.Call(getElementById, "postMarkdownField")
	if !(publishPostForm.Truthy() && postTitleField.Truthy() && postMarkdownField.Truthy()) {
		return nil
	}

	if postTitleField.Get(value).String() == "" {
		alertKey("errorEmptyPostTitleMessage")
		return nil
	}

	if postMarkdownField.Get(value).String() == "" {
		alertKey("errorEmptyPostContentMessage")
		return nil
	}

	target := publishPostForm.Get(action).String()
	publishPostForm.Set(action, convertBlogPreviewUrlToPublish(target))
	publishPostForm.Call(submit)
	return nil
}

func convertBlogPreviewUrlToPublish(url string) string {
	return url[:strings.LastIndexByte(url, '/')+1] + "save"
}

func displayPublishErrorAction(this js.Value, args []js.Value) any {
	alertKey("errorModifiedMarkdownMessage")
	return nil
}

func buildWikiLink(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return "/?error=ErrorTechnicalProblem"
	}

	wikiArg := args[0]
	langArg := args[1]
	title := args[2].String() // always set

	wiki, lang := extractWikiDataFromUrl(js.Global().Get(location).Get(href).String())

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

func displayPasswordHelpAction(this js.Value, args []js.Value) any {
	alertKey("passwordHelpMessage")
	return nil
}

func main() {
	registerValidation()

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

	postTitleField := doc.Call(getElementById, "postTitleField")
	postMarkdownField := doc.Call(getElementById, "postMarkdownField")
	publishPostButton := doc.Call(getElementById, "publishPostButton")
	if postTitleField.Truthy() && postMarkdownField.Truthy() && publishPostButton.Truthy() {
		postTitleField.Set(onchange, js.FuncOf(disablePublishPost))
		postMarkdownField.Set(onchange, js.FuncOf(disablePublishPost))
		publishPostButton.Set(onclick, js.FuncOf(publishPostAction))
	}

	global.Set("buildWikiLink", js.FuncOf(buildWikiLink))

	passwordHelp := doc.Call(getElementById, "passwordHelp")
	if passwordHelp.Truthy() {
		passwordHelp.Set(onclick, js.FuncOf(displayPasswordHelpAction))
	}

	// keep the program active to allow function call from HTML/JavaScript
	<-make(chan struct{})
}
