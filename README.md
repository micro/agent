# Agent

The beginnings of an agent

## Overview

The Micro agent is a service which provides generative AI capabilities. It makes use of the go-micro/genai package to abstract 
the underlying provider. Still a work in progress.

## Usage

Set your genai provider e.g `openai` or `gemini`

```
export MICRO_GENAI=openai
```

Set your API key

```
export MICRO_GENAI_KEY=xxxx
```

Run the agent

```
micro run
```

Query it via the CLI
```
micro agent query --prompt "Tell me about Islam"
```
