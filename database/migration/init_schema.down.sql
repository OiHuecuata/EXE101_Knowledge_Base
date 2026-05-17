-- Drop HNSW index
DROP INDEX IF EXISTS document_chunks_embedding_hnsw_idx;

-- Drop tables
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_sessions;
DROP TABLE IF EXISTS document_chunks;
DROP TABLE IF EXISTS documents;

-- Drop custom types
DROP TYPE IF EXISTS message_role;

-- Drop extension
DROP EXTENSION IF EXISTS vector;