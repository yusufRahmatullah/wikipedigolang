package utils

import (
	"fmt"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/logger"
)

// OnPageOpennedListener is the function that will be called
// if the WebPage success to open an URL
type OnPageOpennedListener func()

// OnEvaluatedListener is the function that will be called
// if the WebPage success to evaluate a JS script
type OnEvaluatedListener func(map[string]interface{})

// OnAsyncEvaluatedListener is the function that will be called
// if the WebPage success to evaluate an async JS script
type OnAsyncEvaluatedListener func()

// OnErrorListener is the function that will be called
// if any error occurred
type OnErrorListener func(err error)

// WebPageWrapper wraps and simplify WebPage
// Don't forget to defer Close method
type WebPageWrapper struct {
	page                      *WebPage
	mLogger                   *logger.Logger
	onPageOpennedListeners    []OnPageOpennedListener
	onEvaluatedListeners      []OnEvaluatedListener
	onErrorListeners          []OnErrorListener
	onAsyncEvaluatedListeners []OnAsyncEvaluatedListener
}

// NewWebPageWrapper instantiate WebPageWrapper instance
func NewWebPageWrapper(mLogger *logger.Logger) *WebPageWrapper {
	page, err := CreateWebPage()
	if err == nil {
		return &WebPageWrapper{page: page, mLogger: mLogger}
	}
	mLogger.Fatal(fmt.Sprintf("Failed to create page. Causes: %v", err))
	return nil
}

func (ww *WebPageWrapper) callOnError(err error) {
	for _, oer := range ww.onErrorListeners {
		oer(err)
	}
}

// OnPageOpenned add OnPageOpennedListener to WebPageWrapper
func (ww *WebPageWrapper) OnPageOpenned(fn OnPageOpennedListener) {
	ww.onPageOpennedListeners = append(ww.onPageOpennedListeners, fn)
}

// OpenURL try to open an URL and run all OnPageOpennedListeners
// if success
func (ww *WebPageWrapper) OpenURL(url string) {
	err := ww.page.Open(url)
	if err == nil {
		for _, opol := range ww.onPageOpennedListeners {
			opol()
		}
	} else {
		ww.mLogger.Fatal(fmt.Sprintf("Failed to open URL: %v", url))
		ww.callOnError(err)
	}
}

// OnEvaluated add OnEvaluatedListener to WebPageWrapper
func (ww *WebPageWrapper) OnEvaluated(fn OnEvaluatedListener) {
	ww.onEvaluatedListeners = append(ww.onEvaluatedListeners, fn)
}

// Evaluate try to evaluate JS script and
// run all OnEvaluatedListeners if success
func (ww *WebPageWrapper) Evaluate(js string) {
	info, err := ww.page.Evaluate(js)
	if err == nil {
		if info == nil {
			ww.mLogger.Fatal("Failed, info is nil")
		} else {
			data := info.(map[string]interface{})
			for _, oel := range ww.onEvaluatedListeners {
				oel(data)
			}
		}
	} else {
		ww.mLogger.Fatal(fmt.Sprintf("Failed to execute JS. Causes: %v", err))
		ww.callOnError(err)
	}
}

// OnAsyncEvaluated add OnAsyncEvaluatedListener to WebPageWrapper
func (ww *WebPageWrapper) OnAsyncEvaluated(fn OnAsyncEvaluatedListener) {
	ww.onAsyncEvaluatedListeners = append(ww.onAsyncEvaluatedListeners, fn)
}

// EvaluateAsync try to evaluate an Async JS script
// should be waiting for receive the data
func (ww *WebPageWrapper) EvaluateAsync(js string) {
	err := ww.page.EvaluateAsync(js, 1*time.Millisecond)
	if err == nil {
		for _, oael := range ww.onAsyncEvaluatedListeners {
			oael()
		}
	} else {
		ww.mLogger.Fatal(fmt.Sprintf("Failed to execute async script. Causes: %v", err))
		ww.callOnError(err)
	}
}

// OnError add OnErrorListeners to WebPageWrapper
func (ww *WebPageWrapper) OnError(fn OnErrorListener) {
	ww.onErrorListeners = append(ww.onErrorListeners, fn)
}

// Close end WebPage by calling the Close method
func (ww *WebPageWrapper) Close() {
	ww.page.Close()
}
