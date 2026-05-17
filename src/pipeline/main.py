import re
import asyncio
from pathlib import Path
from ingest import AsyncIngestionPipeline

PROCESSED_DATA_DIR = Path("./data/processed")

def extract_week_number(file_name: str) -> int:
    match = re.search(r'(?:week|w)\s*(\d+)', file_name, re.IGNORECASE)
    return int(match.group(1)) if match else 1

def chunk_markdown_file(file_path: Path) -> list[dict]:
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    chunks = []
    
    if "## Slide" in content:

        slides = content.split("---")
        for slide in slides:
            slide_content = slide.strip()
            if not slide_content:
                continue
            
            slide_match = re.search(r'## Slide (\d+)', slide_content)
            slide_num = int(slide_match.group(1)) if slide_match else 0
            
            clean_content = re.sub(r'## Slide \d+', '', slide_content).strip()
            if clean_content:
                chunks.append({
                    "content": clean_content,
                    "metadata": {"slide": slide_num}
                })
    
    else:
        paragraphs = content.split("\n\n")
        for idx, para in enumerate(paragraphs):
            para_content = para.strip()
            if para_content and not para_content.startswith("#"):
                chunks.append({
                    "content": para_content,
                    "metadata": {"paragraph_index": idx}
                })
                
    return chunks

async def main():
    if not PROCESSED_DATA_DIR.exists():
        print(f"[ERROR] Processed directory not found at: {PROCESSED_DATA_DIR}. Please run parser.py first.")
        return

    markdown_files = list(PROCESSED_DATA_DIR.glob("*.md"))
    if not markdown_files:
        print(f"[WARNING] No processed .md files found in {PROCESSED_DATA_DIR}")
        return

    print(f"Initializing ingestion pipeline for {len(markdown_files)} files")
    pipeline = AsyncIngestionPipeline()

    for file_path in markdown_files:
        file_name = file_path.name
        week_number = extract_week_number(file_name)
        
        print("\n" + "="*60)
        print(f"Processing: {file_name} | Week: {week_number}")
        print("="*60)

        try:
            chunks = chunk_markdown_file(file_path)
            if not chunks:
                print(f"[WARNING] No chunks extracted from {file_name}. Skipping")
                continue

            await pipeline.execute_async(
                file_name=file_name,
                file_path=str(file_path),
                week_number=week_number,
                chunks=chunks
            )
        except Exception as e:
            print(f"[ERROR] Error processing {file_name}: {e}")
            continue

    print("\nData pipeline execution finished successfully!")

if __name__ == "__main__":
    asyncio.run(main())