#!/usr/bin/env python3
import serial
import serial.tools.list_ports
import time

def send_at_command(port, command, timeout=1):
    try:
        ser = serial.Serial(port, 115200, timeout=timeout)
        ser.write((command + '\r\n').encode())
        time.sleep(0.5)
        response = ser.read(ser.in_waiting).decode('utf-8', errors='ignore')
        ser.close()
        return response
    except Exception as e:
        return None

# Find all USB serial ports
ports = [p.device for p in serial.tools.list_ports.comports() if 'USB' in p.device]
print(f"Found {len(ports)} USB serial ports")

# Group dongles by IMEI
imei_map = {}

for port in ports:
    print(f"\nTesting {port}...")
    
    # Test if it responds to AT commands
    response = send_at_command(port, 'AT')
    if response and 'OK' in response:
        print(f"  {port} responds to AT commands")
        
        # Get IMEI
        imei_response = send_at_command(port, 'AT+CGSN')
        if imei_response:
            lines = imei_response.strip().split('\n')
            for line in lines:
                line = line.strip()
                if line.isdigit() and len(line) == 15:
                    print(f"  IMEI: {line}")
                    if line not in imei_map:
                        imei_map[line] = {'imei': line, 'ports': []}
                    imei_map[line]['ports'].append(port)
                    break

# Now figure out which ports are data and audio
dongles = []
for imei, info in imei_map.items():
    ports_sorted = sorted(info['ports'])
    if len(ports_sorted) >= 2:
        dongle = {
            'imei': imei,
            'data': ports_sorted[0],  # Usually first port is data
            'audio': ports_sorted[1]  # Usually second port is audio
        }
    else:
        # Only one port found, assume it's data
        dongle = {
            'imei': imei,
            'data': ports_sorted[0],
            'audio': None
        }
    dongles.append(dongle)

print(f"\n\nFound {len(dongles)} dongles:")
for i, dongle in enumerate(dongles):
    print(f"\nDongle {i+1}:")
    print(f"  IMEI: {dongle['imei']}")
    print(f"  Data Port: {dongle['data']}")
    print(f"  Audio Port: {dongle['audio'] if dongle['audio'] else 'Not found'}")

# Create dongle.conf configuration
print("\n\nSuggested dongle.conf entries:")
for i, dongle in enumerate(dongles):
    print(f"\n[dongle{i+1}]")
    print(f"imei={dongle['imei']}")
    print(f"data={dongle['data']}")
    if dongle['audio']:
        print(f"audio={dongle['audio']}")
    print(f"context=from-dongle")
