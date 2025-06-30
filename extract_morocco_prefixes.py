#!/usr/bin/env python3

import pdfplumber
import re
import json

def extract_morocco_prefixes(pdf_path):
    """Extract Morocco mobile prefixes from ANRT PDF"""
    
    operators = {
        "ETISALAT AL MAGHRIB": {"name": "Maroc Telecom (IAM)", "prefixes": []},
        "MEDI TELECOM": {"name": "Orange Morocco (Méditel)", "prefixes": []},
        "WANA CORPORATE": {"name": "Inwi", "prefixes": []}
    }
    
    current_operator = None
    mobile_section = False
    
    with pdfplumber.open(pdf_path) as pdf:
        for page_num, page in enumerate(pdf.pages):
            print(f"\nProcessing page {page_num + 1}")
            
            # Extract text from page
            text = page.extract_text()
            if not text:
                continue
                
            lines = text.split('\n')
            
            for line in lines:
                line = line.strip()
                
                # Check for operator name
                for op_key in operators:
                    if op_key in line:
                        current_operator = op_key
                        print(f"Found operator: {current_operator}")
                        break
                
                # Look for mobile numbers section
                if "mobile" in line.lower() or "gsm" in line.lower():
                    mobile_section = True
                elif "fixe" in line.lower() or "fixed" in line.lower():
                    mobile_section = False
                
                # Extract prefixes if we're in mobile section
                if current_operator and mobile_section:
                    # Look for patterns like 060X, 07XX, specific numbers
                    
                    # Pattern for 06XX or 07XX blocks
                    block_matches = re.findall(r'0(6\d{2}|7\d{2})', line)
                    for match in block_matches:
                        prefix = f"212{match}"
                        if prefix not in operators[current_operator]["prefixes"]:
                            operators[current_operator]["prefixes"].append(prefix)
                            print(f"  Added prefix: {prefix}")
                    
                    # Pattern for ranges like 060X (meaning 0600-0609)
                    range_matches = re.findall(r'0(6\d)X|0(7\d)X', line)
                    for match in range_matches:
                        base = match[0] if match[0] else match[1]
                        for i in range(10):
                            prefix = f"212{base}{i}"
                            if prefix not in operators[current_operator]["prefixes"]:
                                operators[current_operator]["prefixes"].append(prefix)
                                print(f"  Added prefix from range: {prefix}")
                    
                    # Pattern for XX ranges like 07XX (meaning 0700-0799)
                    xx_matches = re.findall(r'0([67])XX', line)
                    for match in xx_matches:
                        for tens in range(10):
                            for ones in range(10):
                                prefix = f"212{match}{tens}{ones}"
                                if prefix not in operators[current_operator]["prefixes"]:
                                    operators[current_operator]["prefixes"].append(prefix)
                                    print(f"  Added prefix from XX range: {prefix}")
    
    # Sort prefixes for each operator
    for op in operators.values():
        op["prefixes"].sort()
    
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
            "source": "ANRT Official Document June 2025",
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
    output_path = "/root/e173_go_gateway/data/morocco_mobile_prefixes_correct.json"
    
    print("Extracting Morocco mobile prefixes from ANRT PDF...")
    operators = extract_morocco_prefixes(pdf_path)
    save_to_json(operators, output_path)
    print(f"\nData saved to: {output_path}")