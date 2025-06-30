#!/usr/bin/env python3

import json
import re
from docling.document_converter import DocumentConverter

def extract_morocco_prefixes_from_pdf(pdf_path):
    """Extract Morocco mobile prefixes from ANRT PDF using Docling"""
    
    # Initialize Docling converter
    converter = DocumentConverter()
    
    # Convert PDF
    print(f"Converting PDF: {pdf_path}")
    result = converter.convert(pdf_path)
    
    # Get the document text
    text = result.document.export_to_markdown()
    
    # Save markdown for inspection
    with open("morocco_prefixes.md", "w", encoding="utf-8") as f:
        f.write(text)
    print("Saved extracted text to morocco_prefixes.md")
    
    # Initialize operators
    operators = {
        "ETISALAT AL MAGHRIB": {"name": "Maroc Telecom (IAM)", "prefixes": set()},
        "MEDI TELECOM": {"name": "Orange Morocco (Méditel)", "prefixes": set()},  
        "WANA CORPORATE": {"name": "Inwi", "prefixes": set()}
    }
    
    # Parse the text to extract mobile prefixes
    lines = text.split('\n')
    current_operator = None
    
    print("\nParsing extracted text for mobile prefixes...")
    
    for i, line in enumerate(lines):
        line = line.strip()
        
        # Check for operator names
        for op_key in operators:
            if op_key in line.upper():
                current_operator = op_key
                print(f"\nFound operator: {current_operator}")
                break
        
        # Look for mobile numbers patterns
        if current_operator and ('mobile' in line.lower() or 'gsm' in line.lower() or '06' in line or '07' in line):
            # Extract patterns like 060X, 061X, etc.
            # Pattern for 06XX or 07XX with X notation
            x_patterns = re.findall(r'0([67]\d)X', line)
            for pattern in x_patterns:
                print(f"  Found pattern 0{pattern}X")
                for digit in range(10):
                    prefix = f"212{pattern}{digit}"
                    operators[current_operator]["prefixes"].add(prefix)
            
            # Pattern for XX ranges like 06XX or 07XX
            xx_patterns = re.findall(r'0([67])XX', line)
            for pattern in xx_patterns:
                print(f"  Found pattern 0{pattern}XX")
                for tens in range(10):
                    for ones in range(10):
                        prefix = f"212{pattern}{tens}{ones}"
                        operators[current_operator]["prefixes"].add(prefix)
            
            # Direct 4-digit prefixes like 0610, 0620, etc
            direct_patterns = re.findall(r'0([67]\d{2})', line)
            for pattern in direct_patterns:
                if 'X' not in line:  # Avoid duplicates from X patterns
                    prefix = f"212{pattern}"
                    operators[current_operator]["prefixes"].add(prefix)
                    print(f"  Found direct prefix: 0{pattern}")
    
    # Convert sets to sorted lists
    for op in operators.values():
        op["prefixes"] = sorted(list(op["prefixes"]))
    
    return operators

def save_to_json(operators, output_path):
    """Save extracted data to JSON file"""
    
    data = {
        "operators": operators,
        "metadata": {
            "country": "Morocco",
            "country_code": "212",
            "mobile_prefix_regex": "^212(6[0-9]{2}|7[0-7][0-9])",
            "number_length": 12,
            "source": "ANRT Official Document June 2025 (extracted with Docling)",
            "updated": "2025-06-30"
        }
    }
    
    with open(output_path, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
    
    # Print summary
    print("\n\nSummary:")
    print("========")
    for op_key, op_data in operators.items():
        print(f"{op_data['name']}: {len(op_data['prefixes'])} prefixes")

if __name__ == "__main__":
    pdf_path = "morocco_prefixes.pdf"
    output_path = "/root/e173_go_gateway/data/morocco_mobile_prefixes_docling.json"
    
    print("Extracting Morocco mobile prefixes from ANRT PDF using Docling...")
    operators = extract_morocco_prefixes_from_pdf(pdf_path)
    save_to_json(operators, output_path)
    print(f"\nData saved to: {output_path}")