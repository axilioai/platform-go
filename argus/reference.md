# Reference
## Vision
<details><summary><code>client.Vision.Detect(request) -> *argus.InferenceResponse</code></summary>
<dl>
<dd>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &argus.InferenceRequest{
        Image: "image",
    }
client.Vision.Detect(
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

**confidenceThreshold:** `*float64` — Minimum confidence threshold for detections and OCR results
    
</dd>
</dl>

<dl>
<dd>

**image:** `string` — Base64 encoded image data
    
</dd>
</dl>

<dl>
<dd>

**inferenceType:** `*argus.InferenceType` — Type of inference to perform
    
</dd>
</dl>

<dl>
<dd>

**nmsIouThreshold:** `*float64` — IoU threshold for Non-Maximum Suppression (NMS)
    
</dd>
</dl>

<dl>
<dd>

**ocrEngine:** `*string` — OCR engine to use: 'free' or 'premium'
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Vision.Locate(request) -> *argus.LocateResponse</code></summary>
<dl>
<dd>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &argus.LocateRequest{
        Image: "image",
        Query: "query",
    }
client.Vision.Locate(
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

**image:** `string` — Base64 encoded image (PNG or JPEG)
    
</dd>
</dl>

<dl>
<dd>

**model:** `*string` — VLM model to use; must be one of the models from GET /vision/models. Omit to use the server's configured default. The system prompt is fixed to the element-locator task.
    
</dd>
</dl>

<dl>
<dd>

**query:** `string` — Natural-language target description
    
</dd>
</dl>

<dl>
<dd>

**texts:** `[]*argus.TextElementInput` — Pre-computed OCR text elements. Empty list means Argus skips OCR grounding and asks the VLM to locate from the image alone.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Vision.ListModels() -> *argus.SupportedModelsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

List every model Argus supports, with pricing.

Two families: the curated VLMs served by /vision/locate (per-token
pricing; pass their id as `model`) and the Axilio model line behind
/vision/detect (per-page pricing; selected via `ocr_engine` /
`inference_type` — lite is the free engine, pro is premium).

Public: no API key required (it's a catalog of model names and public
prices, nothing sensitive), so a client can discover supported models
before it holds credentials. The SDK fetches this once, caches it, and
validates find(model=...) locally so a typo fails fast with a clean
error instead of a 400 from /locate.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Vision.ListModels(
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

