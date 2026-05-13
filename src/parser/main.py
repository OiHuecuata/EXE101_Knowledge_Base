import os
import pathlib
from pptx import Presentation
from docx import Document

RAW_DATA_DIR = pathlib.Path("./data/EXE101")
PROCESSED_DATA_DIR = pathlib.Path("./data/processed")

def pptx_parser(file_path):
    try:
        prs = Presentation(file_path)
        md_output = [f"# {file_path.stem}\n"]
        
        for i, slide in enumerate(prs.slides):
            md_output.append(f"## Slide {i+1}")
            
            # Text Extraction
            slide_text = []
            for shape in slide.shapes:
                if hasattr(shape, "text") and shape.text.strip():
                    slide_text.append(shape.text.strip())
            
            if slide_text:
                md_output.append("\n".join(slide_text))
            
            md_output.append("\n---\n") # Line to separate slides
            
        return "\n".join(md_output)
    except Exception as e:
        return f"Error parsing PPTX {file_path.name}: {e}"

def docx_parser(file_path):
    try:
        doc = Document(file_path)
        md_output = [f"# {file_path.stem}\n"]
        
        for para in doc.paragraphs:
            text = para.text.strip()
            if text:
                # Style Header
                md_output.append(text + "\n")
                
        return "\n".join(md_output)
    except Exception as e:
        return f"Error parsing DOCX {file_path.name}: {e}"

def main():

    PROCESSED_DATA_DIR.mkdir(parents=True, exist_ok=True)
    
    file_count = 0
    
    for root, _, files in os.walk(RAW_DATA_DIR):
        for file in files:
            file_path = pathlib.Path(root) / file
            suffix = file_path.suffix.lower()
            
            content = ""
            if suffix == ".pptx":
                print(f"Converting PPTX: {file}")
                content = pptx_parser(file_path)
            elif suffix == ".docx":
                print(f"Converting DOCX: {file}")
                content = docx_parser(file_path)
            else:
                continue

            # Save with Markdown type
            output_file = PROCESSED_DATA_DIR / f"{file_path.stem}.md"
            with open(output_file, "w", encoding="utf-8") as f:
                f.write(content)
            file_count += 1

    print(f"\nExtracted {file_count} files to {PROCESSED_DATA_DIR}")

if __name__ == "__main__":
    main()