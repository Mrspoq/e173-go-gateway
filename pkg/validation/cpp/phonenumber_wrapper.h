#ifndef PHONENUMBER_WRAPPER_H
#define PHONENUMBER_WRAPPER_H

#ifdef __cplusplus
extern "C" {
#endif

// Phone validation result structure
typedef struct {
    int is_valid;
    int is_possible;
    int is_mobile;
    char* formatted_number;
    char* country_code;
    char* national_number;
    char* carrier;
    char* region;
    char* number_type;
    char* error_msg;
} PhoneValidationResult;

// Initialize the library (call once at startup)
int phone_lib_init();

// Validate a phone number
// region_code: default region code (e.g., "US", "NG", "MA")
// phone_number: the phone number to validate
// result: pointer to result structure (caller must free strings)
int validate_phone_number(const char* region_code, const char* phone_number, PhoneValidationResult* result);

// Free result strings
void free_validation_result(PhoneValidationResult* result);

// Check if number is valid for region
int is_valid_number(const char* region_code, const char* phone_number);

// Check if number is valid mobile number
int is_valid_mobile_number(const char* region_code, const char* phone_number);

// Format number in international format
char* format_international(const char* region_code, const char* phone_number);

// Get carrier name for number
char* get_carrier(const char* region_code, const char* phone_number);

// Free string returned by library
void free_string(char* str);

#ifdef __cplusplus
}
#endif

#endif // PHONENUMBER_WRAPPER_H