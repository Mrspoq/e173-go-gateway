#!/usr/bin/env python3
"""
Simple dongle SMS monitor - gets phone numbers and monitors incoming SMS
"""

import serial
import serial.tools.list_ports
import time
import threading
import re
from datetime import datetime

class DongleMonitor:
    def __init__(self):
        self.dongles = []
        self.running = True
        
    def find_dongles(self):
        """Find all connected dongles"""
        print("Scanning for dongles...")
        ports = list(serial.tools.list_ports.comports())
        usb_ports = [p.device for p in ports if 'USB' in p.device]
        
        for port in sorted(usb_ports):
            try:
                ser = serial.Serial(port, 115200, timeout=1)
                ser.write(b'AT\r\n')
                time.sleep(0.2)
                response = ser.read(1000).decode('utf-8', errors='ignore')
                
                if 'OK' in response:
                    # Get IMEI
                    ser.write(b'AT+CGSN\r\n')
                    time.sleep(0.2)
                    imei_response = ser.read(1000).decode('utf-8', errors='ignore')
                    imei = self.extract_number(imei_response, 15)
                    
                    if imei:
                        self.dongles.append({
                            'port': port,
                            'imei': imei,
                            'serial': ser,
                            'phone_number': None
                        })
                        print(f"✓ Found dongle on {port} (IMEI: {imei})")
                else:
                    ser.close()
            except Exception as e:
                pass
        
        print(f"\nFound {len(self.dongles)} working dongles")
        return len(self.dongles) > 0
    
    def extract_number(self, response, length):
        """Extract a numeric string of specific length from response"""
        lines = response.split('\n')
        for line in lines:
            line = line.strip()
            if line.isdigit() and len(line) == length:
                return line
        return None
    
    def get_phone_numbers(self):
        """Get phone numbers from all dongles"""
        print("\nGetting phone numbers...")
        
        for dongle in self.dongles:
            try:
                ser = dongle['serial']
                
                # Method 1: Try CNUM command
                ser.write(b'AT+CNUM\r\n')
                time.sleep(0.5)
                response = ser.read(1000).decode('utf-8', errors='ignore')
                
                # Parse CNUM response
                if '+CNUM:' in response:
                    match = re.search(r'\+CNUM:.*?"(\+?\d+)"', response)
                    if match:
                        dongle['phone_number'] = match.group(1)
                        print(f"Port {dongle['port']}: {dongle['phone_number']}")
                        continue
                
                # Method 2: Try CUSD command to get number
                ser.write(b'AT+CUSD=1,"*#100#"\r\n')
                time.sleep(2)
                response = ser.read(1000).decode('utf-8', errors='ignore')
                
                # Look for phone number pattern
                match = re.search(r'(\+?\d{10,15})', response)
                if match:
                    dongle['phone_number'] = match.group(1)
                    print(f"Port {dongle['port']}: {dongle['phone_number']}")
                else:
                    print(f"Port {dongle['port']}: Number not available (may need SIM PIN)")
                    
            except Exception as e:
                print(f"Error getting number from {dongle['port']}: {e}")
    
    def setup_sms_monitoring(self):
        """Setup SMS text mode and notifications"""
        print("\nSetting up SMS monitoring...")
        
        for dongle in self.dongles:
            try:
                ser = dongle['serial']
                
                # Set SMS to text mode
                ser.write(b'AT+CMGF=1\r\n')
                time.sleep(0.2)
                
                # Enable SMS notifications
                ser.write(b'AT+CNMI=2,1,0,0,0\r\n')
                time.sleep(0.2)
                
                # Clear the buffer
                ser.read(ser.in_waiting)
                
            except Exception as e:
                print(f"Error setting up SMS on {dongle['port']}: {e}")
    
    def monitor_sms(self):
        """Monitor incoming SMS messages"""
        print("\n" + "="*60)
        print("SMS MONITOR STARTED")
        print("="*60)
        print("Waiting for incoming SMS messages...")
        print("Press Ctrl+C to stop\n")
        
        while self.running:
            for dongle in self.dongles:
                try:
                    ser = dongle['serial']
                    if ser.in_waiting:
                        data = ser.read(ser.in_waiting).decode('utf-8', errors='ignore')
                        
                        # Check for SMS notification
                        if '+CMTI:' in data:
                            # Extract SMS index
                            match = re.search(r'\+CMTI:.*?,(\d+)', data)
                            if match:
                                sms_index = match.group(1)
                                self.read_sms(dongle, sms_index)
                        
                except Exception as e:
                    pass
            
            time.sleep(0.1)
    
    def read_sms(self, dongle, index):
        """Read a specific SMS message"""
        try:
            ser = dongle['serial']
            
            # Read the SMS
            ser.write(f'AT+CMGR={index}\r\n'.encode())
            time.sleep(0.5)
            response = ser.read(1000).decode('utf-8', errors='ignore')
            
            # Parse SMS
            if '+CMGR:' in response:
                lines = response.split('\n')
                for i, line in enumerate(lines):
                    if '+CMGR:' in line:
                        # Extract sender info
                        match = re.search(r'"([^"]+)".*?"(\+?\d+)"', line)
                        if match and i + 1 < len(lines):
                            sender = match.group(2)
                            message = lines[i + 1].strip()
                            
                            print("\n" + "="*60)
                            print(f"📱 NEW SMS RECEIVED")
                            print(f"Port: {dongle['port']}")
                            print(f"To: {dongle['phone_number'] or 'Unknown'}")
                            print(f"From: {sender}")
                            print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
                            print(f"Message: {message}")
                            print("="*60 + "\n")
                            
                            # Delete the SMS after reading
                            ser.write(f'AT+CMGD={index}\r\n'.encode())
                            time.sleep(0.2)
                            
        except Exception as e:
            print(f"Error reading SMS: {e}")
    
    def cleanup(self):
        """Close all serial connections"""
        for dongle in self.dongles:
            try:
                dongle['serial'].close()
            except:
                pass
    
    def run(self):
        """Main run loop"""
        try:
            if not self.find_dongles():
                print("No dongles found!")
                return
            
            self.get_phone_numbers()
            self.setup_sms_monitoring()
            self.monitor_sms()
            
        except KeyboardInterrupt:
            print("\n\nStopping SMS monitor...")
        finally:
            self.running = False
            self.cleanup()

if __name__ == "__main__":
    monitor = DongleMonitor()
    monitor.run()
