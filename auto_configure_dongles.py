#!/usr/bin/env python3
"""
Auto-configure all detected dongles with names D01-D07
"""

import serial
import serial.tools.list_ports
import time
import sys

def test_port_for_imei(port):
    """Test if a port responds to AT commands and get IMEI"""
    try:
        ser = serial.Serial(port, 115200, timeout=2)
        # Clear any pending data
        ser.flushInput()
        ser.flushOutput()
        
        # Test AT command
        ser.write(b'AT\r\n')
        time.sleep(0.5)
        response = ser.read(1000).decode('utf-8', errors='ignore')
        
        if 'OK' not in response:
            ser.close()
            return None, None
            
        # Get IMEI
        ser.write(b'AT+CGSN\r\n')
        time.sleep(0.5)
        imei_response = ser.read(1000).decode('utf-8', errors='ignore')
        
        # Extract IMEI
        lines = imei_response.split('\n')
        imei = None
        for line in lines:
            line = line.strip()
            if line.isdigit() and len(line) == 15:
                imei = line
                break
                
        ser.close()
        return port, imei
        
    except Exception as e:
        return None, None

def main():
    print("=== Auto-configuring Dongles ===\n")
    
    # Find all USB ports
    ports = list(serial.tools.list_ports.comports())
    usb_ports = sorted([p.device for p in ports if 'USB' in p.device])
    
    print(f"Found {len(usb_ports)} USB serial ports")
    print("Testing each port for dongle...\n")
    
    # Test each port
    working_dongles = []
    for port in usb_ports:
        print(f"Testing {port}...", end='', flush=True)
        port_device, imei = test_port_for_imei(port)
        if imei:
            print(f" ✓ IMEI: {imei}")
            working_dongles.append((port_device, imei))
        else:
            print(f" ✗ No response")
    
    print(f"\nFound {len(working_dongles)} working dongles")
    
    # Generate dongle.conf content
    config_content = """[general]

interval=15

[defaults]
context=from-dongle
group=0
rxgain=0
txgain=0
autodeletesms=yes
resetdongle=yes
u2diag=-1
usecallingpres=yes
callingpres=allowed_passed_screen
disablesms=no
language=en
smsaspdu=yes
mindtmfgap=45
mindtmfduration=80
mindtmfinterval=200
callwaiting=auto
disable=no
dtmf=relax
init_watchdog=0
notreg_watchdog=0
sms_watchdog=0
dialing_watchdog=0
roaming_watchdog=0
readsms=yes
read_full_sm=no

"""
    
    # Add each dongle with D01-D07 naming
    for idx, (port, imei) in enumerate(working_dongles, 1):
        dongle_name = f"D{idx:02d}"
        config_content += f"""
; Dongle {idx} - IMEI: {imei}
[{dongle_name}]
data={port}
imei={imei}
context=from-dongle
"""
    
    # Write to dongle.conf
    with open('asterisk_private/etc/asterisk/dongle.conf', 'w') as f:
        f.write(config_content)
    
    print(f"\nWritten configuration for {len(working_dongles)} dongles:")
    for idx, (port, imei) in enumerate(working_dongles, 1):
        print(f"  D{idx:02d}: {port} (IMEI: {imei})")
    
    print("\nConfiguration saved to asterisk_private/etc/asterisk/dongle.conf")
    print("\nReload the module with:")
    print("  asterisk -rx 'module reload chan_dongle.so'")

if __name__ == "__main__":
    main()
