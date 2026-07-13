# Reference
## APIKeys
<details><summary><code>client.APIKeys.List() -> *platformgo.APIKeyListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists the API keys for the caller's organization, with optional paging, search, and sort.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.APIKeysListRequest{}
client.APIKeys.List(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**limit:** `*int64` — max items per page
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — pagination offset
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.APIKeys.Create(request) -> *platformgo.APIKeyCreateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Mints a fresh API key for the caller's organization. The plaintext key value is returned exactly once and never stored or returned again.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.APIKeyCreateRequest{
        Name: "name",
    }
client.APIKeys.Create(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**name:** `string` — Human-readable label for the API key.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.APIKeys.Delete(KeyID) -> *platformgo.DeleteAPIKeyOutputBody</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Revokes an API key. Subsequent requests using its value are rejected as unauthorized.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.APIKeysDeleteRequest{
        KeyID: "key_id",
    }
client.APIKeys.Delete(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**keyID:** `string` — API key identifier to delete
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.APIKeys.Regenerate(KeyID) -> *platformgo.APIKeyRegenerateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Rotates the plaintext value for an existing API key, preserving its name and identifier. The previous value is invalidated immediately.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.APIKeysRegenerateRequest{
        KeyID: "key_id",
    }
client.APIKeys.Regenerate(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**keyID:** `string` — API key identifier to regenerate
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Billing
<details><summary><code>client.Billing.GetAutoRecharge() -> *platformgo.SubscriptionAutoRechargeSettingsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the organization's automatic balance top-up configuration and status.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Billing.GetAutoRecharge(
        context.TODO(),
    )
}
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Billing.GetBalance() -> *platformgo.SubscriptionBalanceResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the organization's current credit balance in microdollars plus a display string.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Billing.GetBalance(
        context.TODO(),
    )
}
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Billing.GetHistory() -> *platformgo.BillingHistoryResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Paginated invoice history for the caller's organization with optional filters and sort.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.BillingGetHistoryRequest{}
client.Billing.GetHistory(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**limit:** `*int64` — max items per page
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — pagination offset
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — free-text search across invoice number/description
    
</dd>
</dl>

<dl>
<dd>

**billingCycle:** `*string` — filter by billing cycle
    
</dd>
</dl>

<dl>
<dd>

**sortBy:** `*string` — column to sort by
    
</dd>
</dl>

<dl>
<dd>

**sortOrder:** `*string` — asc or desc
    
</dd>
</dl>

<dl>
<dd>

**status:** `*string` — invoice status filter (lowercase)
    
</dd>
</dl>

<dl>
<dd>

**dateFrom:** `*string` — RFC3339 lower bound for invoice_date
    
</dd>
</dl>

<dl>
<dd>

**dateTo:** `*string` — RFC3339 upper bound for invoice_date (with date_from = an exact calendar month)
    
</dd>
</dl>

<dl>
<dd>

**planName:** `*string` — filter by plan name
    
</dd>
</dl>

<dl>
<dd>

**phoneID:** `*string` — filter to invoices for a dedicated phone
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Billing.GetRentalSubscriptions() -> *platformgo.PhoneRentalSubscriptionListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns every active and pending phone rental subscription owned by the caller's organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Billing.GetRentalSubscriptions(
        context.TODO(),
    )
}
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Billing.GetSubscription() -> *platformgo.SubscriptionResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the caller organization's current plan.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Billing.GetSubscription(
        context.TODO(),
    )
}
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Phones
<details><summary><code>client.Phones.Allocate(request) -> *platformgo.PhoneAllocateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Allocates a phone and opens a session. Omit workflow_id for an interactive lease (drive the phone directly); set it to allocate for a workflow. Pass phone_id to pin a specific dedicated phone. If allocation setup fails the claim is rolled back, so you are never billed for a session that never starts.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhoneAllocateRequest{
        PhoneType: platformgo.PhoneAllocateRequestPhoneTypeAndroid,
    }
