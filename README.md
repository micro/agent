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

### Query the AI

Query it the AI directly

```
micro agent query --question "Tell me about Islam"
```

### Use a command

Make a request that will issue service calls e.g call helloworld service

```
micro agent command --request "Say hello"
```
