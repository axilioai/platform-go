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

Lists the API keys for the caller's organization. Filter by search (name or key preview), created_by_id and a created_after/created_before window; order with order_by ('<field> <asc|desc>'), field one of name, created_by_id, created_at or last_used_at. Defaults to created_at desc.
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

<dl>
<dd>

**search:** `*string` — only keys whose name or key preview contains this text
    
</dd>
</dl>

<dl>
<dd>

**createdByID:** `[]string` — only keys created by one of these users; repeatable
    
</dd>
</dl>

<dl>
<dd>

**createdAfter:** `*time.Time` — only keys created at or after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `*time.Time` — only keys created at or before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is name, created_by_id, created_at or last_used_at. Defaults to created_at desc.
    
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
request := &platformgo.SubscriptionAutoRechargeSettingsUpdate{}
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

**enabled:** `*bool` — Whether auto-recharge is active. Enabling clears any prior payment-failure auto-disable.
    
</dd>
</dl>

<dl>
<dd>

**targetMicrodollars:** `*int64` — Balance the recharge restores to, in microdollars (minimum 5_000_000 = $5.00 while enabled).
    
</dd>
</dl>

<dl>
<dd>

**thresholdMicrodollars:** `*int64` — Recharge when the balance drops below this amount, in microdollars.
    
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

<details><summary><code>client.Billing.GetBalanceStatement() -> *platformgo.SubscriptionBalanceStatement</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

The organization's balance over a period, like a bank statement: the starting balance, credit added, what each product used, credit taken back by refunds and chargebacks, credit that expired, and the ending balance.
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
request := &platformgo.BillingGetBalanceStatementRequest{
        PeriodStartAt: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        PeriodEndAt: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
    }
