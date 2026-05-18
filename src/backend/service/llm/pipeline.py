import asyncio
from backend.config.config import settings
from langchain_core.messages import AIMessage, HumanMessage, SystemMessage
from langchain_openai import ChatOpenAI
from sentence_transformers import SentenceTransformer
from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy import text

if settings.DB_URL:
    db_url = settings.DB_URL
else:
    db_url = f"postgresql+psycopg://{settings.DB_USER}:{settings.DB_PASSWORD}@{settings.DB_HOST}:{settings.DB_PORT}/{settings.DB_NAME}"

if db_url.startswith("postgresql://"):
    db_url = db_url.replace("postgresql://", "postgresql+psycopg://", 1)

engine = create_async_engine(db_url)
embedding_model = SentenceTransformer(settings.EMBEDDING_MODEL_NAME)

llm = ChatOpenAI(
    base_url=settings.OPENROUTER_BASE_URL,
    openai_api_key=settings.OPENROUTER_API_KEY,
    model=settings.LLM_MODEL_NAME,
    temperature=settings.LLM_TEMPERATURE,
    max_tokens=settings.LLM_MAX_TOKENS,
    streaming=True,
)


async def get_relevant_context_async(query: str, limit: int = 3) -> str:
    loop = asyncio.get_running_loop()
    query_vector = await loop.run_in_executor(
        None, embedding_model.encode, query
    )
    query_vector_list = query_vector.tolist()

    async with engine.connect() as conn:
        sql = text("""
            SELECT content 
            FROM knowledge_documents 
            ORDER BY embedding <=> :vector::vector 
            LIMIT :limit;
        """)
        result = await conn.execute(
            sql, {"vector": str(query_vector_list), "limit": limit}
        )
        rows = result.fetchall()
        return "\n---\n".join([row[0] for row in rows]) if rows else ""


async def stream_rag_chat(history_messages: list):
    user_query = history_messages[-1]["content"] if history_messages else ""
    context = await get_relevant_context_async(user_query)

    langchain_messages = [
        SystemMessage(
            content=f"You are an expert tutor for EXE101 course. Use this context to answer:\n{context}"
        )
    ]

    for msg in history_messages:
        if msg["role"] == "user":
            langchain_messages.append(HumanMessage(content=msg["content"]))
        elif msg["role"] == "assistant":
            langchain_messages.append(AIMessage(content=msg["content"]))

    range_stream = llm.astream(langchain_messages)
    async for chunk in range_stream:
        if chunk.content:
            yield f"data:{chunk.content}\n"