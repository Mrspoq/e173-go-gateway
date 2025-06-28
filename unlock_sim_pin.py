#!/usr/bin/env python3
"""
Script to unlock SIM cards with PIN codes via Asterisk CLI
"""

import subprocess
import time
import re

def run_asterisk_command(command):
    """Run command in Asterisk CLI and return output"""
    full_cmd = f"asterisk_private/sbin/asterisk -rx '{command}' -s /root/e173_go_gateway/asterisk_private/var/run/asterisk/asterisk.ctl"
    try:
        result = subprocess.run(full_cmd, shell=True, capture_output=True, text=True, timeout=5)
        return result.stdout.strip()
    except subprocess.TimeoutExpired:
        return "Command timed out"
    except Exception as e:
        return f"Error: {str(e)}"

def check_pin_status(dongle_id):
    """Check if a dongle requires PIN"""
    output = run_asterisk_command(f"dongle show device state {dongle_id}")
    if "PIN" in output:
        return True
    return False

def unlock_with_pin(dongle_id, pin):
    """Try to unlock a dongle with a PIN"""
    print(f"\nTrying PIN {pin} for dongle {dongle_id}...")
    
    # Send PIN unlock command
    cmd = f"dongle cmd {dongle_id} AT+CPIN={pin}"
    result = run_asterisk_command(cmd)
    print(f"PIN unlock result: {result}")
    
    # Wait a bit for the SIM to process
    time.sleep(3)
    
    # Check if unlock was successful
    status = run_asterisk_command(f"dongle show device state {dongle_id}")
    if "Free" in status or "GSM" in status:
        return True
    return False

def main():
    print("=== SIM PIN Unlock Script ===\n")
    
    # Known PIN for Labara
    labara_pin = "2525"
    
    # Get list of all dongles
    output = run_asterisk_command("dongle show devices")
    lines = output.split('\n')
    
    dongles = []
    for line in lines[1:]:  # Skip header
        parts = line.split()
        if len(parts) >= 3:
            dongle_id = parts[0]
            state = parts[2]
            dongles.append((dongle_id, state))
    
    print(f"Found {len(dongles)} dongles")
    
    # Try to unlock each dongle that needs PIN
    for dongle_id, state in dongles:
        print(f"\n--- Checking {dongle_id} (State: {state}) ---")
        
        # Get detailed state
        detailed_state = run_asterisk_command(f"dongle show device state {dongle_id}")
        
        # Check if it needs PIN
        if "PIN" in detailed_state or "Not connec" in state:
            print(f"Dongle {dongle_id} may need PIN unlock")
            
            # Try the Labara PIN
            if unlock_with_pin(dongle_id, labara_pin):
                print(f"✓ Successfully unlocked {dongle_id} with PIN {labara_pin}")
                
                # Get phone number after unlock
                time.sleep(2)
                devices_output = run_asterisk_command("dongle show devices")
                for line in devices_output.split('\n'):
                    if dongle_id in line:
                        print(f"Updated status: {line}")
                        
                # Try to get phone number via SMS
                print("\nTrying to get phone number via SMS...")
                sms_cmd = f"dongle sms {dongle_id} 20344 'BAL'"
                result = run_asterisk_command(sms_cmd)
                print(f"SMS result: {result}")
                
                # Check for incoming SMS
                time.sleep(5)
                print("\nChecking for response SMS...")
                inbox_cmd = f"dongle cmd {dongle_id} AT+CMGL=ALL"
                inbox = run_asterisk_command(inbox_cmd)
                print(f"Inbox check: {inbox[:200]}...")
                
            else:
                print(f"✗ Failed to unlock {dongle_id} with PIN {labara_pin}")
        else:
            print(f"Dongle {dongle_id} doesn't need PIN unlock")
    
    print("\n=== Final Status ===")
    final_status = run_asterisk_command("dongle show devices")
    print(final_status)

if __name__ == "__main__":
    main()
