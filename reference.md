# Reference
## Agents
<details><summary><code>client.Agents.List() -> *platformgo.AgentList</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists the organization's agents, newest first, one page at a time.
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
request := &platformgo.AgentsListRequest{}
client.Agents.List(
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

**offset:** `*int64` — Where the page starts.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Narrow to agents whose name or id contains this, case-insensitively.
    
</dd>
</dl>

<dl>
<dd>

**kind:** `*platformgo.AgentsListRequestKind` — Narrow to agents whose current revision is instructions or code. Omit for both.
    
</dd>
</dl>

<dl>
<dd>

**createdAfter:** `*time.Time` — Only agents created at or after this time (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**createdBefore:** `*time.Time` — Only agents created at or before this time (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**lastRunAfter:** `*time.Time` — Only agents last run at or after this time (RFC 3339).
    
</dd>
</dl>

<dl>
<dd>

**lastRunBefore:** `*time.Time` — Only agents last run at or before this time (RFC 3339).
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.Create(request) -> *platformgo.Agent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates an agent: a task for a phone, as instructions an agent works through a model or as code, with the model, output schema and limits its runs use. Sending instructions or code creates the agent's first revision with it; sending source_run_id saves one of your agent runs as the agent, its task or its script as the first revision. A name already used in the organization is 409; an unsupported model, an invalid output schema, both instructions and code, or a run with no task or script to save is 400; another organization's run is 404.
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
request := &platformgo.AgentCreateRequest{
        Name: "name",
    }
client.Agents.Create(
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

**accessibilityEnabled:** `*bool` — Use accessibility mode for runs (default false). Needs your organization's access to accessibility mode: without it the create is refused with a 403 whose message starts with accessibility_unavailable.
    
</dd>
</dl>

<dl>
<dd>

**captureEnabled:** `*bool` — Capture the media runs produce into the file library (default true).
    
</dd>
</dl>

<dl>
<dd>

**code:** `*string` — Python the agent runs. Creates a code revision. Leave out to add a revision later, or send instructions instead.
    
</dd>
</dl>

<dl>
<dd>

**from:** `*platformgo.AgentCreateRequestFrom` — With source_run_id, what the first revision is. task (default): the run's task, as an instructions revision. script: the script the run left (GET /runs/{run_id}/script), as a code revision. A run with no task, or no script, is 400.
    
</dd>
</dl>

<dl>
<dd>

**instructions:** `*string` — The task for the agent to work through a model. Creates an instructions revision. Leave out to add a revision later, or send code instead.
    
</dd>
</dl>

<dl>
<dd>

**maxCostMicrodollars:** `*int64` — Most a run's model steps may cost, in microdollars (default 10000000).
    
</dd>
</dl>

<dl>
<dd>

**maxSteps:** `*int64` — Most model steps a run may take (default 500).
    
</dd>
</dl>

<dl>
<dd>

**model:** `*string` — The model instruction runs call, one of GET /models. Defaults to the platform default.
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — The agent's name: letters, numbers, hyphens and underscores, unique within the organization.
    
</dd>
</dl>

<dl>
<dd>

**onFailure:** `*platformgo.AgentCreateRequestOnFailure` — With code: what a run does when a step fails. fail (default) fails the run. ask pauses the run for a person (reason step_check_failed), who retries the step or completes the run with its output. repair hands the run to the agent, which finishes the task from the live screen and proposes a fixed revision for a person to promote.
    
</dd>
</dl>

<dl>
<dd>

**outputSchema:** `map[string]any` — The JSON Schema a run's result must satisfy. Leave out to take any result.
    
</dd>
</dl>

<dl>
<dd>

**pauseTimeoutSeconds:** `*int64` — Longest a run waits paused for an answer, in seconds (default 600, 60 to 3600).
    
</dd>
</dl>

<dl>
<dd>

**recordingEnabled:** `*bool` — Record the screen of runs (default false).
    
</dd>
</dl>

<dl>
<dd>

**sourceRunID:** `*string` — Save one of your agent runs as the agent: its first revision comes from the run (see from), and the run's model, output schema and limits stand in for any left out here. Send it without instructions or code. Another organization's run is 404.
    
</dd>
</dl>

<dl>
<dd>

**telemetryEnabled:** `*bool` — Emit telemetry for runs (default true).
    
</dd>
</dl>

<dl>
<dd>

**timeoutSeconds:** `*int64` — Longest a run may take, in seconds (default 3600, at most 3600).
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.Get(AgentID) -> *platformgo.Agent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches one agent. A missing or other-organization id is 404.
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
request := &platformgo.AgentsGetRequest{
        AgentID: "agent_id",
    }
client.Agents.Get(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.Delete(AgentID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes an agent and its revisions. Its runs and sessions keep their history.
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
request := &platformgo.AgentsDeleteRequest{
        AgentID: "agent_id",
    }
client.Agents.Delete(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.Update(AgentID, request) -> *platformgo.Agent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Changes an agent's name, model, output schema, limits or data toggles; omitted fields are left unchanged. Its revisions change through POST /agents/{agent_id}/revisions.
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
request := &platformgo.AgentUpdateRequest{
        AgentID: "agent_id",
    }
client.Agents.Update(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**accessibilityEnabled:** `*bool` — Use accessibility mode for runs. true needs your organization's access to accessibility mode (a 403 whose message starts with accessibility_unavailable without it).
    
</dd>
</dl>

<dl>
<dd>

**captureEnabled:** `*bool` — Capture the media runs produce into the file library.
    
</dd>
</dl>

<dl>
<dd>

**maxCostMicrodollars:** `*int64` — Most a run's model steps may cost, in microdollars.
    
</dd>
</dl>

<dl>
<dd>

**maxSteps:** `*int64` — Most model steps a run may take.
    
</dd>
</dl>

<dl>
<dd>

**model:** `*string` — The model instruction runs call, one of GET /models.
    
</dd>
</dl>

<dl>
<dd>

**name:** `*string` — The agent's new name.
    
</dd>
</dl>

<dl>
<dd>

**outputSchema:** `map[string]any` — The JSON Schema a run's result must satisfy. null clears it, so runs take any result.
    
</dd>
</dl>

<dl>
<dd>

**pauseTimeoutSeconds:** `*int64` — Longest a run waits paused for an answer, in seconds, 60 to 3600.
    
</dd>
</dl>

<dl>
<dd>

**recordingEnabled:** `*bool` — Record the screen of runs.
    
</dd>
</dl>

<dl>
<dd>

**telemetryEnabled:** `*bool` — Emit telemetry for runs.
    
</dd>
</dl>

<dl>
<dd>

**timeoutSeconds:** `*int64` — Longest a run may take, in seconds, at most 3600.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.ListRevisions(AgentID) -> *platformgo.AgentRevisionList</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists an agent's revisions, newest first, one page at a time.
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
request := &platformgo.AgentsListRevisionsRequest{
        AgentID: "agent_id",
    }
client.Agents.ListRevisions(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Page size.
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — Where the page starts.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.CreateRevision(AgentID, request) -> *platformgo.AgentRevision</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Appends a revision to an agent, instructions or code, and makes it the one runs use. Revisions are never changed or removed.
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
request := &platformgo.AgentRevisionCreateRequest{
        AgentID: "agent_id",
        Kind: platformgo.AgentRevisionCreateRequestKindCode,
    }
client.Agents.CreateRevision(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**code:** `*string` — The script, for a code revision.
    
</dd>
</dl>

<dl>
<dd>

**instructions:** `*string` — The task, for an instructions revision.
    
</dd>
</dl>

<dl>
<dd>

**kind:** `*platformgo.AgentRevisionCreateRequestKind` — code needs code; instructions needs instructions.
    
</dd>
</dl>

<dl>
<dd>

**message:** `*string` — A note on the revision.
    
</dd>
</dl>

<dl>
<dd>

**onFailure:** `*platformgo.AgentRevisionCreateRequestOnFailure` — For a code revision, what a run does when a step fails. fail fails the run. ask pauses the run for a person (reason step_check_failed), who retries the step or completes the run with its output. repair hands the run to the agent, which finishes the task and proposes a fixed revision for a person to promote. Defaults to the current revision's, or fail.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.GetRevision(AgentID, RevisionID) -> *platformgo.AgentRevision</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches one of an agent's revisions with its instructions or code. A revision of another agent is 404.
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
request := &platformgo.AgentsGetRevisionRequest{
        AgentID: "agent_id",
        RevisionID: "revision_id",
    }
client.Agents.GetRevision(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**revisionID:** `string` — The revision's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.PromoteRevision(AgentID, RevisionID) -> *platformgo.AgentRevision</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Makes a fix a repair proposed the agent's current revision. 409 when the revision is not a proposed fix waiting to be promoted, or when the agent has moved on from the revision the repaired run ran.
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
request := &platformgo.AgentsPromoteRevisionRequest{
        AgentID: "agent_id",
        RevisionID: "revision_id",
    }
client.Agents.PromoteRevision(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**revisionID:** `string` — The revision's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.ListSchedules(AgentID) -> *platformgo.AgentScheduleList</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists an agent's schedules, or with agent_id - the schedules of every agent, newest first, one page at a time.
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
request := &platformgo.AgentsListSchedulesRequest{
        AgentID: "agent_id",
    }
client.Agents.ListSchedules(
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

**agentID:** `string` — The agent's id, or - for the schedules of every agent.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Page size.
    
</dd>
</dl>

<dl>
<dd>

**offset:** `*int64` — Where the page starts.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.CreateSchedule(AgentID, request) -> *platformgo.AgentSchedule</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a schedule of an agent: each occurrence of cron_expression, read in time_zone, starts a run of the agent with trigger scheduled and the schedule's id, built from run_template as POST /runs builds a saved agent's run. A cron expression or time zone that doesn't parse, an expression that never fires, a revision of another agent, an agent with no revision, or a schedule with no occurrence before its end_at is 400.
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
request := &platformgo.AgentScheduleCreateRequest{
        AgentID: "agent_id",
        CronExpression: "cron_expression",
        TimeZone: "time_zone",
    }
client.Agents.CreateSchedule(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**cronExpression:** `string` — Five-field cron: minute, hour, day of month, month, day of week. For example 0 9 * * 1-5 is 09:00 on weekdays.
    
</dd>
</dl>

<dl>
<dd>

**endAt:** `*time.Time` — No occurrence at or after this time. Leave out to run until deleted.
    
</dd>
</dl>

<dl>
<dd>

**name:** `*string` — A label for the schedule.
    
</dd>
</dl>

<dl>
<dd>

**overlapPolicy:** `*platformgo.AgentScheduleCreateRequestOverlapPolicy` — skip (default) drops an occurrence while the previous run is unfinished; queue starts it anyway.
    
</dd>
</dl>

<dl>
<dd>

**runTemplate:** `*platformgo.AgentScheduleRunTemplate` — What each run starts with.
    
</dd>
</dl>

<dl>
<dd>

**startAt:** `*time.Time` — No occurrence before this time. Leave out to start now.
    
</dd>
</dl>

<dl>
<dd>

**timeZone:** `string` — The IANA time zone cron_expression is read in, for example America/Los_Angeles or UTC.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.GetSchedule(AgentID, ScheduleID) -> *platformgo.AgentSchedule</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches one of an agent's schedules. Another agent's is 404.
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
request := &platformgo.AgentsGetScheduleRequest{
        AgentID: "agent_id",
        ScheduleID: "schedule_id",
    }
client.Agents.GetSchedule(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**scheduleID:** `string` — The schedule's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.DeleteSchedule(AgentID, ScheduleID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes a schedule and cancels its runs that have not started. Its other runs keep their schedule_id.
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
request := &platformgo.AgentsDeleteScheduleRequest{
        AgentID: "agent_id",
        ScheduleID: "schedule_id",
    }
client.Agents.DeleteSchedule(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**scheduleID:** `string` — The schedule's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.UpdateSchedule(AgentID, ScheduleID, request) -> *platformgo.AgentSchedule</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Changes a schedule's name, timing, run template or overlap policy; omitted fields are left unchanged. An active schedule whose timing changed takes its next occurrence from now. Its state changes through :pause and :resume.
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
request := &platformgo.AgentScheduleUpdateRequest{
        AgentID: "agent_id",
        ScheduleID: "schedule_id",
    }
client.Agents.UpdateSchedule(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**scheduleID:** `string` — The schedule's id.
    
</dd>
</dl>

<dl>
<dd>

**cronExpression:** `*string` — Five-field cron.
    
</dd>
</dl>

<dl>
<dd>

**endAt:** `*time.Time` — No occurrence at or after this time. null clears it, so the schedule runs until deleted.
    
</dd>
</dl>

<dl>
<dd>

**name:** `*string` — The schedule's new label.
    
</dd>
</dl>

<dl>
<dd>

**overlapPolicy:** `*platformgo.AgentScheduleUpdateRequestOverlapPolicy` — skip or queue.
    
</dd>
</dl>

<dl>
<dd>

**runTemplate:** `*platformgo.AgentScheduleRunTemplate` — What each run starts with, replaced whole.
    
</dd>
</dl>

<dl>
<dd>

**startAt:** `*time.Time` — No occurrence before this time. null clears it.
    
</dd>
</dl>

<dl>
<dd>

**timeZone:** `*string` — The IANA time zone cron_expression is read in.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.PauseSchedule(AgentID, ScheduleID) -> *platformgo.AgentSchedule</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Stops a schedule's occurrences and cancels its runs that have not started. A run already running goes on. A schedule that is not active is 409.
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
request := &platformgo.AgentsPauseScheduleRequest{
        AgentID: "agent_id",
        ScheduleID: "schedule_id",
    }
client.Agents.PauseSchedule(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**scheduleID:** `string` — The schedule's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.ResumeSchedule(AgentID, ScheduleID) -> *platformgo.AgentSchedule</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Makes a paused schedule active again from its next occurrence after now; occurrences missed while paused do not run. A schedule that is not paused, or has no occurrence left before its end_at, is 409.
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
request := &platformgo.AgentsResumeScheduleRequest{
        AgentID: "agent_id",
        ScheduleID: "schedule_id",
    }
client.Agents.ResumeSchedule(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**scheduleID:** `string` — The schedule's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.RunSchedule(AgentID, ScheduleID) -> *platformgo.RunResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Starts one run of the schedule's template now, whatever the schedule's state and overlap policy, as POST /runs starts a saved agent's run. The run carries the schedule's id; the schedule's occurrences are unchanged. While your organization is at its concurrency limit, the run is refused with 409.
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
request := &platformgo.AgentsRunScheduleRequest{
        AgentID: "agent_id",
        ScheduleID: "schedule_id",
    }
client.Agents.RunSchedule(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**scheduleID:** `string` — The schedule's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Agents.Rollback(AgentID, request) -> *platformgo.AgentRevision</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Copies an earlier revision forward as the agent's newest revision, which becomes the one runs use. History only grows: nothing is removed.
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
request := &platformgo.AgentRollbackRequest{
        AgentID: "agent_id",
        RevisionID: "revision_id",
    }
client.Agents.Rollback(
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

**agentID:** `string` — The agent's id.
    
</dd>
</dl>

<dl>
<dd>

**revisionID:** `string` — The earlier revision to roll back to. It is copied forward as a new revision, which becomes current.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

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

## Models
<details><summary><code>client.Models.List() -> *platformgo.AgentModelList</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists the models an agent's runs can call: supported models that see the screen and call tools, by display name. An agent's model field takes one of their slugs.
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
client.Models.List(
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

Returns one page of the org's phone sessions for the Session Inspector table: active/unbilled sessions pinned on top, terminal history paginated beneath. Covers agents' runs and interactive leases; each row links to a session. Filters: search, agent_id, status.
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

**search:** `*string` — case-insensitive match on phone/session/agent
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `*string` — only sessions of this agent's runs
    
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

Returns one session for the Session Inspector: session lifecycle, phone display fields, the agent's name (for an agent's run) and an inlined presigned recording URL. Works for active and terminal sessions, and for agents' runs and interactive leases. Org-scoped: another org's session reads as not found.
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

Returns a short-lived URL for the session's screen recording, keyed on session_id — so it works for agents' runs and interactive leases alike. Status is "pending" (no URL) when the recording hasn't finished uploading yet, and "off" when the session was not recorded. Org-scoped: another org's session reads as not found.
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

Not available right now: always answers 409 and leaves the phone unchanged.
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

Allocates an Android phone and opens a session you drive directly. Pass phone_id to pin a specific dedicated phone. If allocation setup fails the claim is rolled back, so you are never billed for a session that never starts.
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

<details><summary><code>client.Runs.StreamEvents(RunID) -> platformgo.RunsStreamEventsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Streams an agent run as Server-Sent Events (text/event-stream): the frames of its session's telemetry that belong to the run, the model's text as it streams (model_delta logs), and the run's own events, live as the run runs. Each event's id is its resume cursor, its event name is the frame kind (span or log), and its data is the frame JSON. The run's events are logs: run_status (attribute axilio.run.status: running when the run starts on its phone or goes back to running), run_input_required (it paused for a person; the body is the question, with the input request's id, reason, requested_by, response_schema and expires_at as attributes), run_resumed (a person answered; the body is the action), run_result (how it ended on the phone: axilio.ok, axilio.run.output or axilio.run.error, and the summary as the body), and user_message (a person's message reached the run; the body is the message). The stream's last event is a run_status with the run's final status as the run itself has it (completed, failed or cancelled), so a run the caller cancelled ends cancelled. These events are not archived; the run itself has the result. A session started with telemetry off still streams its run's events and the model's text, and nothing else. A new reader (or one whose cursor is older than the live window) first receives the run's archived frames and then the live ones, deduplicated on span id and phase, followed by an event with no data that carries the cursor the replay ended at. Reconnect with the Last-Event-ID header to resume. The model's text is not archived: a reader that reconnects after it left the live window gets the step's answer from the run's steps instead. A comment line is sent every 15 seconds while the stream is idle. The stream closes after the run's final status; a reconnect after the end with nothing left gets 204. A run that has not reached a phone yet streams only comments until it does. Authenticate with the API key header, or open the run's stream_url, which carries a read-only token for the run's session. Org-scoped: another org's run reads as not found. Readers MUST ignore events with an unknown name, unknown fields, and unknown span_type or log_type values.
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
request := &platformgo.RunsStreamEventsRequest{
        RunID: "run_id",
    }
client.Runs.StreamEvents(
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

**runID:** `string` — Run whose events to stream.
    
</dd>
</dl>

<dl>
<dd>

**token:** `*string` — Read-only token for the run's session, as carried by stream_url. Not needed with the API key header.
    
</dd>
</dl>

<dl>
<dd>

**lastEventID:** `*string` — Resume after this cursor (the id of the last event received).
    
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

Returns paginated recent (non-archived) runs across the caller's organization, every member's runs included. Filters: agent_id, search (run ID substring), status, trigger. Order with order_by ('<field> <asc|desc>'), field one of status, started_at, completed_at, created_at.
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

**agentID:** `*string` — Only runs of this agent.
    
</dd>
</dl>

<dl>
<dd>

**scheduleID:** `*string` — Only runs this schedule started.
    
</dd>
</dl>

<dl>
<dd>

**mode:** `*platformgo.RunsListRequestMode` — Only runs of this mode.
    
</dd>
</dl>

<dl>
<dd>

**deviceClass:** `*platformgo.RunsListRequestDeviceClass` — Only runs that drove this kind of device.
    
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

<details><summary><code>client.Runs.Start(request) -> *platformgo.RunResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Starts a run of a saved agent (agent_id, its current revision and settings) or an ad hoc run (task, with the model, output schema, limits, recording and telemetry given here). A phone_id target must be one of your organization's dedicated phones. While your organization is at its concurrency limit, a run that starts now is refused with 409. Booked with start_at, the run is scheduled instead.
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
request := &platformgo.AgentRunCreateRequest{}
client.Runs.Start(
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

**request:** `*platformgo.AgentRunCreateRequest` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.GetMetrics() -> *platformgo.RunMetricsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns your organization's finished runs over a window as points: one per hour or day in time_zone (and per group with group_by), each with the runs by how they ended, the success rate, the median and 95th percentile duration, what they cost and the model steps they took. A run counts in the bucket it finished in; runs still going are not counted. The window defaults to the last 7 days and may start at most a year back; hour buckets need a window of 7 days or less.
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
request := &platformgo.RunsGetMetricsRequest{}
client.Runs.GetMetrics(
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

**granularity:** `*platformgo.RunsGetMetricsRequestGranularity` — Bucket size: hour or day. hour is allowed for windows of 7 days or less.
    
</dd>
</dl>

<dl>
<dd>

**timeZone:** `*string` — IANA time zone the buckets follow, such as America/Los_Angeles: a day starts at local midnight and an hour on the local hour. An unknown name is an error.
    
</dd>
</dl>

<dl>
<dd>

**groupBy:** `*platformgo.RunsGetMetricsRequestGroupBy` — Split each bucket by agent, revision, status, trigger, mode, schedule or device_class. Omit for one point per bucket.
    
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

<details><summary><code>client.Runs.ListEvents(RunID) -> *platformgo.AgentRunEventList</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists an agent run's events as pages, oldest first: the frames its stream sends (its spans, the model's text and the run's own run_status, run_input_required, run_resumed, run_result and user_message logs) in time order, for a caller that polls rather than holds the stream open. A live run's page always has a next_cursor to poll from; an ended run's last page closes with its final run_status, as the run itself has it, and has none. Pages never repeat a frame, except that a cursor older than the stream's 10-minute window starts again from the run's first event, so a reader that keeps one deduplicates on span_id and phase, as a stream reader does across reconnects. A page keeps a stream entry's frames together, so it can hold a few more than limit. The run's own events and the model's text are live only: past the stream's 10-minute window a page has the archived spans and the closing run_status but no run_result, and a run with telemetry off has only its events while they are live. The run itself always has its result.
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
request := &platformgo.RunsListEventsRequest{
        RunID: "run_id",
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

**runID:** `string` — The run's id.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — The next_cursor of the page before; leave out for the first page.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Page size.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.GetScript(RunID) -> *platformgo.AgentRunScript</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns every cell the run executed, in order, as one Python file and one by one, failed and exploratory cells included. A code run's first cell is its revision's script. The script is kept with the run's data, so it outlives the transcript; when that record is missing, the cells are rebuilt from what remains of the transcript and partial is true. 404 when there is nothing to build it from.
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
request := &platformgo.RunsGetScriptRequest{
        RunID: "run_id",
    }
client.Runs.GetScript(
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

**runID:** `string` — The run's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.ListSteps(RunID) -> *platformgo.AgentRunStepList</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists a run's model steps from its transcript, oldest first: the messages each step sent, the model's answer, links to its screenshots and its cost. Empty once the transcript has passed its retention, or when your organization keeps none.
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
request := &platformgo.RunsListStepsRequest{
        RunID: "run_id",
    }
client.Runs.ListSteps(
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

**runID:** `string` — The run's id.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — The next_cursor of the page before; leave out for the first page.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int64` — Page size.
    
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

Cancels a run that is scheduled, queued, running or paused, scoped to the caller's org. A scheduled or queued run is cancelled before it starts and never runs. A running or paused run is stopped on its phone; an agent-mode run then keeps its phone for keep_alive_seconds, as when a turn ends, so :sendMessage can follow up on the same phone until keep_alive_expires_at. A run that has already reached a terminal state (completed/failed/cancelled) cannot be cancelled and reads as not found. Returns the updated run.
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

<details><summary><code>client.Runs.Pause(RunID) -> *platformgo.RunResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Asks a running run to pause for a person, with reason caller_requested. An agent-mode run pauses at its next model step: a cell already running finishes first, and a run about to finish takes that step instead. A code run pauses at its script's next phone action. Returns the run as it is: it reads paused, with its input_request, once the pause lands, and run.paused fires then. A run that is not running is 409.
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
request := &platformgo.RunsPauseRequest{
        RunID: "run_id",
    }
client.Runs.Pause(
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

**runID:** `string` — The run's id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.Resume(RunID, request) -> *platformgo.RunResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Resumes a paused run. action answer (the default) answers its input request with response, which must satisfy the request's response_schema (422 when it does not). action retry runs a code run's failed script again (only for reason step_check_failed). action complete ends the run as completed with output, which must satisfy the run's output schema (422). Returns the run as it is. With answer or retry it reads running once the edge has the input, and run.resumed fires then; with complete it reads completed once the edge ends it, and run.completed fires then. A run that is not paused, a different input_request_id, or a request that has expired is 409.
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
request := &platformgo.AgentRunResumeRequest{
        RunID: "run_id",
        InputRequestID: "input_request_id",
    }
client.Runs.Resume(
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

**runID:** `string` — The run's id.
    
</dd>
</dl>

<dl>
<dd>

**action:** `*platformgo.AgentRunResumeRequestAction` — answer (the default) hands the run response; not for reason step_check_failed. retry runs the code run's script again in the same session, from where the phone is now; only for a pause with reason step_check_failed. complete ends the run as completed with output, as if it had returned it.
    
</dd>
</dl>

<dl>
<dd>

**inputRequestID:** `string` — The input request answered: the run's input_request.id.
    
</dd>
</dl>

<dl>
<dd>

**message:** `*string` — With answer: a note for the agent with the answer.
    
</dd>
</dl>

<dl>
<dd>

**model:** `*string` — With answer, for an agent-mode run: the model its next steps call, one of GET /models (default the run's). It becomes the run's model, and its current turn's, once the run resumes; each step keeps the model that answered it.
    
</dd>
</dl>

<dl>
<dd>

**output:** `any` — With complete: the run's result, which must satisfy the run's output schema. At most 64 KiB.
    
</dd>
</dl>

<dl>
<dd>

**response:** `any` — With answer: the answer, which must satisfy the input request's response_schema. At most 64 KiB.
    
</dd>
</dl>

<dl>
<dd>

**resultText:** `*string` — With complete: a short account of the result.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.SendMessage(RunID, request) -> *platformgo.RunResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Sends a person's message to an agent-mode run. To a run that is running or paused, the run's model reads it before its next step, as a turn from the person running the task, and the run's stream carries it as a user_message event, its body the message, once the run's phone has it. A message never resumes a paused run: answer the pause with :resume. A run holds at most 20 unread messages, and the edge refuses one past that, or one that arrives as the run finishes, which may then never be read or shown. To a run that ended (completed, failed or cancelled) while it still holds its phone (keep_alive_expires_at is in the future), the message starts the run's next turn on the same phone, and the run goes back to queued. Returns the run as it is. A code run, a queued or scheduled run, a run whose hold has ended, or a run whose phone cannot take messages yet is 409; an empty message, or one over 64 KiB, is 400.
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
request := &platformgo.AgentRunSendMessageRequest{
        RunID: "run_id",
        Message: "message",
    }
client.Runs.SendMessage(
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

**runID:** `string` — The run's id.
    
</dd>
</dl>

<dl>
<dd>

**message:** `string` — What the person says to the run's model. Required, at most 64 KiB.
    
</dd>
</dl>

<dl>
<dd>

**model:** `*string` — The next turn's model, one of GET /models, for a message to a completed or failed run (default the run's). Not taken while the run works.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Runs.BatchCreate(request) -> *platformgo.AgentRunBatchCreateResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Starts up to 1000 runs, each as POST /runs would. All are created or none is; an error names the request it is for. A batch may not name the same dedicated phone twice. The batch is refused with 409 while your organization is at its concurrency limit; its runs beyond the limit wait in the queue for a free slot.
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
request := &platformgo.AgentRunBatchCreateRequest{}
client.Runs.BatchCreate(
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

**requests:** `[]*platformgo.AgentRunCreateRequest` — The runs to start.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Sessions
<details><summary><code>client.Sessions.GetMetrics() -> *platformgo.PhoneSessionMetricsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns your organization's ended phone sessions over a window as points: one per hour or day in time_zone (and per group with group_by), each with the sessions, their phone time in seconds and what it cost. A session counts in the bucket it was billed in, as it ended, dedicated phones' at no cost; sessions still going are not counted. The window defaults to the last 7 days and may start at most a year back; hour buckets need a window of 7 days or less.
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
request := &platformgo.SessionsGetMetricsRequest{}
client.Sessions.GetMetrics(
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

**granularity:** `*platformgo.SessionsGetMetricsRequestGranularity` — Bucket size: hour or day. hour is allowed for windows of 7 days or less.
    
</dd>
</dl>

<dl>
<dd>

**timeZone:** `*string` — IANA time zone the buckets follow, such as America/Los_Angeles: a day starts at local midnight and an hour on the local hour. An unknown name is an error.
    
</dd>
</dl>

<dl>
<dd>

**groupBy:** `*platformgo.SessionsGetMetricsRequestGroupBy` — Split each bucket by allocated_by, device_class or platform. Omit for one point per bucket.
    
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

Paginated, filterable list of inference calls (detect, locate, element_locate, and agent_chat for an agent run's model steps) the caller's user was billed for over a required time window (created_after/created_before). Filters: endpoint, model, session, free-text search. Order with order_by ('<field> <asc|desc>').
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

**endpoint:** `[]*platformgo.UsageListInferencesRequestEndpointItem` — Only inferences of one of these operations; repeatable. agent_chat is an agent run's model step.
    
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

Returns your organization's usage spend over a window as points: one per hour or day in time_zone (and per group with group_by: product, agent, model, api_key, schedule or device_class), each with the total, phone-session and inference spend in microdollars, the phone seconds and tokens it bought, zero when nothing was billed. The window defaults to the last 7 days and may start at most a year back; hour buckets need a window of 7 days or less.
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

**groupBy:** `*platformgo.UsageGetMetricsRequestGroupBy` — Split each bucket by product (session or inference), agent, model, api_key, schedule or device_class. Omit for one point per bucket.
    
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

Paginated, filterable list of the caller's phone sessions over a required time window (allocated_after/allocated_before), each row carrying billing detail: phone-time cost, inference cost, and combined total in microdollars, billing processing status, and allocation source. Filters: session status, billing processing status, agent, allocation source, free-text search. Order with order_by ('<field> <asc|desc>').
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

**agentID:** `*string` — Restrict results to sessions of a single agent's runs.
    
</dd>
</dl>

<dl>
<dd>

**allocatedBy:** `[]*platformgo.UsageListSessionsRequestAllocatedByItem` — Restrict results to the given allocation sources; repeatable.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Filter by session, agent, or phone id substring.
    
</dd>
</dl>

<dl>
<dd>

**orderBy:** `*string` — Sort expression '<field> [asc|desc]', direction defaulting to desc; field is allocated_at, deallocated_at, duration_seconds, cost_microdollars, session_status, processed_status, allocated_by, session_id or agent_id. Defaults to allocated_at desc.
    
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