client.Billing.GetBalanceStatement(
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

**periodStartAt:** `time.Time` — start of the period (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**periodEndAt:** `time.Time` — end of the period, exclusive (RFC 3339)
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Billing.ListBalanceTransactions() -> *platformgo.SubscriptionBalanceTransactionsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

The organization's balance history: every movement of its balance (usage, credit added or granted, expiry, refunds and disputes), newest first, with optional filters.
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
request := &platformgo.BillingListBalanceTransactionsRequest{}
client.Billing.ListBalanceTransactions(
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

**type_:** `[]*platformgo.BillingListBalanceTransactionsRequestTypeItem` — Only transactions of one of these types; repeatable.
    
</dd>
</dl>

<dl>
<dd>

**createdAfter:** `*time.Time` — only transactions at/after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `*time.Time` — only transactions at/before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression 'created_at [asc|desc]', direction defaulting to desc.
    
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

**billingCycle:** `*platformgo.BillingGetHistoryRequestBillingCycle` — filter by billing cycle
    
</dd>
</dl>

<dl>
<dd>

**status:** `*platformgo.BillingGetHistoryRequestStatus` — only invoices at this status
    
</dd>
</dl>

<dl>
<dd>

**invoicedAfter:** `*time.Time` — only invoices created at/after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**invoicedBefore:** `*time.Time` — only invoices created at/before this time (RFC 3339; with invoiced_after, an exact calendar month)
    
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

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is invoiced_at, amount_microdollars, status, plan_name, invoice_number or paid_at.
    
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

Returns one page of the caller organization's phone rental subscriptions, in any lifecycle state unless status filters it (canceled rentals stay listed so past rentals remain visible). Filter by search, status, plan_interval, cancel_at_period_end, payment_status and a current_period_end_after/current_period_end_before window; order with order_by ('<field> <asc|desc>'). The summary block (active_count, combined monthly cost, per-interval breakdown, upcoming charges) covers every active rental and ignores the filters.
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
request := &platformgo.BillingGetRentalSubscriptionsRequest{}
client.Billing.GetRentalSubscriptions(
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

**search:** `*string` — only rentals whose phone name or nickname contains this text
    
</dd>
</dl>

<dl>
<dd>

**status:** `[]*platformgo.BillingGetRentalSubscriptionsRequestStatusItem` — only rentals at one of these lifecycle statuses; repeatable. Omit for every status.
    
</dd>
</dl>

<dl>
<dd>

**planInterval:** `[]*platformgo.BillingGetRentalSubscriptionsRequestPlanIntervalItem` — only rentals on one of these plan intervals; repeatable
    
</dd>
</dl>

<dl>
<dd>

**cancelAtPeriodEnd:** `*platformgo.BillingGetRentalSubscriptionsRequestCancelAtPeriodEnd` — true for rentals set to cancel at period end, false for renewing ones
    
</dd>
</dl>

<dl>
<dd>

**paymentStatus:** `[]*platformgo.BillingGetRentalSubscriptionsRequestPaymentStatusItem` — only rentals in one of these billing states; repeatable
    
</dd>
</dl>

<dl>
<dd>

**currentPeriodEndAfter:** `*time.Time` — only rentals whose current period ends at or after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**currentPeriodEndBefore:** `*time.Time` — only rentals whose current period ends at or before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is phone_name, current_period_end_at, plan_interval or plan_price_microdollars. Defaults to current_period_end_at asc.
    
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
request := &platformgo.SubscriptionUsageAlertSettingsUpdate{}
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

**enabled:** `*bool` — Whether low-balance alerts are active.
    
</dd>
</dl>

<dl>
<dd>

**thresholdMicrodollars:** `*int64` — Alert while the balance sits below this amount, in microdollars. Must be positive while enabled.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Conversations
<details><summary><code>client.Conversations.List() -> *platformgo.ConversationListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one page of your dedicated phones' conversations, newest activity first: one per phone and set of people it texts with, each with its participants and latest message. Filter by phone (phone_id, or phone_number, the phone's own number), and search the participants with search, which matches any part of a person's number with formatting ignored. A text sent to a number and the replies from it are one conversation; the same number on two of your phones is two. Filter by what became of your texts with status, error_code and replied, optionally within a period (period_start_at, period_end_at): they look only at texts you sent in that period (the 7 days before period_end_at when the start is unset, the same window rules as GET /messages/metrics), combine with each other and with the other filters, and list a conversation when it passes every one. status and error_code match a conversation with at least one such text, so they list conversations, not texts: GET /messages/metrics counts both, in totals (texts) and conversations.by_outcome and failures[].conversations (conversations). With only a period, it lists the conversations you texted in it.
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
request := &platformgo.ConversationsListRequest{}
client.Conversations.List(
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

**phoneID:** `*string` — only conversations on this phone
    
</dd>
</dl>

<dl>
<dd>

**phoneNumber:** `*string` — only conversations on the phone with this own number, E.164 (for example +16282788410)
    
</dd>
</dl>

<dl>
<dd>

**status:** `*platformgo.ConversationsListRequestStatus` — only conversations with at least one text of yours, sent in the period, in this status: the same values as a message's status. sent means it left the phone but no delivery report came. A conversation counts once however many of its texts match. More statuses may be added.
    
</dd>
</dl>

<dl>
<dd>

**errorCode:** `*platformgo.ConversationsListRequestErrorCode` — only conversations with at least one text of yours, sent in the period, that failed for this reason: the same codes as a message's error_code. A conversation counts once however many of its texts match. More codes may be added.
    
</dd>
</dl>

<dl>
<dd>

**replied:** `*platformgo.ConversationsListRequestReplied` — true: only conversations you texted in the period where they wrote back within 24 hours of one of those texts (the reply may land after the period ends). false: only the ones you texted that they did not answer that way. The same rule as replied in GET /messages/metrics, so its texted minus replied is the number false lists.
    
</dd>
</dl>

<dl>
<dd>

**periodStartAt:** `*time.Time` — start of the period (RFC 3339) the status, error_code and replied filters look at your texts in. Unset: 7 days before period_end_at. It may be at most a year before now.
    
</dd>
</dl>

<dl>
<dd>

**periodEndAt:** `*time.Time` — end of the period (RFC 3339). Unset: now, and a later time is read as now. Must be after period_start_at.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — free-text search across the conversation's participants: a conversation matches when someone's handle contains this. For a number the digits are compared with formatting ignored, so +1 (415) 555-0123, 4155550123 and 0123 all match +14155550123; other text matches a short code or an alphanumeric sender, case-insensitive.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Conversations.Get(ConversationID) -> *platformgo.Conversation</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one conversation: its phone, its participants and its latest message. Another organization's conversation reads as not found.
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
request := &platformgo.ConversationsGetRequest{
        ConversationID: "conversation_id",
    }
client.Conversations.Get(
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

**conversationID:** `string` — conversation identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Conversations.Delete(ConversationID) -> *platformgo.ConversationDeleted</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes every message in a conversation, as DELETE /messages/{message_id} does for one: the conversation and its messages are removed from every list at once, attachments within 48 hours, and the phone is asked to delete the conversation's messages on its side up to the newest message deleted here (anything it holds that is newer, including texts it has not passed on yet, stays). The response returns at once, with phone_deletion.status pending while the phone has not answered; Axilio keeps checking for up to 7 days, asking the phone again until it confirms (phone_deletion.next_check_at is the next check), and a phone not yet set up to delete messages (not_default_app) is asked once it can. GET /messages/{message_id} on any of the deleted messages shows where the phone's copy stands, and each message fires message.deleted, then message.phone_deleted or message.phone_delete_failed, as DELETE /messages/{message_id} describes. Deleting a conversation does not stop a message in it the phone is still sending. A conversation with no messages left is a 404.
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
request := &platformgo.ConversationsDeleteRequest{
        ConversationID: "conversation_id",
    }
client.Conversations.Delete(
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

**conversationID:** `string` — conversation identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Conversations.ListMessages(ConversationID) -> *platformgo.PhoneMessageListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one page of a conversation's messages, both directions, newest first by default (order_by created_at asc reads it oldest first, as a thread shows it). Each message has its full body.
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
request := &platformgo.ConversationsListMessagesRequest{
        ConversationID: "conversation_id",
    }
client.Conversations.ListMessages(
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

**conversationID:** `string` — conversation identifier
    
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

**orderBy:** `*string` — Sort expression 'created_at [asc|desc]': when the platform stored each message, newest first by default.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Conversations.CreateMessage(ConversationID, request) -> *platformgo.PhoneMessage</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Sends a message into a conversation, from its phone to its participant: the same send as POST /messages, with the same limits, idempotency and webhooks, and the same 202 response: the message, queued. A conversation with more than one other person, or whose other party is a short code or an unknown sender, can't take a reply yet (422).
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
request := &platformgo.ReplyInputBody{
        ConversationID: "conversation_id",
        IdempotencyKey: "idempotency_key",
    }
client.Conversations.CreateMessage(
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

**conversationID:** `string` — conversation identifier
    
</dd>
</dl>

<dl>
<dd>

**idempotencyKey:** `string` — A caller-chosen key; the same key for your organization sends at most once.
    
</dd>
</dl>

<dl>
<dd>

**parts:** `[]*platformgo.PhoneMessageSendPart` — The message content in order. An SMS is one text part; media parts (at most 2 MiB in all) make it an MMS, with any text part first.
    
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

Returns one page of the org's file library, newest first, with the org's standing usage against its storage quota. Every file carries its source (upload, capture or message) and, for a capture or message attachment, its surface and capture state. Filter with source, surface, session_id, phone_id and media_type, among others. Files persist until deleted and any ready file can be delivered to a phone the org holds.
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

**search:** `*string` — filter by filename, case-insensitive substring match
    
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

**phoneID:** `*string` — only files from this phone: attachments on messages it received, and files its sessions captured
    
</dd>
</dl>

<dl>
<dd>

**mimeType:** `*string` — only files of exactly this media type
    
</dd>
</dl>

<dl>
<dd>

**mediaType:** `*platformgo.FilesListRequestMediaType` — only images or only videos, whatever their exact media type
    
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

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is created_at, filename, size_bytes or source (which groups files by source, newest first within each).
    
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

**sizeBytes:** `int64` — Exact size of the upload in bytes, up to 100 MiB (the phone-delivery ceiling); the presigned URL pins it.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Files.Get(FileID) -> *platformgo.FileSummary</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one file from the org's library, with a short-lived read URL once its bytes are ready. A message's media part names its file by id; this is where that id resolves. Another org's file id reads as not found.
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
request := &platformgo.FilesGetRequest{
        FileID: "file_id",
    }
client.Files.Get(
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

**fileID:** `string` — file identifier
    
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

Removes a file from the org's library: the stored object and the library entry go immediately. Copies already on phones are not removed, whether the file was pushed to them or captured from one; removing a file from a phone is a separate operation. This runs the same for an uploaded or a captured file.
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

<details><summary><code>client.Files.Update(FileID, request) -> *platformgo.RenameFileOutputBody</code></summary>
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
client.Files.Update(
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

**search:** `*string` — filter by filename, case-insensitive substring match
    
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

## Messages
<details><summary><code>client.Messages.List() -> *platformgo.PhoneMessageListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one page of every text message your dedicated phones have received and sent, across conversations, newest first, each with its full body and its conversation_id; direction tells them apart. Filter by phone_id, direction and when the platform stored the message (created_after, created_before), for example to collect everything since your last read. To read a conversation, use GET /conversations/{conversation_id}/messages. Only dedicated phones with a current rental receive messages; texts arriving on shared phones, or after a rental lapses, are never stored.
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
request := &platformgo.MessagesListRequest{}
client.Messages.List(
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

**phoneID:** `*string` — only messages on this phone
    
</dd>
</dl>

<dl>
<dd>

**createdAfter:** `*time.Time` — only messages stored at or after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `*time.Time` — only messages stored at or before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**direction:** `*platformgo.MessagesListRequestDirection` — only received (inbound) or sent (outbound) messages
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc. field is created_at (when the platform stored the message), sent_at (when it was sent) or received_at (when the phone received it); messages without that time come last.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Messages.Create(request) -> *platformgo.PhoneMessage</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Sends a message from one of your dedicated phones: an SMS when parts is one text part, an MMS when it carries media from your file library. The phone must have outbound messaging enabled by your administrator. Returns 202 with the message, queued; the outcome arrives as a message.sent or message.failed webhook and is readable at GET /messages/{message_id}. A phone that never received the message within about 15 minutes (it is offline, say) fails it rather than sending it late; one that received it but has not reported back by then makes it unconfirmed until it does. Each phone sends at most 15 messages per 30 minutes. Reusing an idempotency_key returns the original message, as it is now, instead of sending again; reusing it for a different message is a 409.
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
request := &platformgo.SendMessageInputBody{
        IdempotencyKey: "idempotency_key",
        PhoneID: "phone_id",
    }
client.Messages.Create(
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

**idempotencyKey:** `string` — A caller-chosen key; the same key for your organization sends at most once.
    
</dd>
</dl>

<dl>
<dd>

**parts:** `[]*platformgo.PhoneMessageSendPart` — The message content in order. An SMS is one text part; media parts (at most 2 MiB in all) make it an MMS, with any text part first.
    
</dd>
</dl>

<dl>
<dd>

**phoneID:** `string` — The dedicated phone to send from, by its public id.
    
</dd>
</dl>

<dl>
<dd>

**to:** `[]string` — The recipients, as E.164 numbers such as +14155550123. One recipient for now.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Messages.GetMetrics() -> *platformgo.PhoneMessageMetrics</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns your texting over a window (the last 7 days by default, starting at most a year back): texts sent and received, how the sent ones ended up (delivered, left the phone with no delivery report, unconfirmed, failed, still queued) and the delivery rate, how long they took to send and to be delivered, why they failed, your conversations and how many of the ones you texted wrote back within a day, and each phone's share with its busiest 30 minutes, to compare with the send_limit the response carries. The points chart it by the hour or the day (granularity) in time_zone; hour buckets need a window of 7 days or less. Scope it to one phone with phone_id or one conversation with conversation_id. A text counts at the time the platform stored it when sent, and when the phone received it when received.
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
request := &platformgo.MessagesGetMetricsRequest{}
client.Messages.GetMetrics(
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

**periodStartAt:** `*time.Time` — Start of the window (RFC 3339). Defaults to 7 days before period_end_at, and may be at most a year before now.
    
</dd>
</dl>

<dl>
<dd>

**periodEndAt:** `*time.Time` — End of the window (RFC 3339). Defaults to now, and a later time is read as now.
    
</dd>
</dl>

<dl>
<dd>

**granularity:** `*platformgo.MessagesGetMetricsRequestGranularity` — Bucket size: hour or day. hour is allowed for windows of 7 days or less.
    
</dd>
</dl>

<dl>
<dd>

**timeZone:** `*string` — IANA time zone the buckets follow, such as America/Los_Angeles: a day starts at local midnight and an hour on the local hour. An unknown name is an error.
    
</dd>
</dl>

<dl>
<dd>

**phoneID:** `*string` — Only this phone's texts. A phone your organization doesn't rent is not found.
    
</dd>
</dl>

<dl>
<dd>

**conversationID:** `*string` — Only this conversation's texts, and only its phone. Another organization's conversation is not found. When phone_id is also given, it is ignored: the conversation names the phone.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Messages.Get(MessageID) -> *platformgo.PhoneMessage</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns one text message a dedicated phone received or sent, with its full body. A sent message carries its delivery status.
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
request := &platformgo.MessagesGetRequest{
        MessageID: "message_id",
    }
client.Messages.Get(
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

**messageID:** `string` — message identifier
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Messages.Delete(MessageID) -> *platformgo.PhoneMessage</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes a message: it is removed from every list and conversation at once, its attachments are removed from your file library within 48 hours, and the delete is recorded with who made it. The phone is also asked to delete its own copy, asynchronously. The response returns at once with the deleted message as GET /messages/{message_id} returns it (deleted_at set, no body or parts), with phone_deletion.status pending while the phone has not answered. Axilio keeps checking for up to 7 days, asking the phone again until it confirms; phone_deletion.next_check_at is when it next checks (null once the outcome is final). A phone not yet set up to delete messages (not_default_app) is asked once it can, within those 7 days. GET /messages/{message_id} keeps returning the deleted message and shows where the phone's copy stands in phone_deletion. Webhooks: message.deleted fires at once; message.phone_deleted fires when the phone holds no copy, and message.phone_delete_failed when it is final that the copy is still there or may be (ambiguous, not_matched, not_sent, or unconfirmed after the 7 days). A message the phone stored before it became the phone's messaging app can only be deleted on the phone when it is a plain text the phone can match, and only on the first attempt; otherwise the phone answers not_matched. Deleting a message the phone is still sending does not stop the send. Deleting a message that does not exist or is already deleted is a 404.
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
request := &platformgo.MessagesDeleteRequest{
        MessageID: "message_id",
    }
client.Messages.Delete(
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

**messageID:** `string` — message identifier
    
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

**search:** `*string` — free-text search across phone_id, nickname, name, model, location
    
</dd>
</dl>

<dl>
<dd>

**status:** `[]*platformgo.PhonesListRequestStatusItem` — only phones at one of these statuses; repeatable
    
</dd>
</dl>

<dl>
<dd>

**phoneType:** `[]*platformgo.PhonesListRequestPhoneTypeItem` — only phones of one of these platforms; repeatable
    
</dd>
</dl>

<dl>
<dd>

**rentalExpiresAfter:** `*time.Time` — only phones whose rental expires at/after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**rentalExpiresBefore:** `*time.Time` — only phones whose rental expires at/before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**lastHeartbeatAfter:** `*time.Time` — only phones whose last heartbeat is at/after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**lastHeartbeatBefore:** `*time.Time` — only phones whose last heartbeat is at/before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is created_at, rental_expires_at, last_heartbeat_at, status, phone_type or location.
    
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

**limit:** `*int64` — max rows to return
    
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

**status:** `[]*platformgo.PhonesListSessionsRequestStatusItem` — only sessions at one of these statuses; repeatable
    
</dd>
</dl>

<dl>
<dd>

**source:** `[]*platformgo.PhonesListSessionsRequestSourceItem` — only sessions from one of these sources; repeatable
    
</dd>
</dl>

<dl>
<dd>

**dedicated:** `[]*platformgo.PhonesListSessionsRequestDedicatedItem` — only sessions on shared and/or dedicated phones; repeatable
    
</dd>
</dl>

<dl>
<dd>

**allocatedAfter:** `*time.Time` — only sessions allocated at/after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**allocatedBefore:** `*time.Time` — only sessions allocated at/before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**deallocatedAfter:** `*time.Time` — only sessions de-allocated at/after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**deallocatedBefore:** `*time.Time` — only sessions de-allocated at/before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is allocated_at, deallocated_at, status, duration_seconds or source.
    
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

Returns a short-lived URL for the session's screen recording, keyed on session_id — so it works for workflow runs and workflow-less interactive leases alike. Status is "pending" (no URL) when the recording hasn't finished uploading yet, and "off" when the session was not recorded. Org-scoped: another org's session reads as not found.
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

<details><summary><code>client.Phones.SessionLiveViewToken(SessionID) -> *platformgo.PhoneLiveViewTokenResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Mints a fresh live-view token (and hosted viewer URL) for an active session, per the session's allocate-time live_view settings: token-mode sessions re-issue the bearer link allocate returned once, org-mode sessions exchange the caller's identity. The token is view-only when the session was allocated with live_view.view_only or the caller has the viewer role; only members and admins get a token that can drive the phone. Refused for sessions allocated with live_view.disabled, for inactive sessions, and for sessions outside the caller's organization (reads as not found). Sharp edge: the returned URL embeds the token and is a bearer capability — whoever holds it can watch (and, unless the token is view-only, drive) the phone until the session ends, at which point both stop working.
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

**platform:** `*platformgo.PhonesSupportedAppsRequestPlatform` — only apps for this platform
    
</dd>
</dl>

<dl>
<dd>

**category:** `*platformgo.PhonesSupportedAppsRequestCategory` — only apps in this category
    
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

<details><summary><code>client.Phones.Update(PhoneID, request) -> *platformgo.PhoneSummary</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Changes the fields a customer may set on a dedicated phone the caller's org owns: its nickname. Returns the updated phone.
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
request := &platformgo.PhoneUpdatePhoneRequest{
        PhoneID: "phone_id",
        Nickname: "nickname",
    }
client.Phones.Update(
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

**nickname:** `string` — New display name for the phone.
    
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

Sends a library file to a phone the caller's org holds: the phone downloads it over its own connection and inserts it into the media gallery, where app pickers can select it. Accepts any file by id regardless of source (uploaded or captured), and the file must already be ready - finish an uploaded file with POST /files/{file_id}:complete before delivering it. Returns 202 with the delivery record once the phone acknowledges the download started; watch GET /phones/{phone_id}/deliveries or the live preview for completion. Optionally choose the target collection (DCIM / Pictures / Movies).
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

**collection:** `*platformgo.FileDeliveryCreateRequestCollection` — Media collection to insert into on the phone; defaults to pictures for images and Movies for videos.
    
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

<details><summary><code>client.Phones.Deallocate(PhoneID, request) -> *platformgo.PhoneDeallocateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deallocates a phone the caller's org currently holds. The session is billed and the phone is torn down asynchronously. Pass session_id to end that session only: if the phone has already moved on to another session, nothing is released and the call still succeeds.
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
request := &platformgo.PhoneDeallocateRequest{
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

<dl>
<dd>

**sessionID:** `*string` — End this session only. If the phone has already moved on to another session, nothing is released and the call still succeeds. Omit to end whatever session the phone holds.
    
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

**accessibility:** `*bool` — Accessibility mode for this session, which lets locators resolve against the phone's accessibility tree. Defaults to false. true requires a phone that supports it: only such phones are claimed, and a phone_id that does not support it is refused with an accessibility_unavailable conflict. It also needs your organization to have access to accessibility mode; without it the request is refused with a 403 whose message starts with accessibility_unavailable. While on, the accessibility service is visible to apps on the phone.
    
</dd>
</dl>

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

**recording:** `*bool` — Record this session's screen (default false). true stores a video of the session, available from the session's recording once it ends.
    
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

Returns the paginated telemetry frames for a session in the canonical frame envelope: one completed span frame per durable span (operations that never completed appear via their synthesized failed closures) plus log frames, ordered by span start / log time, with response-level billed-cost maps. This is the same envelope the live telemetry event stream (GET /realtime/sessions/{session_id}/events) carries; live and archive differ only in cardinality (start+end frames live, one completed frame here). Org-scoped: another org's session reads as not found. A trace past the organization's telemetry retention window returns an empty list with retention_expired=true; when the retention policy itself cannot be resolved the request fails with a 500 rather than serving frames whose retention state is unknown. Tolerant reader (unified frame contract): consumers MUST ignore frames with an unknown kind, unknown fields within known kinds, and unknown span_type/log_type values (render generically, never error). Generated SDK types surface an unrecognized frame as an explicit UnknownFrame variant carrying the raw JSON, never a silent drop. A live-stream message MAY carry a JSON array of frame objects; consumers MUST accept a single object or an array.
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

Returns paginated recent (non-archived) runs across the caller's organization, every member's runs included. Filters: workflow_id, search (run ID substring), status, trigger. Order with order_by ('<field> <asc|desc>'), field one of status, started_at, completed_at, created_at.
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

**status:** `[]*platformgo.RunsListRequestStatusItem` — Only runs at one of these statuses; repeatable.
    
</dd>
</dl>

<dl>
<dd>

**trigger:** `[]*platformgo.RunsListRequestTriggerItem` — Only runs started by one of these triggers; repeatable.
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is status, started_at, completed_at or created_at. Defaults to created_at desc.
    
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

Returns paginated historic runs across the caller's organization over a required time window (completed_after/completed_before). Use GET /runs for recent (non-archived) runs.
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
        CompletedAfter: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        CompletedBefore: platformgo.MustParseDateTime(
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

**completedAfter:** `time.Time` — Only runs that finished (completed, failed or were cancelled) at or after this time (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**completedBefore:** `time.Time` — Only runs that finished before this time (RFC 3339).
    
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

**status:** `[]*platformgo.RunsListHistoricRequestStatusItem` — Only runs at one of these statuses; repeatable.
    
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

Cancels a run that is scheduled, queued or running, scoped to the caller's org. A scheduled run is cancelled before it starts and never runs. A run that has already reached a terminal state (completed/failed/cancelled) cannot be cancelled and reads as not found. Returns the updated run.
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

**startAt:** `*time.Time` — Start the runs at this time (RFC 3339) instead of now. Must be at least 60 seconds and at most 30 days ahead. The runs are listed as scheduled until then, can be cancelled, and start within about 15 seconds of this time plus the usual wait for a phone. Omit to start now.
    
</dd>
</dl>

<dl>
<dd>

**startTimeoutSeconds:** `*int64` — How long a queued run may wait for a phone before it is auto-cancelled (60-86400). Defaults to 300. For a scheduled run the wait counts from when it is released at start_at, not from booking.
    
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

Paginated, filterable list of inference calls (detect, locate and element_locate) the caller's user was billed for over a required time window (created_after/created_before). Filters: endpoint, model, session, free-text search. Order with order_by ('<field> <asc|desc>').
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
        CreatedAfter: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        CreatedBefore: platformgo.MustParseDateTime(
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

**createdAfter:** `time.Time` — Only inferences created at or after this time (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `time.Time` — Only inferences created before this time (RFC 3339).
    
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

**endpoint:** `[]*platformgo.UsageListInferencesRequestEndpointItem` — Only inferences of one of these operations; repeatable.
    
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

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is created_at, cost_microdollars, latency_ms, endpoint, model or inference_id. Defaults to created_at desc.
    
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

Returns your organization's usage spend over a window as points: one per hour or day in time_zone (and per product with group_by=product), each with the total, phone-session and inference spend in microdollars, zero when nothing was billed. The window defaults to the last 7 days and may start at most a year back; hour buckets need a window of 7 days or less.
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
request := &platformgo.UsageGetMetricsRequest{}
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

**periodStartAt:** `*time.Time` — Start of the window (RFC 3339). Defaults to 7 days before period_end_at, and may be at most a year before now.
    
</dd>
</dl>

<dl>
<dd>

**periodEndAt:** `*time.Time` — End of the window (RFC 3339). Defaults to now, and a later time is read as now.
    
</dd>
</dl>

<dl>
<dd>

**granularity:** `*platformgo.UsageGetMetricsRequestGranularity` — Bucket size: hour or day. hour is allowed for windows of 7 days or less.
    
</dd>
</dl>

<dl>
<dd>

**timeZone:** `*string` — IANA time zone the buckets follow, such as America/Los_Angeles: a day starts at local midnight and an hour on the local hour. An unknown name is an error.
    
</dd>
</dl>

<dl>
<dd>

**groupBy:** `*platformgo.UsageGetMetricsRequestGroupBy` — Split each bucket by product (session or inference). Omit for one point per bucket.
    
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

Paginated, filterable list of the caller's phone sessions over a required time window (allocated_after/allocated_before), each row carrying billing detail: phone-time cost, inference cost, and combined total in microdollars, billing processing status, and allocation source. Filters: session status, billing processing status, workflow, allocation source, free-text search. Order with order_by ('<field> <asc|desc>').
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
        AllocatedAfter: platformgo.MustParseDateTime(
            "2024-01-15T09:30:00Z",
        ),
        AllocatedBefore: platformgo.MustParseDateTime(
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

**allocatedAfter:** `time.Time` — Only sessions allocated at or after this time (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**allocatedBefore:** `time.Time` — Only sessions allocated before this time (RFC 3339).
    
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

**sessionStatus:** `[]*platformgo.UsageListSessionsRequestSessionStatusItem` — Only sessions at one of these lifecycle statuses; repeatable.
    
</dd>
</dl>

<dl>
<dd>

**processedStatus:** `[]*platformgo.UsageListSessionsRequestProcessedStatusItem` — Only sessions at one of these billing processing statuses; repeatable.
    
</dd>
</dl>

<dl>
<dd>

**workflowID:** `*string` — Restrict results to sessions of a single workflow.
    
</dd>
</dl>

<dl>
<dd>

**allocatedBy:** `[]*platformgo.UsageListSessionsRequestAllocatedByItem` — Restrict results to the given allocation sources; repeatable.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filter by session, workflow, or phone id substring.
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is allocated_at, deallocated_at, duration_seconds, cost_microdollars, session_status, processed_status, allocated_by, session_id or workflow_id. Defaults to allocated_at desc.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Webhooks
<details><summary><code>client.Webhooks.ListDeliveries() -> *platformgo.WebhookDeliveryListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists events published for the caller's organization, newest first, one page at a time. Narrow it to one event type (event_type), one destination (endpoint_id) or one message's events (message_id). Pass the response's next_cursor as cursor to fetch the following page.
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
request := &platformgo.WebhooksListDeliveriesRequest{}
client.Webhooks.ListDeliveries(
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

**limit:** `*int64` — Page size.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque cursor from a previous page's next_cursor. Omit for the first page.
    
</dd>
</dl>

<dl>
<dd>

**eventType:** `*string` — Restrict the log to one event type. Omit for all types.
    
</dd>
</dl>

<dl>
<dd>

**endpointID:** `*string` — Restrict the log to deliveries to one destination (endpoint id). Omit for all destinations.
    
</dd>
</dl>

<dl>
<dd>

**messageID:** `*string` — Restrict the log to the events for one message: message.received for a received message, the status events (message.sent, message.delivered, message.failed, message.unconfirmed) for a sent one, and the deletion events (message.deleted, message.phone_deleted, message.phone_delete_failed) for a deleted one. Omit for every event.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.GetDelivery(EventID) -> *platformgo.WebhookDeliveryDetail</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches one event from the delivery log with the payload that was delivered and the endpoints it was routed to. A missing or wrong-organization id returns 404.
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
request := &platformgo.WebhooksGetDeliveryRequest{
        EventID: "event_id",
    }
client.Webhooks.GetDelivery(
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

**eventID:** `string` — Event identifier.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.ListAttempts(EventID) -> *platformgo.WebhookAttemptListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists the delivery attempts for one event across the organization's endpoints, newest first, including automatic retries and manual replays. Pages the same way as the delivery log.
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
request := &platformgo.WebhooksListAttemptsRequest{
        EventID: "event_id",
    }
client.Webhooks.ListAttempts(
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

**eventID:** `string` — Event identifier.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Page size.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque cursor from a previous page's next_cursor. Omit for the first page.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.ReplayDelivery(EventID, request) -> *platformgo.WebhookReplayResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Re-drives one event to one endpoint. Used to recover a dead-lettered delivery once the endpoint is healthy again. The replay is queued, not delivered synchronously; the response identifies it, and its attempt appears under the event's attempts with manual set.
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
request := &platformgo.WebhookReplayRequest{
        EventID: "event_id",
        EndpointID: "endpoint_id",
    }
client.Webhooks.ReplayDelivery(
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

**eventID:** `string` — Event identifier to replay.
    
</dd>
</dl>

<dl>
<dd>

**endpointID:** `string` — Endpoint to replay the event to.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.SendTestEvent(request) -> *platformgo.WebhookTestEventResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Publishes a simulated event to the caller's organization so an endpoint can be verified end to end. Name an event_type from the catalog to simulate it with that type's sample payload (default webhook.test), and an endpoint_id to deliver to that endpoint only (a missing or wrong-organization id returns 404; an endpoint not subscribed to the type is a validation error rather than a silent no-op). Omit the body to deliver webhook.test to every endpoint subscribed to it. The event is queued, not delivered synchronously; the response carries the event id to look up in the delivery log.
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
request := &platformgo.WebhookTestEventRequest{}
client.Webhooks.SendTestEvent(
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

**endpointID:** `*string` — Deliver the test event to this endpoint only. It must be subscribed to the event type. Omit to deliver to every endpoint subscribed to the type.
    
</dd>
</dl>

<dl>
<dd>

**eventType:** `*string` — The event type to simulate, from the event-type catalog; the event carries that type's sample payload. Omit for webhook.test.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.ListEndpoints() -> *platformgo.WebhookEndpointListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists every webhook endpoint registered for the caller's organization. The collection is bounded, so it is returned whole rather than paged. Signing secrets are never included.
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
client.Webhooks.ListEndpoints(
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

<details><summary><code>client.Webhooks.CreateEndpoint(request) -> *platformgo.WebhookEndpointCreateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Registers an HTTPS endpoint to receive signed events for the caller's organization, subscribed to a set of event types (or all). The signing secret is returned exactly once.
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
request := &platformgo.WebhookEndpointCreateRequest{}
client.Webhooks.CreateEndpoint(
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

**config:** `map[string]string` — Non-secret configuration for a BYOS destination (e.g. queue_url + region for aws_sqs).
    
</dd>
</dl>

<dl>
<dd>

**credentials:** `map[string]string` — Customer secrets for a BYOS destination (e.g. AWS key/secret). Never stored or logged by Axilio; passed to the delivery gateway, which encrypts them.
    
</dd>
</dl>

<dl>
<dd>

**eventTypes:** `[]string` — Event types to subscribe to. Empty means all event types.
    
</dd>
</dl>

<dl>
<dd>

**type_:** `*platformgo.WebhookEndpointCreateRequestType` — Destination type: webhook (default when omitted), aws_sqs, aws_s3, or aws_kinesis.
    
</dd>
</dl>

<dl>
<dd>

**url:** `*string` — HTTPS endpoint that signed events are delivered to (webhook destinations).
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.GetEndpoint(EndpointID) -> *platformgo.WebhookEndpoint</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches one webhook endpoint by id. A missing or wrong-organization id returns 404.
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
request := &platformgo.WebhooksGetEndpointRequest{
        EndpointID: "endpoint_id",
    }
client.Webhooks.GetEndpoint(
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

**endpointID:** `string` — Endpoint identifier.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.DeleteEndpoint(EndpointID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Removes a webhook endpoint. Subsequent events are no longer delivered to it. A missing or wrong-organization id returns 404.
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
request := &platformgo.WebhooksDeleteEndpointRequest{
        EndpointID: "endpoint_id",
    }
client.Webhooks.DeleteEndpoint(
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

**endpointID:** `string` — Endpoint identifier.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.ListEndpointAttempts(EndpointID) -> *platformgo.WebhookEndpointAttemptListResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists every delivery attempt made to one endpoint across the organization's events, newest first, each with its outcome, the HTTP status the endpoint returned, its try number and the event it delivered. Filter by outcome (status) or event type. Pages the same way as the delivery log. A missing or wrong-organization endpoint id returns 404.
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
request := &platformgo.WebhooksListEndpointAttemptsRequest{
        EndpointID: "endpoint_id",
    }
client.Webhooks.ListEndpointAttempts(
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

**endpointID:** `string` — Endpoint identifier.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Page size.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Opaque cursor from a previous page's next_cursor. Omit for the first page.
    
</dd>
</dl>

<dl>
<dd>

**status:** `*platformgo.WebhooksListEndpointAttemptsRequestStatus` — Keep only attempts with this outcome. Omit for both.
    
</dd>
</dl>

<dl>
<dd>

**eventType:** `*string` — Keep only attempts of events of this type. Omit for all types.
    
</dd>
</dl>

<dl>
<dd>

**deliveredAfter:** `*time.Time` — Only attempts that completed at or after this time (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**deliveredBefore:** `*time.Time` — Only attempts that completed at or before this time (RFC 3339).
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.RotateEndpointSecret(EndpointID, request) -> *platformgo.WebhookRotateSecretResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Issues a new signing secret for the endpoint, returned exactly once. With mode graceful (the default) the previous secret keeps verifying for 24 hours (previous_secret_invalid_at in the response) so the new secret can be rolled out across a deploy without dropping deliveries in flight. With mode immediate the previous secret stops verifying at once, for a secret that may be compromised. A missing or wrong-organization id returns 404.
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
request := &platformgo.WebhookRotateSecretRequest{
        EndpointID: "endpoint_id",
    }
client.Webhooks.RotateEndpointSecret(
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

**endpointID:** `string` — Endpoint identifier.
    
</dd>
</dl>

<dl>
<dd>

**mode:** `*platformgo.WebhookRotateSecretRequestMode` — graceful: the previous secret keeps verifying for 24 hours so the new one can be rolled out across a deploy. immediate: the previous secret stops verifying at once, for a secret that may be compromised. Omit for graceful.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Webhooks.ListEventTypes() -> *platformgo.WebhookEventTypeCatalogResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists the event types an endpoint can subscribe to, with a description of each, so a caller knows what to subscribe to before registering an endpoint.
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
client.Webhooks.ListEventTypes(
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

<details><summary><code>client.Webhooks.GetMetrics() -> *platformgo.WebhookDeliveryMetricsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns delivery volume, error rate and destination response time over a window as points: one per hour or day in time_zone, for every destination together, one (endpoint_id), or each one (group_by=endpoint_id). Counts are zero and latencies null for a bucket with no attempts. The window defaults to the last 7 days and may start at most a year back; hour buckets need a window of 7 days or less.
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
request := &platformgo.WebhooksGetMetricsRequest{}
client.Webhooks.GetMetrics(
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

**periodStartAt:** `*time.Time` — Start of the window (RFC 3339). Defaults to 7 days before period_end_at, and may be at most a year before now.
    
</dd>
</dl>

<dl>
<dd>

**periodEndAt:** `*time.Time` — End of the window (RFC 3339). Defaults to now, and a later time is read as now.
    
</dd>
</dl>

<dl>
<dd>

**granularity:** `*platformgo.WebhooksGetMetricsRequestGranularity` — Bucket size: hour or day. hour is allowed for windows of 7 days or less.
    
</dd>
</dl>

<dl>
<dd>

**timeZone:** `*string` — IANA time zone the buckets follow, such as America/Los_Angeles: a day starts at local midnight and an hour on the local hour. An unknown name is an error.
    
</dd>
</dl>

<dl>
<dd>

**groupBy:** `*platformgo.WebhooksGetMetricsRequestGroupBy` — Split each bucket by destination, one series per destination that had attempts (the list-view sparklines). Omit for one point per bucket.
    
</dd>
</dl>

<dl>
<dd>

**endpointID:** `*string` — Only this destination's attempts.
    
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

Paginated list of workflows in the caller's org, with optional search, status, platform, and created/last-run date filters via query params. Order with order_by ('<field> <asc|desc>'), field a path into the listed item such as workflow.name.
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

**status:** `[]*platformgo.WorkflowsListRequestStatusItem` — only workflows at one of these statuses; repeatable
    
</dd>
</dl>

<dl>
<dd>

**platform:** `[]*platformgo.WorkflowsListRequestPlatformItem` — only workflows for one of these platforms; repeatable
    
</dd>
</dl>

<dl>
<dd>

**createdAfter:** `*time.Time` — only workflows created at/after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `*time.Time` — only workflows created at/before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**lastRunAfter:** `*time.Time` — only workflows last run at/after this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**lastRunBefore:** `*time.Time` — only workflows last run at/before this time (RFC 3339)
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is workflow.id, workflow.name, workflow.platform, workflow.status, workflow.last_run_at, workflow.created_at or workflow.updated_at. Defaults to workflow.created_at desc.
    
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

**accessibility:** `*bool` — Accessibility mode for this workflow's runs (default false). true requires a phone that supports it for every run dispatched through the scheduler, and your organization's access to accessibility mode: without it the save is refused with a 403 whose message starts with accessibility_unavailable.
    
</dd>
</dl>

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

**recording:** `*bool` — Record the screen of this workflow's runs (default false). true stores a video of every run dispatched through the scheduler.
    
</dd>
</dl>

<dl>
<dd>

**telemetry:** `*bool` — Emit telemetry for this workflow's runs (default true). false suppresses both live frames and durable trace storage.
    
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

**accessibility:** `*bool` — Accessibility mode for this workflow's runs. true requires a phone that supports it for every run dispatched through the scheduler, and your organization's access to accessibility mode: without it the update is refused with a 403 whose message starts with accessibility_unavailable. false turns it off. Omit to leave it unchanged.
    
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

**recording:** `*bool` — Record the screen of this workflow's runs. true stores a video of every run dispatched through the scheduler; false stops recording. Omit to leave it unchanged. New workflows default to false.
    
</dd>
</dl>

<dl>
<dd>

**status:** `*platformgo.WorkflowUpdateRequestStatus` — Updated lifecycle status.
    
</dd>
</dl>

<dl>
<dd>

**telemetry:** `*bool` — Emit telemetry for this workflow's runs (default true). false suppresses both live frames and durable trace storage.
    
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

Returns the source code of the workflow's current revision, or of revision_id when given, scoped to the caller's org. A revision_id that is not one of this workflow's revisions reads as not found.
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

<dl>
<dd>

**revisionID:** `*string` — Revision whose source to return. Defaults to the workflow's current revision.
    
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

<details><summary><code>client.Workflows.CreateRevision(WorkflowID, request) -> *platformgo.WorkflowCreateRevisionResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Saves source as the workflow's new current revision. Hash-deduplicates against the current revision (no_op is true and nothing changes when the source is unchanged). Source is capped at 256KB.
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
request := &platformgo.WorkflowCreateRevisionRequest{
        WorkflowID: "workflow_id",
        Source: "source",
    }
client.Workflows.CreateRevision(
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

<details><summary><code>client.Workflows.RestoreRevision(WorkflowID, RevisionID) -> *platformgo.WorkflowCreateRevisionResponse</code></summary>
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
request := &platformgo.WorkflowsRestoreRevisionRequest{
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

**revisionID:** `string` — revision to restore
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Workflows.GetRunStats(WorkflowID) -> *platformgo.RunStatsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the total run count and success rate of the caller's own runs of the given workflow. A workflow outside the caller's organization reads as not found.
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
request := &platformgo.WorkflowsGetRunStatsRequest{
        WorkflowID: "workflow_id",
    }
client.Workflows.GetRunStats(
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

