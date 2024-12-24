#!/bin/sh

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <language> <source_file_path>"
    exit 1
fi

language="$1"
source_file="$2"

if [ ! -f "$source_file" ]; then
    echo "Error: Source file '$source_file' does not exist!" >&2
    exit 1
fi

# C++ function to compile and run the program
run_cpp() {
    executable="a.out"

    # Compile the C++ program and direct any errors to stderr
    g++ "$source_file" -o "$executable" 2>&1
    if [ $? -ne 0 ]; then
        echo "Compilation failed. Please check the error messages above." >&2
        exit 1
    fi

    # Run the compiled executable, redirecting both stdout and stderr to the terminal
    ./$executable 2>&1

    if [ $? -ne 0 ]; then
        echo "Runtime error occurred. Check the output above for details." >&2
    fi

    # Clean up the executable
    rm "$executable"
}

# Go function to run the program
run_go() {
    # Run the Go program and capture both stdout and stderr
    go run "$source_file" 2>&1
    if [ $? -ne 0 ]; then
        echo "Runtime error occurred while running the Go program." >&2
    fi
}

# Check the programming language and call the appropriate function
case "$language" in
    cpp)
        run_cpp
        ;;
    golang)
        run_go
        ;;
    *)
        echo "Error: Unsupported language '$language'. Please use 'cpp' or 'go'." >&2
        exit 1
        ;;
esac
