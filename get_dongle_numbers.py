#!/usr/bin/env python3
"""
Simple script to get phone numbers from USB dongles for OTP purposes
"""

import serial
import serial.tools.list_ports
import time
import re

def send_at_command(port, command, timeout=5):
    """Send AT command and return response"""
    try:
        ser = serial.Serial(port, 115200, timeout=timeout)
        ser.flushInput()
        ser.flushOutput()
        
        # Send command
        ser.write(f"{command}\r\n".encode())
        time.sleep(1)
        
        # Read response
        response = ser.read(1000).decode('utf-8', errors='ignore')
        ser.close()
        return response.strip()
    except Exception as e:
        return f"Error: {str(e)}"

def get_operator_from_imsi(imsi):
    """Identify operator from IMSI"""
    if not imsi or len(imsi) < 6:
        return "Unknown"
    
    # Common Spanish operators (MCC 214)
    mcc_mnc = imsi[:5]
    operators = {
        "21401": "Vodafone ES",
        "21403": "Orange ES", 
        "21404": "Yoigo ES",
        "21405": "Movistar ES",
        "21406": "Vodafone ES",
        "21407": "Movistar ES",
        "21408": "Euskaltel",
        "21409": "Orange ES",
        "21415": "BT España",
        "21416": "Telecable",
        "21417": "R Cable",
        "21418": "ONO",
        "21419": "Simyo",
        "21420": "Fonyou",
        "21421": "Jazztel",
        "21422": "Vectone",
        "21423": "Barablu",
        "21424": "Eroski Móvil"
    }
    
    return operators.get(mcc_mnc, f"Unknown ({mcc_mnc})")

def get_phone_number_commands():
    """Return different AT commands to get phone number based on operator"""
    return [
        "AT+CNUM",           # Standard command
        "AT+CPBS=\"ON\"",    # Phone book selection
        "AT+CPBR=1",         # Read phonebook entry 1
        "AT*CPBS=\"ON\"",    # Alternative phonebook
        "AT*CPBR=1",         # Alternative read
        "AT+CUSD=1,\"*135#\"", # USSD for some operators
        "AT+CUSD=1,\"*123#\"", # Alternative USSD
        "AT+CUSD=1,\"*100#\"", # Another USSD
    ]

def main():
    print("=== Dongle Phone Number Detection ===\n")
    
    # Find USB serial ports
    ports = list(serial.tools.list_ports.comports())
    usb_ports = [p.device for p in ports if 'USB' in p.device]
    
    print(f"Found {len(usb_ports)} USB serial ports")
    
    dongles = []
    
    # Test each port
    for port in sorted(usb_ports):
        print(f"\nTesting {port}...")
        
        # Test if it responds to AT
        response = send_at_command(port, "AT")
        if "OK" not in response:
            print(f"  No AT response")
            continue
            
        print(f"  ✓ Responds to AT commands")
        
        # Get IMEI
        imei_response = send_at_command(port, "AT+CGSN")
        imei = None
        for line in imei_response.split('\n'):
            line = line.strip()
            if line.isdigit() and len(line) == 15:
                imei = line
                break
        
        # Get IMSI  
        imsi_response = send_at_command(port, "AT+CIMI")
        imsi = None
        for line in imsi_response.split('\n'):
            line = line.strip()
            if line.isdigit() and len(line) >= 15:
                imsi = line
                break
        
        operator = get_operator_from_imsi(imsi) if imsi else "Unknown"
        
        # Try to get phone number
        phone_number = "Unknown"
        for cmd in get_phone_number_commands():
            try:
                response = send_at_command(port, cmd, timeout=10)
                
                # Look for phone number patterns
                phone_patterns = [
                    r'\+34\d{9}',  # Spanish format
                    r'34\d{9}',    # Without +
                    r'6\d{8}',     # Mobile without country code
                    r'7\d{8}',     # Alternative mobile
                ]
                
                for pattern in phone_patterns:
                    matches = re.findall(pattern, response)
                    if matches:
                        phone_number = matches[0]
                        if not phone_number.startswith('+'):
                            if phone_number.startswith('34'):
                                phone_number = '+' + phone_number
                            elif phone_number.startswith('6') or phone_number.startswith('7'):
                                phone_number = '+34' + phone_number
                        break
                
                if phone_number != "Unknown":
                    break
                    
            except Exception as e:
                continue
        
        dongle_info = {
            'port': port,
            'imei': imei or "Unknown",
            'imsi': imsi or "Unknown", 
            'operator': operator,
            'number': phone_number
        }
        
        dongles.append(dongle_info)
        
        print(f"  IMEI: {imei}")
        print(f"  IMSI: {imsi}")
        print(f"  Operator: {operator}")
        print(f"  Number: {phone_number}")
    
    print(f"\n=== Summary ===")
    print(f"Found {len(dongles)} working dongles:\n")
    
    for i, dongle in enumerate(dongles, 1):
        print(f"Dongle {i}:")
        print(f"  Port: {dongle['port']}")
        print(f"  Operator: {dongle['operator']}")
        print(f"  Number: {dongle['number']}")
        print(f"  IMEI: {dongle['imei']}")
        print(f"  IMSI: {dongle['imsi']}")
        print()
    
    if dongles:
        print("To check for SMS messages, you can use:")
        print("python3 check_sms.py")

if __name__ == "__main__":
    main()
