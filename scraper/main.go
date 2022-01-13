package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func main() {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// navigate to a page, wait for an element, click
	var start string
	var cap string
	var max string
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(`https://app.unicrypt.network/amm/pancake-v2/ilo/0x1a1039810162266238a569c5141D12c8FdE779Ed`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.ScrollIntoView(`footer`),
		chromedp.WaitVisible(`div.v-tabs`, chromedp.ByQuery),
		chromedp.EvaluateAsDevTools(`$x("/html/body/div/div[1]/main/div/div[2]/div/div[2]/div[2]/div[3]/div[2]/div[3]/div[1]/div/div[2]/div/div[3]")[0].click()`, nil),

		chromedp.Text(`//div[contains(., 'Start block')][not(@class)]`, &start, chromedp.NodeVisible, chromedp.BySearch),
		chromedp.Text(`/html/body/div/div[1]/main/div/div[2]/div/div[2]/div[2]/div[3]/div[2]/div[3]/div[2]/div[2]/div/div[4]/div/div[1]`, &cap, chromedp.NodeVisible, chromedp.BySearch),
		chromedp.Text(`/html/body/div/div[1]/main/div/div[2]/div/div[2]/div[2]/div[3]/div[2]/div[3]/div[2]/div[2]/div/div[5]/div/div[1]`, &max, chromedp.NodeVisible, chromedp.BySearch),
	)
	if err != nil {
		log.Fatal(err)
	}

	// parts := strings.Split(start, " ")
	// start = parts[0]
	cap = strings.TrimRight(cap, " BN")
	max = strings.TrimRight(max, " BN")

	fmt.Println(start + "\n" + cap + "\n" + max)
}
