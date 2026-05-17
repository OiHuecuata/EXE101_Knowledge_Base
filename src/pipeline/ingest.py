import psycopg
from pgvector.psycopg import register_vector_async
from backend.config.config import settings
from pipeline.embedding import EmbeddingEngine

class AsyncIngestionPipeline:
    def __init__(self):
        self.embed_engine = EmbeddingEngine()
        self.db_url = settings.DB_URL or f"postgresql://{settings.DB_USER}:{settings.DB_PASSWORD}@{settings.DB_HOST}:{settings.DB_PORT}/{settings.DB_NAME}"

    async def get_or_create_document(self, aconn, file_name: str, file_path: str, week_number: int) -> int:
        async with aconn.cursor() as acur:
            await acur.execute(
                """
                INSERT INTO documents (file_name, file_path, week_number)
                VALUES (%s, %s, %s)
                ON CONFLICT (file_name) 
                DO UPDATE SET file_path = EXCLUDED.file_path, week_number = EXCLUDED.week_number
                RETURNING id;
                """,
                (file_name, file_path, week_number)
            )
            row = await acur.fetchone()
            return row[0]

    async def execute_async(self, file_name: str, file_path: str, week_number: int, chunks: list[dict]):
        if not chunks:
            print(f"[WARNING] No data chunks provided for {file_name}.")
            return

        print(f"Starting async embedding generation for {len(chunks)} chunks")
        texts = [chunk["content"] for chunk in chunks]
        vectors = await self.embed_engine.get_embeddings_async(texts)

        print("--> Establishing async connection to PostgreSQL...")
        async with await psycopg.AsyncConnection.connect(self.db_url) as aconn:
            await register_vector_async(aconn)
            
            document_id = await self.get_or_create_document(aconn, file_name, file_path, week_number)
            
            async with aconn.cursor() as acur:
                await acur.execute("DELETE FROM document_chunks WHERE document_id = %s;", (document_id,))
                
                insert_query = """
                    INSERT INTO document_chunks (document_id, chunk_index, content, embedding, metadata)
                    VALUES (%s, %s, %s, %s, %s);
                """
                
                print(f"Inserting batch chunks into document_chunks for document ID: {document_id}")
                for i, chunk in enumerate(chunks):
                    await acur.execute(insert_query, (
                        document_id,
                        i,
                        chunk.get("content"),
                        vectors[i],
                        psycopg.types.json.Jsonb(chunk.get("metadata", {}))
                    ))
                
                await aconn.commit()
                
        print(f"Successfully ingested {file_name} into database!")