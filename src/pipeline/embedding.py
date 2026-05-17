from huggingface_hub import AsyncInferenceClient
from backend.config.config import settings

class EmbeddingEngine:
    def __init__(self):
        print(f"Initializing API Embedding Engine: {settings.EMBEDDING_MODEL_NAME}")
        self.client = AsyncInferenceClient(
            model=settings.EMBEDDING_MODEL_NAME,
            token=settings.HF_TOKEN
        )

    async def get_embeddings_async(self, texts: list[str]) -> list[list[float]]:
        if not texts:
            return []

        try:
            return await self.client.feature_extraction(texts)
        
        except Exception as e:
            print(f"[ERROR] Embedding API failed: {str(e)}")
            raise Exception(f"Embedding API Error: {str(e)}")