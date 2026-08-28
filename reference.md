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

Lists the API keys for the caller's organization, with optional paging (limit + offset). Ordered newest-first.
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

<details><summary><code>client.Billing.UpdateAutoRecharge(request) -> *platformgo.SubscriptionAutoRechargeSettingsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Configures automatic balance top-up: when the balance drops below the threshold, a one-off invoice charges the saved payment method to restore it to the target. Requires an admin-role caller: an API key carries its creator's organization role, so the key must belong to an org admin. This only tunes when the saved payment method is charged — adding funds or payment methods stays in the dashboard.
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
request := &platformgo.SubscriptionAutoRechargeSettings{
        Enabled: true,
        TargetCents: int64(1000000),
        ThresholdCents: int64(1000000),
    }
client.Billing.UpdateAutoRecharge(
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

**enabled:** `bool` — Whether auto-recharge is active.
    
</dd>
</dl>

<dl>
<dd>

**targetCents:** `int64` — Balance the recharge restores to, in cents (minimum 500 = $5.00). The charge is target minus current balance.
    
</dd>
</dl>

<dl>
<dd>

**thresholdCents:** `int64` — Recharge when the balance drops below this amount, in cents.
    
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

<details><summary><code>client.Billing.DownloadInvoice(InvoiceID) -> *platformgo.BillingHistoryInvoiceDownloadResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a temporary PDF download URL for an invoice the caller's org owns. The URL expires; re-request it rather than storing it.
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
request := &platformgo.BillingDownloadInvoiceRequest{
        InvoiceID: "invoice_id",
    }
client.Billing.DownloadInvoice(
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

**invoiceID:** `string` — Billing history item ID to download.
    
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

Returns every phone rental subscription owned by the caller's organization, in any lifecycle state (active or canceled) - canceled rentals are included so past rentals stay visible. Each item's status field tells them apart. The summary block (active_count, combined monthly cost) counts only active rentals.
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

Returns the caller organization's active subscription plan. Only an active subscription is returned; if the org has none, this responds 404. (The status field is therefore always 'active' here.)
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

<details><summary><code>client.Billing.GetUsageAlerts() -> *platformgo.SubscriptionUsageAlertSettingsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the organization's balance-alert configuration and any currently open alerts (low balance, negative balance).
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
client.Billing.GetUsageAlerts(
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

<details><summary><code>client.Billing.UpdateUsageAlerts(request) -> *platformgo.SubscriptionUsageAlertSettingsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Configures the low-balance alert: while enabled and the balance sits below the threshold, an alert stays open and surfaces in the dashboard. Negative-balance alerts are always on. Requires an admin-role caller: an API key carries its creator's organization role, so the key must belong to an org admin.
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
request := &platformgo.SubscriptionUsageAlertSettings{
        Enabled: true,
        ThresholdCents: int64(1000000),
    }
client.Billing.UpdateUsageAlerts(
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

**enabled:** `bool` — Whether low-balance alerts are active. Negative-balance alerts are always on.
    
</dd>
</dl>

<dl>
<dd>

**thresholdCents:** `int64` — Alert while the balance sits below this amount, in cents. Must be positive.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Files
<details><summary><code>client.Files.List() -> *platformgo.FileListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one page of the org's file library, newest first, with the org's standing usage against its storage quota. Every file carries its source (upload or capture) and, for a capture, its surface and capture state. Filter with source, surface and session_id. Files persist until deleted and any ready file can be delivered to a phone the org holds.
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
request := &platformgo.FilesListRequest{}
client.Files.List(
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

**q:** `*string` — filter by filename, case-insensitive substring match
    
</dd>
</dl>

<dl>
<dd>

**source:** `*platformgo.FilesListRequestSource` — only files of this source
    
</dd>
</dl>

<dl>
<dd>

**surface:** `*platformgo.FilesListRequestSurface` — only captures off this surface
    
</dd>
</dl>

<dl>
<dd>

**sessionID:** `*string` — only files captured by this session
    
</dd>
</dl>

<dl>
<dd>

**mimeType:** `*string` — only files of exactly this media type
    
</dd>
</dl>

<dl>
<dd>

**minSizeBytes:** `*int64` — only files at least this many bytes (0 = no bound)
    
</dd>
</dl>

<dl>
<dd>

**maxSizeBytes:** `*int64` — only files at most this many bytes (0 = no bound)
    
</dd>
</dl>

<dl>
<dd>

**createdAfter:** `*time.Time` — only files registered at or after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `*time.Time` — only files registered at or before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**sort:** `*platformgo.FilesListRequestSort` — field to sort by; source groups uploads and captures
    
</dd>
</dl>

<dl>
<dd>

**order:** `*platformgo.FilesListRequestOrder` — sort direction
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Files.Create(request) -> *platformgo.FileUploadResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Registers an image or video in the org's file library and returns a presigned S3 URL to upload the bytes to. PUT the raw file to upload_url with the declared Content-Type and Content-Length headers before the URL expires, then call the complete endpoint to make it ready. The registered file has source=upload. Uploads are capped per file and per org by total size.
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
request := &platformgo.FileCreateRequest{
        Filename: "filename",
        MimeType: "mime_type",
        SizeBytes: int64(1000000),
    }
client.Files.Create(
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

**filename:** `string` — Display name for the file (also its name on the phone).
    
</dd>
</dl>

<dl>
<dd>

**mimeType:** `string` — MIME type of the upload; must be an allowed image or video type.
    
</dd>
</dl>

<dl>
<dd>

**sizeBytes:** `int64` — Exact size of the upload in bytes, up to 1 GiB; the presigned URL pins it.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Files.Delete(FileID) -> *platformgo.DeleteFileOutputBody</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Removes a file from the org's library and everywhere it was delivered: the stored object and the library entry go immediately, and every phone holding a copy is scheduled to remove it (removal is confirmed per phone and retried until it lands). The response reports how many phones that recall reaches. This runs the same for an uploaded or a captured file; a capture's source phone keeps its own session copy, which belongs to the session, not the library.
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
request := &platformgo.FilesDeleteRequest{
        FileID: "file_id",
    }
client.Files.Delete(
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

**fileID:** `string` — file identifier to delete
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Files.Rename(FileID, request) -> *platformgo.RenameFileOutputBody</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates the file's display name. Metadata only: storage is keyed by id, so the object never moves and existing URLs keep working, and past deliveries keep the name they were sent under. Only an uploaded file can be renamed; a captured file's name is part of its provenance, so renaming one is rejected.
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
request := &platformgo.FileRenameRequest{
        FileID: "file_id",
        Filename: "filename",
    }
client.Files.Rename(
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

**fileID:** `string` — file identifier to rename
    
</dd>
</dl>

<dl>
<dd>

**filename:** `string` — New display name for the file.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Files.Complete(FileID) -> *platformgo.CompleteFileOutputBody</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Call after PUTting the bytes to the upload URL. Verifies the object landed at the declared size and type, checks the content really is the media it claims to be, and moves the file to ready so it can be delivered. Idempotent: completing an already-ready file just returns it. Applies to source=upload files; a captured file finalizes on its own from the phone's report.
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
request := &platformgo.FilesCompleteRequest{
        FileID: "file_id",
    }
client.Files.Complete(
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

**fileID:** `string` — file identifier to finalize
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Files.PhonesSessionFiles(SessionID) -> *platformgo.FileListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the files this session captured off its phone, newest first — the direct answer to "what did this session capture". Every row is source=capture. Rows appear at detection, before the bytes finish moving, so a caller waiting on a file watches its capture state progress rather than an empty list.
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
request := &platformgo.PhonesSessionFilesRequest{
        SessionID: "session_id",
    }
client.Files.PhonesSessionFiles(
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

**sessionID:** `string` — session whose captures to list
    
</dd>
</dl>

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

**q:** `*string` — filter by filename, case-insensitive substring match
    
</dd>
</dl>

<dl>
<dd>

**mimeType:** `*string` — only files of exactly this media type
    
</dd>
</dl>

<dl>
<dd>

**minSizeBytes:** `*int64` — only files at least this many bytes (0 = no bound)
    
</dd>
</dl>

<dl>
<dd>

**maxSizeBytes:** `*int64` — only files at most this many bytes (0 = no bound)
    
</dd>
</dl>

<dl>
<dd>

**createdAfter:** `*time.Time` — only files registered at or after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `*time.Time` — only files registered at or before this time (RFC 3339)
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Organization
<details><summary><code>client.Organization.Get() -> *platformgo.OrganizationDescriptorResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a descriptor of the organization the API key belongs to: name, slug, id, active plan name. Organization management (rename, members, invitations) is deliberately not part of the public API: it is dashboard-only so an API key can never modify the organization that issued it.
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
client.Organization.Get(
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

<details><summary><code>client.Organization.ListInvitations() -> *platformgo.OrganizationInvitationSummaryListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns every outstanding invitation for the caller's organization. Requires the admin role, matching the dashboard. Organization management (rename, members, invitations) is deliberately not part of the public API: it is dashboard-only so an API key can never modify the organization that issued it.
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
client.Organization.ListInvitations(
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

<details><summary><code>client.Organization.ListMembers() -> *platformgo.OrganizationMemberSummaryListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns every member of the caller's organization with email, role, and join date. Organization management (rename, members, invitations) is deliberately not part of the public API: it is dashboard-only so an API key can never modify the organization that issued it.
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
client.Organization.ListMembers(
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
<details><summary><code>client.Phones.List() -> *platformgo.PhonePrivateListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the caller org's full dedicated (private/rented) phone inventory - every state, not just the free ones: busy phones in an active session, offline/inactive phones, and phones in maintenance are all included, so this is the endpoint to discover a phone_id you can pin via POST /phones:allocate. Each phone's current_session_id and status reflect its live state. include_expired=true also keeps rentals past their rental_expires_at so users can see what they used to own. Filter by status/phone_type and paginate; the response total is the full match count. Pass ownership=dedicated.
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
request := &platformgo.PhonesListRequest{
        Ownership: platformgo.PhonesListRequestOwnershipDedicated,
    }
client.Phones.List(
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

**ownership:** `*platformgo.PhonesListRequestOwnership` — Which phones to return. 'dedicated' = the org's dedicated/rented inventory in every state (busy, offline, and maintenance phones included).
    
</dd>
</dl>

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

**status:** `[]string` — filter by phone status (active/inactive/maintenance/suspended); case-insensitive
    
</dd>
</dl>

<dl>
<dd>

**type_:** `[]string` — filter by phone type (iphone/android); case-insensitive
    
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

<details><summary><code>client.Phones.Availability() -> *platformgo.PhoneAvailabilityResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a capacity snapshot to check before POST /phones:allocate: shared-pool availability broken down by phone type and by location, plus the caller org's dedicated phones with how many are idle (claimable right now). Optional phone_type and location filters narrow every count. Advisory only - availability can change between this read and an allocate, so allocation remains the authority and can still refuse.
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
request := &platformgo.PhonesAvailabilityRequest{}
client.Phones.Availability(
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

**phoneType:** `*platformgo.PhonesAvailabilityRequestPhoneType` — Only count phones of this platform.
    
</dd>
</dl>

<dl>
<dd>

**location:** `*string` — Only count phones at this location slug.
    
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

**sort:** `*string` — sort column: started|ended|status|duration|source (default started)
    
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

<details><summary><code>client.Phones.SessionLiveViewToken(SessionID) -> *platformgo.PhoneLiveViewTokenResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Mints a fresh live-view token (and hosted viewer URL) for an active session, per the session's allocate-time live_view settings: token-mode sessions re-issue the bearer link allocate returned once, org-mode sessions exchange the caller's identity. Refused for sessions allocated with live_view.disabled, for inactive sessions, and for sessions outside the caller's organization (reads as not found). Sharp edge: the returned URL embeds the token and is a bearer capability — whoever holds it can watch (and, unless the session was allocated view-only, drive) the phone until the session ends, at which point both stop working.
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
request := &platformgo.PhonesSessionLiveViewTokenRequest{
        SessionID: "session_id",
    }
client.Phones.SessionLiveViewToken(
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

<details><summary><code>client.Phones.SessionTelemetryToken(SessionID) -> *platformgo.PhoneTelemetryTokenResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Mints a fresh telemetry_url for an active session — the same read-only WebSocket URL POST /phones:allocate returns once (live trace spans + output logs for exactly this session, 3h token). This is the telemetry/frames leg, distinct from the live-view (video) token: it can only watch the trace, never the screen, and never drive the phone. The stream's end frame is the session-end signal. Refused for inactive sessions — an ended session's telemetry is served by GET /phones/sessions/{session_id}/frames — and for sessions outside the caller's organization (reads as not found).
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
request := &platformgo.PhonesSessionTelemetryTokenRequest{
        SessionID: "session_id",
    }
client.Phones.SessionTelemetryToken(
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

<details><summary><code>client.Phones.ListDeliveries(PhoneID) -> *platformgo.FileDeliveryListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the phone's file delivery records, newest first: which library files were sent to it and where each stands (dispatching / dispatched / delivered / failed). Org-scoped: another org's phone reads as not found.
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
request := &platformgo.PhonesListDeliveriesRequest{
        PhoneID: "phone_id",
    }
client.Phones.ListDeliveries(
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

**phoneID:** `string` — phone to list deliveries for
    
</dd>
</dl>

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

<details><summary><code>client.Phones.CreateDelivery(PhoneID, request) -> *platformgo.FilePushResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Sends a library file to a phone the caller's org holds: the phone downloads it over its own connection and inserts it into the media gallery, where app pickers can select it. Accepts any file by id regardless of source (uploaded or captured), and the file must already be ready - finish an uploaded file with POST /files/{file_id}/complete before delivering it. Returns 202 with the delivery record once the phone acknowledges the download started; watch GET /phones/{phone_id}/deliveries or the live preview for completion. Optionally choose the target collection (DCIM / Pictures / Movies).
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
request := &platformgo.FileDeliveryCreateRequest{
        PhoneID: "phone_id",
        FileID: "file_id",
    }
client.Phones.CreateDelivery(
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

**phoneID:** `string` — target phone_id
    
</dd>
</dl>

<dl>
<dd>

**collection:** `*platformgo.FileDeliveryCreateRequestCollection` — Media collection to insert into on the phone; defaults to Pictures for images and Movies for videos.
    
</dd>
</dl>

<dl>
<dd>

**fileID:** `string` — Library file to deliver; accepts any file id regardless of source.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.GetDelivery(PhoneID, DeliveryID) -> *platformgo.FileDeliverySummary</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a single delivery by id and its current status. Poll this to wait on a specific push: the list endpoint pages the newest records and can drop a delivery that ages past the page on a busy phone.
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
request := &platformgo.PhonesGetDeliveryRequest{
        PhoneID: "phone_id",
        DeliveryID: "delivery_id",
    }
client.Phones.GetDelivery(
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

**phoneID:** `string` — phone the delivery belongs to
    
</dd>
</dl>

<dl>
<dd>

**deliveryID:** `string` — delivery to fetch
    
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

<details><summary><code>client.Phones.Preview(PhoneID) -> *platformgo.PhonePreviewResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns a short-lived URL for the phone's current screen preview — a rolling JPEG refreshed every few seconds while the phone is paired, available with or without an active session. Poll this endpoint and swap the image; every call mints a fresh URL. Status is "pending" when no preview exists yet. Authorized to the org that owns the phone or currently holds its active session; any other org reads as not found.
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
request := &platformgo.PhonesPreviewRequest{
        PhoneID: "phone_id",
    }
client.Phones.Preview(
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

**phoneID:** `string` — Phone identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Phones.Deallocate(PhoneID) -> *platformgo.PhoneDeallocateResponse</code></summary>
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

**phoneID:** `string` — phone identifier to deallocate
    
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

<details><summary><code>client.Phones.Allocate(request) -> *platformgo.PhoneAllocateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Allocates an Android phone and opens a session. Omit workflow_id for an interactive lease (drive the phone directly); set it to allocate for a workflow. Pass phone_id to pin a specific dedicated phone. If allocation setup fails the claim is rolled back, so you are never billed for a session that never starts.
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

**capture:** `*bool` — Capture media this session produces on the phone into the org's file library (default true). false disables capture for this session entirely.
    
</dd>
</dl>

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

**pool:** `*platformgo.PhoneAllocateRequestPool` — Which pool to draw the phone from. Omit for shared. 'dedicated' claims any idle phone your organization rents; combine with phone_id to pin a specific one.
    
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

**telemetry:** `*bool` — Emit this session's telemetry (default true). false suppresses telemetry entirely - no live trace stream and no durable trace store; the session dashboard shows a 'telemetry disabled' state.
    
</dd>
</dl>

<dl>
<dd>

**ttl:** `*platformgo.PhoneSessionTTLOptions` — Idle timeout for this session. Omit for no idle timeout: the session runs until the 1-hour max-session cap.
    
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

## Runs
<details><summary><code>client.Runs.SessionsListFrames(SessionID) -> *platformgo.RunSessionFramesResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the paginated telemetry frames for a session in the canonical frame envelope: one completed span frame per durable span (operations that never completed appear via their synthesized failed closures) plus log frames, ordered by span start / log time, with response-level billed-cost maps. This is the same envelope the live telemetry WebSocket streams; live and archive differ only in cardinality (start+end frames live, one completed frame here). Org-scoped: another org's session reads as not found. A trace past the organization's telemetry retention window returns an empty list with retention_expired=true; when the retention policy itself cannot be resolved the request fails with a 500 rather than serving frames whose retention state is unknown. Tolerant reader (unified frame contract): consumers MUST ignore frames with an unknown kind, unknown fields within known kinds, and unknown span_type/log_type values (render generically, never error). Generated SDK types surface an unrecognized frame as an explicit UnknownFrame variant carrying the raw JSON, never a silent drop. A live-stream message MAY carry a JSON array of frame objects; consumers MUST accept a single object or an array.
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
request := &platformgo.SessionsListFramesRequest{
        SessionID: "session_id",
    }
client.Runs.SessionsListFrames(
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

**sessionID:** `string` — Session whose frames to return.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Maximum number of frames to return (1-1000).
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — Pagination offset.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.List() -> *platformgo.RunListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns paginated recent (non-archived) runs the caller started - scoped to their own user within the org, not every member's runs. Filters: workflow_id, search (run ID substring), status_filter, trigger_filter. Order with order_by ('<field> <asc|desc>'), field one of status, started_at, completed_at, created_at.
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
request := &platformgo.RunsListRequest{}
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

**workflowID:** `*string` — Filter results to a single workflow.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Maximum number of runs to return per page (1-500).
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — Pagination offset.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filter by run id substring.
    
</dd>
</dl>

<dl>
<dd>

**statusFilter:** `[]*platformgo.RunsListRequestStatusFilterItem` — Restrict results to the given run statuses.
    
</dd>
</dl>

<dl>
<dd>

**triggerFilter:** `[]*platformgo.RunsListRequestTriggerFilterItem` — Restrict results to the given triggers.
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> <asc|desc>'; field is one of status, started_at, completed_at, created_at. Defaults to created_at desc.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.ListHistoric() -> *platformgo.RunHistoryResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns paginated historic runs for the caller's user over a required time window (start_date/end_date). Use GET /runs for recent (non-archived) runs.
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
request := &platformgo.RunsListHistoricRequest{
        StartDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        EndDate: platformgo.MustParseDateTime(
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

**startDate:** `time.Time` — Beginning of the query time window (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**endDate:** `time.Time` — End of the query time window (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**workflowID:** `*string` — Filter results to a single workflow.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Maximum number of runs to return (1-500).
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — Pagination offset.
    
</dd>
</dl>

<dl>
<dd>

**statusFilter:** `[]*platformgo.RunsListHistoricRequestStatusFilterItem` — Restrict results to the given run statuses (case-insensitive).
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filter by run id or workflow id substring.
    
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

<details><summary><code>client.Runs.Cancel(RunID) -> *platformgo.RunResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Cancels a run that is still queued or running, scoped to the caller's org. A run that has already reached a terminal state (completed/failed/cancelled) cannot be cancelled and reads as not found. Returns the updated run.
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
        Runs: []*platformgo.RunConfig{
            &platformgo.RunConfig{},
        },
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

**runs:** `[]*platformgo.RunConfig` — Per-run variable configurations. One run is created per entry; 1-1000 entries per request.
    
</dd>
</dl>

<dl>
<dd>

**startTimeoutSeconds:** `*int64` — How long a queued run may wait for a phone before it is auto-cancelled (60-86400). Defaults to 300.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Skill
<details><summary><code>client.Skill.GetSkill() -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

The canonical Axilio agent skill: markdown instructions that teach a coding agent to drive a phone through the axilio CLI, then hand back the equivalent SDK script. Single source of truth — the dashboard's "Get agent prompt" action and `axilio init` both fetch this instead of carrying their own copy, so a change here reaches both without a separate deploy.
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
client.Skill.GetSkill(
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

## Usage
<details><summary><code>client.Usage.ListInferences() -> *platformgo.UsageInferencesResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Paginated, filterable list of inference calls (detect + locate) the caller's user was billed for over a required time window (start_date/end_date). Filters: endpoint, model, session, free-text search. Order with order_by ('<field> <asc|desc>').
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
request := &platformgo.UsageListInferencesRequest{
        StartDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        EndDate: platformgo.MustParseDateTime(
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

**startDate:** `time.Time` — Beginning of the inferences query window (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**endDate:** `time.Time` — End of the inferences query window (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Number of inferences per page (1-100).
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — Pagination offset.
    
</dd>
</dl>

<dl>
<dd>

**endpointFilter:** `[]string` — Restrict results to the given vision endpoints ('detect'/'locate').
    
</dd>
</dl>

<dl>
<dd>

**model:** `*string` — Restrict results to a single model name.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filter by inference (event) id substring.
    
</dd>
</dl>

<dl>
<dd>

**sessionID:** `*string` — Restrict results to inferences that ran under one phone session.
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> <asc|desc>'; field one of created_at, cost_microdollars, latency_ms, endpoint, model, inference_id. Defaults to created_at desc.
    
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

Returns infrastructure cost and compute-minute summaries for the caller's user over a date range, plus per-bucket chart data. Granularity is hourly (≤24h window) or daily. Pass the window and granularity as query params.
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

<details><summary><code>client.Usage.ListSessions() -> *platformgo.UsageSessionsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Paginated, filterable list of the caller's phone sessions over a required time window (start_date/end_date), each row carrying billing detail: phone-time cost, inference cost, and combined total in microdollars, billing processing status, and allocation source. Filters: session status, billing processing status, workflow, allocation source, free-text search. Order with order_by ('<field> <asc|desc>').
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
request := &platformgo.UsageListSessionsRequest{
        StartDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        EndDate: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
    }
client.Usage.ListSessions(
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

**startDate:** `time.Time` — Beginning of the sessions query window (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**endDate:** `time.Time` — End of the sessions query window (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Number of sessions per page (1-100).
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — Pagination offset.
    
</dd>
</dl>

<dl>
<dd>

**sessionStatusFilter:** `[]string` — Restrict results to the given session lifecycle statuses.
    
</dd>
</dl>

<dl>
<dd>

**processedStatusFilter:** `[]string` — Restrict results to the given billing processing statuses.
    
</dd>
</dl>

<dl>
<dd>

**workflowID:** `*string` — Restrict results to sessions of a single workflow.
    
</dd>
</dl>

<dl>
<dd>

**allocatedBy:** `[]string` — Restrict results to the given allocation sources.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filter by session, workflow, or phone id substring.
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> <asc|desc>'; field one of allocated_at, deallocated_at, duration, cost_microdollars, session_status, processed_status, allocated_by, session_id, workflow_id. Defaults to allocated_at desc.
    
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

**search:** `*string` — free-text search across workflow name or ID substring
    
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

**capture:** `*bool` — Capture media this workflow's runs produce on the phone into the org's file library (default true). false disables capture for every run dispatched through the scheduler.
    
</dd>
</dl>

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

**ocrEngine:** `*platformgo.WorkflowCreateRequestOcrEngine` — OCR backend to use.
    
</dd>
</dl>

<dl>
<dd>

**platform:** `*platformgo.WorkflowCreateRequestPlatform` — Target OS platform.
    
</dd>
</dl>

<dl>
<dd>

**recording:** `*bool` — Record this workflow's runs (default true). false suppresses video recording and the rolling thumbnail entirely, for every run dispatched through the scheduler.
    
</dd>
</dl>

<dl>
<dd>

**telemetry:** `*bool` — Persist telemetry spans for this workflow's runs (default true). false skips the durable trace store; the live telemetry stream still works while a run is active.
    
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

**capture:** `*bool` — Capture media this workflow's runs produce on the phone into the org's file library (default true). false disables capture for every run dispatched through the scheduler.
    
</dd>
</dl>

<dl>
<dd>

**name:** `*string` — Updated workflow name.
    
</dd>
</dl>

<dl>
<dd>

**ocrEngine:** `*platformgo.WorkflowUpdateRequestOcrEngine` — Updated OCR backend selection.
    
</dd>
</dl>

<dl>
<dd>

**platform:** `*platformgo.WorkflowUpdateRequestPlatform` — Updated target platform.
    
</dd>
</dl>

<dl>
<dd>

**recording:** `*bool` — Record this workflow's runs (default true). false suppresses video recording and the rolling thumbnail entirely, for every run dispatched through the scheduler.
    
</dd>
</dl>

<dl>
<dd>

**status:** `*platformgo.WorkflowUpdateRequestStatus` — Updated lifecycle status.
    
</dd>
</dl>

<dl>
<dd>

**telemetry:** `*bool` — Persist telemetry spans for this workflow's runs (default true). false skips the durable trace store; the live telemetry stream still works while a run is active.
    
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

