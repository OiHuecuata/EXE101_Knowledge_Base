from typing import List
from urllib.parse import urlparse

import uvicorn
from backend.service.llm.pipeline import stream_rag_chat
from fastapi import FastAPI
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

app = FastAPI()


class MessageSchema(BaseModel):
    role: str
    content: str


class ChatRequest(BaseModel):
    session_id: str
    messages: List[MessageSchema]


@app.post("/api/v1/chat/stream")
async def chat_stream(payload: ChatRequest):
    history_list = [msg.model_dump() for msg in payload.messages]
    return StreamingResponse(
        stream_rag_chat(history_list), media_type="text/event-stream"
    )


if __name__ == "__main__":
    from backend.config.config import settings

    try:
        parsed_url = urlparse(settings.PYTHON_LLM_SERVICE_URL)
        port = parsed_url.port if parsed_url.port else 5000
    except Exception:
        port = 5000

    uvicorn.run(
        "backend.service.llm.server:app",
        host="127.0.0.1",
        port=port,
        reload=True,
    )