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

package fronthelper

import (
	"strings"
	"syscall/js"
)

const (
	Action           = "action"
	ClassList        = "classList"
	Document         = "document"
	GetAttribute     = "getAttribute"
	GetElementById   = "getElementById"
	Href             = "href"
	Location         = "location"
	Onchange         = "onchange"
	Onclick          = "onclick"
	QuerySelectorAll = "querySelectorAll"
	Submit           = "submit"
	TextContent      = "textContent"
	Toggle           = "toggle"
	Value            = "value"
)

const (
	validateFormAttrName        = "validate-form"
	validateFormAttrSelector    = "[" + validateFormAttrName + "]"
	requiredMessageAttrName     = "required-message"
	requiredMessageAttrSelector = "[" + requiredMessageAttrName + "]"
	forbiddenValuesAttrName     = "forbidden-values"
	confirmFieldAttrName        = "confirm-field"
	confirmFieldAttrSelector    = "[" + confirmFieldAttrName + "]"
	confirmMessageAttrName      = "confirm-message"

	displayMessageAttrName     = "display-message"
	displayMessageAttrSelector = "[" + displayMessageAttrName + "]"
)

func AlertKey(messageSpanId string) {
	global := js.Global()
	jsAlert := global.Get("alert")
	errorMessageSpan := global.Get(Document).Call(GetElementById, messageSpanId)
	if jsAlert.Truthy() && errorMessageSpan.Truthy() {
		jsAlert.Invoke(errorMessageSpan.Get(TextContent).String())
	}
}

func TruthyOnclick(button js.Value, action func(js.Value, []js.Value) any) {
	if button.Truthy() {
		button.Set(Onclick, js.FuncOf(action))
	}
}

func RegisterValidationRules() {
	doc := js.Global().Get(Document)
	buttons := doc.Call(QuerySelectorAll, validateFormAttrSelector)
	size := buttons.Length()
	if size == 0 {
		return
	}

	for i := 0; i < size; i++ {
		button := buttons.Index(i)
		formToValidate := doc.Call(GetElementById, button.Call(GetAttribute, validateFormAttrName).String())
		if formToValidate.Truthy() {
			button.Set(Onclick, js.FuncOf(func(_ js.Value, _ []js.Value) any {
				fields := formToValidate.Call(QuerySelectorAll, requiredMessageAttrSelector)
				size := fields.Length()
				for i := 0; i < size; i++ {
					field := fields.Index(i)

					fieldValue := field.Get(Value).String()
					if fieldValue == "" {
						AlertKey(field.Call(GetAttribute, requiredMessageAttrName).String())
						return nil
					}

					if forbiddenValues := field.Call(GetAttribute, forbiddenValuesAttrName).String(); forbiddenValues != "" {
						for _, forbiddenValue := range strings.Split(forbiddenValues, ",") {
							if strings.EqualFold(fieldValue, strings.TrimSpace(forbiddenValue)) {
								AlertKey(field.Call(GetAttribute, requiredMessageAttrName).String())
								return nil
							}
						}
					}
				}

				fields = formToValidate.Call(QuerySelectorAll, confirmFieldAttrSelector)
				size = fields.Length()
				for i := 0; i < size; i++ {
					field := fields.Index(i)

					field2 := doc.Call(GetElementById, field.Call(GetAttribute, confirmFieldAttrName).String())
					if !field2.Truthy() {
						return nil
					}

					if field.Get(Value).String() != field2.Get(Value).String() {
						AlertKey(field.Call(GetAttribute, confirmMessageAttrName).String())
						return nil
					}
				}

				formToValidate.Call(Submit)
				return nil
			}))
		}
	}
}

func RegisterDisplayMessageAction() {
	doc := js.Global().Get(Document)
	buttons := doc.Call(QuerySelectorAll, displayMessageAttrSelector)
	size := buttons.Length()
	if size == 0 {
		return
	}

	for i := 0; i < size; i++ {
		button := buttons.Index(i)
		messageId := button.Call(GetAttribute, displayMessageAttrName).String()
		button.Set(Onclick, js.FuncOf(func(_ js.Value, _ []js.Value) any {
			AlertKey(messageId)
			return nil
		}))
	}
}

// keep the program active to allow function call from HTML/JavaScript
func KeepRunning() {
	<-make(chan struct{})
}
