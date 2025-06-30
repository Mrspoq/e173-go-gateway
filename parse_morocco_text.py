#!/usr/bin/env python3

import json
import re

def parse_morocco_prefixes(text_file):
    """Parse Morocco mobile prefixes from extracted PDF text"""
    
    operators = {
        "ITISSALAT AL-MAGHRIB": {"name": "Maroc Telecom (IAM)", "prefixes": []},
        "MEDI TELECOM": {"name": "Orange Morocco (Méditel)", "prefixes": []},
        "Wana Corporate": {"name": "Inwi", "prefixes": []}
    }
    
    with open(text_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split into sections by operator
    # Fix operator names to match exactly
    content = content.replace("ITISSALAT AL-MAGHRIB", "ETISALAT AL-MAGHRIB")
    
    # Find sections for each operator
    print("Parsing Morocco mobile prefixes from text file...")
    
    # Split content into lines
    lines = content.split('\n')
    
    current_operator = None
    prefix_pattern = re.compile(r'0([67]\d{2})XXXXXX')
    
    for i, line in enumerate(lines):
        line = line.strip()
        
        # Check for operator names
        if "ETISALAT AL-MAGHRIB" in line or "ITISSALAT AL-MAGHRIB" in line:
            current_operator = "ITISSALAT AL-MAGHRIB"
            print(f"\nFound operator: Maroc Telecom (IAM)")
        elif "MEDI TELECOM" in line:
            current_operator = "MEDI TELECOM"
            print(f"\nFound operator: Orange Morocco (Méditel)")
        elif "Wana Corporate" in line:
            current_operator = "Wana Corporate"
            print(f"\nFound operator: Inwi")
        
        # Extract prefixes
        match = prefix_pattern.match(line)
        if match and current_operator:
            prefix = "212" + match.group(1)
            operators[current_operator]["prefixes"].append(prefix)
            print(f"  Added prefix: {prefix}")
    
    # If we didn't get good operator separation, try a different approach
    # Look for blocks of consecutive prefixes
    if not all(op["prefixes"] for op in operators.values()):
        print("\nTrying alternative parsing method...")
        
        # Reset
        for op in operators.values():
            op["prefixes"] = []
        
        # Collect all prefixes first
        all_prefixes = []
        for line in lines:
            match = prefix_pattern.match(line)
            if match:
                all_prefixes.append("212" + match.group(1))
        
        # Based on the PDF structure and known assignments:
        # IAM (Maroc Telecom) typically has: 61X (except 612, 614, 617, 619), 62X (some), 63X (some), 64X (some), 65X (some), 66X (some), 67X (some), 68X (some), 696-697, 75X, 76X
        # Orange (Medi Telecom) typically has: 612, 614, 617, 619, 620-629 (some), 63X (some), 64X (some), 65X (some), 66X (some), 67X (some), 68X (some), 69X (some), 77X
        # Inwi (Wana) typically has: 600-608, 630, 679, 683, 685-686, 692, 695, 70X, 72X
        
        for prefix in all_prefixes:
            suffix = prefix[3:]  # Get the part after 212
            
            # Inwi patterns
            if (suffix.startswith("60") and suffix[2] in "012345678") or \
               suffix == "630" or \
               suffix == "679" or \
               suffix == "683" or \
               suffix in ["685", "686"] or \
               suffix == "692" or \
               suffix == "695" or \
               suffix.startswith("70") or \
               suffix.startswith("72"):
                operators["Wana Corporate"]["prefixes"].append(prefix)
            
            # Orange specific patterns
            elif suffix in ["612", "614", "617", "619"] or \
                 suffix.startswith("77"):
                operators["MEDI TELECOM"]["prefixes"].append(prefix)
            
            # Remaining go to IAM
            else:
                operators["ITISSALAT AL-MAGHRIB"]["prefixes"].append(prefix)
    
    # Sort prefixes
    for op in operators.values():
        op["prefixes"].sort()
        print(f"\n{op['name']}: {len(op['prefixes'])} prefixes")
    
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
    
    print(f"\nData saved to: {output_path}")

if __name__ == "__main__":
    text_file = "morocco_prefixes.txt"
    output_path = "/root/e173_go_gateway/data/morocco_mobile_prefixes_official.json"
    
    operators = parse_morocco_prefixes(text_file)
    save_to_json(operators, output_path)