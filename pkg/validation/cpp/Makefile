CXX = g++
CXXFLAGS = -fPIC -O3 -Wall -std=c++11
LDFLAGS = -shared
LIBS = -lphonenumber -lgeocoding -lprotobuf -lboost_system -lboost_thread -licui18n -licuuc -licudata

TARGET = libphonenumber_wrapper.so
SOURCES = phonenumber_wrapper.cpp
OBJECTS = $(SOURCES:.cpp=.o)

all: $(TARGET)

$(TARGET): $(OBJECTS)
	$(CXX) $(LDFLAGS) -o $@ $^ $(LIBS)

%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

clean:
	rm -f $(OBJECTS) $(TARGET)

install: $(TARGET)
	cp $(TARGET) /usr/local/lib/
	ldconfig

.PHONY: all clean install