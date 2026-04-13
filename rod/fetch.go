package rod

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/proto"
	readability "github.com/go-shiori/go-readability"
)

type Output int

const (
	OutputMarkdown Output = iota
	OutputText
)

type FetchOption struct {
	Timeout   time.Duration
	IdleWait  time.Duration
	MaxLength int
	UserAgent string
	Output    Output
}

const (
	defaultTimeout   = 30 * time.Second
	defaultIdleWait  = 5 * time.Second
	defaultMaxLength = 100 << 10
)

func Fetch(ctx context.Context, href string, opt *FetchOption) (string, error) {
	if opt == nil {
		opt = &FetchOption{}
	}
	timeout := opt.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	idle := opt.IdleWait
	if idle == 0 {
		idle = defaultIdleWait
	}
	maxLen := opt.MaxLength
	if maxLen == 0 {
		maxLen = defaultMaxLength
	}

	parsed, err := url.Parse(href)
	if err != nil {
		return "", fmt.Errorf("url.Parse: %w", err)
	}
	if parsed.Scheme == "" || !strings.Contains(parsed.Hostname(), ".") {
		return "", fmt.Errorf("invalid url: %s", href)
	}

	b, err := ensureBrowser(opt.UserAgent)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	page, err := b.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return "", fmt.Errorf("browser.Page: %w", err)
	}
	defer func() { _ = page.Close() }()

	page = page.Context(ctx)
	if err := page.Navigate(href); err != nil {
		return "", fmt.Errorf("page.Navigate: %w", err)
	}
	if err := page.WaitLoad(); err != nil {
		return "", fmt.Errorf("page.WaitLoad: %w", err)
	}
	_ = page.WaitIdle(idle)

	htmlSrc, err := page.HTML()
	if err != nil {
		return "", fmt.Errorf("page.HTML: %w", err)
	}

	article, err := readability.FromReader(strings.NewReader(htmlSrc), parsed)
	if err != nil {
		return "", fmt.Errorf("readability: %w", err)
	}

	var text string
	switch opt.Output {
	case OutputText:
		text = strings.TrimSpace(article.TextContent)
	default:
		md, err := HTMLToMarkdown(article.Content, href)
		if err != nil {
			return "", fmt.Errorf("HTMLToMarkdown: %w", err)
		}
		text = md
	}
	if text == "" {
		return "", fmt.Errorf("empty content")
	}
	if maxLen > 0 && len(text) > maxLen {
		text = text[:maxLen]
	}
	return text, nil
}
