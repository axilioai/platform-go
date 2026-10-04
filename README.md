# Axilio Go Library

[![fern shield](https://img.shields.io/badge/%F0%9F%8C%BF-Built%20with%20Fern-brightgreen)](https://buildwithfern.com?utm_source=github&utm_medium=github&utm_campaign=readme&utm_source=Axilio%2FGo)

The Axilio Go library provides convenient access to the Axilio APIs from Go: a
typed REST client for phones, workflows, runs, files, usage and billing, plus a
hand-written mobile driver (`drivers/mobile`) that drives an allocated phone over
its device control channel.

## Table of Contents

- [Quickstart](#quickstart)
  - [Accessibility mode](#accessibility-mode)
- [Reference](#reference)
- [Usage](#usage)
- [Environments](#environments)
- [Errors](#errors)
- [Request Options](#request-options)
- [Advanced](#advanced)
  - [Response Headers](#response-headers)
  - [Retries](#retries)
  - [Timeouts](#timeouts)
  - [Explicit Null](#explicit-null)
- [Contributing](#contributing)

## Quickstart

The task-first walkthrough lives in the
[Go quickstart](https://docs.axilio.ai/go-quickstart): install, authenticate,
allocate a real Android phone, inspect its screen, and release it. The same
flow in one file:

```bash
go get github.com/axilioai/platform-go@latest
export AXILIO_API_KEY=axl_your_key_here
```

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "os"

    platformgo "github.com/axilioai/platform-go"
    "github.com/axilioai/platform-go/client"
    "github.com/axilioai/platform-go/drivers/mobile"
    "github.com/axilioai/platform-go/option"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    apiKey := os.Getenv("AXILIO_API_KEY")
    if apiKey == "" {
        return errors.New("AXILIO_API_KEY is not set")
    }

    ctx := context.Background()
    c := client.NewClient(option.WithAPIKey(apiKey))

    // Allocate a phone from the shared pool. Pass PhoneID to pin a dedicated one.
    session, err := c.Phones.Allocate(ctx, &platformgo.PhoneAllocateRequest{
        PhoneType: platformgo.PhoneAllocateRequestPhoneTypeAndroid,
    })
    if err != nil {
        return fmt.Errorf("allocate phone: %w", err)
    }
    // Release the phone on every return path, including errors below.
    defer func() {
        if _, err := c.Phones.Deallocate(context.Background(), &platformgo.PhonesDeallocateRequest{
            PhoneID: session.PhoneID,
        }); err != nil {
            log.Printf("release phone: %v", err)
        }
    }()

    if session.ControlURL == nil {
        return errors.New("allocation returned no control URL")
    }

    driver := mobile.ConnectRemote(*session.ControlURL)
    defer func() { _ = driver.Close() }()

    screen, err := driver.Observe()
    if err != nil {
        return fmt.Errorf("observe screen: %w", err)
    }
    fmt.Printf("Found %d text regions and %d icons\n", len(screen.Texts), len(screen.Icons))

    png, err := driver.Screenshot()
    if err != nil {
        return fmt.Errorf("capture screenshot: %w", err)
    }
    if err := os.WriteFile("screen.png", png, 0o600); err != nil {
        return fmt.Errorf("save screenshot: %w", err)
    }
    fmt.Println("Saved screen.png")
    return nil
}
```

The driver is built around locators: `driver.GetByText("Settings", mobile.Exact())`
and `driver.Locator(mobile.Query("the blue Continue button"))` each describe a
target without touching the device. Calling `loc.Tap()`, `loc.Fill(text)`,
`loc.Press(mobile.KeyEnter)`, `loc.WaitFor(mobile.StateVisible)`,
`loc.BoundingBox()`, `loc.Text()` or `loc.Count()` sends it: the device
resolves the locator, auto-waits until it's actionable, and acts, all in one
round trip. `Count` is the exception: it reports how many targets match the
current screen right now, zero included, and never waits, so use `WaitFor`
(not `Count`) to wait for something to appear. Raw input (`Tap`, `Swipe`,
`TypeText`, `KeyPress`, `Press`) by literal coordinate stays on the driver,
and `Observe` still returns a `Screen` of OCR/icon data you can filter
locally with `Screen.FindText`/`FindAllText`.

Resolution options (which model, which OCR engine) belong to the locator, not
the action: pass `mobile.Model("...")` or `mobile.OCREngine("premium")`
alongside the selector options (`mobile.Text`, `mobile.Exact`, `mobile.Query`,
`mobile.Role`, ...) to `driver.GetByText` or `driver.Locator`. An action or query on the locator
(`Tap`/`Fill`/`Press`/`WaitFor`/`BoundingBox`/`Text`/`Count`) takes a
`mobile.ActionOption`, not a `mobile.CallOption`: the only one is
`mobile.WithTimeout`, and it's a compile error to pass `mobile.WithOCREngine`
(a `CallOption`, for `Observe`) to a locator action instead of setting
`mobile.OCREngine` on the locator. The resolution follows the locator that
resolves, else the driver's `WithDefaultModel`/`WithDefaultOCREngine`/`WithDefaultStrategy`, else
it's left off the wire so the server's own default applies. Refining a
locator (`Nth`, `First`, `Within`, `Has`, `Filter`) keeps the receiver's own
resolution options; the locator passed into `Within`/`Has` only ever
contributes its selector fields, since one call resolves the whole locator
and the outer locator's options govern it. If that inner locator carries its
own `Model`/`OCREngine`/`Strategy` (set on itself, not inherited from a driver default),
`Within`/`Has` record a build error on the result instead of silently
dropping them: every action or query on it (and on anything further refined
from it) fails locally with a `CodeInvalidArgs` `*Error` naming what to do,
and nothing is sent. `driver.Press(key, ...)` (no locator) takes no
resolution options at all.

Without the accessibility tree, a plain `mobile.Text` locator is matched by OCR; a locator that also carries
`mobile.Query`, `Within`, `Has` or `Nth` is instead resolved by one
vision-model call, with a prompt composed from the whole locator, so `Nth` on
a query-based locator works. `Count` is the exception: it needs a plain text
locator, and answers a `CodeInvalidArgs` error for one that also carries
`Query`, `Within` or `Has`, since counting needs every independent match and
a vision-model call only resolves a single target per prompt.

### Accessibility mode

Sessions run with the phone's accessibility tree on by default.
`PhoneAllocateRequest.Accessibility` defaults to true when left nil, and true
requires a phone that supports it: only such phones are claimed, and a
`PhoneID` that does not support it is refused with a conflict
(`IsAccessibilityUnavailable`). Set it to `platformgo.Bool(false)` to allocate
any phone with the tree off. The response's `Accessibility` field reports the
value. While the tree is on, the accessibility service is visible to apps on
the phone.

```go
session, err := c.Phones.Allocate(ctx, &platformgo.PhoneAllocateRequest{
    PhoneType: platformgo.PhoneAllocateRequestPhoneTypeAndroid,
    PhoneID:   platformgo.String("your-dedicated-phone-id"),
    // Accessibility is on by default; platformgo.Bool(false) allows any phone.
})
if platformgo.IsAccessibilityUnavailable(err) {
    // The named phone does not support accessibility mode. Allocate it with
    // Accessibility: platformgo.Bool(false), or drop PhoneID.
}
```

Workflows take the same setting: `WorkflowCreateRequest.Accessibility`
defaults to true, and `WorkflowUpdateRequest.Accessibility` leaves the
current value unchanged when nil.

With the tree on, literal selectors resolve against it on the device:

```go
driver.GetByRole("button", mobile.Name("Log in")).Tap()
driver.GetByID("com.example.app:id/login").Tap()
driver.GetByRole("textbox", mobile.Name("Email")).Fill("me@example.com")
driver.Locator(mobile.Query("the log in button"), mobile.Strategy(mobile.StrategyVision)).Tap()
```

The tree-only options are `Role`, `Name`, `ID`, `States`, `Value`,
`WindowID`, `NodeID`, `AndroidClassName` and `AndroidPackageName` (`Exact`
applies to `Text`, `Name` and `Value`). On a session whose tree is off they
answer `IsStrategyUnavailable`; they are never turned into a model prompt.
`mobile.Strategy` (or the driver-wide `mobile.WithDefaultStrategy`) picks the
resolver: `StrategyAuto` (the default: the tree when it is on, vision
otherwise), `StrategyVision` or `StrategyAccessibility`. `WaitFor` also
accepts `StateEnabled`, which needs the tree.

`driver.Accessibility()` reads the tree itself:

```go
tree, err := driver.Accessibility().Snapshot()            // every window, layout nodes dropped
tree, err = driver.Accessibility().Snapshot(mobile.WithWindow(id), mobile.WithInterestingOnly(false))
nodes, err := driver.Accessibility().Query(mobile.AXQuery{Role: "button"})
state, err := driver.Accessibility().State()              // Enabled, Toggleable
err = driver.Accessibility().Disable()                    // returns once the phone confirms
err = driver.Accessibility().Enable()
```

`Partial(nodeID, fetchRelatives)` and `Children(nodeID)` read around one
node; a node id also works as a locator (`mobile.NodeID(id)`). Errors to
expect: `IsStrategyUnavailable` (the tree is off), `IsTreeUnavailable` (a
system dialog such as a permission prompt covers the app; vision still sees
it) and `IsStaleNode` (the node from an earlier snapshot is gone).
`Enable` and `Disable` work only where `State().Toggleable` is true.

## Reference

A full reference for this library is available [here](./reference.md).

## Usage

Instantiate and use the client with the following:

```go
package example

import (
    context "context"

    platformgo "github.com/axilioai/platform-go"
    client "github.com/axilioai/platform-go/client"
    option "github.com/axilioai/platform-go/option"
)

func do() {
    client := client.NewClient(
        option.WithAPIKey(
            "<value>",
        ),
    )
    request := &platformgo.APIKeyCreateRequest{
        Name: "name",
    }
    client.APIKeys.Create(
        context.TODO(),
        request,
    )
}
```

## Environments

You can choose between different environments by using the `option.WithBaseURL` option. You can configure any arbitrary base
URL, which is particularly useful in test environments.

```go
client := client.NewClient(
    option.WithBaseURL(api.Environments.Default),
)
```

## Errors

Structured error types are returned from API calls that return non-success status codes. These errors are compatible
with the `errors.Is` and `errors.As` APIs, so you can access the error like so:

```go
response, err := client.APIKeys.Create(...)
if err != nil {
    var apiError *core.APIError
    if errors.As(err, apiError) {
        // Do something with the API error ...
    }
    return err
}
```

## Request Options

A variety of request options are included to adapt the behavior of the library, which includes configuring
authorization tokens, or providing your own instrumented `*http.Client`.

These request options can either be
specified on the client so that they're applied on every request, or for an individual request, like so:

> Providing your own `*http.Client` is recommended. Otherwise, the `http.DefaultClient` will be used,
> and your client will wait indefinitely for a response (unless the per-request, context-based timeout
> is used).

```go
// Specify default options applied on every request.
client := client.NewClient(
    option.WithToken("<YOUR_API_KEY>"),
    option.WithHTTPClient(
        &http.Client{
            Timeout: 5 * time.Second,
        },
    ),
)

// Specify options for an individual request.
response, err := client.APIKeys.Create(
    ...,
    option.WithToken("<YOUR_API_KEY>"),
)
```

## Advanced

### Response Headers

You can access the raw HTTP response data by using the `WithRawResponse` field on the client. This is useful
when you need to examine the response headers received from the API call. (When the endpoint is paginated,
the raw HTTP response data will be included automatically in the Page response object.)

```go
response, err := client.APIKeys.WithRawResponse.Create(...)
if err != nil {
    return err
}
fmt.Printf("Got response headers: %v", response.Header)
fmt.Printf("Got status code: %d", response.StatusCode)
```

### Retries

The SDK is instrumented with automatic retries with exponential backoff. A request will be retried as long
as the request is deemed retryable and the number of retry attempts has not grown larger than the configured
retry limit (default: 2).

Which status codes are retried depends on the `retryStatusCodes` generator configuration:

**`legacy`** (current default): retries on
- [408](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/408) (Timeout)
- [429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) (Too Many Requests)
- [5XX](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#server_error_responses) (All server errors, including 500)

**`recommended`**: retries on
- [408](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/408) (Timeout)
- [429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) (Too Many Requests)
- [502](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/502) (Bad Gateway)
- [503](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/503) (Service Unavailable)
- [504](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/504) (Gateway Timeout)

If the `Retry-After` header is present in the response, the SDK will prioritize respecting its value exactly
over the default exponential backoff.

Use the `option.WithMaxAttempts` option to configure this behavior for the entire client or an individual request:

```go
client := client.NewClient(
    option.WithMaxAttempts(1),
)

response, err := client.APIKeys.Create(
    ...,
    option.WithMaxAttempts(1),
)
```

### Timeouts

Setting a timeout for each individual request is as simple as using the standard context library. Setting a one second timeout for an individual API call looks like the following:

```go
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel()

response, err := client.APIKeys.Create(ctx, ...)
```

### Explicit Null

If you want to send the explicit `null` JSON value through an optional parameter, you can use the setters\
that come with every object. Calling a setter method for a property will flip a bit in the `explicitFields`
bitfield for that setter's object; during serialization, any property with a flipped bit will have its
omittable status stripped, so zero or `nil` values will be sent explicitly rather than omitted altogether:

```go
type ExampleRequest struct {
    // An optional string parameter.
    Name *string `json:"name,omitempty" url:"-"`

    // Private bitmask of fields set to an explicit value and therefore not to be omitted
    explicitFields *big.Int `json:"-" url:"-"`
}

request := &ExampleRequest{}
request.SetName(nil)

response, err := client.APIKeys.Create(ctx, request, ...)
```

## Contributing

While we value open-source contributions to this SDK, this library is generated programmatically.
Additions made directly to this library would have to be moved over to our generation code,
otherwise they would be overwritten upon the next generated release. Feel free to open a PR as
a proof of concept, but know that we will not be able to merge it as-is. We suggest opening
an issue first to discuss with us!

This README is hand-maintained and survives regeneration (see CONTRIBUTING.md),
so contributions to it are always welcome.
