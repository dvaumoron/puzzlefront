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

func loginSubmitAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	loginForm := doc.Call(getElementById, "loginForm")
	loginField := doc.Call(getElementById, "loginField")
	passwordField := doc.Call(getElementById, "passwordField")
	if !(loginForm.Truthy() && loginField.Truthy() && passwordField.Truthy()) {
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

	loginForm.Call(submit)
	return nil
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

func saveRoleAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	editRoleForm := doc.Call(getElementById, "editRoleForm")
	editRoleNameField := doc.Call(getElementById, "editRoleNameField")
	if !(editRoleForm.Truthy() && editRoleNameField.Truthy()) {
		return nil
	}

	roleName := editRoleNameField.Get(value).String()
	if roleName == "" || strings.EqualFold(roleName, "new") {
		alertKey("errorBadRoleNameMessage")
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

func previewPostAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	previewPostForm := doc.Call(getElementById, "previewPostForm")
	postTitleField := doc.Call(getElementById, "postTitleField")
	postMarkdownField := doc.Call(getElementById, "postMarkdownField")
	if !(previewPostForm.Truthy() && postTitleField.Truthy() && postMarkdownField.Truthy()) {
		return nil
	}

	if postTitleField.Get(value).String() == "" {
		alertKey("errorEmptyPostTitleMessage")
		return nil
	}

	markdown := postMarkdownField.Get(value).String()
	if markdown == "" {
		alertKey("errorEmptyPostContentMessage")
		return nil
	}

	defaultCommentSpan := doc.Call(getElementById, "unmodifiedMessage")
	if defaultCommentSpan.Truthy() && defaultCommentSpan.Get(textContent).String() == markdown {
		alertKey("errorEmptyPostContentMessage")
		return nil
	}

	previewPostForm.Call(submit)
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

func commentAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	commentForm := doc.Call(getElementById, "commentForm")
	commentField := doc.Call(getElementById, "commentField")
	if !(commentForm.Truthy() && commentField.Truthy()) {
		return nil
	}

	comment := commentField.Get(value).String()
	if comment == "" {
		alertKey("errorEmptyCommentMessage")
		return nil
	}

	defaultCommentSpan := doc.Call(getElementById, "unmodifiedComment")
	if defaultCommentSpan.Truthy() && defaultCommentSpan.Get(textContent).String() == comment {
		alertKey("errorEmptyCommentMessage")
		return nil
	}

	commentForm.Call(submit)
	return nil
}

func changeLoginAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	changeLoginForm := doc.Call(getElementById, "changeLoginForm")
	loginField := doc.Call(getElementById, "loginField")
	passwordField := doc.Call(getElementById, "passwordField")
	if !(changeLoginForm.Truthy() && loginField.Truthy() && passwordField.Truthy()) {
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

	changeLoginForm.Call(submit)
	return nil
}

func changePasswordAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	changePasswordForm := doc.Call(getElementById, "changePasswordForm")
	passwordField := doc.Call(getElementById, "changePasswordField")
	confirmPasswordField := doc.Call(getElementById, "confirmChangePasswordField")
	if !(changePasswordForm.Truthy() && passwordField.Truthy() && confirmPasswordField.Truthy()) {
		return nil
	}

	if passwordField.Get(value).String() == "" {
		alertKey("errorEmptyPasswordMessage")
		return nil
	}

	if passwordField.Get(value).String() == confirmPasswordField.Get(value).String() {
		changePasswordForm.Call(submit)
	} else {
		alertKey("errorWrongConfimPasswordMessage")
	}
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

func createThreadAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	createThreadForm := doc.Call(getElementById, "createThreadForm")
	threadTitleField := doc.Call(getElementById, "threadTitleField")
	forumMessageField := doc.Call(getElementById, "forumMessageField")
	if !(createThreadForm.Truthy() && threadTitleField.Truthy() && forumMessageField.Truthy()) {
		return nil
	}

	if threadTitleField.Get(value).String() == "" {
		alertKey("errorEmptyThreadTitleMessage")
		return nil
	}

	if forumMessageField.Get(value).String() == "" {
		alertKey("errorEmptyForumMessageMessage")
		return nil
	}

	createThreadForm.Call(submit)
	return nil
}

func createForumMessageAction(this js.Value, args []js.Value) any {
	doc := js.Global().Get(document)
	createForumMessageForm := doc.Call(getElementById, "createForumMessageForm")
	forumMessageField := doc.Call(getElementById, "forumMessageField")
	if !(createForumMessageForm.Truthy() && forumMessageField.Truthy()) {
		return nil
	}

	message := forumMessageField.Get(value).String()
	if message == "" {
		alertKey("errorEmptyForumMessageMessage")
		return nil
	}

	defaultCommentSpan := doc.Call(getElementById, "unmodifiedMessage")
	if defaultCommentSpan.Truthy() && defaultCommentSpan.Get(textContent).String() == message {
		alertKey("errorEmptyForumMessageMessage")
		return nil
	}

	createForumMessageForm.Call(submit)
	return nil
}

func main() {
	global := js.Global()
	doc := global.Get(document)

	loginSubmitButton := doc.Call(getElementById, "loginSubmitButton")
	if loginSubmitButton.Truthy() {
		loginSubmitButton.Set(onclick, js.FuncOf(loginSubmitAction))
	}

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

	previewPostButton := doc.Call(getElementById, "previewPostButton")
	if previewPostButton.Truthy() {
		previewPostButton.Set(onclick, js.FuncOf(previewPostAction))
	}

	postTitleField := doc.Call(getElementById, "postTitleField")
	postMarkdownField := doc.Call(getElementById, "postMarkdownField")
	publishPostButton := doc.Call(getElementById, "publishPostButton")
	if postTitleField.Truthy() && postMarkdownField.Truthy() && publishPostButton.Truthy() {
		postTitleField.Set(onchange, js.FuncOf(disablePublishPost))
		postMarkdownField.Set(onchange, js.FuncOf(disablePublishPost))
		publishPostButton.Set(onclick, js.FuncOf(publishPostAction))
	}

	commentButton := doc.Call(getElementById, "commentButton")
	if commentButton.Truthy() {
		commentButton.Set(onclick, js.FuncOf(commentAction))
	}

	changeLoginButton := doc.Call(getElementById, "changeLoginButton")
	if changeLoginButton.Truthy() {
		changeLoginButton.Set(onclick, js.FuncOf(changeLoginAction))
	}

	changePasswordButton := doc.Call(getElementById, "changePasswordButton")
	if changePasswordButton.Truthy() {
		changePasswordButton.Set(onclick, js.FuncOf(changePasswordAction))
	}

	global.Set("buildWikiLink", js.FuncOf(buildWikiLink))

	createThreadButton := doc.Call(getElementById, "createThreadButton")
	if createThreadButton.Truthy() {
		createThreadButton.Set(onclick, js.FuncOf(createThreadAction))
	}

	createForumMessageButton := doc.Call(getElementById, "createForumMessageButton")
	if createForumMessageButton.Truthy() {
		createForumMessageButton.Set(onclick, js.FuncOf(createForumMessageAction))
	}

	// keep the program active to allow function call from HTML/JavaScript
	<-make(chan struct{})
}
