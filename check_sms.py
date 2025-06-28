#!/usr/bin/env python3
"""
Simple script to check SMS messages on USB dongles
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
        response = ser.read(2000).decode('utf-8', errors='ignore')
        ser.close()
        return response.strip()
    except Exception as e:
        return f"Error: {str(e)}"

def parse_sms_list(response):
    """Parse SMS list from AT+CMGL response"""
    messages = []
    lines = response.split('\n')
    
    i = 0
    while i < len(lines):
        line = lines[i].strip()
        # Look for +CMGL: header
        if line.startswith('+CMGL:'):
            try:
                # Parse header
                parts = line.split(',')
                index = parts[0].split(':')[1].strip()
                status = parts[1].strip('"') if len(parts) > 1 else "Unknown"
                sender = parts[2].strip('"') if len(parts) > 2 else "Unknown"
                
                # Get message content (next line)
                i += 1
                if i < len(lines):
                    content = lines[i].strip()
                    if not content.startswith('+CMGL:') and not content.startswith('OK'):
                        messages.append({
                            'index': index,
                            'status': status,
                            'sender': sender,
                            'content': content
                        })
            except:
                pass
        i += 1
    
    return messages

def main():
    print("=== SMS Check on Dongles ===\n")
    
    # Find USB serial ports
    ports = list(serial.tools.list_ports.comports())
    usb_ports = [p.device for p in ports if 'USB' in p.device]
    
    print(f"Found {len(usb_ports)} USB serial ports")
    
    working_dongles = []
    
    # Test each port
    for port in sorted(usb_ports):
        response = send_at_command(port, "AT")
        if "OK" in response:
            working_dongles.append(port)
    
    print(f"Found {len(working_dongles)} working dongles\n")
    
    for i, port in enumerate(working_dongles, 1):
        print(f"\n=== Dongle {i} on {port} ===")
        
        # Set SMS text mode
        send_at_command(port, "AT+CMGF=1")
        
        # Check for unread messages
        print("\nUnread messages:")
        response = send_at_command(port, 'AT+CMGL="REC UNREAD"', timeout=10)
        unread = parse_sms_list(response)
        
        if unread:
            for msg in unread:
                print(f"  From: {msg['sender']}")
                print(f"  Message: {msg['content']}")
                print(f"  Index: {msg['index']}")
                print()
        else:
            print("  No unread messages")
        
        # Check all messages
        print("\nAll messages:")
        response = send_at_command(port, 'AT+CMGL="ALL"', timeout=10)
        all_msgs = parse_sms_list(response)
        
        if all_msgs:
            for msg in all_msgs:
                print(f"  [{msg['status']}] From: {msg['sender']}")
                print(f"  Message: {msg['content']}")
                print(f"  Index: {msg['index']}")
                print()
        else:
            print("  No messages")
        
        # Show delete command
        if all_msgs:
            print(f"\nTo delete a message, use: AT+CMGD={msg['index']}")
            print(f"To delete all messages, use: AT+CMGD=1,4")

if __name__ == "__main__":
    main()
