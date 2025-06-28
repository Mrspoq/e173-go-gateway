#!/usr/bin/env python3
"""
Dongle Troubleshooting Script
Helps diagnose and fix issues with USB dongles in Asterisk
"""

import serial
import serial.tools.list_ports
import time
import subprocess
import sys
import os

def run_command(cmd):
    """Run a shell command and return the output"""
    try:
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
        return result.stdout.strip()
    except Exception as e:
        return f"Error: {str(e)}"

def test_port(port):
    """Test if a port responds to AT commands"""
    try:
        ser = serial.Serial(port, 115200, timeout=2)
        ser.write(b'AT\r\n')
        time.sleep(0.5)
        response = ser.read(1000).decode('utf-8', errors='ignore')
        ser.close()
        return 'OK' in response, response
    except Exception as e:
        return False, str(e)

def get_imei(port):
    """Get IMEI from a working port"""
    try:
        ser = serial.Serial(port, 115200, timeout=2)
        ser.write(b'AT+CGSN\r\n')
        time.sleep(0.5)
        response = ser.read(1000).decode('utf-8', errors='ignore')
        ser.close()
        
        # Extract IMEI from response
        lines = response.split('\n')
        for line in lines:
            line = line.strip()
            if line.isdigit() and len(line) == 15:
                return line
        return None
    except:
        return None

def check_port_usage(port):
    """Check if a port is in use by another process"""
    cmd = f"lsof {port} 2>/dev/null"
    output = run_command(cmd)
    if output and output != "Error:":
        return True, output
    return False, None

def get_asterisk_dongles():
    """Get current dongle status from Asterisk"""
    cmd = "asterisk_private/sbin/asterisk -rx 'dongle show devices' -s /root/e173_go_gateway/asterisk_private/var/run/asterisk/asterisk.ctl"
    return run_command(cmd)

def main():
    print("=== Dongle Troubleshooting Tool ===\n")
    
    # 1. List all USB serial ports
    print("1. Detecting USB serial ports...")
    ports = list(serial.tools.list_ports.comports())
    usb_ports = [p for p in ports if 'USB' in p.device]
    
    print(f"Found {len(usb_ports)} USB serial ports:")
    for port in usb_ports:
        print(f"  - {port.device}: {port.description}")
    
    # 2. Test each port
    print("\n2. Testing ports for AT command response...")
    working_ports = []
    for port in usb_ports:
        print(f"\nTesting {port.device}...")
        
        # Check if port is in use
        in_use, usage = check_port_usage(port.device)
        if in_use:
            print(f"  ⚠️  Port is in use by: {usage.split()[0] if usage else 'unknown'}")
        
        # Test AT command
        responds, response = test_port(port.device)
        if responds:
            print(f"  ✓ Responds to AT commands")
            imei = get_imei(port.device)
            if imei:
                print(f"  ✓ IMEI: {imei}")
                working_ports.append((port.device, imei))
            else:
                print(f"  ⚠️  Could not retrieve IMEI")
        else:
            print(f"  ✗ No AT response: {response[:50]}...")
    
    # 3. Check Asterisk dongle status
    print("\n3. Current Asterisk dongle status:")
    asterisk_status = get_asterisk_dongles()
    print(asterisk_status)
    
    # 4. Compare findings
    print("\n4. Analysis:")
    
    # Parse Asterisk dongles
    asterisk_dongles = []
    if asterisk_status:
        lines = asterisk_status.split('\n')
        for line in lines[1:]:  # Skip header
            parts = line.split()
            if len(parts) >= 15:
                dongle_id = parts[0]
                imei = parts[14] if parts[14].isdigit() else None
                state = parts[2]
                asterisk_dongles.append((dongle_id, imei, state))
    
    print(f"\nFound {len(working_ports)} working ports via direct test")
    print(f"Asterisk shows {len(asterisk_dongles)} configured dongles")
    
    # 5. Recommendations
    print("\n5. Recommendations:")
    
    # Find working ports not in Asterisk
    asterisk_imeis = [d[1] for d in asterisk_dongles if d[1]]
    for port, imei in working_ports:
        if imei not in asterisk_imeis:
            print(f"\n✓ Port {port} with IMEI {imei} is working but not configured in Asterisk")
            print(f"  Add to dongle.conf:")
            print(f"  [dongle_{imei[-4:]}]")
            print(f"  data={port}")
            print(f"  imei={imei}")
            print(f"  context=from-dongle")
    
    # Find Asterisk dongles not working
    working_imeis = [imei for _, imei in working_ports]
    for dongle_id, imei, state in asterisk_dongles:
        if imei and imei not in working_imeis and state != "Free":
            print(f"\n⚠️  Dongle {dongle_id} (IMEI: {imei}) is configured but not responding")
            print(f"  Current state: {state}")
            print(f"  Try: asterisk -rx 'dongle reset {dongle_id}'")
    
    # 6. USB reset option
    print("\n\n6. USB Reset Option:")
    print("If dongles are not responding, you can try resetting the USB subsystem:")
    print("  sudo modprobe -r option usb_wwan")
    print("  sudo modprobe option usb_wwan")
    print("\nOr reset specific USB devices:")
    print("  usbreset /dev/bus/usb/XXX/YYY  (use lsusb to find XXX/YYY)")

if __name__ == "__main__":
    main()