client.Phones.Allocate(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**liveView:** `*platformgo.PhoneLiveViewOptions` — Hosted live-view options for this session; omit for the defaults (token auth, interactive, enabled).
    
</dd>
</dl>

<dl>
<dd>

**name:** `*string` — Optional session label (letters, numbers, dots, hyphens, underscores; max 64). Unique among the org's active sessions - allocating with a name already in use returns a conflict.
    
</dd>
</dl>

<dl>
<dd>

**phoneID:** `*string` — PhoneID pins allocation to a specific device (for dedicated devices).
    
</dd>
</dl>

<dl>
<dd>

**phoneType:** `*platformgo.PhoneAllocateRequestPhoneType` — Category of device to allocate.
    
</dd>
</dl>

<dl>
<dd>

**recording:** `*bool` — Record this session's screen (default true). false suppresses the video recording and rolling thumbnail entirely - no screen content is ever written.
    
</dd>
</dl>

<dl>
<dd>

**tags:** `map[string]string` — Optional key->value labels for organizing sessions (max 50 tags; keys up to 40 chars, values up to 128).
    
</dd>
</dl>

<dl>
<dd>

**telemetry:** `*bool` — Persist this session's telemetry spans (default true). false skips the durable trace store; the live telemetry stream still works while the session runs.
    
</dd>
</dl>

<dl>
<dd>

**ttl:** `*platformgo.PhoneSessionTTLOptions` — Idle-timeout override for this session; omit for the defaults (inactive after 5 min, close 10 min later).
    
</dd>
</dl>

<dl>
<dd>

**workflowID:** `*string` — Workflow requesting allocation; nil for an interactive lease.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.SupportedApps() -> *platformgo.PhoneSupportedAppsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the apps the platform supports orchestration for, optionally filtered by platform and category.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesSupportedAppsRequest{}
client.Phones.SupportedApps(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**platform:** `*string` — filter by platform
    
</dd>
</dl>

<dl>
<dd>

**category:** `*string` — filter by app category
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.Available() -> *platformgo.PhoneAvailableListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns ACTIVE unallocated phones the caller's org can claim, optionally filtered by phone type (iphone/android). Counts by type are included alongside the list.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesAvailableRequest{}
client.Phones.Available(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**deviceType:** `*string` — filter by device type (iphone/android)
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.Deallocate() -> *platformgo.PhoneDeallocateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deallocates a phone the caller's org currently holds. The session is billed and the phone is torn down asynchronously.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesDeallocateRequest{
        PhoneID: "phone_id",
    }
client.Phones.Deallocate(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**phoneID:** `string` — device identifier to deallocate
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.Mine() -> *platformgo.PhonePrivateListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns private and rented phones owned by the caller's org. include_expired=true keeps rentals past their rental_expires_at in the result so users can see what they used to own.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesMineRequest{}
client.Phones.Mine(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**includeExpired:** `*bool` — include rented devices whose rental window has expired
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` 
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` 
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — free-text search across nickname, name, model, location
    
</dd>
</dl>

<dl>
<dd>

**status:** `[]string` — filter by phone status (ACTIVE/INACTIVE/MAINTENANCE/SUSPENDED)
    
</dd>
</dl>

<dl>
<dd>

**type_:** `[]string` — filter by phone type (IPHONE/ANDROID)
    
</dd>
</dl>

<dl>
<dd>

**rentalExpiresAfter:** `*string` — only phones whose rental expires at/after this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**rentalExpiresBefore:** `*string` — only phones whose rental expires at/before this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**lastActiveAfter:** `*string` — only phones last seen at/after this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**lastActiveBefore:** `*string` — only phones last seen at/before this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**sort:** `*string` — sort column (created|rental_expires|last_active|status|type|location)
    
</dd>
</dl>

<dl>
<dd>

**order:** `*string` — sort direction (asc|desc)
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.ListSessions() -> *platformgo.PhoneSessionListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one page of the org's phone sessions for the Session Inspector table: active/unbilled sessions pinned on top, terminal history paginated beneath. Covers workflow runs and workflow-less interactive leases; each row links to a session. Filters: search, workflow_id, status.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesListSessionsRequest{}
client.Phones.ListSessions(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**limit:** `*int64` — max rows to return (default 50, max 100)
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — rows to skip for pagination
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — case-insensitive match on phone/session/workflow
    
</dd>
</dl>

<dl>
<dd>

**workflowID:** `*string` — only sessions for this workflow
    
</dd>
</dl>

<dl>
<dd>

**status:** `[]string` — filter by session status (ACTIVE/COMPLETED/CANCELLED/EXPIRED); repeatable
    
</dd>
</dl>

<dl>
<dd>

**source:** `[]string` — filter by source: workflow and/or interactive; repeatable
    
</dd>
</dl>

<dl>
<dd>

**dedicated:** `[]string` — filter by type: shared and/or dedicated; repeatable
    
</dd>
</dl>

<dl>
<dd>

**startedAfter:** `*string` — only sessions started at/after this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**startedBefore:** `*string` — only sessions started at/before this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**endedAfter:** `*string` — only sessions de-allocated at/after this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**endedBefore:** `*string` — only sessions de-allocated at/before this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**sort:** `*string` — sort column: started|ended|status|duration (default started)
    
</dd>
</dl>

<dl>
<dd>

**order:** `*string` — sort direction: asc|desc (default desc)
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.ActiveSessions() -> *platformgo.PhoneActiveSessionsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one page of the organization's currently-active phone sessions joined with phone + workflow display fields: in-flight runs, workflow-less interactive leases, and dedicated phones in use. Paginated via limit (default 25, max 100) + offset; the response total is the full active count.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesActiveSessionsRequest{}
client.Phones.ActiveSessions(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**limit:** `*int64` — max rows to return (default 25, max 100)
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — rows to skip for pagination
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — case-insensitive match on phone name/nickname/id, session id, or workflow name
    
</dd>
</dl>

<dl>
<dd>

**dedicated:** `*string` — filter by ownership: 'dedicated' or 'shared'
    
</dd>
</dl>

<dl>
<dd>

**source:** `*string` — filter by source: 'workflow' or 'interactive'
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.GetSession(SessionID) -> *platformgo.PhoneSessionDetailResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one session for the Session Inspector: session lifecycle + phone display fields + workflow name (when tied to one) + an inlined presigned recording URL. Works for active and terminal sessions, and for workflow runs and workflow-less interactive leases. Org-scoped: another org's session reads as not found.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesGetSessionRequest{
        SessionID: "session_id",
    }
client.Phones.GetSession(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**sessionID:** `string` — Phone session identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.SessionRecording(SessionID) -> *platformgo.PhoneSessionRecordingResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a short-lived URL for the session's screen recording, keyed on session_id — so it works for workflow runs and workflow-less interactive leases alike. Status is "pending" (no URL) when the recording hasn't finished uploading yet. Org-scoped: another org's session reads as not found.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesSessionRecordingRequest{
        SessionID: "session_id",
    }
client.Phones.SessionRecording(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**sessionID:** `string` — Phone session identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.SessionThumbnail(SessionID) -> *platformgo.PhoneSessionThumbnailResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a short-lived URL for the session's current screen thumbnail — a rolling JPEG refreshed every few seconds while the session is active. Poll this endpoint and swap the image; every call mints a fresh URL. Status is "pending" (no URL) before the first frame lands or after the session ends. Org-scoped: another org's session reads as not found.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesSessionThumbnailRequest{
        SessionID: "session_id",
    }
client.Phones.SessionThumbnail(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**sessionID:** `string` — Phone session identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.Get(PhoneID) -> *platformgo.PhoneSummary</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a single phone by its identifier.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesGetRequest{
        PhoneID: "phone_id",
    }
client.Phones.Get(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**phoneID:** `string` — device identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.Nickname(PhoneID, request) -> *platformgo.PhoneSummary</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Sets the human-readable display name on a private phone the caller's org owns. Returns the updated phone summary.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhoneUpdateNicknameRequest{
        PhoneID: "phone_id",
        Nickname: "nickname",
    }
client.Phones.Nickname(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**phoneID:** `string` — device identifier
    
</dd>
</dl>

<dl>
<dd>

**nickname:** `string` — New display name for the device.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.Wipe(PhoneID) -> *platformgo.PhoneSuccessResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Requests an on-demand factory reset of a private phone the caller's org owns. Requires the phone to be ACTIVE and not currently allocated. Sets the phone to MAINTENANCE while the wipe is carried out.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.PhonesWipeRequest{
        PhoneID: "phone_id",
    }
client.Phones.Wipe(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**phoneID:** `string` — device identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Runs
<details><summary><code>client.Runs.List(request) -> *platformgo.RunListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns paginated recent runs for the caller's org. Filters: workflow, search text, status, trigger; sort by any of the columns listed in RunSortField.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.RunListRequest{
        Limit: int64(1000000),
        Offset: int64(1000000),
    }
client.Runs.List(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**limit:** `int64` — Maximum number of runs to return per page.
    
</dd>
</dl>

<dl>
<dd>

**offset:** `int64` — Pagination offset.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filters by run ID substring.
    
</dd>
</dl>

<dl>
<dd>

**sortBy:** `[]*platformgo.RunSortSpec` — Ordered list of sort specs; first entry is primary.
    
</dd>
</dl>

<dl>
<dd>

**statusFilter:** `[]*platformgo.RunListRequestStatusFilterItem` — StatusFilter restricts results to runs in the given statuses.
    
</dd>
</dl>

<dl>
<dd>

**triggerFilter:** `[]*platformgo.RunListRequestTriggerFilterItem` — TriggerFilter restricts results to runs with the given triggers.
    
</dd>
</dl>

<dl>
<dd>

**workflowID:** `*string` — Filters results to a single workflow.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.ListEvents(request) -> *platformgo.RunEventsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns paginated run events for a session, filtered by session_id.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.RunEventsRequest{
        Limit: int64(1000000),
        Offset: int64(1000000),
        SessionID: "session_id",
    }
client.Runs.ListEvents(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**eventTypes:** `[]string` — EventTypes restricts results to specific event type codes (RUN_STARTED / OUTPUT_LOG / SDK_CALL_COMPLETED / etc.).
    
</dd>
</dl>

<dl>
<dd>

**limit:** `int64` — Maximum number of events to return.
    
</dd>
</dl>

<dl>
<dd>

**offset:** `int64` — Pagination offset.
    
</dd>
</dl>

<dl>
<dd>

**sessionID:** `string` — Filters events to a specific device session (formerly allocation_id; W6-2).
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.ListHistoric(request) -> *platformgo.RunHistoryResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns paginated historic runs for the caller's user. Use POST /runs for recent (non-archived) runs.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.RunHistoryRequest{
        EndDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        Limit: int64(1000000),
        Offset: int64(1000000),
        StartDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
    }
client.Runs.ListHistoric(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**endDate:** `time.Time` — End of the query time window.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `int64` — Maximum number of runs to return.
    
</dd>
</dl>

<dl>
<dd>

**offset:** `int64` — Pagination offset.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filters by run ID, workflow ID, or device ID.
    
</dd>
</dl>

<dl>
<dd>

**startDate:** `time.Time` — Beginning of the query time window.
    
</dd>
</dl>

<dl>
<dd>

**statusFilter:** `[]string` — StatusFilter restricts results to runs in the given statuses.
    
</dd>
</dl>

<dl>
<dd>

**workflowID:** `*string` — Filters results to a single workflow.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.Stats(WorkflowID) -> *platformgo.RunStatsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns total run count + success rate for the given workflow, scoped to the caller's user.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.RunsStatsRequest{
        WorkflowID: "workflow_id",
    }
client.Runs.Stats(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.Get(RunID) -> *platformgo.RunResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one run by ID, scoped to the caller's organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.RunsGetRequest{
        RunID: "run_id",
    }
client.Runs.Get(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**runID:** `string` — run identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.Cancel(RunID) -> *platformgo.RunSuccessResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Transitions a run to CANCELLED, scoped to the caller's org.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.RunsCancelRequest{
        RunID: "run_id",
    }
client.Runs.Cancel(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**runID:** `string` — run identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.Create(WorkflowID, request) -> *platformgo.RunCreateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates one or more runs against the given workflow and queues them for execution. Pre-flight checks: balance sufficient, concurrency limit, workflow exists. Runs that fail to queue are marked FAILED immediately so they stop counting toward the concurrency limit.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.RunCreateRequest{
        WorkflowID: "workflow_id",
        Count: int64(1000000),
    }
client.Runs.Create(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow to create runs for
    
</dd>
</dl>

<dl>
<dd>

**count:** `int64` — Number of runs to create.
    
</dd>
</dl>

<dl>
<dd>

**runs:** `[]*platformgo.RunConfig` — Per-run variable configurations.
    
</dd>
</dl>

<dl>
<dd>

**startTimeoutSeconds:** `*int64` — How long a queued run may wait for a phone before it is auto-cancelled.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Usage
<details><summary><code>client.Usage.ListInferences(request) -> *platformgo.UsageInferencesResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Paginated, filterable list of inference calls (/infer + /locate) the caller's user was billed for. Filters: date range, endpoint, free-text search. Ordered by call time DESC.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.UsageInferencesRequest{
        EndDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        StartDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
    }
client.Usage.ListInferences(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**endDate:** `time.Time` — End of the inferences query window.
    
</dd>
</dl>

<dl>
<dd>

**endpointFilter:** `[]string` — Restricts results to the given vision endpoints ('detect'/'locate').
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Number of inferences per page.
    
</dd>
</dl>

<dl>
<dd>

**model:** `*string` — Model restricts results to a single model name.
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — Pagination offset.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filters by inference (event) ID substring.
    
</dd>
</dl>

<dl>
<dd>

**sessionID:** `*string` — Restricts results to inferences that ran under one phone session.
    
</dd>
</dl>

<dl>
<dd>

**sortBy:** `[]*platformgo.UsageInferenceSortSpec` — Ordered list of sort specs; first entry is primary.
    
</dd>
</dl>

<dl>
<dd>

**startDate:** `time.Time` — Beginning of the inferences query window.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Usage.GetMetrics() -> *platformgo.UsageMetricsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns infrastructure cost and compute-minute summaries for the caller's user over a date range, plus per-bucket chart data. Granularity is hourly (≤24h window) or daily. Use POST /usage/metrics if you need richer body params; this endpoint takes query params only.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.UsageGetMetricsRequest{
        StartDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        EndDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
    }
client.Usage.GetMetrics(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**startDate:** `time.Time` — start of reporting window (RFC3339)
    
</dd>
</dl>

<dl>
<dd>

**endDate:** `time.Time` — end of reporting window (RFC3339)
    
</dd>
</dl>

<dl>
<dd>

**granularity:** `*platformgo.UsageGetMetricsRequestGranularity` — bucket resolution
    
</dd>
</dl>

<dl>
<dd>

**timezone:** `*string` — IANA timezone for bucketing periods (e.g., America/Los_Angeles)
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Workflows
<details><summary><code>client.Workflows.List() -> *platformgo.WorkflowListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Paginated list of workflows in the caller's org, with optional search, status, platform, and created/last-run date filters via query params.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowsListRequest{}
client.Workflows.List(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**limit:** `*int64` 
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` 
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — free-text search across workflow name
    
</dd>
</dl>

<dl>
<dd>

**status:** `[]string` — filter by workflow status (lowercase)
    
</dd>
</dl>

<dl>
<dd>

**platform:** `[]string` — filter by device platform (lowercase)
    
</dd>
</dl>

<dl>
<dd>

**createdAfter:** `*string` — only workflows created at/after this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `*string` — only workflows created at/before this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**lastRunAfter:** `*string` — only workflows last run at/after this RFC3339 time
    
</dd>
</dl>

<dl>
<dd>

**lastRunBefore:** `*string` — only workflows last run at/before this RFC3339 time
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.Create(request) -> *platformgo.WorkflowCreateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a workflow in the caller's org. Name must match ^[A-Za-z0-9_-]+$ and be unique within the org. Pass code to save the workflow's first code revision atomically with it; omit it to create an empty workflow and add code later. Returns the workflow_id (plus revision_id and revision when code was provided).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowCreateRequest{
        Name: "name",
    }
client.Workflows.Create(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**code:** `*string` — Optional Python source for the workflow's first revision, saved atomically with the workflow when provided.
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — Human-readable workflow name.
    
</dd>
</dl>

<dl>
<dd>

**ocrEngine:** `*string` — OCR backend to use.
    
</dd>
</dl>

<dl>
<dd>

**platform:** `*string` — Target OS platform.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.Get(WorkflowID) -> *platformgo.WorkflowResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a single workflow, scoped to the caller's org (workflows in other orgs return 404).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowsGetRequest{
        WorkflowID: "workflow_id",
    }
client.Workflows.Get(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.Update(WorkflowID, request) -> *platformgo.WorkflowResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Applies a partial update (name, platform, status, ocr_engine). Org-scoped — workflows in other orgs return 404.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowUpdateRequest{
        WorkflowID: "workflow_id",
    }
client.Workflows.Update(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>

<dl>
<dd>

**name:** `*string` — Updated workflow name.
    
</dd>
</dl>

<dl>
<dd>

**ocrEngine:** `*string` — Updated OCR backend selection.
    
</dd>
</dl>

<dl>
<dd>

**platform:** `*string` — Updated target platform.
    
</dd>
</dl>

<dl>
<dd>

**status:** `*string` — Updated lifecycle status.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.Delete(WorkflowID) -> *platformgo.MessageOutputBody</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes a workflow. Org-scoped — workflows in other orgs return 404.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowsDeleteRequest{
        WorkflowID: "workflow_id",
    }
client.Workflows.Delete(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.GetCode(WorkflowID) -> *platformgo.WorkflowGetCodeResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the source code of the workflow's current revision, scoped to the caller's org.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowsGetCodeRequest{
        WorkflowID: "workflow_id",
    }
client.Workflows.GetCode(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.SaveCode(WorkflowID, request) -> *platformgo.WorkflowSaveCodeResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Persists a new revision of the workflow's code. Hash-deduplicates against the current revision (no-op if unchanged). Source is capped at 256KB.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowSaveCodeRequest{
        WorkflowID: "workflow_id",
        Source: "source",
    }
client.Workflows.SaveCode(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>

<dl>
<dd>

**message:** `*string` — Optional commit-style note.
    
</dd>
</dl>

<dl>
<dd>

**source:** `string` — Python source the user typed.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.RestoreRevision(WorkflowID, request) -> *platformgo.WorkflowSaveCodeResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a new revision with the source of the named revision. Does NOT dedup against the current revision so the action is auditable in the revision history.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowRestoreRevisionRequest{
        WorkflowID: "workflow_id",
        RevisionID: "revision_id",
    }
client.Workflows.RestoreRevision(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>

<dl>
<dd>

**revisionID:** `string` — Revision to restore.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.ListRevisions(WorkflowID) -> *platformgo.WorkflowListRevisionsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns revision metadata (id, number, author, message, bytes, sha256, created_at) for the workflow, in reverse-chronological order. Use the `before` cursor to paginate older revisions.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowsListRevisionsRequest{
        WorkflowID: "workflow_id",
    }
client.Workflows.ListRevisions(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` 
    
</dd>
</dl>

<dl>
<dd>

**before:** `*int64` — cursor: return revisions older than this revision number
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.GetRevision(WorkflowID, RevisionID) -> *platformgo.WorkflowRevisionDetail</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a single revision including its full source. Defense-in-depth: the revision must belong to the named workflow or it's treated as missing.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &platformgo.WorkflowsGetRevisionRequest{
        WorkflowID: "workflow_id",
        RevisionID: "revision_id",
    }
client.Workflows.GetRevision(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**workflowID:** `string` — workflow identifier
    
</dd>
</dl>

<dl>
<dd>

**revisionID:** `string` — revision identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

