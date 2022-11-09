build:
	mkdir -p build	
	conan install . -if build -s build_type=Release
	cmake -Bbuild -S. -DCMAKE_BUILD_TYPE=Release
	cmake --build build
test:
deploy: